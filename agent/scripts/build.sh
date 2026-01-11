#!/bin/bash

# Cross-platform build script for the agent

set -e

echo "Building cross-platform agent..."

# Create build directory
mkdir -p dist

# Build for different targets
TARGETS=(
    "x86_64-unknown-linux-gnu"
    "x86_64-pc-windows-gnu" 
    "x86_64-apple-darwin"
    "aarch64-apple-darwin"
)

for target in "${TARGETS[@]}"; do
    echo "Building for $target..."
    
    # Install target if not already installed
    rustup target add $target 2>/dev/null || true
    
    # Build
    cargo build --release --target $target
    
    # Copy binary to dist directory
    if [[ $target == *"windows"* ]]; then
        cp target/$target/release/agent.exe dist/agent-$target.exe
    else
        cp target/$target/release/agent dist/agent-$target
    fi
    
    echo "✓ Built for $target"
done

echo "All builds completed successfully!"
echo "Binaries available in dist/ directory:"
ls -la dist/