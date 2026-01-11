use crate::config::AgentConfig;
use crate::control_plane::ControlPlane;
use crate::data_plane::DataPlane;
use crate::error::{AgentError, AgentResult, RetryManager, ErrorRecovery};
use crate::health::{HealthChecker, ControlPlaneHealthCheck, DiskSpaceHealthCheck, MemoryHealthCheck};
use crate::metrics::{MetricRegistry, SystemMetricsCollector, AgentMetrics};
use crate::platform::PlatformInfo;
use crate::security::SecurityValidator;
use crate::server::AgentServer;
use crate::types::AgentInfo;
use anyhow::Result;
use chrono::Utc;
use log::{info, error, warn};
use std::sync::Arc;
use std::time::Duration;
use tokio::signal;
use uuid::Uuid;

pub struct Agent {
    config: AgentConfig,
    agent_info: AgentInfo,
    control_plane: Option<ControlPlane>,
    data_plane: Option<DataPlane>,
    health_checker: Arc<HealthChecker>,
    metric_registry: Arc<MetricRegistry>,
    security_validator: Arc<SecurityValidator>,
    retry_manager: Arc<RetryManager>,
    server: Option<AgentServer>,
    shutdown_signal: Option<tokio::sync::oneshot::Receiver<()>>,
}

impl Agent {
    pub async fn new(config: AgentConfig) -> AgentResult<Self> {
        let platform_info = PlatformInfo::new()
            .map_err(|e| AgentError::PlatformError(e.to_string()))?;
        
        let agent_id = config.agent.id.unwrap_or_else(|| Uuid::new_v4());
        
        let agent_info = AgentInfo {
            id: agent_id,
            hostname: platform_info.hostname,
            platform: platform_info.os,
            arch: platform_info.arch,
            version: env!("CARGO_PKG_VERSION").to_string(),
            last_heartbeat: Utc::now(),
        };

        // Initialize core components
        let metric_registry = Arc::new(MetricRegistry::new());
        let health_checker = Arc::new(HealthChecker::new(agent_info.clone()));
        let security_validator = Arc::new(SecurityValidator::new(
            crate::security::SecurityConfig::default()
        ));
        let retry_manager = Arc::new(RetryManager::new(ErrorRecovery::default()));

        // Register health checks
        if config.control_plane.enabled {
            health_checker.register_check(Box::new(
                ControlPlaneHealthCheck::new(config.control_plane.server_url.clone())
            ));
        }
        
        health_checker.register_check(Box::new(
            DiskSpaceHealthCheck::new("/".to_string(), 80.0)
        ));
        
        health_checker.register_check(Box::new(
            MemoryHealthCheck::new(85.0)
        ));

        let control_plane = if config.control_plane.enabled {
            Some(ControlPlane::new(
                config.control_plane.clone(), 
                agent_info.clone(),
                security_validator.clone(),
                retry_manager.clone(),
            ))
        } else {
            None
        };

        let data_plane = if config.data_plane.enabled {
            Some(DataPlane::new(
                config.data_plane.clone(),
                metric_registry.clone(),
            ))
        } else {
            None
        };

        // Initialize HTTP server if enabled
        let server = if config.agent.enable_http_server.unwrap_or(false) {
            let port = config.agent.http_port.unwrap_or(8080);
            Some(AgentServer::new(
                port,
                health_checker.clone(),
                metric_registry.clone(),
                agent_info.clone(),
            ))
        } else {
            None
        };

        Ok(Self {
            config,
            agent_info,
            control_plane,
            data_plane,
            health_checker,
            metric_registry,
            security_validator,
            retry_manager,
            server,
            shutdown_signal: None,
        })
    }

