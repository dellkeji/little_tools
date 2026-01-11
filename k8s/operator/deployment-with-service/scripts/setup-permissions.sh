#!/bin/bash

# Script to set up correct file permissions for the Service Operator project

set -e

echo "Setting up file permissions for Service Operator project..."

# Make all shell scripts executable
find scripts/ -name "*.sh" -type f -exec chmod +x {} \;

# Make sure the main directories have correct permissions
chmod 755 api/ controllers/ config/ deploy/ docs/ examples/ scripts/ .github/

# Set correct permissions for configuration files
find config/ -type f -name "*.yaml" -exec chmod 644 {} \;
find deploy/ -type f -name "*.yaml" -exec chmod 644 {} \;
find examples/ -type f -name "*.yaml" -exec chmod 644 {} \;

# Set permissions for Go source files
find . -name "*.go" -type f -exec chmod 644 {} \;

# Set permissions for documentation files
find docs/ -name "*.md" -type f -exec chmod 644 {} \;
chmod 644 *.md

# Set permissions for configuration files
chmod 644 go.mod go.sum Dockerfile Makefile PROJECT LICENSE .gitignore .golangci.yml

# Set permissions for GitHub workflow files
find .github/ -type f -exec chmod 644 {} \;

# Set permissions for Helm chart files
find deploy/helm/ -type f -exec chmod 644 {} \;

echo "File permissions set up successfully!"

# Display executable files
echo ""
echo "Executable files:"
find . -type f -executable -name "*.sh" | sort