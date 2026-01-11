#!/bin/bash

# This script generates a standalone installation script for releases

set -e

VERSION=${1:-"latest"}
OUTPUT_FILE=${2:-"install.sh"}

cat > "$OUTPUT_FILE" << 'EOF'
#!/bin/bash

# Service Operator Installation Script
# This script installs Service Operator on a Kubernetes cluster

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
VERSION="latest"
NAMESPACE="service-operator-system"
INSTALL_METHOD="manifests"
GITHUB_REPO="example/service-operator"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}[INSTALL]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking prerequisites..."
    
    if ! command_exists kubectl; then
        print_error "kubectl is not installed or not in PATH"
        print_error "Please install kubectl: https://kubernetes.io/docs/tasks/tools/"
        exit 1
    fi
    
    # Check if kubectl can connect to cluster
    if ! kubectl cluster-info >/dev/null 2>&1; then
        print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    print_status "Prerequisites check passed"
}

# Get latest version from GitHub
get_latest_version() {
    if [ "$VERSION" = "latest" ]; then
        print_status "Fetching latest version..."
        VERSION=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [ -z "$VERSION" ]; then
            print_error "Failed to fetch latest version"
            exit 1
        fi
        print_status "Latest version: $VERSION"
    fi
}

# Install using manifests
install_manifests() {
    print_header "Installing Service Operator using manifests..."
    
    local manifest_url="https://github.com/$GITHUB_REPO/releases/download/$VERSION/install.yaml"
    
    print_status "Downloading manifests from $manifest_url"
    if ! curl -sSL "$manifest_url" | kubectl apply -f -; then
        print_error "Failed to install manifests"
        exit 1
    fi
    
    print_status "Manifests installed successfully"
}

# Install using Helm
install_helm() {
    print_header "Installing Service Operator using Helm..."
    
    if ! command_exists helm; then
        print_error "helm is not installed or not in PATH"
        print_error "Please install Helm: https://helm.sh/docs/intro/install/"
        exit 1
    fi
    
    local chart_url="https://github.com/$GITHUB_REPO/releases/download/$VERSION/service-operator-${VERSION#v}.tgz"
    
    print_status "Installing Helm chart from $chart_url"
    if ! helm install service-operator "$chart_url" --namespace "$NAMESPACE" --create-namespace; then
        print_error "Failed to install Helm chart"
        exit 1
    fi
    
    print_status "Helm chart installed successfully"
}

# Verify installation
verify_installation() {
    print_header "Verifying installation..."
    
    # Wait for deployment to be ready
    print_status "Waiting for operator deployment to be ready..."
    if ! kubectl wait --for=condition=available --timeout=300s deployment/service-operator-controller-manager -n "$NAMESPACE"; then
        print_error "Operator deployment failed to become ready"
        exit 1
    fi
    
    # Check if CRD is available
    if kubectl get crd services.apps.example.com >/dev/null 2>&1; then
        print_status "CRD services.apps.example.com is available"
    else
        print_error "CRD services.apps.example.com not found"
        exit 1
    fi
    
    print_status "Installation verification completed successfully"
}

# Show post-installation information
show_post_install_info() {
    print_header "Installation Complete!"
    echo ""
    echo "Service Operator has been installed successfully!"
    echo ""
    echo "Version: $VERSION"
    echo "Namespace: $NAMESPACE"
    echo "Installation Method: $INSTALL_METHOD"
    echo ""
    echo "Next steps:"
    echo ""
    echo "1. Create a Service resource:"
    echo "   kubectl apply -f https://raw.githubusercontent.com/$GITHUB_REPO/main/examples/sample-service.yaml"
    echo ""
    echo "2. Check the status:"
    echo "   kubectl get services.apps.example.com"
    echo ""
    echo "3. View operator logs:"
    echo "   kubectl logs -n $NAMESPACE deployment/service-operator-controller-manager -f"
    echo ""
    echo "For more information, visit: https://github.com/$GITHUB_REPO"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -v, --version VERSION        Version to install (default: latest)"
    echo "  -n, --namespace NAMESPACE    Namespace to install operator (default: service-operator-system)"
    echo "  -m, --method METHOD         Installation method: manifests or helm (default: manifests)"
    echo "  -r, --repo REPO             GitHub repository (default: example/service-operator)"
    echo "  -h, --help                  Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Install latest version with manifests"
    echo "  $0 -v v0.1.0 -m helm                # Install specific version with Helm"
    echo "  $0 -n my-namespace                   # Install in custom namespace"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -m|--method)
            INSTALL_METHOD="$2"
            shift 2
            ;;
        -r|--repo)
            GITHUB_REPO="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Validate installation method
if [[ "$INSTALL_METHOD" != "manifests" && "$INSTALL_METHOD" != "helm" ]]; then
    print_error "Invalid installation method: $INSTALL_METHOD"
    print_error "Supported methods: manifests, helm"
    exit 1
fi

# Main installation process
main() {
    print_header "Service Operator Installation"
    print_status "Version: $VERSION"
    print_status "Namespace: $NAMESPACE"
    print_status "Method: $INSTALL_METHOD"
    print_status "Repository: $GITHUB_REPO"
    echo ""
    
    check_prerequisites
    get_latest_version
    
    case "$INSTALL_METHOD" in
        "manifests")
            install_manifests
            ;;
        "helm")
            install_helm
            ;;
    esac
    
    verify_installation
    show_post_install_info
}

# Run main function
main
EOF

chmod +x "$OUTPUT_FILE"

echo "Generated installation script: $OUTPUT_FILE"
echo "Usage: ./$OUTPUT_FILE [options]"