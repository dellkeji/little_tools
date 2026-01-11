#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
NAMESPACE="service-operator-system"
FORCE=false

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
    
    # Check if kubectl can connect to cluster
    if ! kubectl cluster-info >/dev/null 2>&1; then
        print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
        exit 1
    fi
    
    print_status "Prerequisites check passed"
}

# Function to check for existing Service resources
check_existing_resources() {
    if [ "$FORCE" = true ]; then
        return
    fi
    
    print_status "Checking for existing Service resources..."
    
    if kubectl get services.apps.example.com --all-namespaces --no-headers 2>/dev/null | grep -q .; then
        print_warning "Found existing Service resources:"
        kubectl get services.apps.example.com --all-namespaces
        print_warning ""
        print_warning "Uninstalling the operator will leave these resources orphaned."
        print_warning "Consider deleting them first or use --force to proceed anyway."
        
        read -p "Do you want to continue? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_status "Uninstallation cancelled"
            exit 0
        fi
    else
        print_status "No existing Service resources found"
    fi
}

# Function to undeploy operator
undeploy_operator() {
    print_status "Undeploying Service Operator..."
    
    if kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
        make undeploy
        print_status "Service Operator undeployed successfully"
    else
        print_warning "Namespace $NAMESPACE not found, skipping operator undeployment"
    fi
}

# Function to uninstall CRDs
uninstall_crds() {
    print_status "Uninstalling Custom Resource Definitions..."
    
    if kubectl get crd services.apps.example.com >/dev/null 2>&1; then
        make uninstall
        print_status "CRDs uninstalled successfully"
    else
        print_warning "CRD services.apps.example.com not found, skipping CRD uninstallation"
    fi
}

# Function to cleanup remaining resources
cleanup_remaining() {
    print_status "Cleaning up remaining resources..."
    
    # Remove any remaining configmaps, secrets, etc. in the operator namespace
    if kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
        print_status "Removing remaining resources in namespace $NAMESPACE..."
        kubectl delete all --all -n "$NAMESPACE" --ignore-not-found=true
        kubectl delete configmaps --all -n "$NAMESPACE" --ignore-not-found=true
        kubectl delete secrets --all -n "$NAMESPACE" --ignore-not-found=true
        
        # Wait a bit for resources to be deleted
        sleep 5
        
        # Delete the namespace if it's empty
        if ! kubectl get all -n "$NAMESPACE" --no-headers 2>/dev/null | grep -q .; then
            kubectl delete namespace "$NAMESPACE" --ignore-not-found=true
            print_status "Namespace $NAMESPACE deleted"
        else
            print_warning "Namespace $NAMESPACE still contains resources, not deleting"
        fi
    fi
}

# Function to verify uninstallation
verify_uninstallation() {
    print_status "Verifying uninstallation..."
    
    # Check if CRD is removed
    if kubectl get crd services.apps.example.com >/dev/null 2>&1; then
        print_warning "CRD services.apps.example.com still exists"
    else
        print_status "CRD successfully removed"
    fi
    
    # Check if namespace is removed
    if kubectl get namespace "$NAMESPACE" >/dev/null 2>&1; then
        print_warning "Namespace $NAMESPACE still exists"
    else
        print_status "Namespace successfully removed"
    fi
    
    print_status "Uninstallation verification completed"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -n, --namespace NAMESPACE    Namespace where operator is installed (default: service-operator-system)"
    echo "  -f, --force                 Force uninstallation without checking for existing resources"
    echo "  -h, --help                 Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                          # Uninstall with default settings"
    echo "  $0 -f                       # Force uninstallation"
    echo "  $0 -n my-namespace          # Uninstall from custom namespace"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -f|--force)
            FORCE=true
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

# Main uninstallation process
main() {
    print_status "Starting Service Operator uninstallation..."
    print_status "Namespace: $NAMESPACE"
    
    check_prerequisites
    check_existing_resources
    undeploy_operator
    uninstall_crds
    cleanup_remaining
    verify_uninstallation
    
    print_status "Service Operator uninstallation completed successfully!"
}

# Run main function
main