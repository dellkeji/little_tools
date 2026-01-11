# 故障排除指南

本文档提供了 Service Operator 常见问题的诊断和解决方案。

## 常见问题

### 1. Operator 无法启动

#### 症状
- Operator Pod 处于 CrashLoopBackOff 状态
- Pod 日志显示启动错误

#### 可能原因和解决方案

**原因 1: CRD 未安装**
```bash
# 检查 CRD 是否存在
kubectl get crd services.apps.example.com

# 如果不存在，安装 CRD
make install
# 或
kubectl apply -f config/crd/bases/apps.example.com_services.yaml
```

**原因 2: RBAC 权限不足**
```bash
# 检查 ServiceAccount 权限
kubectl auth can-i create deployments --as=system:serviceaccount:service-operator-system:service-operator-controller-manager

# 检查 ClusterRole 是否存在
kubectl get clusterrole service-operator-manager-role

# 重新应用 RBAC 配置
kubectl apply -f config/rbac/
```

**原因 3: 镜像拉取失败**
```bash
# 检查镜像是否存在
docker pull service-operator:latest

# 检查 Pod 事件
kubectl describe pod -n service-operator-system

# 如果是私有镜像，创建 ImagePullSecret
kubectl create secret docker-registry regcred \
  --docker-server=<your-registry-server> \
  --docker-username=<your-name> \
  --docker-password=<your-pword> \
  --docker-email=<your-email> \
  -n service-operator-system
```

### 2. Service 资源创建失败

#### 症状
- Service 资源状态为 Pending
- 相关的 Kubernetes 资源未创建

#### 诊断步骤

```bash
# 查看 Service 资源状态
kubectl get services.apps.example.com <service-name> -o yaml

# 查看 Service 资源事件
kubectl describe services.apps.example.com <service-name>

# 查看 Operator 日志
kubectl logs -n service-operator-system deployment/service-operator-controller-manager -f

# 查看集群事件
kubectl get events --sort-by=.metadata.creationTimestamp
```

#### 常见解决方案

**问题 1: 镜像不存在**
```yaml
# 确保镜像名称正确
spec:
  image: nginx:1.21  # 确保标签存在
```

**问题 2: 资源配置错误**
```yaml
# 检查资源配置格式
spec:
  resources:
    requests:
      cpu: "100m"     # 注意单位
      memory: "128Mi"  # 注意单位
    limits:
      cpu: "500m"
      memory: "512Mi"
```

**问题 3: 命名空间不存在**
```bash
# 创建命名空间
kubectl create namespace <namespace-name>
```

### 3. Ingress 无法访问

#### 症状
- Ingress 创建成功但无法访问
- 返回 404 或 502 错误

#### 诊断步骤

```bash
# 检查 Ingress 状态
kubectl get ingress <service-name>
kubectl describe ingress <service-name>

# 检查 Ingress Controller
kubectl get pods -n ingress-nginx
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller

# 检查 Service 端点
kubectl get endpoints <service-name>

# 测试 Service 连通性
kubectl port-forward service/<service-name> 8080:80
curl http://localhost:8080
```

#### 解决方案

**问题 1: Ingress Controller 未安装**
```bash
# 安装 nginx-ingress
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml
```

**问题 2: DNS 解析问题**
```bash
# 检查 DNS 解析
nslookup <ingress-host>

# 或者修改本地 hosts 文件进行测试
echo "<ingress-ip> <ingress-host>" >> /etc/hosts
```

**问题 3: TLS 证书问题**
```bash
# 检查证书 Secret
kubectl get secret <tls-secret-name> -o yaml

# 如果使用 cert-manager，检查证书状态
kubectl get certificate
kubectl describe certificate <cert-name>
```

### 4. ConfigMap 配置未生效

#### 症状
- ConfigMap 创建成功但应用未读取到配置
- 应用使用默认配置

#### 诊断步骤

