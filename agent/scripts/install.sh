#!/bin/bash

# Installation script for the cross-platform agent

set -e

INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/agent"
SERVICE_DIR="/etc/systemd/system"
USER="agent"

echo "Installing Cross-Platform Agent..."

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root (use sudo)" 
   exit 1
fi

# Create agent user
if ! id "$USER" &>/dev/null; then
    echo "Creating agent user..."
    useradd -r -s /bin/false -d /var/lib/agent $USER
fi

# Create directories
echo "Creating directories..."
mkdir -p $CONFIG_DIR
mkdir -p /var/lib/agent
mkdir -p /var/log/agent

# Copy binary
echo "Installing binary..."
cp agent $INSTALL_DIR/agent
chmod +x $INSTALL_DIR/agent

# Generate default config
echo "Generating default configuration..."
$INSTALL_DIR/agent config -o $CONFIG_DIR/config.toml

# Set permissions
chown -R $USER:$USER /var/lib/agent
chown -R $USER:$USER /var/log/agent
chown root:$USER $CONFIG_DIR/config.toml
chmod 640 $CONFIG_DIR/config.toml

# Create systemd service
echo "Creating systemd service..."
cat > $SERVICE_DIR/agent.service << EOF
[Unit]
Description=Cross Platform Agent
After=network.target
Wants=network.target

[Service]
Type=simple
User=$USER
Group=$USER
ExecStart=$INSTALL_DIR/agent start -c $CONFIG_DIR/config.toml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=agent

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/agent /var/log/agent

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and enable service
echo "Enabling service..."
systemctl daemon-reload
systemctl enable agent.service

echo "Installation completed successfully!"
echo ""
echo "Next steps:"
echo "1. Edit configuration: $CONFIG_DIR/config.toml"
echo "2. Start the service: systemctl start agent"
echo "3. Check status: systemctl status agent"
echo "4. View logs: journalctl -u agent -f"