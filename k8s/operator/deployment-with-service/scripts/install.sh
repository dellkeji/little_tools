#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
NAMESPACE="service-operator-system"
IMAGE="service-operator:latest"
SKIP_BUILD=false

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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    if ! command_exists kubectl; then
        print_error "kubectl is not installed or not in PATH"
        exit 1
    fi
    
    if ! command_exists docker; then
        print_error "docker is not installed or not in PATH"
        exit 1
    fi
    
    # Check if kubectl can connect to cluster
    if ! kubectl cluster-info >/dev/null 2>&1; then
        print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    print_status "Prerequisites check passed"
}

# Function to build and push image
build_image() {
    if [ "$SKIP_BUILD" = true ]; then
        print_status "Skipping image build"
        return
    fi
    
    print_status "Building Docker image: $IMAGE"
    make docker-build IMG="$IMAGE"
    
    print_status "Image built successfully"
}

# Function to install CRDs
install_crds() {
    print_status "Installing Custom Resource Definitions..."
    make install
    
    # Wait for CRDs to be established
    print_status "Waiting for CRDs to be established..."
    kubectl wait --for condition=established --timeout=60s crd/services.apps.example.com
    
    print_status "CRDs installed successfully"
}

# Function to deploy operator
deploy_operator() {
    print_status "Deploying Service Operator..."
    make deploy IMG="$IMAGE"
    
    # Wait for deployment to be ready
    print_status "Waiting for operator deployment to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment/service-operator-controller-manager -n "$NAMESPACE"
    
    print_status "Service Operator deployed successfully"
}

# Function to verify installation
verify_installation() {
    print_status "Verifying installation..."
    
    # Check if operator pod is running
    if kubectl get pods -n "$NAMESPACE" -l control-plane=controller-manager --no-headers | grep -q Running; then
        print_status "Operator pod is running"
    else
        print_error "Operator pod is not running"
        kubectl get pods -n "$NAMESPACE" -l control-plane=controller-manager
        exit 1
    fi
    
    # Check if CRD is available
    if kubectl get crd services.apps.example.com >/dev/null 2>&1; then
        print_status "CRD is available"
    else
        print_error "CRD is not available"
        exit 1
    fi
    
    print_status "Installation verification completed successfully"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -n, --namespace NAMESPACE    Namespace to install operator (default: service-operator-system)"
    echo "  -i, --image IMAGE           Docker image for operator (default: service-operator:latest)"
    echo "  -s, --skip-build           Skip building Docker image"
    echo "  -h, --help                 Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Install with default settings"
    echo "  $0 -i myregistry/service-operator:v1.0.0  # Use custom image"
    echo "  $0 -s                                # Skip building image"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -i|--image)
            IMAGE="$2"
            shift 2
            ;;
        -s|--skip-build)
            SKIP_BUILD=true
            shift
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

# Main installation process
main() {
    print_status "Starting Service Operator installation..."
    print_status "Namespace: $NAMESPACE"
    print_status "Image: $IMAGE"
    
    check_prerequisites
    build_image
    install_crds
    deploy_operator
    verify_installation
    
    print_status "Service Operator installation completed successfully!"
    print_status ""
    print_status "Next steps:"
    print_status "1. Create a Service resource: kubectl apply -f examples/sample-service.yaml"
    print_status "2. Check the status: kubectl get services.apps.example.com"
    print_status "3. View operator logs: kubectl logs -n $NAMESPACE deployment/service-operator-controller-manager"
}

# Run main function
main