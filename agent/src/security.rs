use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
use std::path::Path;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SecurityConfig {
    pub allowed_commands: HashSet<String>,
    pub allowed_paths: HashSet<String>,
    pub max_file_size: u64,
    pub enable_command_whitelist: bool,
    pub enable_path_restriction: bool,
}

impl Default for SecurityConfig {
    fn default() -> Self {
        let mut allowed_commands = HashSet::new();
        allowed_commands.insert("ls".to_string());
        allowed_commands.insert("dir".to_string());
        allowed_commands.insert("ps".to_string());
        allowed_commands.insert("systemctl".to_string());
        allowed_commands.insert("service".to_string());
        
        let mut allowed_paths = HashSet::new();
        allowed_paths.insert("/tmp".to_string());
        allowed_paths.insert("/var/log".to_string());
        allowed_paths.insert("C:\\temp".to_string());
        allowed_paths.insert("C:\\logs".to_string());

        Self {
            allowed_commands,
            allowed_paths,
            max_file_size: 100 * 1024 * 1024, // 100MB
            enable_command_whitelist: true,
            enable_path_restriction: true,
        }
    }
}

pub struct SecurityValidator {
    config: SecurityConfig,
}

impl SecurityValidator {
    pub fn new(config: SecurityConfig) -> Self {
        Self { config }
    }

    pub fn validate_command(&self, command: &str) -> Result<()> {
        if !self.config.enable_command_whitelist {
            return Ok(());
        }

        if self.config.allowed_commands.contains(command) {
            Ok(())
        } else {
            Err(anyhow::anyhow!("Command '{}' is not allowed", command))
        }
    }

    pub fn validate_path(&self, path: &str) -> Result<()> {
        if !self.config.enable_path_restriction {
            return Ok(());
        }

        let path = Path::new(path);
        let path_str = path.to_string_lossy();

        for allowed_path in &self.config.allowed_paths {
            if path_str.starts_with(allowed_path) {
                return Ok(());
            }
        }

        Err(anyhow::anyhow!("Path '{}' is not allowed", path_str))
    }

    pub fn validate_file_size(&self, size: u64) -> Result<()> {
        if size > self.config.max_file_size {
            Err(anyhow::anyhow!(
                "File size {} exceeds maximum allowed size {}",
                size,
                self.config.max_file_size
            ))
        } else {
            Ok(())
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_command_validation() {
        let validator = SecurityValidator::new(SecurityConfig::default());
        
        assert!(validator.validate_command("ls").is_ok());
        assert!(validator.validate_command("rm").is_err());
    }

    #[test]
    fn test_path_validation() {
        let validator = SecurityValidator::new(SecurityConfig::default());
        
        assert!(validator.validate_path("/tmp/test.txt").is_ok());
        assert!(validator.validate_path("/etc/passwd").is_err());
    }
}