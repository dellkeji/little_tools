use crate::config::ControlPlaneConfig;
use crate::error::{AgentError, AgentResult, RetryManager};
use crate::security::SecurityValidator;
use crate::types::{Command, CommandResult, CommandType, AgentInfo};
use crate::platform::{execute_command, deploy_file, write_config};
use anyhow::Result;
use log::{info, error, warn};
use reqwest::Client;
use serde_json::Value;
use std::sync::Arc;
use std::time::{Duration, Instant};
use tokio::time::sleep;
use uuid::Uuid;

#[derive(Clone)]
pub struct ControlPlane {
    config: ControlPlaneConfig,
    client: Client,
    agent_info: AgentInfo,
    security_validator: Arc<SecurityValidator>,
    retry_manager: Arc<RetryManager>,
}

impl ControlPlane {
    pub fn new(
        config: ControlPlaneConfig, 
        agent_info: AgentInfo,
        security_validator: Arc<SecurityValidator>,
        retry_manager: Arc<RetryManager>,
    ) -> Self {
        let mut headers = reqwest::header::HeaderMap::new();
        
        if let Some(api_key) = &config.api_key {
            headers.insert(
                "Authorization",
                format!("Bearer {}", api_key).parse().unwrap(),
            );
        }
        
        let client = Client::builder()
            .default_headers(headers)
            .timeout(Duration::from_secs(30))
            .build()
            .unwrap();

        Self {
            config,
            client,
            agent_info,
            security_validator,
            retry_manager,
        }
    }

    pub async fn start(&self) -> AgentResult<()> {
        info!("Starting control plane");
        
        // Register agent with control server using retry logic
        self.retry_manager.execute_with_retry(|| {
            let rt = tokio::runtime::Handle::current();
            rt.block_on(self.register_agent())
                .map_err(|e| AgentError::ControlPlaneError(e.to_string()))
        }).await?;
        
        // Start command polling loop
        let mut interval = tokio::time::interval(Duration::from_secs(self.config.poll_interval));
        
        loop {
            interval.tick().await;
            
            // Skip if circuit breaker is open
            if self.retry_manager.is_circuit_open() {
                warn!("Control plane circuit breaker is open, skipping operations");
                tokio::time::sleep(Duration::from_secs(30)).await;
                continue;
            }
            
            if let Err(e) = self.poll_commands().await {
                error!("Error polling commands: {}", e);
            }
            
            if let Err(e) = self.send_heartbeat().await {
                error!("Error sending heartbeat: {}", e);
            }
        }
    }

    async fn register_agent(&self) -> Result<()> {
        let url = format!("{}/api/agents/register", self.config.server_url);
        
        let response = self.client
            .post(&url)
            .json(&self.agent_info)
            .send()
            .await?;

        if response.status().is_success() {
            info!("Agent registered successfully");
        } else {
            warn!("Failed to register agent: {}", response.status());
        }

        Ok(())
    }

    async fn send_heartbeat(&self) -> Result<()> {
        let url = format!("{}/api/agents/{}/heartbeat", self.config.server_url, self.agent_info.id);
        
        let response = self.client
            .post(&url)
            .json(&self.agent_info)
            .send()
            .await?;

        if !response.status().is_success() {
            warn!("Failed to send heartbeat: {}", response.status());
        }

        Ok(())
    }

    async fn poll_commands(&self) -> Result<()> {
        let url = format!("{}/api/agents/{}/commands", self.config.server_url, self.agent_info.id);
        
        let response = self.client
            .get(&url)
            .send()
            .await?;

        if response.status().is_success() {
            let commands: Vec<Command> = response.json().await?;
            
            for command in commands {
                info!("Received command: {:?}", command.id);
                
                let result = self.execute_command(command.clone()).await;
                
                if let Err(e) = self.send_command_result(&command.id, result).await {
                    error!("Failed to send command result: {}", e);
                }
            }
        }

        Ok(())
    }

