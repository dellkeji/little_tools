#!/bin/bash

# End-to-end test script for Service Operator

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
NAMESPACE="e2e-test"
SERVICE_NAME="test-service"
TIMEOUT=300

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
    echo -e "${BLUE}[E2E]${NC} $1"
}

# Function to cleanup resources
cleanup() {
    print_header "Cleaning up test resources..."
    
    # Delete test service
    kubectl delete services.apps.example.com "$SERVICE_NAME" -n "$NAMESPACE" --ignore-not-found=true
    
    # Delete test namespace
    kubectl delete namespace "$NAMESPACE" --ignore-not-found=true
    
    print_status "Cleanup completed"
}

# Set up cleanup trap
trap cleanup EXIT

# Function to wait for condition
wait_for_condition() {
    local resource=$1
    local condition=$2
    local timeout=${3:-$TIMEOUT}
    
    print_status "Waiting for $resource to be $condition (timeout: ${timeout}s)..."
    
    if kubectl wait --for=condition="$condition" --timeout="${timeout}s" "$resource" -n "$NAMESPACE"; then
        print_status "$resource is $condition"
        return 0
    else
        print_error "$resource failed to become $condition within ${timeout}s"
        return 1
    fi
}

# Function to test basic service creation
test_basic_service() {
    print_header "Testing basic service creation..."
    
    # Create test service
    cat <<EOF | kubectl apply -f -
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: $SERVICE_NAME
  namespace: $NAMESPACE
spec:
  image: nginx:1.21
  replicas: 2
  port: 80
  serviceType: ClusterIP
EOF
    
    # Wait for service to be ready
    wait_for_condition "services.apps.example.com/$SERVICE_NAME" "Ready"
    
    # Verify deployment was created
    if kubectl get deployment "$SERVICE_NAME" -n "$NAMESPACE" >/dev/null 2>&1; then
        print_status "Deployment created successfully"
    else
        print_error "Deployment was not created"
        return 1
    fi
    
    # Verify service was created
    if kubectl get service "$SERVICE_NAME" -n "$NAMESPACE" >/dev/null 2>&1; then
        print_status "Service created successfully"
    else
        print_error "Service was not created"
        return 1
    fi
    
    # Check replicas
    local ready_replicas
    ready_replicas=$(kubectl get services.apps.example.com "$SERVICE_NAME" -n "$NAMESPACE" -o jsonpath='{.status.readyReplicas}')
    if [ "$ready_replicas" = "2" ]; then
        print_status "Correct number of replicas are ready: $ready_replicas"
    else
        print_error "Expected 2 ready replicas, got: $ready_replicas"
        return 1
    fi
    
    print_status "Basic service test passed"
}

# Function to test service with config
test_service_with_config() {
    print_header "Testing service with configuration..."
    
    # Update service with config
    cat <<EOF | kubectl apply -f -
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: $SERVICE_NAME
  namespace: $NAMESPACE
spec:
  image: nginx:1.21
  replicas: 2
  port: 80
  serviceType: ClusterIP
  configData:
    nginx.conf: |
      server {
          listen 80;
          location / {
              return 200 'Hello from E2E test!';
              add_header Content-Type text/plain;
          }
      }
  env:
    - name: TEST_ENV
      value: "e2e-test"
EOF
    
    # Wait for deployment to be updated
    kubectl rollout status deployment/"$SERVICE_NAME" -n "$NAMESPACE" --timeout="${TIMEOUT}s"
    
    # Verify configmap was created
    if kubectl get configmap "$SERVICE_NAME-config" -n "$NAMESPACE" >/dev/null 2>&1; then
        print_status "ConfigMap created successfully"
    else
        print_error "ConfigMap was not created"
        return 1
    fi
    
    # Verify environment variable is set
    local env_value
    env_value=$(kubectl get deployment "$SERVICE_NAME" -n "$NAMESPACE" -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="TEST_ENV")].value}')
    if [ "$env_value" = "e2e-test" ]; then
        print_status "Environment variable set correctly: TEST_ENV=$env_value"
    else
        print_error "Environment variable not set correctly, got: $env_value"
        return 1
    fi
    
    print_status "Service with config test passed"
}

