#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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
    echo -e "${BLUE}[SETUP]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking development prerequisites..."
    
    local missing_tools=()
    
    if ! command_exists go; then
        missing_tools+=("go")
    fi
    
    if ! command_exists docker; then
        missing_tools+=("docker")
    fi
    
    if ! command_exists kubectl; then
        missing_tools+=("kubectl")
    fi
    
    if ! command_exists kind; then
        print_warning "kind not found - will provide installation instructions"
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        print_error "Missing required tools: ${missing_tools[*]}"
        print_error "Please install them before running this script"
        exit 1
    fi
    
    print_status "Prerequisites check passed"
}

# Setup development cluster with kind
setup_kind_cluster() {
    if command_exists kind; then
        print_header "Setting up kind cluster for development..."
        
        # Check if cluster already exists
        if kind get clusters | grep -q "service-operator-dev"; then
            print_status "Kind cluster 'service-operator-dev' already exists"
            kubectl cluster-info --context kind-service-operator-dev
        else
            print_status "Creating kind cluster 'service-operator-dev'..."
            
            # Create kind config
            cat <<EOF > /tmp/kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF
            
            kind create cluster --name service-operator-dev --config /tmp/kind-config.yaml
            rm /tmp/kind-config.yaml
            
            # Install nginx ingress controller
            print_status "Installing nginx ingress controller..."
            kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
            
            # Wait for ingress controller to be ready
            kubectl wait --namespace ingress-nginx \
                --for=condition=ready pod \
                --selector=app.kubernetes.io/component=controller \
                --timeout=90s
        fi
        
        # Set kubectl context
        kubectl config use-context kind-service-operator-dev
        print_status "Kind cluster setup completed"
    else
        print_warning "kind not found. To install kind, run:"
        print_warning "  # On macOS with Homebrew:"
        print_warning "  brew install kind"
        print_warning "  # On Linux:"
        print_warning "  curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64"
        print_warning "  chmod +x ./kind"
        print_warning "  sudo mv ./kind /usr/local/bin/kind"
    fi
}

# Install development tools
install_dev_tools() {
    print_header "Installing development tools..."
    
    # Install controller-gen, kustomize, etc.
    print_status "Installing kubebuilder tools..."
    make controller-gen
    make kustomize
    make envtest
    
    print_status "Development tools installed"
}

# Setup Go modules
setup_go_modules() {
    print_header "Setting up Go modules..."
    
    if [ ! -f "go.sum" ]; then
        print_status "Downloading Go dependencies..."
        go mod download
        go mod tidy
    else
        print_status "Go modules already set up"
    fi
}

# Generate code and manifests
generate_code() {
    print_header "Generating code and manifests..."
    
    print_status "Generating deepcopy code..."
    make generate
    
    print_status "Generating CRD manifests..."
    make manifests
    
    print_status "Code generation completed"
}

# Run tests
run_tests() {
    print_header "Running tests..."
    
    print_status "Running unit tests..."
    make test
    
    print_status "Tests completed successfully"
}

# Setup development environment
setup_dev_env() {
    print_header "Setting up development environment..."
    
    # Create local bin directory
    mkdir -p bin
    
    # Install CRDs to the cluster
    print_status "Installing CRDs to development cluster..."
    make install
    
    print_status "Development environment setup completed"
}

# Show development commands
show_dev_commands() {
    print_header "Development Commands Reference"
    echo ""
    echo "Common development tasks:"
    echo ""
    echo "  make run                    # Run operator locally against the cluster"
    echo "  make test                   # Run tests"
    echo "  make build                  # Build the operator binary"
    echo "  make docker-build           # Build Docker image"
    echo "  make install                # Install CRDs"
    echo "  make deploy                 # Deploy operator to cluster"
    echo "  make undeploy               # Remove operator from cluster"
    echo "  make uninstall              # Remove CRDs"
    echo ""
    echo "Testing with examples:"
    echo ""
    echo "  kubectl apply -f examples/sample-service.yaml"
    echo "  kubectl get services.apps.example.com"
    echo "  kubectl describe services.apps.example.com sample-web-service"
    echo ""
    echo "Viewing logs:"
    echo ""
    echo "  # If running locally with 'make run':"
    echo "  # Logs will appear in the terminal"
    echo ""
    echo "  # If deployed to cluster:"
    echo "  kubectl logs -n service-operator-system deployment/service-operator-controller-manager"
    echo ""
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --skip-cluster             Skip kind cluster setup"
    echo "  --skip-tests               Skip running tests"
    echo "  -h, --help                Show this help message"
    echo ""
    echo "This script sets up a complete development environment for the Service Operator."
}

# Parse command line arguments
SKIP_CLUSTER=false
SKIP_TESTS=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-cluster)
            SKIP_CLUSTER=true
            shift
            ;;
        --skip-tests)
            SKIP_TESTS=true
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

# Main setup process
main() {
    print_header "Service Operator Development Setup"
    print_status "This script will set up a complete development environment"
    echo ""
    
    check_prerequisites
    
    if [ "$SKIP_CLUSTER" = false ]; then
        setup_kind_cluster
    fi
    
    setup_go_modules
    install_dev_tools
    generate_code
    
    if [ "$SKIP_TESTS" = false ]; then
        run_tests
    fi
    
    setup_dev_env
    show_dev_commands
    
    print_status "Development setup completed successfully!"
    print_status ""
    print_status "You can now start developing:"
    print_status "1. Run 'make run' to start the operator locally"
    print_status "2. In another terminal, apply examples: 'kubectl apply -f examples/sample-service.yaml'"
    print_status "3. Check the results: 'kubectl get services.apps.example.com'"
}

# Run main function
main