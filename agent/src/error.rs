use serde::{Deserialize, Serialize};
use std::fmt;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AgentError {
    ConfigError(String),
    NetworkError(String),
    SecurityError(String),
    PlatformError(String),
    CommandError(String),
    DataPlaneError(String),
    ControlPlaneError(String),
    HealthCheckError(String),
    ServerError(String),
    ValidationError(String),
    TimeoutError(String),
    AuthenticationError(String),
    PermissionError(String),
    ResourceError(String),
}

impl fmt::Display for AgentError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            AgentError::ConfigError(msg) => write!(f, "Configuration error: {}", msg),
            AgentError::NetworkError(msg) => write!(f, "Network error: {}", msg),
            AgentError::SecurityError(msg) => write!(f, "Security error: {}", msg),
            AgentError::PlatformError(msg) => write!(f, "Platform error: {}", msg),
            AgentError::CommandError(msg) => write!(f, "Command error: {}", msg),
            AgentError::DataPlaneError(msg) => write!(f, "Data plane error: {}", msg),
            AgentError::ControlPlaneError(msg) => write!(f, "Control plane error: {}", msg),
            AgentError::HealthCheckError(msg) => write!(f, "Health check error: {}", msg),
            AgentError::ServerError(msg) => write!(f, "Server error: {}", msg),
            AgentError::ValidationError(msg) => write!(f, "Validation error: {}", msg),
            AgentError::TimeoutError(msg) => write!(f, "Timeout error: {}", msg),
            AgentError::AuthenticationError(msg) => write!(f, "Authentication error: {}", msg),
            AgentError::PermissionError(msg) => write!(f, "Permission error: {}", msg),
            AgentError::ResourceError(msg) => write!(f, "Resource error: {}", msg),
        }
    }
}

impl std::error::Error for AgentError {}

impl From<std::io::Error> for AgentError {
    fn from(err: std::io::Error) -> Self {
        AgentError::PlatformError(err.to_string())
    }
}

impl From<reqwest::Error> for AgentError {
    fn from(err: reqwest::Error) -> Self {
        AgentError::NetworkError(err.to_string())
    }
}

impl From<serde_json::Error> for AgentError {
    fn from(err: serde_json::Error) -> Self {
        AgentError::ConfigError(err.to_string())
    }
}

impl From<toml::de::Error> for AgentError {
    fn from(err: toml::de::Error) -> Self {
        AgentError::ConfigError(err.to_string())
    }
}

impl From<anyhow::Error> for AgentError {
    fn from(err: anyhow::Error) -> Self {
        AgentError::PlatformError(err.to_string())
    }
}

pub type AgentResult<T> = Result<T, AgentError>;

// Error recovery strategies
#[derive(Debug, Clone)]
pub struct ErrorRecovery {
    pub max_retries: u32,
    pub retry_delay_ms: u64,
    pub exponential_backoff: bool,
    pub circuit_breaker_threshold: u32,
}

impl Default for ErrorRecovery {
    fn default() -> Self {
        Self {
            max_retries: 3,
            retry_delay_ms: 1000,
            exponential_backoff: true,
            circuit_breaker_threshold: 5,
        }
    }
}

pub struct RetryManager {
    config: ErrorRecovery,
    failure_count: std::sync::Arc<std::sync::Mutex<u32>>,
}

impl RetryManager {
    pub fn new(config: ErrorRecovery) -> Self {
        Self {
            config,
            failure_count: std::sync::Arc::new(std::sync::Mutex::new(0)),
        }
    }

    pub async fn execute_with_retry<F, T, E>(&self, mut operation: F) -> Result<T, E>
    where
        F: FnMut() -> Result<T, E>,
        E: std::fmt::Debug,
    {
        let mut attempts = 0;
        let mut delay = self.config.retry_delay_ms;

        loop {
            match operation() {
                Ok(result) => {
                    // Reset failure count on success
                    *self.failure_count.lock().unwrap() = 0;
                    return Ok(result);
                }
                Err(err) => {
                    attempts += 1;
                    *self.failure_count.lock().unwrap() += 1;

                    if attempts >= self.config.max_retries {
                        log::error!("Operation failed after {} attempts: {:?}", attempts, err);
                        return Err(err);
                    }

                    log::warn!("Operation failed (attempt {}), retrying in {}ms: {:?}", 
                              attempts, delay, err);

                    tokio::time::sleep(tokio::time::Duration::from_millis(delay)).await;

                    if self.config.exponential_backoff {
                        delay *= 2;
                    }
                }
            }
        }
    }

    pub fn is_circuit_open(&self) -> bool {
        *self.failure_count.lock().unwrap() >= self.config.circuit_breaker_threshold
    }

    pub fn reset_circuit(&self) {
        *self.failure_count.lock().unwrap() = 0;
    }
}