# Function to test service scaling
test_service_scaling() {
    print_header "Testing service scaling..."
    
    # Scale up to 3 replicas
    kubectl patch services.apps.example.com "$SERVICE_NAME" -n "$NAMESPACE" -p '{"spec":{"replicas":3}}'
    
    # Wait for scaling to complete
    kubectl rollout status deployment/"$SERVICE_NAME" -n "$NAMESPACE" --timeout="${TIMEOUT}s"
    
    # Verify replica count
    local ready_replicas
    ready_replicas=$(kubectl get services.apps.example.com "$SERVICE_NAME" -n "$NAMESPACE" -o jsonpath='{.status.readyReplicas}')
    if [ "$ready_replicas" = "3" ]; then
        print_status "Scaling up successful: $ready_replicas replicas"
    else
        print_error "Scaling failed, expected 3 replicas, got: $ready_replicas"
        return 1
    fi
    
    # Scale down to 1 replica
    kubectl patch services.apps.example.com "$SERVICE_NAME" -n "$NAMESPACE" -p '{"spec":{"replicas":1}}'
    
    # Wait for scaling to complete
    kubectl rollout status deployment/"$SERVICE_NAME" -n "$NAMESPACE" --timeout="${TIMEOUT}s"
    
    # Verify replica count
    ready_replicas=$(kubectl get services.apps.example.com "$SERVICE_NAME" -n "$NAMESPACE" -o jsonpath='{.status.readyReplicas}')
    if [ "$ready_replicas" = "1" ]; then
        print_status "Scaling down successful: $ready_replicas replica"
    else
        print_error "Scaling failed, expected 1 replica, got: $ready_replicas"
        return 1
    fi
    
    print_status "Service scaling test passed"
}

# Function to test service connectivity
test_service_connectivity() {
    print_header "Testing service connectivity..."
    
    # Create a test pod to test connectivity
    kubectl run test-pod --image=curlimages/curl:latest --rm -i --restart=Never -n "$NAMESPACE" -- \
        curl -f -s "http://$SERVICE_NAME.$NAMESPACE.svc.cluster.local" || {
        print_error "Service connectivity test failed"
        return 1
    }
    
    print_status "Service connectivity test passed"
}

# Function to test service deletion
test_service_deletion() {
    print_header "Testing service deletion..."
    
    # Delete the service
    kubectl delete services.apps.example.com "$SERVICE_NAME" -n "$NAMESPACE"
    
    # Wait for resources to be cleaned up
    local timeout=60
    local count=0
    while [ $count -lt $timeout ]; do
        if ! kubectl get deployment "$SERVICE_NAME" -n "$NAMESPACE" >/dev/null 2>&1 && \
           ! kubectl get service "$SERVICE_NAME" -n "$NAMESPACE" >/dev/null 2>&1 && \
           ! kubectl get configmap "$SERVICE_NAME-config" -n "$NAMESPACE" >/dev/null 2>&1; then
            print_status "All resources cleaned up successfully"
            break
        fi
        sleep 1
        ((count++))
    done
    
    if [ $count -eq $timeout ]; then
        print_error "Resources were not cleaned up within ${timeout}s"
        return 1
    fi
    
    print_status "Service deletion test passed"
}

# Function to run all tests
run_tests() {
    print_header "Starting E2E tests..."
    
    # Create test namespace
    kubectl create namespace "$NAMESPACE" || true
    
    # Run tests
    test_basic_service
    test_service_with_config
    test_service_scaling
    test_service_connectivity
    test_service_deletion
    
    print_header "All E2E tests passed!"
}

# Function to check prerequisites
check_prerequisites() {
    print_header "Checking prerequisites..."
    
    # Check if kubectl is available
    if ! command -v kubectl >/dev/null 2>&1; then
        print_error "kubectl is not installed or not in PATH"
        exit 1
    fi
    
    # Check if cluster is accessible
    if ! kubectl cluster-info >/dev/null 2>&1; then
        print_error "Cannot connect to Kubernetes cluster"
        exit 1
    fi
    
    # Check if Service Operator is installed
    if ! kubectl get crd services.apps.example.com >/dev/null 2>&1; then
        print_error "Service Operator CRD not found. Please install Service Operator first."
        exit 1
    fi
    
    # Check if operator is running
    if ! kubectl get deployment service-operator-controller-manager -n service-operator-system >/dev/null 2>&1; then
        print_error "Service Operator is not running. Please install Service Operator first."
        exit 1
    fi
    
    print_status "Prerequisites check passed"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -n, --namespace NAMESPACE    Test namespace (default: e2e-test)"
    echo "  -s, --service SERVICE        Test service name (default: test-service)"
    echo "  -t, --timeout TIMEOUT        Timeout in seconds (default: 300)"
    echo "  -h, --help                   Show this help message"
    echo ""
    echo "This script runs end-to-end tests for Service Operator."
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -s|--service)
            SERVICE_NAME="$2"
            shift 2
            ;;
        -t|--timeout)
            TIMEOUT="$2"
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

# Main execution
main() {
    print_header "Service Operator E2E Tests"
    print_status "Test namespace: $NAMESPACE"
    print_status "Service name: $SERVICE_NAME"
    print_status "Timeout: ${TIMEOUT}s"
    echo ""
    
    check_prerequisites
    run_tests
    
    print_status "E2E tests completed successfully!"
}

# Run main function
main