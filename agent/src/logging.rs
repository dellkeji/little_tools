use crate::config::AgentSettings;
use crate::error::{AgentError, AgentResult};
use log::LevelFilter;
use std::fs::OpenOptions;
use std::io::Write;
use std::path::Path;

pub struct LoggingConfig {
    pub level: LevelFilter,
    pub file_path: Option<String>,
    pub max_files: u32,
    pub max_file_size_mb: u32,
    pub enable_console: bool,
    pub enable_json: bool,
}

impl From<&AgentSettings> for LoggingConfig {
    fn from(settings: &AgentSettings) -> Self {
        let level = match settings.log_level.as_deref().unwrap_or("info") {
            "trace" => LevelFilter::Trace,
            "debug" => LevelFilter::Debug,
            "info" => LevelFilter::Info,
            "warn" => LevelFilter::Warn,
            "error" => LevelFilter::Error,
            "off" => LevelFilter::Off,
            _ => LevelFilter::Info,
        };

        Self {
            level,
            file_path: Some("agent.log".to_string()),
            max_files: settings.max_log_files.unwrap_or(10),
            max_file_size_mb: settings.log_file_size_mb.unwrap_or(100),
            enable_console: true,
            enable_json: false,
        }
    }
}

pub fn init_logging(config: LoggingConfig) -> AgentResult<()> {
    let mut builder = env_logger::Builder::new();
    builder.filter_level(config.level);

    if config.enable_json {
        builder.format(|buf, record| {
            let timestamp = chrono::Utc::now().to_rfc3339();
            writeln!(
                buf,
                r#"{{"timestamp":"{}","level":"{}","target":"{}","message":"{}"}}"#,
                timestamp,
                record.level(),
                record.target(),
                record.args()
            )
        });
    } else {
        builder.format(|buf, record| {
            let timestamp = chrono::Local::now().format("%Y-%m-%d %H:%M:%S%.3f");
            writeln!(
                buf,
                "[{}] {} [{}:{}] {}",
                timestamp,
                record.level(),
                record.file().unwrap_or("unknown"),
                record.line().unwrap_or(0),
                record.args()
            )
        });
    }

    // Setup file logging if configured
    if let Some(file_path) = config.file_path {
        setup_file_logging(&file_path, config.max_file_size_mb, config.max_files)?;
    }

    builder.init();
    Ok(())
}

fn setup_file_logging(
    file_path: &str,
    max_size_mb: u32,
    max_files: u32,
) -> AgentResult<()> {
    // Create log directory if it doesn't exist
    if let Some(parent) = Path::new(file_path).parent() {
        std::fs::create_dir_all(parent)
            .map_err(|e| AgentError::ConfigError(format!("Failed to create log directory: {}", e)))?;
    }

    // Rotate logs if current file is too large
    if let Ok(metadata) = std::fs::metadata(file_path) {
        let size_mb = metadata.len() / (1024 * 1024);
        if size_mb >= max_size_mb as u64 {
            rotate_log_files(file_path, max_files)?;
        }
    }

    Ok(())
}

fn rotate_log_files(base_path: &str, max_files: u32) -> AgentResult<()> {
    // Move existing log files
    for i in (1..max_files).rev() {
        let old_path = if i == 1 {
            base_path.to_string()
        } else {
            format!("{}.{}", base_path, i - 1)
        };
        
        let new_path = format!("{}.{}", base_path, i);
        
        if Path::new(&old_path).exists() {
            std::fs::rename(&old_path, &new_path)
                .map_err(|e| AgentError::ConfigError(format!("Failed to rotate log file: {}", e)))?;
        }
    }

    Ok(())
}

// Structured logging macros
#[macro_export]
macro_rules! log_with_context {
    ($level:expr, $context:expr, $($arg:tt)*) => {
        log::log!($level, "[{}] {}", $context, format!($($arg)*));
    };
}

#[macro_export]
macro_rules! info_with_context {
    ($context:expr, $($arg:tt)*) => {
        $crate::log_with_context!(log::Level::Info, $context, $($arg)*);
    };
}

#[macro_export]
macro_rules! warn_with_context {
    ($context:expr, $($arg:tt)*) => {
        $crate::log_with_context!(log::Level::Warn, $context, $($arg)*);
    };
}

#[macro_export]
macro_rules! error_with_context {
    ($context:expr, $($arg:tt)*) => {
        $crate::log_with_context!(log::Level::Error, $context, $($arg)*);
    };
}

// Audit logging for security events
pub struct AuditLogger {
    file_path: String,
}

impl AuditLogger {
    pub fn new(file_path: String) -> Self {
        Self { file_path }
    }

    pub fn log_command_execution(
        &self,
        user: &str,
        command: &str,
        args: &[String],
        success: bool,
    ) -> AgentResult<()> {
        let entry = serde_json::json!({
            "timestamp": chrono::Utc::now().to_rfc3339(),
            "event_type": "command_execution",
            "user": user,
            "command": command,
            "args": args,
            "success": success,
            "pid": std::process::id(),
        });

        self.write_audit_entry(&entry)
    }

    pub fn log_file_access(
        &self,
        user: &str,
        file_path: &str,
        operation: &str,
        success: bool,
    ) -> AgentResult<()> {
        let entry = serde_json::json!({
            "timestamp": chrono::Utc::now().to_rfc3339(),
            "event_type": "file_access",
            "user": user,
            "file_path": file_path,
            "operation": operation,
            "success": success,
            "pid": std::process::id(),
        });

        self.write_audit_entry(&entry)
    }

    pub fn log_authentication(
        &self,
        user: &str,
        method: &str,
        success: bool,
        source_ip: Option<&str>,
    ) -> AgentResult<()> {
        let entry = serde_json::json!({
            "timestamp": chrono::Utc::now().to_rfc3339(),
            "event_type": "authentication",
            "user": user,
            "method": method,
            "success": success,
            "source_ip": source_ip,
            "pid": std::process::id(),
        });

        self.write_audit_entry(&entry)
    }

    fn write_audit_entry(&self, entry: &serde_json::Value) -> AgentResult<()> {
        let mut file = OpenOptions::new()
            .create(true)
            .append(true)
            .open(&self.file_path)
            .map_err(|e| AgentError::ConfigError(format!("Failed to open audit log: {}", e)))?;

        writeln!(file, "{}", entry)
            .map_err(|e| AgentError::ConfigError(format!("Failed to write audit log: {}", e)))?;

        file.flush()
            .map_err(|e| AgentError::ConfigError(format!("Failed to flush audit log: {}", e)))?;

        Ok(())
    }
}