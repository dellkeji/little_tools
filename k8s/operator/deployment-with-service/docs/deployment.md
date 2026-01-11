# 部署指南

本文档介绍如何在不同环境中部署 Service Operator。

## 部署方式

Service Operator 支持多种部署方式：

1. **Kustomize 部署** - 使用 Kubernetes 原生的 Kustomize
2. **Helm 部署** - 使用 Helm Chart
3. **脚本部署** - 使用提供的安装脚本
4. **手动部署** - 直接应用 YAML 文件

## 前置条件

- Kubernetes 集群 v1.20+
- kubectl 配置正确
- 集群管理员权限（用于创建 CRD 和 ClusterRole）

## 方式一：使用安装脚本（推荐）

### 快速安装

```bash
# 克隆项目
git clone <repository-url>
cd service-operator

# 运行安装脚本
chmod +x scripts/install.sh
./scripts/install.sh
```

### 自定义安装

```bash
# 使用自定义镜像
./scripts/install.sh -i myregistry/service-operator:v1.0.0

# 跳过镜像构建（如果镜像已存在）
./scripts/install.sh -s

# 查看帮助
./scripts/install.sh -h
```

### 卸载

```bash
chmod +x scripts/uninstall.sh
./scripts/uninstall.sh
```

## 方式二：使用 Helm

### 安装 Helm Chart

```bash
# 基本安装
chmod +x scripts/helm-install.sh
./scripts/helm-install.sh

# 自定义安装
./scripts/helm-install.sh -n my-namespace -r my-operator -t v1.0.0

# 使用自定义配置文件
./scripts/helm-install.sh -f my-values.yaml

# 升级现有安装
./scripts/helm-install.sh --upgrade -t v1.1.0
```

### 手动 Helm 安装

```bash
# 添加并安装 Chart
helm install service-operator deploy/helm/service-operator \
  --namespace service-operator-system \
  --create-namespace \
  --set image.tag=latest

# 升级
helm upgrade service-operator deploy/helm/service-operator \
  --namespace service-operator-system \
  --set image.tag=v1.0.0

# 卸载
helm uninstall service-operator -n service-operator-system
```

### Helm 配置选项

创建自定义的 `values.yaml` 文件：

```yaml
# values.yaml
image:
  repository: myregistry/service-operator
  tag: v1.0.0
  pullPolicy: Always

replicaCount: 2

resources:
  limits:
    cpu: 1000m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi

# 启用监控
monitoring:
  serviceMonitor:
    enabled: true
    namespace: monitoring
    interval: 30s

# 节点选择器
nodeSelector:
  kubernetes.io/os: linux

# 容忍度
tolerations:
- key: "node-role.kubernetes.io/master"
  operator: "Exists"
  effect: "NoSchedule"
```

然后使用：

```bash
helm install service-operator deploy/helm/service-operator \
  -f values.yaml \
  --namespace service-operator-system \
  --create-namespace
```

## 方式三：使用 Kustomize

### 基本部署

```bash
# 构建并应用
kubectl apply -k config/default

# 或者分步骤
make install  # 安装 CRD
make deploy   # 部署 Operator
```

### 自定义部署

创建自定义的 kustomization.yaml：

```yaml
# my-kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: my-operator-system

resources:
- github.com/example/service-operator/config/default

images:
- name: controller
  newName: myregistry/service-operator
  newTag: v1.0.0

patchesStrategicMerge:
- |-
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: controller-manager
    namespace: system
  spec:
    replicas: 2
    template:
      spec:
        containers:
        - name: manager
          resources:
            limits:
              cpu: 1000m
              memory: 256Mi
            requests:
              cpu: 100m
              memory: 128Mi
```

应用自定义配置：

```bash
kubectl apply -k my-kustomization.yaml
```

## 方式四：手动部署

### 1. 安装 CRD

```bash
kubectl apply -f config/crd/bases/apps.example.com_services.yaml
```

### 2. 创建命名空间

```bash
kubectl create namespace service-operator-system
```

### 3. 应用 RBAC

```bash
kubectl apply -f config/rbac/
```

### 4. 部署 Operator

```bash
kubectl apply -f config/manager/
```

## 验证部署

### 检查 Operator 状态

```bash
# 检查 Pod 状态
kubectl get pods -n service-operator-system

# 检查 Deployment 状态
kubectl get deployment -n service-operator-system

# 查看日志
kubectl logs -n service-operator-system deployment/service-operator-controller-manager -f
```

### 检查 CRD

```bash
# 检查 CRD 是否安装
kubectl get crd services.apps.example.com

# 查看 CRD 详情
kubectl describe crd services.apps.example.com
```

### 测试功能

