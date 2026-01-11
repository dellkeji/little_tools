#!/bin/bash

# Script to validate Kubernetes manifests and Helm charts

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
    echo -e "${BLUE}[VALIDATE]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking prerequisites..."
    
    local missing_tools=()
    
    if ! command_exists kubectl; then
        missing_tools+=("kubectl")
    fi
    
    if ! command_exists helm; then
        missing_tools+=("helm")
    fi
    
    if ! command_exists kustomize; then
        print_warning "kustomize not found, will try to use kubectl kustomize"
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        print_error "Missing required tools: ${missing_tools[*]}"
        exit 1
    fi
    
    print_status "Prerequisites check passed"
}

# Validate YAML syntax
validate_yaml_syntax() {
    print_header "Validating YAML syntax..."
    
    local error_count=0
    
    # Find all YAML files
    while IFS= read -r -d '' file; do
        if ! kubectl apply --dry-run=client -f "$file" >/dev/null 2>&1; then
            print_error "Invalid YAML syntax in: $file"
            ((error_count++))
        else
            print_status "Valid YAML: $file"
        fi
    done < <(find config/ examples/ deploy/ -name "*.yaml" -type f -print0)
    
    if [ $error_count -eq 0 ]; then
        print_status "All YAML files have valid syntax"
    else
        print_error "Found $error_count files with invalid YAML syntax"
        return 1
    fi
}

# Validate Kubernetes manifests
validate_k8s_manifests() {
    print_header "Validating Kubernetes manifests..."
    
    local temp_dir
    temp_dir=$(mktemp -d)
    
    # Validate CRD
    print_status "Validating CRD..."
    if kubectl apply --dry-run=server -f config/crd/bases/apps.example.com_services.yaml >/dev/null 2>&1; then
        print_status "CRD validation passed"
    else
        print_error "CRD validation failed"
        return 1
    fi
    
    # Validate RBAC
    print_status "Validating RBAC..."
    if kustomize build config/rbac > "$temp_dir/rbac.yaml" 2>/dev/null || kubectl kustomize config/rbac > "$temp_dir/rbac.yaml" 2>/dev/null; then
        if kubectl apply --dry-run=server -f "$temp_dir/rbac.yaml" >/dev/null 2>&1; then
            print_status "RBAC validation passed"
        else
            print_error "RBAC validation failed"
            return 1
        fi
    else
        print_error "Failed to build RBAC manifests"
        return 1
    fi
    
    # Validate manager
    print_status "Validating manager..."
    if kustomize build config/manager > "$temp_dir/manager.yaml" 2>/dev/null || kubectl kustomize config/manager > "$temp_dir/manager.yaml" 2>/dev/null; then
        if kubectl apply --dry-run=server -f "$temp_dir/manager.yaml" >/dev/null 2>&1; then
            print_status "Manager validation passed"
        else
            print_error "Manager validation failed"
            return 1
        fi
    else
        print_error "Failed to build manager manifests"
        return 1
    fi
    
    # Validate default configuration
    print_status "Validating default configuration..."
    if kustomize build config/default > "$temp_dir/default.yaml" 2>/dev/null || kubectl kustomize config/default > "$temp_dir/default.yaml" 2>/dev/null; then
        if kubectl apply --dry-run=server -f "$temp_dir/default.yaml" >/dev/null 2>&1; then
            print_status "Default configuration validation passed"
        else
            print_error "Default configuration validation failed"
            return 1
        fi
    else
        print_error "Failed to build default configuration"
        return 1
    fi
    
    # Cleanup
    rm -rf "$temp_dir"
    
    print_status "Kubernetes manifests validation completed"
}

