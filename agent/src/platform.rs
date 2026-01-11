use std::process::{Command, Output};
use anyhow::Result;
use serde_json::Value;

pub struct PlatformInfo {
    pub os: String,
    pub arch: String,
    pub hostname: String,
}

impl PlatformInfo {
    pub fn new() -> Result<Self> {
        Ok(Self {
            os: std::env::consts::OS.to_string(),
            arch: std::env::consts::ARCH.to_string(),
            hostname: hostname::get()?.to_string_lossy().to_string(),
        })
    }
}

pub struct SystemMetrics {
    pub cpu_usage: f64,
    pub memory_usage: f64,
    pub disk_usage: f64,
}

pub fn get_system_metrics() -> Result<SystemMetrics> {
    #[cfg(target_os = "windows")]
    {
        get_windows_metrics()
    }
    #[cfg(target_os = "linux")]
    {
        get_linux_metrics()
    }
    #[cfg(target_os = "macos")]
    {
        get_macos_metrics()
    }
}

#[cfg(target_os = "windows")]
fn get_windows_metrics() -> Result<SystemMetrics> {
    // CPU usage via wmic
    let cpu_output = Command::new("wmic")
        .args(&["cpu", "get", "loadpercentage", "/value"])
        .output()?;
    
    let cpu_str = String::from_utf8_lossy(&cpu_output.stdout);
    let cpu_usage = parse_windows_cpu(&cpu_str)?;

    // Memory usage via wmic
    let mem_output = Command::new("wmic")
        .args(&["OS", "get", "TotalVisibleMemorySize,FreePhysicalMemory", "/value"])
        .output()?;
    
    let mem_str = String::from_utf8_lossy(&mem_output.stdout);
    let memory_usage = parse_windows_memory(&mem_str)?;

    Ok(SystemMetrics {
        cpu_usage,
        memory_usage,
        disk_usage: 0.0, // Simplified for demo
    })
}

#[cfg(target_os = "linux")]
fn get_linux_metrics() -> Result<SystemMetrics> {
    // CPU usage from /proc/stat
    let cpu_usage = get_linux_cpu_usage()?;
    
    // Memory usage from /proc/meminfo
    let memory_usage = get_linux_memory_usage()?;
    
    Ok(SystemMetrics {
        cpu_usage,
        memory_usage,
        disk_usage: 0.0, // Simplified for demo
    })
}

#[cfg(target_os = "macos")]
fn get_macos_metrics() -> Result<SystemMetrics> {
    // Use system_profiler and vm_stat for macOS
    let cpu_output = Command::new("top")
        .args(&["-l", "1", "-n", "0"])
        .output()?;
    
    let cpu_str = String::from_utf8_lossy(&cpu_output.stdout);
    let cpu_usage = parse_macos_cpu(&cpu_str)?;

    let mem_output = Command::new("vm_stat")
        .output()?;
    
    let mem_str = String::from_utf8_lossy(&mem_output.stdout);
    let memory_usage = parse_macos_memory(&mem_str)?;

    Ok(SystemMetrics {
        cpu_usage,
        memory_usage,
        disk_usage: 0.0, // Simplified for demo
    })
}

pub fn execute_command(cmd: &str, args: &[&str]) -> Result<Output> {
    let output = Command::new(cmd)
        .args(args)
        .output()?;
    Ok(output)
}

pub fn deploy_file(source: &str, destination: &str) -> Result<()> {
    std::fs::copy(source, destination)?;
    Ok(())
}

pub fn write_config(path: &str, content: &Value) -> Result<()> {
    let content_str = serde_json::to_string_pretty(content)?;
    std::fs::write(path, content_str)?;
    Ok(())
}

// Helper functions for parsing platform-specific output
#[cfg(target_os = "windows")]
fn parse_windows_cpu(output: &str) -> Result<f64> {
    for line in output.lines() {
        if line.starts_with("LoadPercentage=") {
            let value = line.split('=').nth(1).unwrap_or("0");
            return Ok(value.parse::<f64>().unwrap_or(0.0));
        }
    }
    Ok(0.0)
}

#[cfg(target_os = "windows")]
fn parse_windows_memory(output: &str) -> Result<f64> {
    let mut total = 0u64;
    let mut free = 0u64;
    
    for line in output.lines() {
        if line.starts_with("TotalVisibleMemorySize=") {
            total = line.split('=').nth(1).unwrap_or("0").parse().unwrap_or(0);
        } else if line.starts_with("FreePhysicalMemory=") {
            free = line.split('=').nth(1).unwrap_or("0").parse().unwrap_or(0);
        }
    }
    
    if total > 0 {
        Ok(((total - free) as f64 / total as f64) * 100.0)
    } else {
        Ok(0.0)
    }
}

#[cfg(target_os = "linux")]
fn get_linux_cpu_usage() -> Result<f64> {
    // Simplified CPU usage calculation
    let stat = std::fs::read_to_string("/proc/stat")?;
    // Parse first line for CPU stats
    // This is a simplified implementation
    Ok(50.0) // Placeholder
}

#[cfg(target_os = "linux")]
fn get_linux_memory_usage() -> Result<f64> {
    let meminfo = std::fs::read_to_string("/proc/meminfo")?;
    let mut total = 0u64;
    let mut available = 0u64;
    
    for line in meminfo.lines() {
        if line.starts_with("MemTotal:") {
            total = line.split_whitespace().nth(1).unwrap_or("0").parse().unwrap_or(0);
        } else if line.starts_with("MemAvailable:") {
            available = line.split_whitespace().nth(1).unwrap_or("0").parse().unwrap_or(0);
        }
    }
    
    if total > 0 {
        Ok(((total - available) as f64 / total as f64) * 100.0)
    } else {
        Ok(0.0)
    }
}

#[cfg(target_os = "macos")]
fn parse_macos_cpu(output: &str) -> Result<f64> {
    // Parse top output for CPU usage
    // This is a simplified implementation
    Ok(30.0) // Placeholder
}

#[cfg(target_os = "macos")]
fn parse_macos_memory(output: &str) -> Result<f64> {
    // Parse vm_stat output for memory usage
    // This is a simplified implementation
    Ok(60.0) // Placeholder
}