    async fn execute_command(&self, command: Command) -> CommandResult {
        let start_time = Instant::now();
        
        let (success, output, error) = match command.command_type {
            CommandType::Execute => self.handle_execute_command(&command.payload).await,
            CommandType::Deploy => self.handle_deploy_command(&command.payload).await,
            CommandType::Configure => self.handle_configure_command(&command.payload).await,
            CommandType::Monitor => self.handle_monitor_command(&command.payload).await,
            CommandType::Stop => self.handle_stop_command(&command.payload).await,
        };

        let execution_time = start_time.elapsed().as_millis() as u64;

        CommandResult {
            command_id: command.id,
            success,
            output,
            error,
            execution_time,
        }
    }

    async fn handle_execute_command(&self, payload: &Value) -> (bool, String, Option<String>) {
        let cmd = payload.get("command").and_then(|v| v.as_str()).unwrap_or("");
        let args: Vec<&str> = payload.get("args")
            .and_then(|v| v.as_array())
            .map(|arr| arr.iter().filter_map(|v| v.as_str()).collect())
            .unwrap_or_default();

        // Security validation
        if let Err(e) = self.security_validator.validate_command(cmd) {
            return (false, String::new(), Some(format!("Security validation failed: {}", e)));
        }

        match execute_command(cmd, &args) {
            Ok(output) => {
                let stdout = String::from_utf8_lossy(&output.stdout);
                let stderr = String::from_utf8_lossy(&output.stderr);
                
                if output.status.success() {
                    (true, stdout.to_string(), None)
                } else {
                    (false, stdout.to_string(), Some(stderr.to_string()))
                }
            }
            Err(e) => (false, String::new(), Some(e.to_string())),
        }
    }

    async fn handle_deploy_command(&self, payload: &Value) -> (bool, String, Option<String>) {
        let source = payload.get("source").and_then(|v| v.as_str()).unwrap_or("");
        let destination = payload.get("destination").and_then(|v| v.as_str()).unwrap_or("");

        // Security validation for paths
        if let Err(e) = self.security_validator.validate_path(source) {
            return (false, String::new(), Some(format!("Source path validation failed: {}", e)));
        }
        
        if let Err(e) = self.security_validator.validate_path(destination) {
            return (false, String::new(), Some(format!("Destination path validation failed: {}", e)));
        }

        match deploy_file(source, destination) {
            Ok(_) => (true, format!("File deployed from {} to {}", source, destination), None),
            Err(e) => (false, String::new(), Some(e.to_string())),
        }
    }

    async fn handle_configure_command(&self, payload: &Value) -> (bool, String, Option<String>) {
        let path = payload.get("path").and_then(|v| v.as_str()).unwrap_or("");
        let config = payload.get("config").unwrap_or(&Value::Null);

        // Security validation for configuration path
        if let Err(e) = self.security_validator.validate_path(path) {
            return (false, String::new(), Some(format!("Configuration path validation failed: {}", e)));
        }

        match write_config(path, config) {
            Ok(_) => (true, format!("Configuration written to {}", path), None),
            Err(e) => (false, String::new(), Some(e.to_string())),
        }
    }

    async fn handle_monitor_command(&self, payload: &Value) -> (bool, String, Option<String>) {
        // Handle monitoring configuration updates
        info!("Updating monitoring configuration: {:?}", payload);
        (true, "Monitoring configuration updated".to_string(), None)
    }

    async fn handle_stop_command(&self, payload: &Value) -> (bool, String, Option<String>) {
        let process_name = payload.get("process").and_then(|v| v.as_str()).unwrap_or("");
        
        #[cfg(target_os = "windows")]
        let result = execute_command("taskkill", &["/F", "/IM", process_name]);
        
        #[cfg(unix)]
        let result = execute_command("pkill", &[process_name]);

        match result {
            Ok(output) => {
                if output.status.success() {
                    (true, format!("Process {} stopped", process_name), None)
                } else {
                    let stderr = String::from_utf8_lossy(&output.stderr);
                    (false, String::new(), Some(stderr.to_string()))
                }
            }
            Err(e) => (false, String::new(), Some(e.to_string())),
        }
    }

    async fn send_command_result(&self, command_id: &Uuid, result: CommandResult) -> Result<()> {
        let url = format!("{}/api/commands/{}/result", self.config.server_url, command_id);
        
        let response = self.client
            .post(&url)
            .json(&result)
            .send()
            .await?;

        if response.status().is_success() {
            info!("Command result sent successfully for command: {}", command_id);
        } else {
            warn!("Failed to send command result: {}", response.status());
        }

        Ok(())
    }
}