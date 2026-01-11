# 快速开始指南

本指南将帮助您快速上手 Service Operator，从安装到创建第一个服务。

## 概述

Service Operator 是一个 Kubernetes Operator，它简化了完整服务的部署和管理。通过单个自定义资源，您可以管理：

- **Deployment** - 应用程序的 Pod 管理
- **ConfigMap** - 配置文件管理
- **Service** - 网络服务暴露
- **Ingress** - 外部访问配置

## 前置条件

在开始之前，请确保您有：

- Kubernetes 集群 (v1.20+)
- kubectl 命令行工具
- 集群管理员权限
- Docker (如果需要构建镜像)

## 第一步：安装 Service Operator

### 方式一：使用安装脚本（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd service-operator

# 运行安装脚本
chmod +x scripts/install.sh
./scripts/install.sh
```

### 方式二：使用 Helm

```bash
# 使用 Helm 安装
chmod +x scripts/helm-install.sh
./scripts/helm-install.sh
```

### 方式三：手动安装

```bash
# 安装 CRD
kubectl apply -f config/crd/bases/apps.example.com_services.yaml

# 部署 Operator
kubectl apply -f config/rbac/
kubectl apply -f config/manager/
```

## 第二步：验证安装

```bash
# 检查 Operator 是否运行
kubectl get pods -n service-operator-system

# 检查 CRD 是否安装
kubectl get crd services.apps.example.com

# 查看 Operator 日志
kubectl logs -n service-operator-system deployment/service-operator-controller-manager
```

预期输出：
```
NAME                                                READY   STATUS    RESTARTS   AGE
service-operator-controller-manager-xxxxxxxxx-xxxxx   2/2     Running   0          1m
```

## 第三步：创建第一个服务

### 创建简单的 Web 服务

创建文件 `my-first-service.yaml`：

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: my-web-app
  namespace: default
spec:
  image: nginx:1.21
  replicas: 2
  port: 80
  serviceType: ClusterIP
```

应用配置：

```bash
kubectl apply -f my-first-service.yaml
```

### 查看创建的资源

```bash
# 查看 Service 资源状态
kubectl get services.apps.example.com

# 查看详细信息
kubectl describe services.apps.example.com my-web-app

# 查看创建的 Kubernetes 资源
kubectl get deployments,services,pods -l app=my-web-app
```

预期输出：
```
NAME        IMAGE        REPLICAS   READY   PHASE   AGE
my-web-app  nginx:1.21   2          2       Ready   1m
```

## 第四步：测试服务

```bash
# 端口转发测试服务
kubectl port-forward service/my-web-app 8080:80

# 在另一个终端测试
curl http://localhost:8080
```

## 第五步：添加配置文件

更新服务以包含自定义配置：

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: my-web-app
  namespace: default
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
              return 200 'Hello from Service Operator!';
              add_header Content-Type text/plain;
          }
      }
  env:
    - name: ENVIRONMENT
      value: "development"
```

应用更新：

```bash
kubectl apply -f my-first-service.yaml
```

查看配置是否生效：

```bash
# 检查 ConfigMap
kubectl get configmap my-web-app-config -o yaml

# 测试更新后的服务
kubectl port-forward service/my-web-app 8080:80
curl http://localhost:8080
```

## 第六步：添加外部访问

如果您的集群有 Ingress Controller，可以添加外部访问：

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: my-web-app
  namespace: default
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
              return 200 'Hello from Service Operator!';
              add_header Content-Type text/plain;
          }
      }
  env:
    - name: ENVIRONMENT
      value: "development"
  ingress:
    enabled: true
    host: my-app.local
    path: /
    annotations:
      kubernetes.io/ingress.class: nginx
```

应用配置：

```bash
kubectl apply -f my-first-service.yaml
```

如果使用 kind 或本地集群，添加 hosts 条目：

```bash
echo "127.0.0.1 my-app.local" | sudo tee -a /etc/hosts
```

然后访问：

```bash
curl http://my-app.local
```

## 第七步：监控服务状态

```bash
# 实时查看服务状态
kubectl get services.apps.example.com my-web-app -w

# 查看服务事件
kubectl get events --field-selector involvedObject.name=my-web-app

# 查看 Pod 日志
kubectl logs -l app=my-web-app
```

## 常用操作

### 扩缩容

```bash
# 扩容到 5 个副本
kubectl patch services.apps.example.com my-web-app -p '{"spec":{"replicas":5}}'

# 缩容到 1 个副本
kubectl patch services.apps.example.com my-web-app -p '{"spec":{"replicas":1}}'
```

### 更新镜像

```bash
# 更新到新版本
kubectl patch services.apps.example.com my-web-app -p '{"spec":{"image":"nginx:1.22"}}'
```

### 更新配置

```bash
# 更新环境变量
kubectl patch services.apps.example.com my-web-app -p '{"spec":{"env":[{"name":"ENVIRONMENT","value":"production"}]}}'
```

### 查看状态

```bash
# 查看服务 URL（如果启用了 Ingress）
kubectl get services.apps.example.com my-web-app -o jsonpath='{.status.url}'

# 查看就绪副本数
kubectl get services.apps.example.com my-web-app -o jsonpath='{.status.readyReplicas}'
```

## 更多示例

### 数据库服务

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: postgres-db
  namespace: default
spec:
  image: postgres:13
  replicas: 1
  port: 5432
  serviceType: ClusterIP
  env:
    - name: POSTGRES_DB
      value: "myapp"
    - name: POSTGRES_USER
      value: "user"
    - name: POSTGRES_PASSWORD
      value: "password"
  resources:
    requests:
      cpu: "250m"
      memory: "512Mi"
    limits:
      cpu: "1000m"
      memory: "2Gi"
```

### 带 TLS 的 Web 服务

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: secure-web-app
  namespace: default
spec:
  image: nginx:1.21
  replicas: 3
  port: 80
  serviceType: ClusterIP
  resources:
    requests:
      cpu: "100m"
      memory: "128Mi"
    limits:
      cpu: "500m"
      memory: "512Mi"
  ingress:
    enabled: true
    host: secure-app.example.com
    path: /
    annotations:
      kubernetes.io/ingress.class: nginx
      cert-manager.io/cluster-issuer: letsencrypt-prod
    tls:
      enabled: true
      secretName: secure-app-tls
```

## 清理资源

当您完成测试后，可以清理创建的资源：

```bash
# 删除服务资源
kubectl delete services.apps.example.com my-web-app

# 卸载 Operator
./scripts/uninstall.sh
```

## 下一步

现在您已经成功创建了第一个服务，可以：

1. 查看 [API 参考文档](api-reference.md) 了解所有可用选项
2. 浏览 [示例目录](../examples/) 查看更多用例
3. 阅读 [部署指南](deployment.md) 了解生产环境部署
4. 查看 [故障排除指南](troubleshooting.md) 解决常见问题

## 获取帮助

如果遇到问题：

1. 查看 [故障排除指南](troubleshooting.md)
2. 搜索现有的 [GitHub Issues](../../issues)
3. 创建新的 Issue 报告问题
4. 参与社区讨论

欢迎使用 Service Operator！🚀