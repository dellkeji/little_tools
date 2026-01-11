---
name: Bug report
about: Create a report to help us improve
title: '[BUG] '
labels: 'bug'
assignees: ''

---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Apply the following Service resource: '...'
2. Run command '....'
3. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Actual behavior**
A clear and concise description of what actually happened.

**Service Resource**
```yaml
# Please provide the Service resource that caused the issue
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: example
spec:
  # ... your configuration
```

**Environment (please complete the following information):**
- Kubernetes version: [e.g. v1.28.0]
- Service Operator version: [e.g. v0.1.0]
- Installation method: [e.g. Helm, Kustomize, script]
- Operating System: [e.g. Ubuntu 20.04]
- Container Runtime: [e.g. containerd, docker]

**Logs**
```
# Operator logs
kubectl logs -n service-operator-system deployment/service-operator-controller-manager

# Service resource status
kubectl describe services.apps.example.com <service-name>

# Related resource status
kubectl get deployments,services,configmaps,ingresses -l app=<service-name>
```

**Additional context**
Add any other context about the problem here, such as:
- Screenshots
- Related Issues
- Workarounds you've tried
- Impact on your application

**Checklist**
- [ ] I have searched existing issues to ensure this is not a duplicate
- [ ] I have provided all the requested information
- [ ] I have tested with the latest version of Service Operator
- [ ] I have included relevant logs and resource definitions