```bash
# 检查 ConfigMap 内容
kubectl get configmap <service-name>-config -o yaml

# 检查 Pod 中的挂载
kubectl describe pod <pod-name>
kubectl exec <pod-name> -- ls -la /etc/config

# 检查应用日志
kubectl logs <pod-name>
```

#### 解决方案

**问题 1: 挂载路径错误**
- 确保应用从正确路径读取配置文件
- 检查 volumeMount 配置

**问题 2: 配置格式错误**
- 验证配置文件语法
- 检查字符编码

**问题 3: 应用需要重启**
```bash
# 重启 Deployment
kubectl rollout restart deployment <service-name>
```

### 5. 资源更新不生效

#### 症状
- 修改 Service 资源后，相关资源未更新
- Deployment 使用旧的配置

#### 诊断步骤

```bash
# 检查 Service 资源的 generation
kubectl get services.apps.example.com <service-name> -o jsonpath='{.metadata.generation}'

# 检查 status 中的 observedGeneration
kubectl get services.apps.example.com <service-name> -o jsonpath='{.status.conditions[0].observedGeneration}'

# 查看 Operator 日志中的 reconcile 事件
kubectl logs -n service-operator-system deployment/service-operator-controller-manager | grep "Reconciling"
```

#### 解决方案

**问题 1: Operator 未收到更新事件**
```bash
# 重启 Operator
kubectl rollout restart deployment/service-operator-controller-manager -n service-operator-system
```

**问题 2: 资源被外部修改**
```bash
# 检查资源的 ownerReferences
kubectl get deployment <service-name> -o yaml | grep -A 10 ownerReferences

# 如果 ownerReferences 丢失，删除资源让 Operator 重新创建
kubectl delete deployment <service-name>
```

## 性能问题

### 1. Operator 响应缓慢

#### 症状
- 资源创建或更新延迟很高
- Operator CPU/内存使用率高

#### 诊断步骤

```bash
# 检查 Operator 资源使用
kubectl top pod -n service-operator-system

# 检查 Operator 配置
kubectl get deployment service-operator-controller-manager -n service-operator-system -o yaml

# 查看工作队列指标
kubectl port-forward -n service-operator-system service/service-operator-controller-manager-metrics-service 8443:8443
curl -k https://localhost:8443/metrics | grep workqueue
```

#### 解决方案

**增加资源限制**
```yaml
resources:
  limits:
    cpu: 1000m
    memory: 512Mi
  requests:
    cpu: 200m
    memory: 256Mi
```

**调整并发设置**
```yaml
args:
- --max-concurrent-reconciles=5
```

### 2. 大量资源管理

#### 症状
- 管理大量 Service 资源时性能下降
- 内存使用持续增长

#### 解决方案

**启用分片**
```yaml
args:
- --enable-leader-election
- --leader-election-namespace=service-operator-system
```

**优化缓存**
```yaml
args:
- --cache-sync-timeout=10m
```

## 调试技巧

### 1. 启用详细日志

```bash
# 临时启用 debug 日志
kubectl patch deployment service-operator-controller-manager \
  -n service-operator-system \
  -p '{"spec":{"template":{"spec":{"containers":[{"name":"manager","args":["--zap-log-level=debug"]}]}}}}'

# 恢复正常日志级别
kubectl patch deployment service-operator-controller-manager \
  -n service-operator-system \
  -p '{"spec":{"template":{"spec":{"containers":[{"name":"manager","args":["--zap-log-level=info"]}]}}}}'
```

### 2. 使用事件查看器

```bash
# 实时查看事件
kubectl get events --watch

# 查看特定资源的事件
kubectl get events --field-selector involvedObject.name=<service-name>
```

### 3. 资源状态检查脚本

创建一个检查脚本 `debug-service.sh`：

