use crate::config::{DataPlaneConfig, CollectorConfig, ExporterConfig};
use crate::error::{AgentError, AgentResult};
use crate::metrics::MetricRegistry;
use crate::types::{MetricData, LogData};
use crate::platform::{get_system_metrics, execute_command};
use anyhow::Result;
use log::{info, error, warn};
use reqwest::Client;
use std::collections::HashMap;
use std::sync::Arc;
use std::time::Duration;
use tokio::sync::mpsc;
use tokio::time::sleep;
use chrono::Utc;

pub struct DataPlane {
    config: DataPlaneConfig,
    metric_sender: mpsc::Sender<MetricData>,
    log_sender: mpsc::Sender<LogData>,
    metric_registry: Arc<MetricRegistry>,
}

impl DataPlane {
    pub fn new(config: DataPlaneConfig, metric_registry: Arc<MetricRegistry>) -> Self {
        let (metric_sender, _) = mpsc::channel(config.buffer_size);
        let (log_sender, _) = mpsc::channel(config.buffer_size);

        Self {
            config,
            metric_sender,
            log_sender,
            metric_registry,
        }
    }

    pub async fn start(&mut self) -> AgentResult<()> {
        info!("Starting data plane");

        let (metric_tx, metric_rx) = mpsc::channel(self.config.buffer_size);
        let (log_tx, log_rx) = mpsc::channel(self.config.buffer_size);

        self.metric_sender = metric_tx.clone();
        self.log_sender = log_tx.clone();

        // Start collectors with error handling
        for collector_config in &self.config.collectors {
            if collector_config.enabled {
                if let Err(e) = self.start_collector(collector_config.clone(), metric_tx.clone(), log_tx.clone()).await {
                    error!("Failed to start collector {}: {}", collector_config.name, e);
                    // Continue with other collectors instead of failing completely
                }
            }
        }

        // Start exporters with error handling
        for exporter_config in &self.config.exporters {
            if exporter_config.enabled {
                if let Err(e) = self.start_exporter(exporter_config.clone(), metric_rx, log_rx).await {
                    error!("Failed to start exporter {}: {}", exporter_config.name, e);
                    // Continue with other exporters instead of failing completely
                }
            }
        }

        Ok(())
    }

    async fn start_collector(
        &self,
        config: CollectorConfig,
        metric_tx: mpsc::Sender<MetricData>,
        log_tx: mpsc::Sender<LogData>,
    ) -> AgentResult<()> {
        let collector_name = config.name.clone();
        let metric_registry = self.metric_registry.clone();
        
        tokio::spawn(async move {
            info!("Starting collector: {}", collector_name);
            
            let result = match config.collector_type.as_str() {
                "system" => {
                    SystemCollector::new(config, metric_tx, metric_registry).start().await
                }
                "log" => {
                    LogCollector::new(config, log_tx).start().await
                }
                "custom" => {
                    CustomCollector::new(config, metric_tx, log_tx).start().await
                }
                _ => {
                    error!("Unknown collector type: {}", config.collector_type);
                    Err(AgentError::DataPlaneError(format!("Unknown collector type: {}", config.collector_type)))
                }
            };

            if let Err(e) = result {
                error!("Collector {} failed: {}", collector_name, e);
            }
        });

        Ok(())
    }

    async fn start_exporter(
        &self,
        config: ExporterConfig,
        mut metric_rx: mpsc::Receiver<MetricData>,
        mut log_rx: mpsc::Receiver<LogData>,
    ) -> AgentResult<()> {
        let exporter_name = config.name.clone();
        
        tokio::spawn(async move {
            info!("Starting exporter: {}", exporter_name);
            
            let result = match config.exporter_type.as_str() {
                "http" => {
                    HttpExporter::new(config).start(metric_rx, log_rx).await
                }
                "file" => {
                    FileExporter::new(config).start(metric_rx, log_rx).await
                }
                _ => {
                    error!("Unknown exporter type: {}", config.exporter_type);
                    Err(AgentError::DataPlaneError(format!("Unknown exporter type: {}", config.exporter_type)))
                }
            };

            if let Err(e) = result {
                error!("Exporter {} failed: {}", exporter_name, e);
            }
        });

        Ok(())
    }
}

// System Metrics Collector
struct SystemCollector {
    config: CollectorConfig,
    metric_tx: mpsc::Sender<MetricData>,
    metric_registry: Arc<MetricRegistry>,
}

impl SystemCollector {
    fn new(config: CollectorConfig, metric_tx: mpsc::Sender<MetricData>, metric_registry: Arc<MetricRegistry>) -> Self {
        Self { config, metric_tx, metric_registry }
    }

