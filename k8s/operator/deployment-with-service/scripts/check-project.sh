#!/bin/bash

# Script to check project completeness and consistency

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0
WARNING_CHECKS=0

# Function to print colored output
print_status() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((PASSED_CHECKS++))
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
    ((WARNING_CHECKS++))
}

print_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((FAILED_CHECKS++))
}

print_header() {
    echo -e "${BLUE}[CHECK]${NC} $1"
}

# Function to check if file exists
check_file() {
    local file=$1
    local description=$2
    ((TOTAL_CHECKS++))
    
    if [ -f "$file" ]; then
        print_status "$description: $file"
    else
        print_error "$description missing: $file"
    fi
}

# Function to check if directory exists
check_directory() {
    local dir=$1
    local description=$2
    ((TOTAL_CHECKS++))
    
    if [ -d "$dir" ]; then
        print_status "$description: $dir"
    else
        print_error "$description missing: $dir"
    fi
}

# Function to check if file is executable
check_executable() {
    local file=$1
    local description=$2
    ((TOTAL_CHECKS++))
    
    if [ -f "$file" ] && [ -x "$file" ]; then
        print_status "$description executable: $file"
    elif [ -f "$file" ]; then
        print_error "$description not executable: $file"
    else
        print_error "$description missing: $file"
    fi
}

# Check project structure
check_project_structure() {
    print_header "Checking project structure..."
    
    # Core directories
    check_directory "api" "API directory"
    check_directory "api/v1" "API v1 directory"
    check_directory "controllers" "Controllers directory"
    check_directory "config" "Config directory"
    check_directory "config/crd" "CRD config directory"
    check_directory "config/rbac" "RBAC config directory"
    check_directory "config/manager" "Manager config directory"
    check_directory "docs" "Documentation directory"
    check_directory "examples" "Examples directory"
    check_directory "scripts" "Scripts directory"
    check_directory "deploy/helm" "Helm deployment directory"
    check_directory ".github" "GitHub directory"
    check_directory ".github/workflows" "GitHub workflows directory"
}

# Check core files
check_core_files() {
    print_header "Checking core files..."
    
    # Go files
    check_file "main.go" "Main Go file"
    check_file "go.mod" "Go module file"
    check_file "go.sum" "Go sum file"
    check_file "api/v1/service_types.go" "Service types"
    check_file "api/v1/groupversion_info.go" "Group version info"
    check_file "api/v1/zz_generated.deepcopy.go" "Generated deepcopy"
    check_file "controllers/service_controller.go" "Service controller"
    
    # Build files
    check_file "Makefile" "Makefile"
    check_file "Dockerfile" "Dockerfile"
    check_file ".dockerignore" "Docker ignore file"
    check_file "PROJECT" "Kubebuilder project file"
    
    # Configuration files
    check_file ".gitignore" "Git ignore file"
    check_file ".golangci.yml" "Golangci-lint config"
}

# Check documentation
check_documentation() {
    print_header "Checking documentation..."
    
    check_file "README.md" "Main README"
    check_file "CHANGELOG.md" "Changelog"
    check_file "CONTRIBUTING.md" "Contributing guide"
    check_file "LICENSE" "License file"
    check_file "SECURITY.md" "Security policy"
    check_file "PROJECT_STRUCTURE.md" "Project structure"
    check_file "VERSION" "Version file"
    
    # Documentation files
    check_file "docs/getting-started.md" "Getting started guide"
    check_file "docs/api-reference.md" "API reference"
    check_file "docs/development.md" "Development guide"
    check_file "docs/deployment.md" "Deployment guide"
    check_file "docs/architecture.md" "Architecture documentation"
    check_file "docs/troubleshooting.md" "Troubleshooting guide"
    check_file "docs/faq.md" "FAQ"
    check_file "docs/project-summary.md" "Project summary"
}

# Check scripts
check_scripts() {
    print_header "Checking scripts..."
    
    check_executable "scripts/install.sh" "Install script"
    check_executable "scripts/uninstall.sh" "Uninstall script"
    check_executable "scripts/dev-setup.sh" "Development setup script"
    check_executable "scripts/helm-install.sh" "Helm install script"
    check_executable "scripts/e2e-test.sh" "E2E test script"
    check_executable "scripts/validate-manifests.sh" "Manifest validation script"
    check_executable "scripts/generate-install-script.sh" "Install script generator"
    check_executable "scripts/setup-permissions.sh" "Permissions setup script"
    check_executable "scripts/check-project.sh" "Project check script"
}

# Check Kubernetes manifests
check_k8s_manifests() {
    print_header "Checking Kubernetes manifests..."
    
    # CRD
    check_file "config/crd/bases/apps.example.com_services.yaml" "Service CRD"
    check_file "config/crd/kustomization.yaml" "CRD kustomization"
    
    # RBAC
    check_file "config/rbac/role.yaml" "RBAC role"
    check_file "config/rbac/role_binding.yaml" "RBAC role binding"
    check_file "config/rbac/service_account.yaml" "Service account"
    check_file "config/rbac/leader_election_role.yaml" "Leader election role"
    check_file "config/rbac/leader_election_role_binding.yaml" "Leader election role binding"
    check_file "config/rbac/kustomization.yaml" "RBAC kustomization"
    
    # Manager
    check_file "config/manager/manager.yaml" "Manager deployment"
    check_file "config/manager/kustomization.yaml" "Manager kustomization"
    check_file "config/manager/controller_manager_config.yaml" "Controller manager config"
    
    # Default
    check_file "config/default/kustomization.yaml" "Default kustomization"
}