```bash
#!/bin/bash

SERVICE_NAME=$1
NAMESPACE=${2:-default}

if [ -z "$SERVICE_NAME" ]; then
    echo "Usage: $0 <service-name> [namespace]"
    exit 1
fi

echo "=== Service Resource ==="
kubectl get services.apps.example.com $SERVICE_NAME -n $NAMESPACE -o yaml

echo -e "\n=== Related Resources ==="
kubectl get deployment,service,configmap,ingress -l app=$SERVICE_NAME -n $NAMESPACE

echo -e "\n=== Pod Status ==="
kubectl get pods -l app=$SERVICE_NAME -n $NAMESPACE

echo -e "\n=== Events ==="
kubectl get events --field-selector involvedObject.name=$SERVICE_NAME -n $NAMESPACE --sort-by=.metadata.creationTimestamp

echo -e "\n=== Operator Logs (last 50 lines) ==="
kubectl logs -n service-operator-system deployment/service-operator-controller-manager --tail=50 | grep $SERVICE_NAME
```

使用方法：
```bash
chmod +x debug-service.sh
./debug-service.sh my-service default
```

## 监控和告警

### 1. 关键指标监控

如果使用 Prometheus，监控以下指标：

```yaml
# Operator 健康状态
up{job="service-operator"}

# 控制器队列长度
workqueue_depth{name="service"}

# 调谐频率
controller_runtime_reconcile_total{controller="service"}

# 调谐错误率
rate(controller_runtime_reconcile_errors_total{controller="service"}[5m])

# 调谐延迟
histogram_quantile(0.95, controller_runtime_reconcile_time_seconds_bucket{controller="service"})
```

### 2. 告警规则

```yaml
groups:
- name: service-operator
  rules:
  - alert: ServiceOperatorDown
    expr: up{job="service-operator"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Service Operator is down"
      
  - alert: ServiceOperatorHighErrorRate
    expr: rate(controller_runtime_reconcile_errors_total{controller="service"}[5m]) > 0.1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Service Operator has high error rate"
      
  - alert: ServiceOperatorHighLatency
    expr: histogram_quantile(0.95, controller_runtime_reconcile_time_seconds_bucket{controller="service"}) > 30
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Service Operator has high reconcile latency"
```

## 获取帮助

### 1. 收集诊断信息

在报告问题时，请收集以下信息：

```bash
# 创建诊断信息收集脚本
cat > collect-diagnostics.sh << 'EOF'
#!/bin/bash

NAMESPACE="service-operator-system"
OUTPUT_DIR="service-operator-diagnostics-$(date +%Y%m%d-%H%M%S)"

mkdir -p $OUTPUT_DIR

echo "Collecting Service Operator diagnostics..."

# 基本信息
kubectl version > $OUTPUT_DIR/kubectl-version.txt
kubectl get nodes -o wide > $OUTPUT_DIR/nodes.txt

# CRD 信息
kubectl get crd services.apps.example.com -o yaml > $OUTPUT_DIR/crd.yaml

# Operator 信息
kubectl get all -n $NAMESPACE > $OUTPUT_DIR/operator-resources.txt
kubectl describe deployment service-operator-controller-manager -n $NAMESPACE > $OUTPUT_DIR/operator-deployment.txt
kubectl logs -n $NAMESPACE deployment/service-operator-controller-manager --tail=1000 > $OUTPUT_DIR/operator-logs.txt

# Service 资源
kubectl get services.apps.example.com --all-namespaces -o yaml > $OUTPUT_DIR/service-resources.yaml

# 事件
kubectl get events --all-namespaces --sort-by=.metadata.creationTimestamp > $OUTPUT_DIR/events.txt

echo "Diagnostics collected in $OUTPUT_DIR/"
tar -czf $OUTPUT_DIR.tar.gz $OUTPUT_DIR/
echo "Archive created: $OUTPUT_DIR.tar.gz"
EOF

chmod +x collect-diagnostics.sh
./collect-diagnostics.sh
```

### 2. 社区支持

- GitHub Issues: 报告 bug 和功能请求
- 文档: 查看最新文档和示例
- 社区论坛: 获取社区帮助

### 3. 企业支持

如果需要企业级支持，请联系维护团队获取：
- 优先级支持
- 定制化开发
- 培训服务
- SLA 保证