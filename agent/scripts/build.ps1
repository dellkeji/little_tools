# Cross-platform build script for Windows

Write-Host "Building cross-platform agent..." -ForegroundColor Green

# Create build directory
New-Item -ItemType Directory -Force -Path "dist" | Out-Null

# Build targets
$targets = @(
    "x86_64-unknown-linux-gnu",
    "x86_64-pc-windows-gnu", 
    "x86_64-apple-darwin",
    "aarch64-apple-darwin"
)

foreach ($target in $targets) {
    Write-Host "Building for $target..." -ForegroundColor Yellow
    
    # Install target if not already installed
    rustup target add $target 2>$null
    
    # Build
    cargo build --release --target $target
    
    # Copy binary to dist directory
    if ($target -like "*windows*") {
        Copy-Item "target\$target\release\agent.exe" "dist\agent-$target.exe"
    } else {
        Copy-Item "target\$target\release\agent" "dist\agent-$target"
    }
    
    Write-Host "✓ Built for $target" -ForegroundColor Green
}

Write-Host "All builds completed successfully!" -ForegroundColor Green
Write-Host "Binaries available in dist\ directory:" -ForegroundColor Cyan
Get-ChildItem dist\