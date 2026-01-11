use crate::types::MetricData;
use chrono::Utc;
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::time::{Duration, Instant};

#[derive(Debug, Clone)]
pub struct MetricRegistry {
    metrics: Arc<Mutex<HashMap<String, MetricEntry>>>,
}

#[derive(Debug, Clone)]
struct MetricEntry {
    value: f64,
    labels: HashMap<String, String>,
    last_updated: Instant,
}

impl MetricRegistry {
    pub fn new() -> Self {
        Self {
            metrics: Arc::new(Mutex::new(HashMap::new())),
        }
    }

    pub fn record_metric(&self, name: &str, value: f64, labels: HashMap<String, String>) {
        let mut metrics = self.metrics.lock().unwrap();
        metrics.insert(
            name.to_string(),
            MetricEntry {
                value,
                labels,
                last_updated: Instant::now(),
            },
        );
    }

    pub fn get_all_metrics(&self) -> Vec<MetricData> {
        let metrics = self.metrics.lock().unwrap();
        let now = Utc::now();
        
        metrics
            .iter()
            .map(|(name, entry)| MetricData {
                name: name.clone(),
                value: entry.value,
                labels: entry.labels.clone(),
                timestamp: now,
            })
            .collect()
    }

    pub fn cleanup_old_metrics(&self, max_age: Duration) {
        let mut metrics = self.metrics.lock().unwrap();
        let cutoff = Instant::now() - max_age;
        
        metrics.retain(|_, entry| entry.last_updated > cutoff);
    }

    pub fn get_metric_count(&self) -> usize {
        self.metrics.lock().unwrap().len()
    }
}

// Built-in system metrics
pub struct SystemMetricsCollector {
    registry: MetricRegistry,
}

impl SystemMetricsCollector {
    pub fn new(registry: MetricRegistry) -> Self {
        Self { registry }
    }

    pub async fn collect_all(&self) -> anyhow::Result<()> {
        self.collect_cpu_metrics().await?;
        self.collect_memory_metrics().await?;
        self.collect_disk_metrics().await?;
        self.collect_network_metrics().await?;
        Ok(())
    }

    async fn collect_cpu_metrics(&self) -> anyhow::Result<()> {
        let metrics = crate::platform::get_system_metrics()?;
        
        let mut labels = HashMap::new();
        labels.insert("type".to_string(), "system".to_string());
        
        self.registry.record_metric("cpu_usage_percent", metrics.cpu_usage, labels);
        Ok(())
    }

    async fn collect_memory_metrics(&self) -> anyhow::Result<()> {
        let metrics = crate::platform::get_system_metrics()?;
        
        let mut labels = HashMap::new();
        labels.insert("type".to_string(), "system".to_string());
        
        self.registry.record_metric("memory_usage_percent", metrics.memory_usage, labels);
        Ok(())
    }

    async fn collect_disk_metrics(&self) -> anyhow::Result<()> {
        let metrics = crate::platform::get_system_metrics()?;
        
        let mut labels = HashMap::new();
        labels.insert("type".to_string(), "system".to_string());
        
        self.registry.record_metric("disk_usage_percent", metrics.disk_usage, labels);
        Ok(())
    }

    async fn collect_network_metrics(&self) -> anyhow::Result<()> {
        // Placeholder for network metrics
        let mut labels = HashMap::new();
        labels.insert("type".to_string(), "network".to_string());
        
        self.registry.record_metric("network_bytes_sent", 0.0, labels.clone());
        self.registry.record_metric("network_bytes_received", 0.0, labels);
        Ok(())
    }
}

// Performance metrics for the agent itself
pub struct AgentMetrics {
    registry: MetricRegistry,
    start_time: Instant,
}

impl AgentMetrics {
    pub fn new(registry: MetricRegistry) -> Self {
        Self {
            registry,
            start_time: Instant::now(),
        }
    }

    pub fn record_command_execution(&self, duration: Duration, success: bool) {
        let mut labels = HashMap::new();
        labels.insert("type".to_string(), "command".to_string());
        labels.insert("success".to_string(), success.to_string());
        
        self.registry.record_metric(
            "command_execution_duration_ms",
            duration.as_millis() as f64,
            labels,
        );
    }

    pub fn record_data_export(&self, count: usize, exporter: &str) {
        let mut labels = HashMap::new();
        labels.insert("exporter".to_string(), exporter.to_string());
        
        self.registry.record_metric("data_exported_count", count as f64, labels);
    }

    pub fn record_uptime(&self) {
        let uptime = self.start_time.elapsed().as_secs() as f64;
        let labels = HashMap::new();
        
        self.registry.record_metric("agent_uptime_seconds", uptime, labels);
    }

    pub fn record_memory_usage(&self) {
        // Get current process memory usage (simplified)
        let mut labels = HashMap::new();
        labels.insert("type".to_string(), "agent".to_string());
        
        // This would need platform-specific implementation
        self.registry.record_metric("agent_memory_bytes", 0.0, labels);
    }
}