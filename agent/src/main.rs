use clap::{Parser, Subcommand};
use log::{info, error};
use std::path::PathBuf;

mod agent;
mod config;
mod control_plane;
mod data_plane;
mod error;
mod health;
mod logging;
mod metrics;
mod platform;
mod security;
mod server;
mod types;

use agent::Agent;
use config::AgentConfig;
use error::{AgentError, AgentResult};
use logging::{LoggingConfig, init_logging};

#[derive(Parser)]
#[command(name = "agent")]
#[command(about = "Cross-platform agent for remote control and monitoring")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Start the agent
    Start {
        /// Configuration file path
        #[arg(short, long, default_value = "config.toml")]
        config: PathBuf,
    },
    /// Generate default configuration
    Config {
        /// Output configuration file path
        #[arg(short, long, default_value = "config.toml")]
        output: PathBuf,
    },
}

#[tokio::main]
async fn main() -> AgentResult<()> {
    let cli = Cli::parse();
    
    match cli.command {
        Commands::Start { config } => {
            // Load configuration first
            let config = AgentConfig::load(&config)
                .map_err(|e| AgentError::ConfigError(e.to_string()))?;
            
            // Initialize logging
            let logging_config = LoggingConfig::from(&config.agent);
            init_logging(logging_config)?;
            
            info!("Starting agent with config: {:?}", config);
            
            // Create and run agent
            let mut agent = Agent::new(config).await?;
            
            // Handle the result of running the agent
            match agent.run().await {
                Ok(_) => {
                    info!("Agent stopped gracefully");
                    Ok(())
                }
                Err(e) => {
                    error!("Agent failed: {}", e);
                    Err(e)
                }
            }
        }
        Commands::Config { output } => {
            // Initialize basic logging for config generation
            env_logger::init();
            
            info!("Generating default config to: {:?}", output);
            AgentConfig::generate_default(&output)
                .map_err(|e| AgentError::ConfigError(e.to_string()))?;
            println!("Default configuration generated at: {:?}", output);
            Ok(())
        }
    }
}