```bash
# 应用示例
kubectl apply -f examples/sample-service.yaml

# 检查创建的资源
kubectl get services.apps.example.com
kubectl get deployments,services,configmaps

# 查看详细状态
kubectl describe services.apps.example.com sample-web-service
```

## 生产环境部署

### 安全配置

1. **使用非 root 用户**：
   ```yaml
   securityContext:
     runAsNonRoot: true
     runAsUser: 65532
   ```

2. **限制权限**：
   ```yaml
   securityContext:
     allowPrivilegeEscalation: false
     capabilities:
       drop:
       - ALL
     readOnlyRootFilesystem: true
   ```

3. **网络策略**：
   ```yaml
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: service-operator-netpol
     namespace: service-operator-system
   spec:
     podSelector:
       matchLabels:
         control-plane: controller-manager
     policyTypes:
     - Ingress
     - Egress
     ingress:
     - from:
       - namespaceSelector: {}
       ports:
       - protocol: TCP
         port: 8443
     egress:
     - to: []
       ports:
       - protocol: TCP
         port: 443
       - protocol: TCP
         port: 6443
   ```

### 高可用配置

```yaml
# 多副本部署
spec:
  replicas: 3
  
# Pod 反亲和性
affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchExpressions:
          - key: control-plane
            operator: In
            values:
            - controller-manager
        topologyKey: kubernetes.io/hostname

# 资源限制
resources:
  limits:
    cpu: 1000m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 128Mi

# 健康检查
livenessProbe:
  httpGet:
    path: /healthz
    port: 8081
  initialDelaySeconds: 15
  periodSeconds: 20
readinessProbe:
  httpGet:
    path: /readyz
    port: 8081
  initialDelaySeconds: 5
  periodSeconds: 10
```

### 监控配置

如果使用 Prometheus Operator：

```yaml
# ServiceMonitor
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: service-operator-metrics
  namespace: service-operator-system
spec:
  selector:
    matchLabels:
      control-plane: controller-manager
  endpoints:
  - port: https
    scheme: https
    path: /metrics
    interval: 30s
    bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    tlsConfig:
      insecureSkipVerify: true
```

### 日志配置

```yaml
# 配置结构化日志
args:
- --zap-log-level=info
- --zap-encoder=json
- --zap-time-encoding=iso8601
```

## 升级

### 升级 Operator

```bash
# 使用脚本升级
./scripts/install.sh -i service-operator:v1.1.0

# 使用 Helm 升级
helm upgrade service-operator deploy/helm/service-operator \
  --set image.tag=v1.1.0

# 使用 Kustomize 升级
# 更新 kustomization.yaml 中的镜像标签，然后：
kubectl apply -k config/default
```

### 升级 CRD

```bash
# 更新 CRD
make install

# 或手动应用
kubectl apply -f config/crd/bases/apps.example.com_services.yaml
```

## 故障排除

### 常见问题

1. **CRD 未安装**：
   ```bash
   kubectl get crd services.apps.example.com
   # 如果不存在，运行：
   make install
   ```

2. **权限问题**：
   ```bash
   kubectl auth can-i create services.apps.example.com
   kubectl describe clusterrole service-operator-manager-role
   ```

3. **Pod 启动失败**：
   ```bash
   kubectl describe pod -n service-operator-system
   kubectl logs -n service-operator-system deployment/service-operator-controller-manager
   ```

4. **镜像拉取失败**：
   ```bash
   # 检查镜像是否存在
   docker pull service-operator:latest
   
   # 检查 ImagePullSecrets
   kubectl get pods -n service-operator-system -o yaml | grep imagePullSecrets
   ```

### 调试模式

启用详细日志：

```bash
# 修改 Deployment 参数
kubectl patch deployment service-operator-controller-manager \
  -n service-operator-system \
  -p '{"spec":{"template":{"spec":{"containers":[{"name":"manager","args":["--zap-log-level=debug","--health-probe-bind-address=:8081","--metrics-bind-address=127.0.0.1:8080","--leader-elect"]}]}}}}'
```

## 卸载

### 完全卸载

```bash
# 使用脚本卸载
./scripts/uninstall.sh

# 使用 Helm 卸载
helm uninstall service-operator -n service-operator-system

# 使用 Kustomize 卸载
kubectl delete -k config/default

# 手动卸载
kubectl delete -f config/manager/
kubectl delete -f config/rbac/
kubectl delete -f config/crd/bases/apps.example.com_services.yaml
kubectl delete namespace service-operator-system
```

### 保留数据卸载

如果要保留 Service 资源：

```bash
# 只删除 Operator，保留 CRD
kubectl delete -f config/manager/
kubectl delete -f config/rbac/
# 不删除 CRD
```

这样现有的 Service 资源会保留，但不会被管理，直到重新安装 Operator。