# Validate examples
validate_examples() {
    print_header "Validating examples..."
    
    local error_count=0
    
    for example in examples/*.yaml; do
        if [ -f "$example" ]; then
            print_status "Validating example: $example"
            if kubectl apply --dry-run=client -f "$example" >/dev/null 2>&1; then
                print_status "Example validation passed: $example"
            else
                print_error "Example validation failed: $example"
                ((error_count++))
            fi
        fi
    done
    
    if [ $error_count -eq 0 ]; then
        print_status "All examples validation passed"
    else
        print_error "Found $error_count invalid examples"
        return 1
    fi
}

# Validate Helm chart
validate_helm_chart() {
    print_header "Validating Helm chart..."
    
    local chart_dir="deploy/helm/service-operator"
    
    # Lint Helm chart
    print_status "Linting Helm chart..."
    if helm lint "$chart_dir"; then
        print_status "Helm chart lint passed"
    else
        print_error "Helm chart lint failed"
        return 1
    fi
    
    # Template Helm chart
    print_status "Templating Helm chart..."
    local temp_dir
    temp_dir=$(mktemp -d)
    
    if helm template test-release "$chart_dir" --output-dir "$temp_dir" >/dev/null 2>&1; then
        print_status "Helm chart templating passed"
        
        # Validate templated manifests
        print_status "Validating templated manifests..."
        if kubectl apply --dry-run=client -f "$temp_dir/service-operator/templates/" >/dev/null 2>&1; then
            print_status "Templated manifests validation passed"
        else
            print_error "Templated manifests validation failed"
            rm -rf "$temp_dir"
            return 1
        fi
    else
        print_error "Helm chart templating failed"
        rm -rf "$temp_dir"
        return 1
    fi
    
    # Test with different values
    print_status "Testing Helm chart with custom values..."
    cat > "$temp_dir/test-values.yaml" << EOF
replicaCount: 2
image:
  tag: test
resources:
  limits:
    cpu: 1000m
    memory: 256Mi
monitoring:
  serviceMonitor:
    enabled: true
EOF
    
    if helm template test-release "$chart_dir" -f "$temp_dir/test-values.yaml" --output-dir "$temp_dir/custom" >/dev/null 2>&1; then
        print_status "Helm chart with custom values passed"
    else
        print_error "Helm chart with custom values failed"
        rm -rf "$temp_dir"
        return 1
    fi
    
    # Cleanup
    rm -rf "$temp_dir"
    
    print_status "Helm chart validation completed"
}

# Validate API schema
validate_api_schema() {
    print_header "Validating API schema..."
    
    # Check if CRD has proper OpenAPI schema
    local crd_file="config/crd/bases/apps.example.com_services.yaml"
    
    if grep -q "openAPIV3Schema" "$crd_file"; then
        print_status "CRD contains OpenAPI v3 schema"
    else
        print_error "CRD missing OpenAPI v3 schema"
        return 1
    fi
    
    # Check required fields
    if grep -q "required:" "$crd_file"; then
        print_status "CRD has required field validation"
    else
        print_warning "CRD missing required field validation"
    fi
    
    # Check default values
    if grep -q "default:" "$crd_file"; then
        print_status "CRD has default values"
    else
        print_warning "CRD missing default values"
    fi
    
    print_status "API schema validation completed"
}

# Validate RBAC permissions
validate_rbac_permissions() {
    print_header "Validating RBAC permissions..."
    
    local rbac_file="config/rbac/role.yaml"
    
    # Check if all required permissions are present
    local required_resources=("deployments" "services" "configmaps" "ingresses")
    local missing_resources=()
    
    for resource in "${required_resources[@]}"; do
        if ! grep -q "$resource" "$rbac_file"; then
            missing_resources+=("$resource")
        fi
    done
    
    if [ ${#missing_resources[@]} -eq 0 ]; then
        print_status "All required RBAC permissions are present"
    else
        print_error "Missing RBAC permissions for: ${missing_resources[*]}"
        return 1
    fi
    
    # Check for overly broad permissions
    if grep -q "resources: \[\"*\"\]" "$rbac_file"; then
        print_warning "Found wildcard resource permissions - consider being more specific"
    fi
    
    if grep -q "verbs: \[\"*\"\]" "$rbac_file"; then
        print_warning "Found wildcard verb permissions - consider being more specific"
    fi
    
    print_status "RBAC permissions validation completed"
}

# Main validation function
run_validation() {
    print_header "Starting manifest validation..."
    
    local failed_checks=()
    
    # Run all validation checks
    if ! validate_yaml_syntax; then
        failed_checks+=("YAML syntax")
    fi
    
    if ! validate_k8s_manifests; then
        failed_checks+=("Kubernetes manifests")
    fi
    
    if ! validate_examples; then
        failed_checks+=("Examples")
    fi
    
    if ! validate_helm_chart; then
        failed_checks+=("Helm chart")
    fi
    
    if ! validate_api_schema; then
        failed_checks+=("API schema")
    fi
    
    if ! validate_rbac_permissions; then
        failed_checks+=("RBAC permissions")
    fi
    
    # Report results
    if [ ${#failed_checks[@]} -eq 0 ]; then
        print_header "All validation checks passed! ✅"
        return 0
    else
        print_error "The following validation checks failed:"
        for check in "${failed_checks[@]}"; do
            echo "  - $check"
        done
        return 1
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --yaml-only              Validate only YAML syntax"
    echo "  --k8s-only              Validate only Kubernetes manifests"
    echo "  --examples-only         Validate only examples"
    echo "  --helm-only             Validate only Helm chart"
    echo "  --api-only              Validate only API schema"
    echo "  --rbac-only             Validate only RBAC permissions"
    echo "  -h, --help              Show this help message"
    echo ""
    echo "This script validates all Kubernetes manifests and Helm charts in the project."
}

# Parse command line arguments
case "${1:-}" in
    --yaml-only)
        check_prerequisites
        validate_yaml_syntax
        ;;
    --k8s-only)
        check_prerequisites
        validate_k8s_manifests
        ;;
    --examples-only)
        check_prerequisites
        validate_examples
        ;;
    --helm-only)
        check_prerequisites
        validate_helm_chart
        ;;
    --api-only)
        validate_api_schema
        ;;
    --rbac-only)
        validate_rbac_permissions
        ;;
    -h|--help)
        show_usage
        exit 0
        ;;
    "")
        check_prerequisites
        run_validation
        ;;
    *)
        print_error "Unknown option: $1"
        show_usage
        exit 1
        ;;
esac