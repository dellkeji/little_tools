# Windows installation script for the cross-platform agent

param(
    [string]$InstallPath = "C:\Program Files\Agent",
    [string]$ServiceName = "CrossPlatformAgent"
)

Write-Host "Installing Cross-Platform Agent..." -ForegroundColor Green

# Check if running as administrator
if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Error "This script must be run as Administrator"
    exit 1
}

# Create installation directory
Write-Host "Creating installation directory..." -ForegroundColor Yellow
New-Item -ItemType Directory -Force -Path $InstallPath | Out-Null
New-Item -ItemType Directory -Force -Path "$InstallPath\config" | Out-Null
New-Item -ItemType Directory -Force -Path "$InstallPath\logs" | Out-Null

# Copy binary
Write-Host "Installing binary..." -ForegroundColor Yellow
Copy-Item "agent.exe" "$InstallPath\agent.exe"

# Generate default config
Write-Host "Generating default configuration..." -ForegroundColor Yellow
& "$InstallPath\agent.exe" config -o "$InstallPath\config\config.toml"

# Install as Windows service using sc command
Write-Host "Installing Windows service..." -ForegroundColor Yellow
$servicePath = "`"$InstallPath\agent.exe`" start -c `"$InstallPath\config\config.toml`""

# Remove existing service if it exists
sc.exe delete $ServiceName 2>$null

# Create new service
sc.exe create $ServiceName binPath= $servicePath start= auto DisplayName= "Cross Platform Agent"

if ($LASTEXITCODE -eq 0) {
    Write-Host "Service installed successfully!" -ForegroundColor Green
    
    # Set service description
    sc.exe description $ServiceName "Cross-platform agent for remote control and monitoring"
    
    Write-Host "Installation completed successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Cyan
    Write-Host "1. Edit configuration: $InstallPath\config\config.toml"
    Write-Host "2. Start the service: sc start $ServiceName"
    Write-Host "3. Check status: sc query $ServiceName"
    Write-Host "4. View logs in: $InstallPath\logs\"
} else {
    Write-Error "Failed to install service"
    exit 1
}