    async fn start(self) -> AgentResult<()> {
        let interval = self.config.config
            .get("interval")
            .and_then(|v| v.as_u64())
            .unwrap_or(30);

        let mut interval_timer = tokio::time::interval(Duration::from_secs(interval));

        loop {
            interval_timer.tick().await;

            if let Err(e) = self.collect_metrics().await {
                error!("Error collecting system metrics: {}", e);
            }
        }
    }

    async fn collect_metrics(&self) -> Result<()> {
        let metrics = get_system_metrics()?;
        let timestamp = Utc::now();
        let mut labels = HashMap::new();
        labels.insert("collector".to_string(), self.config.name.clone());

        // CPU metric
        let cpu_metric = MetricData {
            name: "system_cpu_usage".to_string(),
            value: metrics.cpu_usage,
            labels: labels.clone(),
            timestamp,
        };

        // Memory metric
        let memory_metric = MetricData {
            name: "system_memory_usage".to_string(),
            value: metrics.memory_usage,
            labels: labels.clone(),
            timestamp,
        };

        // Disk metric
        let disk_metric = MetricData {
            name: "system_disk_usage".to_string(),
            value: metrics.disk_usage,
            labels: labels.clone(),
            timestamp,
        };

        // Send metrics
        if let Err(e) = self.metric_tx.send(cpu_metric).await {
            error!("Failed to send CPU metric: {}", e);
        }

        if let Err(e) = self.metric_tx.send(memory_metric).await {
            error!("Failed to send memory metric: {}", e);
        }

        if let Err(e) = self.metric_tx.send(disk_metric).await {
            error!("Failed to send disk metric: {}", e);
        }

        Ok(())
    }
}

// Log Collector
struct LogCollector {
    config: CollectorConfig,
    log_tx: mpsc::Sender<LogData>,
}

impl LogCollector {
    fn new(config: CollectorConfig, log_tx: mpsc::Sender<LogData>) -> Self {
        Self { config, log_tx }
    }

    async fn start(self) {
        info!("Log collector started for: {}", self.config.name);
        
        // Simplified log collection - in real implementation, you'd use file watchers
        let mut interval = tokio::time::interval(Duration::from_secs(10));
        
        loop {
            interval.tick().await;
            
            if let Err(e) = self.collect_logs().await {
                error!("Error collecting logs: {}", e);
            }
        }
    }

    async fn collect_logs(&self) -> Result<()> {
        // Simplified log collection
        let log_data = LogData {
            level: "INFO".to_string(),
            message: "Sample log message".to_string(),
            source: self.config.name.clone(),
            timestamp: Utc::now(),
            labels: HashMap::new(),
        };

        if let Err(e) = self.log_tx.send(log_data).await {
            error!("Failed to send log data: {}", e);
        }

        Ok(())
    }
}

// Custom Collector
struct CustomCollector {
    config: CollectorConfig,
    metric_tx: mpsc::Sender<MetricData>,
    log_tx: mpsc::Sender<LogData>,
}

impl CustomCollector {
    fn new(
        config: CollectorConfig,
        metric_tx: mpsc::Sender<MetricData>,
        log_tx: mpsc::Sender<LogData>,
    ) -> Self {
        Self {
            config,
            metric_tx,
            log_tx,
        }
    }

    async fn start(self) {
        info!("Custom collector started for: {}", self.config.name);
        
        let command = self.config.config
            .get("command")
            .and_then(|v| v.as_str())
            .unwrap_or("");

        let interval = self.config.config
            .get("interval")
            .and_then(|v| v.as_u64())
            .unwrap_or(60);

        let mut interval_timer = tokio::time::interval(Duration::from_secs(interval));

        loop {
            interval_timer.tick().await;

            if let Err(e) = self.execute_custom_command(command).await {
                error!("Error executing custom command: {}", e);
            }
        }
    }

    async fn execute_custom_command(&self, command: &str) -> Result<()> {
        let parts: Vec<&str> = command.split_whitespace().collect();
        if parts.is_empty() {
            return Ok(());
        }

        let cmd = parts[0];
        let args = &parts[1..];

        match execute_command(cmd, args) {
            Ok(output) => {
                let stdout = String::from_utf8_lossy(&output.stdout);
                
                // Parse output as metric (simplified)
                if let Ok(value) = stdout.trim().parse::<f64>() {
                    let metric = MetricData {
                        name: format!("custom_{}", self.config.name),
                        value,
                        labels: HashMap::new(),
                        timestamp: Utc::now(),
                    };

                    if let Err(e) = self.metric_tx.send(metric).await {
                        error!("Failed to send custom metric: {}", e);
                    }
                }
            }
            Err(e) => {
                error!("Custom command failed: {}", e);
            }
        }

        Ok(())
    }
}

