use crate::types::AgentInfo;
use anyhow::Result;
use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::time::{Duration, Instant};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HealthStatus {
    pub status: HealthState,
    pub timestamp: DateTime<Utc>,
    pub checks: HashMap<String, CheckResult>,
    pub uptime: u64,
    pub version: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum HealthState {
    Healthy,
    Degraded,
    Unhealthy,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CheckResult {
    pub status: HealthState,
    pub message: String,
    pub last_check: DateTime<Utc>,
    pub duration_ms: u64,
}

pub struct HealthChecker {
    checks: Arc<Mutex<HashMap<String, Box<dyn HealthCheck + Send + Sync>>>>,
    start_time: Instant,
    agent_info: AgentInfo,
}

pub trait HealthCheck {
    fn name(&self) -> &str;
    fn check(&self) -> Result<CheckResult>;
}

impl HealthChecker {
    pub fn new(agent_info: AgentInfo) -> Self {
        Self {
            checks: Arc::new(Mutex::new(HashMap::new())),
            start_time: Instant::now(),
            agent_info,
        }
    }

    pub fn register_check(&self, check: Box<dyn HealthCheck + Send + Sync>) {
        let mut checks = self.checks.lock().unwrap();
        checks.insert(check.name().to_string(), check);
    }

    pub async fn get_health_status(&self) -> HealthStatus {
        let mut check_results = HashMap::new();
        let mut overall_status = HealthState::Healthy;

        let checks = self.checks.lock().unwrap();
        for (name, check) in checks.iter() {
            let start = Instant::now();
            let result = match check.check() {
                Ok(mut result) => {
                    result.duration_ms = start.elapsed().as_millis() as u64;
                    result
                }
                Err(e) => CheckResult {
                    status: HealthState::Unhealthy,
                    message: format!("Check failed: {}", e),
                    last_check: Utc::now(),
                    duration_ms: start.elapsed().as_millis() as u64,
                },
            };

            // Update overall status based on individual check results
            match result.status {
                HealthState::Unhealthy => overall_status = HealthState::Unhealthy,
                HealthState::Degraded if matches!(overall_status, HealthState::Healthy) => {
                    overall_status = HealthState::Degraded
                }
                _ => {}
            }

            check_results.insert(name.clone(), result);
        }

        HealthStatus {
            status: overall_status,
            timestamp: Utc::now(),
            checks: check_results,
            uptime: self.start_time.elapsed().as_secs(),
            version: self.agent_info.version.clone(),
        }
    }
}

// Built-in health checks
pub struct ControlPlaneHealthCheck {
    server_url: String,
    client: reqwest::Client,
}

impl ControlPlaneHealthCheck {
    pub fn new(server_url: String) -> Self {
        Self {
            server_url,
            client: reqwest::Client::new(),
        }
    }
}

impl HealthCheck for ControlPlaneHealthCheck {
    fn name(&self) -> &str {
        "control_plane_connectivity"
    }

    fn check(&self) -> Result<CheckResult> {
        let rt = tokio::runtime::Runtime::new()?;
        
        let result = rt.block_on(async {
            let response = self
                .client
                .get(&format!("{}/health", self.server_url))
                .timeout(Duration::from_secs(5))
                .send()
                .await;

            match response {
                Ok(resp) if resp.status().is_success() => CheckResult {
                    status: HealthState::Healthy,
                    message: "Control plane is reachable".to_string(),
                    last_check: Utc::now(),
                    duration_ms: 0,
                },
                Ok(resp) => CheckResult {
                    status: HealthState::Degraded,
                    message: format!("Control plane returned status: {}", resp.status()),
                    last_check: Utc::now(),
                    duration_ms: 0,
                },
                Err(e) => CheckResult {
                    status: HealthState::Unhealthy,
                    message: format!("Cannot reach control plane: {}", e),
                    last_check: Utc::now(),
                    duration_ms: 0,
                },
            }
        });

        Ok(result)
    }
}

pub struct DiskSpaceHealthCheck {
    path: String,
    threshold_percent: f64,
}

impl DiskSpaceHealthCheck {
    pub fn new(path: String, threshold_percent: f64) -> Self {
        Self {
            path,
            threshold_percent,
        }
    }
}

impl HealthCheck for DiskSpaceHealthCheck {
    fn name(&self) -> &str {
        "disk_space"
    }

    fn check(&self) -> Result<CheckResult> {
        // Simplified disk space check
        let usage = get_disk_usage(&self.path)?;
        
        let (status, message) = if usage > 95.0 {
            (HealthState::Unhealthy, format!("Disk usage critical: {:.1}%", usage))
        } else if usage > self.threshold_percent {
            (HealthState::Degraded, format!("Disk usage high: {:.1}%", usage))
        } else {
            (HealthState::Healthy, format!("Disk usage normal: {:.1}%", usage))
        };

        Ok(CheckResult {
            status,
            message,
            last_check: Utc::now(),
            duration_ms: 0,
        })
    }
}

pub struct MemoryHealthCheck {
    threshold_percent: f64,
}

impl MemoryHealthCheck {
    pub fn new(threshold_percent: f64) -> Self {
        Self { threshold_percent }
    }
}

impl HealthCheck for MemoryHealthCheck {
    fn name(&self) -> &str {
        "memory_usage"
    }

    fn check(&self) -> Result<CheckResult> {
        let metrics = crate::platform::get_system_metrics()?;
        let usage = metrics.memory_usage;
        
        let (status, message) = if usage > 95.0 {
            (HealthState::Unhealthy, format!("Memory usage critical: {:.1}%", usage))
        } else if usage > self.threshold_percent {
            (HealthState::Degraded, format!("Memory usage high: {:.1}%", usage))
        } else {
            (HealthState::Healthy, format!("Memory usage normal: {:.1}%", usage))
        };

        Ok(CheckResult {
            status,
            message,
            last_check: Utc::now(),
            duration_ms: 0,
        })
    }
}

// Platform-specific disk usage function
fn get_disk_usage(path: &str) -> Result<f64> {
    #[cfg(target_os = "windows")]
    {
        use std::ffi::CString;
        use winapi::um::fileapi::GetDiskFreeSpaceExA;
        
        let path_c = CString::new(path)?;
        let mut free_bytes = 0u64;
        let mut total_bytes = 0u64;
        
        unsafe {
            if GetDiskFreeSpaceExA(
                path_c.as_ptr(),
                &mut free_bytes,
                &mut total_bytes,
                std::ptr::null_mut(),
            ) != 0 {
                let used_bytes = total_bytes - free_bytes;
                Ok((used_bytes as f64 / total_bytes as f64) * 100.0)
            } else {
                Err(anyhow::anyhow!("Failed to get disk space"))
            }
        }
    }
    
    #[cfg(unix)]
    {
        use std::ffi::CString;
        use std::mem;
        
        let path_c = CString::new(path)?;
        let mut statvfs: libc::statvfs = unsafe { mem::zeroed() };
        
        let result = unsafe { libc::statvfs(path_c.as_ptr(), &mut statvfs) };
        
        if result == 0 {
            let total = statvfs.f_blocks * statvfs.f_frsize;
            let free = statvfs.f_bavail * statvfs.f_frsize;
            let used = total - free;
            Ok((used as f64 / total as f64) * 100.0)
        } else {
            Err(anyhow::anyhow!("Failed to get disk space"))
        }
    }
}