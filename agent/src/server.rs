use crate::health::{HealthChecker, HealthStatus};
use crate::metrics::MetricRegistry;
use crate::types::AgentInfo;
use anyhow::Result;
use log::{info, error};
use serde_json::json;
use std::convert::Infallible;
use std::sync::Arc;
use tokio::net::TcpListener;
use warp::{Filter, Reply};

pub struct AgentServer {
    port: u16,
    health_checker: Arc<HealthChecker>,
    metric_registry: Arc<MetricRegistry>,
    agent_info: AgentInfo,
}

impl AgentServer {
    pub fn new(
        port: u16,
        health_checker: Arc<HealthChecker>,
        metric_registry: Arc<MetricRegistry>,
        agent_info: AgentInfo,
    ) -> Self {
        Self {
            port,
            health_checker,
            metric_registry,
            agent_info,
        }
    }

    pub async fn start(&self) -> Result<()> {
        info!("Starting agent HTTP server on port {}", self.port);

        let health_checker = self.health_checker.clone();
        let metric_registry = self.metric_registry.clone();
        let agent_info = self.agent_info.clone();

        // Health endpoint
        let health = warp::path("health")
            .and(warp::get())
            .and(with_health_checker(health_checker.clone()))
            .and_then(health_handler);

        // Metrics endpoint
        let metrics = warp::path("metrics")
            .and(warp::get())
            .and(with_metric_registry(metric_registry.clone()))
            .and_then(metrics_handler);

        // Info endpoint
        let info = warp::path("info")
            .and(warp::get())
            .and(with_agent_info(agent_info))
            .and_then(info_handler);

        // Status endpoint (combined health + metrics + info)
        let status = warp::path("status")
            .and(warp::get())
            .and(with_health_checker(health_checker))
            .and(with_metric_registry(metric_registry))
            .and_then(status_handler);

        let routes = health
            .or(metrics)
            .or(info)
            .or(status)
            .with(warp::cors().allow_any_origin());

        warp::serve(routes)
            .run(([0, 0, 0, 0], self.port))
            .await;

        Ok(())
    }
}

fn with_health_checker(
    health_checker: Arc<HealthChecker>,
) -> impl Filter<Extract = (Arc<HealthChecker>,), Error = Infallible> + Clone {
    warp::any().map(move || health_checker.clone())
}

fn with_metric_registry(
    metric_registry: Arc<MetricRegistry>,
) -> impl Filter<Extract = (Arc<MetricRegistry>,), Error = Infallible> + Clone {
    warp::any().map(move || metric_registry.clone())
}

fn with_agent_info(
    agent_info: AgentInfo,
) -> impl Filter<Extract = (AgentInfo,), Error = Infallible> + Clone {
    warp::any().map(move || agent_info.clone())
}

async fn health_handler(
    health_checker: Arc<HealthChecker>,
) -> Result<impl Reply, Infallible> {
    let health_status = health_checker.get_health_status().await;
    let status_code = match health_status.status {
        crate::health::HealthState::Healthy => warp::http::StatusCode::OK,
        crate::health::HealthState::Degraded => warp::http::StatusCode::OK,
        crate::health::HealthState::Unhealthy => warp::http::StatusCode::SERVICE_UNAVAILABLE,
    };

    Ok(warp::reply::with_status(
        warp::reply::json(&health_status),
        status_code,
    ))
}

async fn metrics_handler(
    metric_registry: Arc<MetricRegistry>,
) -> Result<impl Reply, Infallible> {
    let metrics = metric_registry.get_all_metrics();
    Ok(warp::reply::json(&json!({
        "metrics": metrics,
        "count": metrics.len()
    })))
}

async fn info_handler(agent_info: AgentInfo) -> Result<impl Reply, Infallible> {
    Ok(warp::reply::json(&agent_info))
}

async fn status_handler(
    health_checker: Arc<HealthChecker>,
    metric_registry: Arc<MetricRegistry>,
) -> Result<impl Reply, Infallible> {
    let health_status = health_checker.get_health_status().await;
    let metrics = metric_registry.get_all_metrics();
    
    let status = json!({
        "health": health_status,
        "metrics_count": metrics.len(),
        "timestamp": chrono::Utc::now()
    });

    let status_code = match health_status.status {
        crate::health::HealthState::Healthy => warp::http::StatusCode::OK,
        crate::health::HealthState::Degraded => warp::http::StatusCode::OK,
        crate::health::HealthState::Unhealthy => warp::http::StatusCode::SERVICE_UNAVAILABLE,
    };

    Ok(warp::reply::with_status(
        warp::reply::json(&status),
        status_code,
    ))
}

// Prometheus metrics format endpoint
pub async fn prometheus_metrics_handler(
    metric_registry: Arc<MetricRegistry>,
) -> Result<impl Reply, Infallible> {
    let metrics = metric_registry.get_all_metrics();
    let mut prometheus_output = String::new();

    for metric in metrics {
        // Convert to Prometheus format
        let labels = if metric.labels.is_empty() {
            String::new()
        } else {
            let label_pairs: Vec<String> = metric.labels
                .iter()
                .map(|(k, v)| format!("{}=\"{}\"", k, v))
                .collect();
            format!("{{{}}}", label_pairs.join(","))
        };

        prometheus_output.push_str(&format!(
            "{}{} {} {}\n",
            metric.name,
            labels,
            metric.value,
            metric.timestamp.timestamp_millis()
        ));
    }

    Ok(warp::reply::with_header(
        prometheus_output,
        "content-type",
        "text/plain; version=0.0.4; charset=utf-8",
    ))
}