    pub async fn run(&mut self) -> AgentResult<()> {
        info!("Starting agent: {} ({})", self.agent_info.hostname, self.agent_info.id);
        info!("Platform: {} {}", self.agent_info.platform, self.agent_info.arch);

        // Setup graceful shutdown
        let (shutdown_tx, shutdown_rx) = tokio::sync::oneshot::channel();
        self.shutdown_signal = Some(shutdown_rx);

        // Setup signal handlers
        tokio::spawn(async move {
            let _ = signal::ctrl_c().await;
            info!("Received shutdown signal");
            let _ = shutdown_tx.send(());
        });

        // Start HTTP server if configured
        if let Some(ref server) = self.server {
            let server_clone = server.clone();
            tokio::spawn(async move {
                if let Err(e) = server_clone.start().await {
                    error!("HTTP server error: {}", e);
                }
            });
        }

        // Start system metrics collection
        let system_collector = SystemMetricsCollector::new(self.metric_registry.clone());
        let metric_registry_clone = self.metric_registry.clone();
        tokio::spawn(async move {
            let mut interval = tokio::time::interval(Duration::from_secs(30));
            loop {
                interval.tick().await;
                if let Err(e) = system_collector.collect_all().await {
                    error!("Failed to collect system metrics: {}", e);
                }
                
                // Cleanup old metrics
                metric_registry_clone.cleanup_old_metrics(Duration::from_secs(3600));
            }
        });

        // Start agent metrics collection
        let agent_metrics = AgentMetrics::new(self.metric_registry.clone());
        tokio::spawn(async move {
            let mut interval = tokio::time::interval(Duration::from_secs(60));
            loop {
                interval.tick().await;
                agent_metrics.record_uptime();
                agent_metrics.record_memory_usage();
            }
        });

        // Start data plane with error handling
        if let Some(ref mut data_plane) = self.data_plane {
            info!("Starting data plane...");
            let retry_manager = self.retry_manager.clone();
            
            let data_plane_future = retry_manager.execute_with_retry(|| {
                // This would need to be adapted for the actual data plane start method
                Ok(())
            });

            if let Err(e) = data_plane_future.await {
                error!("Failed to start data plane after retries: {:?}", e);
                return Err(AgentError::DataPlaneError("Failed to start data plane".to_string()));
            }

            // Start data plane in background
            tokio::spawn(async move {
                if let Err(e) = data_plane.start().await {
                    error!("Data plane error: {}", e);
                }
            });
        }

        // Start control plane with error handling
        if let Some(ref mut control_plane) = self.control_plane {
            info!("Starting control plane...");
            
            let control_plane_clone = control_plane.clone();
            tokio::spawn(async move {
                loop {
                    match control_plane_clone.start().await {
                        Ok(_) => break,
                        Err(e) => {
                            error!("Control plane error: {}, retrying in 30s...", e);
                            tokio::time::sleep(Duration::from_secs(30)).await;
                        }
                    }
                }
            });
        }

        // Main event loop
        let mut health_check_interval = tokio::time::interval(Duration::from_secs(60));
        
        loop {
            tokio::select! {
                // Handle shutdown signal
                _ = &mut self.shutdown_signal.as_mut().unwrap() => {
                    info!("Shutting down agent gracefully...");
                    self.shutdown().await?;
                    break;
                }
                
                // Periodic health checks
                _ = health_check_interval.tick() => {
                    let health_status = self.health_checker.get_health_status().await;
                    match health_status.status {
                        crate::health::HealthState::Unhealthy => {
                            warn!("Agent health is unhealthy: {:?}", health_status.checks);
                        }
                        crate::health::HealthState::Degraded => {
                            warn!("Agent health is degraded: {:?}", health_status.checks);
                        }
                        _ => {}
                    }
                }
                
                // Handle other periodic tasks
                _ = tokio::time::sleep(Duration::from_secs(1)) => {
                    // Keep the main loop alive
                }
            }
        }

        Ok(())
    }

    async fn shutdown(&mut self) -> AgentResult<()> {
        info!("Initiating graceful shutdown...");

        // Stop data plane
        if let Some(ref mut data_plane) = self.data_plane {
            info!("Stopping data plane...");
            // Add shutdown method to data plane
        }

        // Stop control plane
        if let Some(ref mut control_plane) = self.control_plane {
            info!("Stopping control plane...");
            // Add shutdown method to control plane
        }

        // Final health check and metrics export
        let health_status = self.health_checker.get_health_status().await;
        info!("Final health status: {:?}", health_status.status);

        let metrics = self.metric_registry.get_all_metrics();
        info!("Exported {} metrics before shutdown", metrics.len());

        info!("Agent shutdown completed");
        Ok(())
    }

    pub fn get_agent_info(&self) -> &AgentInfo {
        &self.agent_info
    }

    pub fn get_config(&self) -> &AgentConfig {
        &self.config
    }

    pub fn get_health_checker(&self) -> Arc<HealthChecker> {
        self.health_checker.clone()
    }

    pub fn get_metric_registry(&self) -> Arc<MetricRegistry> {
        self.metric_registry.clone()
    }

    pub fn get_security_validator(&self) -> Arc<SecurityValidator> {
        self.security_validator.clone()
    }
}