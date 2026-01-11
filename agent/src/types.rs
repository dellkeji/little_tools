use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use chrono::{DateTime, Utc};
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AgentInfo {
    pub id: Uuid,
    pub hostname: String,
    pub platform: String,
    pub arch: String,
    pub version: String,
    pub last_heartbeat: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Command {
    pub id: Uuid,
    pub command_type: CommandType,
    pub payload: serde_json::Value,
    pub timeout: Option<u64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum CommandType {
    Execute,
    Deploy,
    Configure,
    Monitor,
    Stop,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CommandResult {
    pub command_id: Uuid,
    pub success: bool,
    pub output: String,
    pub error: Option<String>,
    pub execution_time: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MetricData {
    pub name: String,
    pub value: f64,
    pub labels: HashMap<String, String>,
    pub timestamp: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LogData {
    pub level: String,
    pub message: String,
    pub source: String,
    pub timestamp: DateTime<Utc>,
    pub labels: HashMap<String, String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MonitorConfig {
    pub metrics: Vec<MetricConfig>,
    pub logs: Vec<LogConfig>,
    pub interval: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MetricConfig {
    pub name: String,
    pub metric_type: MetricType,
    pub command: Option<String>,
    pub file_path: Option<String>,
    pub labels: HashMap<String, String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum MetricType {
    SystemCpu,
    SystemMemory,
    SystemDisk,
    SystemNetwork,
    ProcessCount,
    Custom,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LogConfig {
    pub name: String,
    pub source_type: LogSourceType,
    pub path: String,
    pub pattern: Option<String>,
    pub labels: HashMap<String, String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum LogSourceType {
    File,
    Command,
    SystemLog,
}