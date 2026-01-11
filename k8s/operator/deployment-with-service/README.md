# Service Operator

[![Go Report Card](https://goreportcard.com/badge/github.com/example/service-operator)](https://goreportcard.com/report/github.com/example/service-operator)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-v1.20+-blue.svg)](https://kubernetes.io/)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)

一个用 Go 编写的 Kubernetes Operator，用于简化和自动化完整服务的部署和管理。通过单个自定义资源 (CRD)，您可以轻松管理 Deployment、ConfigMap、Service 和 Ingress 等多个 Kubernetes 资源。

## ✨ 功能特性

- 🚀 **统一管理**: 通过单个 CRD 管理 Deployment、ConfigMap、Service 和 Ingress
- 🔄 **自动化部署**: 自动创建和管理相关的 Kubernetes 资源
- ⚙️ **配置管理**: 支持通过 ConfigMap 管理应用配置文件
- 🌐 **Ingress 支持**: 可选的 Ingress 配置，支持 TLS 和自定义注解
- 📊 **资源管理**: 支持设置 CPU 和内存资源请求和限制
- 📈 **状态监控**: 实时监控部署状态和就绪副本数
- 🔒 **安全性**: 遵循 Kubernetes 安全最佳实践
- 📚 **丰富文档**: 完整的文档和示例

## 🚀 快速开始

### 前置条件

- Kubernetes 集群 (v1.20+)
- kubectl 配置正确
- 集群管理员权限

### 一键安装

```bash
# 克隆项目
git clone <repository-url>
cd service-operator

# 运行安装脚本
chmod +x scripts/install.sh
./scripts/install.sh
```

### 创建第一个服务

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: my-web-service
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
      value: "production"
  resources:
    requests:
      cpu: "100m"
      memory: "128Mi"
    limits:
      cpu: "500m"
      memory: "512Mi"
  ingress:
    enabled: true
    host: my-app.example.com
    path: /
    annotations:
      kubernetes.io/ingress.class: nginx
```

```bash
# 应用配置
kubectl apply -f my-service.yaml

# 查看状态
kubectl get services.apps.example.com
```

详细的快速开始指南请参考 [Getting Started](docs/getting-started.md)。

## 📖 文档

| 文档 | 描述 |
|------|------|
| [快速开始](docs/getting-started.md) | 从安装到创建第一个服务的完整指南 |
| [API 参考](docs/api-reference.md) | 完整的 API 规范和字段说明 |
| [部署指南](docs/deployment.md) | 生产环境部署的最佳实践 |
| [开发指南](docs/development.md) | 开发环境设置和贡献指南 |
| [架构设计](docs/architecture.md) | 系统架构和设计原理 |
| [故障排除](docs/troubleshooting.md) | 常见问题和解决方案 |

## 🏗️ 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                       │
│                                                             │
│  ┌─────────────────┐    ┌─────────────────────────────────┐ │
│  │  Service CRD    │    │      Service Operator          │ │
│  │                 │    │                                 │ │
│  │  ┌───────────┐  │    │  ┌─────────────────────────────┐│ │
│  │  │ Service   │◄─┼────┼──┤   Service Controller        ││ │
│  │  │ Resource  │  │    │  │                             ││ │
│  │  └───────────┘  │    │  │  - Reconcile Loop           ││ │
│  └─────────────────┘    │  │  - Resource Management      ││ │
│                         │  │  - Status Updates           ││ │
│  ┌─────────────────┐    │  └─────────────────────────────┘│ │
│  │ Managed         │    └─────────────────────────────────┘ │
│  │ Resources       │                                        │
│  │                 │                                        │
│  │ ┌─────────────┐ │                                        │
│  │ │ Deployment  │ │                                        │
│  │ └─────────────┘ │                                        │
│  │ ┌─────────────┐ │                                        │
│  │ │ ConfigMap   │ │                                        │
│  │ └─────────────┘ │                                        │
│  │ ┌─────────────┐ │                                        │
│  │ │ Service     │ │                                        │
│  │ └─────────────┘ │                                        │
│  │ ┌─────────────┐ │                                        │
│  │ │ Ingress     │ │                                        │
│  │ └─────────────┘ │                                        │
│  └─────────────────┘                                        │
└─────────────────────────────────────────────────────────────┘
```

## 🛠️ 安装方式

### 方式一：使用安装脚本（推荐）

```bash
./scripts/install.sh
```

### 方式二：使用 Helm

```bash
./scripts/helm-install.sh
```

### 方式三：使用 Kustomize

```bash
make install  # 安装 CRD
make deploy   # 部署 Operator
```

## 📋 示例

### 基础 Web 服务

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: web-app
spec:
  image: nginx:1.21
  replicas: 3
  port: 80
```

### 数据库服务

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: postgres-db
spec:
  image: postgres:13
  replicas: 1
  port: 5432
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

### 带 Ingress 的微服务

```yaml
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: api-service
spec:
  image: myapp/api:v1.0.0
  replicas: 5
  port: 8080
  configData:
    application.yml: |
      server:
        port: 8080
      database:
        url: jdbc:postgresql://postgres-db:5432/myapp
  ingress:
    enabled: true
    host: api.example.com
    path: /api
    annotations:
      kubernetes.io/ingress.class: nginx
      cert-manager.io/cluster-issuer: letsencrypt-prod
    tls:
      enabled: true
      secretName: api-tls
```

更多示例请查看 [examples](examples/) 目录。

## 🔧 开发

### 设置开发环境

```bash
# 自动设置开发环境
./scripts/dev-setup.sh
```

### 常用命令

```bash
# 构建
make build

# 运行测试
make test

# 本地运行
make run

# 构建镜像
make docker-build IMG=service-operator:dev

# 生成代码
make generate manifests
```

## 📊 监控

Service Operator 提供 Prometheus 指标和健康检查端点：

- **指标端点**: `:8080/metrics`
- **健康检查**: `:8081/healthz`
- **就绪检查**: `:8081/readyz`

## 🤝 贡献

我们欢迎各种形式的贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何参与项目。

### 贡献者

感谢所有为项目做出贡献的人！

## 📄 许可证

本项目采用 [Apache License 2.0](LICENSE) 许可证。

## 🆘 获取帮助

- 📖 查看 [文档](docs/)
- 🐛 报告 [Issues](../../issues)
- 💬 参与 [Discussions](../../discussions)
- 📧 联系维护者

## ⭐ Star History

如果这个项目对您有帮助，请给我们一个 Star！

---

**Service Operator** - 让 Kubernetes 服务部署变得简单 🚀