#!/bin/bash
# Package script for Daily Report Tool using Fyne
# This script creates a packaged Windows application

echo "Packaging Daily Report Tool with Fyne..."

# Check if fyne command is available
if ! command -v fyne &> /dev/null; then
    echo "Fyne CLI not found. Installing..."
    go install fyne.io/fyne/v2/cmd/fyne@latest
fi

# Package the application
if [ -f icon.png ]; then
    fyne package -os windows -icon icon.png -name "Daily Report Tool" -appID com.dailyreport.tool
else
    echo "Warning: icon.png not found, packaging without icon"
    fyne package -os windows -name "Daily Report Tool" -appID com.dailyreport.tool
fi

if [ $? -eq 0 ]; then
    echo "Package successful!"
else
    echo "Package failed!"
    exit 1
fi
