use serde::{Deserialize, Serialize};
use std::path::Path;
use anyhow::Result;
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AgentConfig {
    pub agent: AgentSettings,
    pub control_plane: ControlPlaneConfig,
    pub data_plane: DataPlaneConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AgentSettings {
    pub id: Option<Uuid>,
    pub name: String,
    pub tags: Vec<String>,
    pub heartbeat_interval: u64,
    pub command_timeout: u64,
    pub enable_http_server: Option<bool>,
    pub http_port: Option<u16>,
    pub log_level: Option<String>,
    pub max_log_files: Option<u32>,
    pub log_file_size_mb: Option<u32>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ControlPlaneConfig {
    pub enabled: bool,
    pub server_url: String,
    pub api_key: Option<String>,
    pub poll_interval: u64,
    pub max_concurrent_commands: usize,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DataPlaneConfig {
    pub enabled: bool,
    pub collectors: Vec<CollectorConfig>,
    pub exporters: Vec<ExporterConfig>,
    pub buffer_size: usize,
    pub flush_interval: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CollectorConfig {
    pub name: String,
    pub collector_type: String,
    pub config: serde_json::Value,
    pub enabled: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExporterConfig {
    pub name: String,
    pub exporter_type: String,
    pub endpoint: String,
    pub headers: std::collections::HashMap<String, String>,
    pub batch_size: usize,
    pub enabled: bool,
}

impl AgentConfig {
    pub fn load<P: AsRef<Path>>(path: P) -> Result<Self> {
        let content = std::fs::read_to_string(path)?;
        let config: AgentConfig = toml::from_str(&content)?;
        Ok(config)
    }

    pub fn generate_default<P: AsRef<Path>>(path: P) -> Result<()> {
        let default_config = Self::default();
        let toml_content = toml::to_string_pretty(&default_config)?;
        std::fs::write(path, toml_content)?;
        Ok(())
    }
}

impl Default for AgentConfig {
    fn default() -> Self {
        let mut headers = std::collections::HashMap::new();
        headers.insert("Content-Type".to_string(), "application/json".to_string());

        Self {
            agent: AgentSettings {
                id: Some(Uuid::new_v4()),
                name: "default-agent".to_string(),
                tags: vec!["production".to_string()],
                heartbeat_interval: 30,
                command_timeout: 300,
                enable_http_server: Some(true),
                http_port: Some(8080),
                log_level: Some("info".to_string()),
                max_log_files: Some(10),
                log_file_size_mb: Some(100),
            },
            control_plane: ControlPlaneConfig {
                enabled: true,
                server_url: "http://localhost:8080".to_string(),
                api_key: None,
                poll_interval: 10,
                max_concurrent_commands: 5,
            },
            data_plane: DataPlaneConfig {
                enabled: true,
                collectors: vec![
                    CollectorConfig {
                        name: "system_metrics".to_string(),
                        collector_type: "system".to_string(),
                        config: serde_json::json!({
                            "interval": 30,
                            "metrics": ["cpu", "memory", "disk"]
                        }),
                        enabled: true,
                    }
                ],
                exporters: vec![
                    ExporterConfig {
                        name: "http_exporter".to_string(),
                        exporter_type: "http".to_string(),
                        endpoint: "http://localhost:8081/metrics".to_string(),
                        headers,
                        batch_size: 100,
                        enabled: true,
                    }
                ],
                buffer_size: 1000,
                flush_interval: 60,
            },
        }
    }
}