# Check Helm chart
check_helm_chart() {
    print_header "Checking Helm chart..."
    
    check_file "deploy/helm/service-operator/Chart.yaml" "Helm Chart.yaml"
    check_file "deploy/helm/service-operator/values.yaml" "Helm values.yaml"
    check_file "deploy/helm/service-operator/templates/_helpers.tpl" "Helm helpers"
    check_file "deploy/helm/service-operator/templates/deployment.yaml" "Helm deployment template"
    check_file "deploy/helm/service-operator/templates/rbac.yaml" "Helm RBAC template"
    check_file "deploy/helm/service-operator/templates/service.yaml" "Helm service template"
    check_file "deploy/helm/service-operator/templates/serviceaccount.yaml" "Helm service account template"
    check_file "deploy/helm/service-operator/crds/apps.example.com_services.yaml" "Helm CRD"
}

# Check examples
check_examples() {
    print_header "Checking examples..."
    
    check_file "examples/sample-service.yaml" "Sample service example"
    check_file "examples/database-service.yaml" "Database service example"
    check_file "examples/microservice-with-ingress.yaml" "Microservice with ingress example"
    check_file "examples/complete-application.yaml" "Complete application example"
    check_file "examples/namespace.yaml" "Namespace example"
}

# Check GitHub templates
check_github_templates() {
    print_header "Checking GitHub templates..."
    
    check_file ".github/workflows/ci.yml" "CI workflow"
    check_file ".github/workflows/release.yml" "Release workflow"
    check_file ".github/ISSUE_TEMPLATE/bug_report.md" "Bug report template"
    check_file ".github/ISSUE_TEMPLATE/feature_request.md" "Feature request template"
    check_file ".github/PULL_REQUEST_TEMPLATE.md" "Pull request template"
}

# Check tests
check_tests() {
    print_header "Checking tests..."
    
    check_file "controllers/service_controller_test.go" "Controller tests"
    check_file "controllers/suite_test.go" "Test suite"
}

# Check file permissions
check_permissions() {
    print_header "Checking file permissions..."
    
    ((TOTAL_CHECKS++))
    local script_count=0
    local executable_count=0
    
    # Count shell scripts
    while IFS= read -r -d '' file; do
        ((script_count++))
        if [ -x "$file" ]; then
            ((executable_count++))
        fi
    done < <(find scripts/ -name "*.sh" -type f -print0 2>/dev/null)
    
    if [ $script_count -eq $executable_count ] && [ $script_count -gt 0 ]; then
        print_status "All shell scripts are executable ($executable_count/$script_count)"
    elif [ $script_count -gt 0 ]; then
        print_error "Some shell scripts are not executable ($executable_count/$script_count)"
    else
        print_warning "No shell scripts found"
    fi
}

# Check consistency
check_consistency() {
    print_header "Checking consistency..."
    
    # Check if VERSION file matches
    ((TOTAL_CHECKS++))
    if [ -f "VERSION" ]; then
        local version
        version=$(cat VERSION)
        if grep -q "$version" CHANGELOG.md; then
            print_status "Version consistency: VERSION file matches CHANGELOG.md"
        else
            print_warning "Version in VERSION file not found in CHANGELOG.md"
        fi
    else
        print_error "VERSION file missing"
    fi
    
    # Check if all examples have proper apiVersion
    ((TOTAL_CHECKS++))
    local invalid_examples=0
    for example in examples/*.yaml; do
        if [ -f "$example" ] && ! grep -q "apiVersion: apps.example.com/v1" "$example"; then
            ((invalid_examples++))
        fi
    done
    
    if [ $invalid_examples -eq 0 ]; then
        print_status "All examples use correct apiVersion"
    else
        print_error "$invalid_examples examples have incorrect apiVersion"
    fi
}

# Generate report
generate_report() {
    print_header "Project Check Summary"
    echo ""
    echo "Total checks: $TOTAL_CHECKS"
    echo -e "${GREEN}Passed: $PASSED_CHECKS${NC}"
    echo -e "${YELLOW}Warnings: $WARNING_CHECKS${NC}"
    echo -e "${RED}Failed: $FAILED_CHECKS${NC}"
    echo ""
    
    local success_rate
    success_rate=$((PASSED_CHECKS * 100 / TOTAL_CHECKS))
    
    if [ $FAILED_CHECKS -eq 0 ]; then
        echo -e "${GREEN}✅ Project check completed successfully! (${success_rate}% passed)${NC}"
        return 0
    else
        echo -e "${RED}❌ Project check found issues. (${success_rate}% passed)${NC}"
        return 1
    fi
}

# Main function
main() {
    print_header "Service Operator Project Check"
    echo ""
    
    check_project_structure
    check_core_files
    check_documentation
    check_scripts
    check_k8s_manifests
    check_helm_chart
    check_examples
    check_github_templates
    check_tests
    check_permissions
    check_consistency
    
    echo ""
    generate_report
}

# Run main function
main