// HTTP Exporter
struct HttpExporter {
    config: ExporterConfig,
    client: Client,
}

impl HttpExporter {
    fn new(config: ExporterConfig) -> Self {
        let mut headers = reqwest::header::HeaderMap::new();
        
        for (key, value) in &config.headers {
            if let (Ok(name), Ok(val)) = (key.parse(), value.parse()) {
                headers.insert(name, val);
            }
        }

        let client = Client::builder()
            .default_headers(headers)
            .timeout(Duration::from_secs(30))
            .build()
            .unwrap();

        Self { config, client }
    }

    async fn start(
        self,
        mut metric_rx: mpsc::Receiver<MetricData>,
        mut log_rx: mpsc::Receiver<LogData>,
    ) {
        let mut metrics_buffer = Vec::new();
        let mut logs_buffer = Vec::new();
        
        let mut flush_interval = tokio::time::interval(Duration::from_secs(60));

        loop {
            tokio::select! {
                metric = metric_rx.recv() => {
                    if let Some(metric) = metric {
                        metrics_buffer.push(metric);
                        
                        if metrics_buffer.len() >= self.config.batch_size {
                            self.export_metrics(&metrics_buffer).await;
                            metrics_buffer.clear();
                        }
                    }
                }
                log = log_rx.recv() => {
                    if let Some(log) = log {
                        logs_buffer.push(log);
                        
                        if logs_buffer.len() >= self.config.batch_size {
                            self.export_logs(&logs_buffer).await;
                            logs_buffer.clear();
                        }
                    }
                }
                _ = flush_interval.tick() => {
                    if !metrics_buffer.is_empty() {
                        self.export_metrics(&metrics_buffer).await;
                        metrics_buffer.clear();
                    }
                    
                    if !logs_buffer.is_empty() {
                        self.export_logs(&logs_buffer).await;
                        logs_buffer.clear();
                    }
                }
            }
        }
    }

    async fn export_metrics(&self, metrics: &[MetricData]) {
        let response = self.client
            .post(&self.config.endpoint)
            .json(metrics)
            .send()
            .await;

        match response {
            Ok(resp) => {
                if resp.status().is_success() {
                    info!("Exported {} metrics", metrics.len());
                } else {
                    warn!("Failed to export metrics: {}", resp.status());
                }
            }
            Err(e) => {
                error!("Error exporting metrics: {}", e);
            }
        }
    }

    async fn export_logs(&self, logs: &[LogData]) {
        let logs_endpoint = format!("{}/logs", self.config.endpoint);
        
        let response = self.client
            .post(&logs_endpoint)
            .json(logs)
            .send()
            .await;

        match response {
            Ok(resp) => {
                if resp.status().is_success() {
                    info!("Exported {} logs", logs.len());
                } else {
                    warn!("Failed to export logs: {}", resp.status());
                }
            }
            Err(e) => {
                error!("Error exporting logs: {}", e);
            }
        }
    }
}

// File Exporter
struct FileExporter {
    config: ExporterConfig,
}

impl FileExporter {
    fn new(config: ExporterConfig) -> Self {
        Self { config }
    }

    async fn start(
        self,
        mut metric_rx: mpsc::Receiver<MetricData>,
        mut log_rx: mpsc::Receiver<LogData>,
    ) {
        info!("File exporter started, writing to: {}", self.config.endpoint);
        
        loop {
            tokio::select! {
                metric = metric_rx.recv() => {
                    if let Some(metric) = metric {
                        self.write_metric(&metric).await;
                    }
                }
                log = log_rx.recv() => {
                    if let Some(log) = log {
                        self.write_log(&log).await;
                    }
                }
            }
        }
    }

    async fn write_metric(&self, metric: &MetricData) {
        let content = format!("{}\n", serde_json::to_string(metric).unwrap_or_default());
        
        if let Err(e) = tokio::fs::write(&self.config.endpoint, content).await {
            error!("Failed to write metric to file: {}", e);
        }
    }

    async fn write_log(&self, log: &LogData) {
        let log_file = format!("{}.log", self.config.endpoint);
        let content = format!("{}\n", serde_json::to_string(log).unwrap_or_default());
        
        if let Err(e) = tokio::fs::write(&log_file, content).await {
            error!("Failed to write log to file: {}", e);
        }
    }
}