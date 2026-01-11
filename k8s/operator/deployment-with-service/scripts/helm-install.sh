#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
NAMESPACE="service-operator-system"
RELEASE_NAME="service-operator"
CHART_PATH="deploy/helm/service-operator"
VALUES_FILE=""
IMAGE_TAG="latest"
DRY_RUN=false
UPGRADE=false

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
    echo -e "${BLUE}[HELM]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking prerequisites..."
    
    if ! command_exists helm; then
        print_error "helm is not installed or not in PATH"
        print_error "Please install Helm: https://helm.sh/docs/intro/install/"
        exit 1
    fi
    
    if ! command_exists kubectl; then
        print_error "kubectl is not installed or not in PATH"
        exit 1
    fi
    
    # Check if kubectl can connect to cluster
    if ! kubectl cluster-info >/dev/null 2>&1; then
        print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    # Check if chart directory exists
    if [ ! -d "$CHART_PATH" ]; then
        print_error "Chart directory not found: $CHART_PATH"
        exit 1
    fi
    
    print_status "Prerequisites check passed"
}

# Create namespace if it doesn't exist
create_namespace() {
    if ! kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
        print_status "Creating namespace: $NAMESPACE"
        kubectl create namespace "$NAMESPACE"
    else
        print_status "Namespace $NAMESPACE already exists"
    fi
}

# Install or upgrade the chart
install_chart() {
    print_header "Installing Service Operator with Helm..."
    
    local helm_args=(
        "$RELEASE_NAME"
        "$CHART_PATH"
        "--namespace" "$NAMESPACE"
        "--create-namespace"
        "--set" "image.tag=$IMAGE_TAG"
    )
    
    if [ -n "$VALUES_FILE" ]; then
        helm_args+=("--values" "$VALUES_FILE")
    fi
    
    if [ "$DRY_RUN" = true ]; then
        helm_args+=("--dry-run" "--debug")
    fi
    
    if [ "$UPGRADE" = true ]; then
        print_status "Upgrading release $RELEASE_NAME..."
        helm upgrade --install "${helm_args[@]}"
    else
        print_status "Installing release $RELEASE_NAME..."
        helm install "${helm_args[@]}"
    fi
}

# Verify installation
verify_installation() {
    if [ "$DRY_RUN" = true ]; then
        print_status "Dry run completed successfully"
        return
    fi
    
    print_header "Verifying installation..."
    
    # Wait for deployment to be ready
    print_status "Waiting for deployment to be ready..."
    kubectl wait --for=condition=available --timeout=300s \
        deployment/"$RELEASE_NAME-service-operator-controller-manager" \
        -n "$NAMESPACE"
    
    # Check if pods are running
    print_status "Checking pod status..."
    kubectl get pods -n "$NAMESPACE" -l "app.kubernetes.io/instance=$RELEASE_NAME"
    
    # Check if CRD is installed
    if kubectl get crd services.apps.example.com >/dev/null 2>&1; then
        print_status "CRD services.apps.example.com is available"
    else
        print_warning "CRD services.apps.example.com not found"
    fi
    
    print_status "Installation verification completed"
}

# Show post-installation information
show_post_install_info() {
    if [ "$DRY_RUN" = true ]; then
        return
    fi
    
    print_header "Post-Installation Information"
    echo ""
    echo "Service Operator has been installed successfully!"
    echo ""
    echo "Release Name: $RELEASE_NAME"
    echo "Namespace: $NAMESPACE"
    echo "Chart Version: $(helm list -n $NAMESPACE -o json | jq -r ".[] | select(.name==\"$RELEASE_NAME\") | .chart")"
    echo ""
    echo "Useful commands:"
    echo ""
    echo "  # Check operator status"
    echo "  kubectl get pods -n $NAMESPACE"
    echo ""
    echo "  # View operator logs"
    echo "  kubectl logs -n $NAMESPACE deployment/$RELEASE_NAME-service-operator-controller-manager -f"
    echo ""
    echo "  # Create a sample service"
    echo "  kubectl apply -f examples/sample-service.yaml"
    echo ""
    echo "  # List services"
    echo "  kubectl get services.apps.example.com"
    echo ""
    echo "  # Upgrade the operator"
    echo "  helm upgrade $RELEASE_NAME $CHART_PATH -n $NAMESPACE"
    echo ""
    echo "  # Uninstall the operator"
    echo "  helm uninstall $RELEASE_NAME -n $NAMESPACE"
    echo ""
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -n, --namespace NAMESPACE    Namespace to install operator (default: service-operator-system)"
    echo "  -r, --release RELEASE        Helm release name (default: service-operator)"
    echo "  -f, --values VALUES_FILE     Values file for Helm chart"
    echo "  -t, --tag IMAGE_TAG         Docker image tag (default: latest)"
    echo "  --dry-run                   Perform a dry run without installing"
    echo "  --upgrade                   Upgrade existing installation"
    echo "  -h, --help                  Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Install with default settings"
    echo "  $0 -n my-namespace -r my-operator    # Install in custom namespace with custom release name"
    echo "  $0 -f my-values.yaml                 # Install with custom values file"
    echo "  $0 --upgrade -t v1.0.0               # Upgrade to specific image tag"
    echo "  $0 --dry-run                         # Dry run to see what would be installed"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -r|--release)
            RELEASE_NAME="$2"
            shift 2
            ;;
        -f|--values)
            VALUES_FILE="$2"
            shift 2
            ;;
        -t|--tag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --upgrade)
            UPGRADE=true
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
    print_header "Service Operator Helm Installation"
    print_status "Release: $RELEASE_NAME"
    print_status "Namespace: $NAMESPACE"
    print_status "Image Tag: $IMAGE_TAG"
    if [ -n "$VALUES_FILE" ]; then
        print_status "Values File: $VALUES_FILE"
    fi
    if [ "$DRY_RUN" = true ]; then
        print_status "Mode: Dry Run"
    elif [ "$UPGRADE" = true ]; then
        print_status "Mode: Upgrade"
    else
        print_status "Mode: Install"
    fi
    echo ""
    
    check_prerequisites
    
    if [ "$DRY_RUN" = false ]; then
        create_namespace
    fi
    
    install_chart
    verify_installation
    show_post_install_info
    
    print_status "Helm installation completed successfully!"
}

# Run main function
main