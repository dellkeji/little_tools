#!/bin/bash
# Build script for Daily Report Tool
# This script builds the Windows executable with console window hidden

echo "Building Daily Report Tool..."

# Build the executable with windowsgui flag to hide console
go build -ldflags="-H windowsgui" -o daily-report.exe ./cmd/daily-report

if [ $? -eq 0 ]; then
    echo "Build successful! Executable: daily-report.exe"
else
    echo "Build failed!"
    exit 1
fi
