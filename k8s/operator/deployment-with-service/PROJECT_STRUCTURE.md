# 项目结构

本文档描述了 Service Operator 项目的完整目录结构和文件说明。

```
service-operator/
├── api/                                    # API 定义
│   └── v1/                                # v1 版本 API
│       ├── groupversion_info.go           # API 组版本信息
│       ├── service_types.go               # Service CRD 类型定义
│       └── zz_generated.deepcopy.go       # 自动生成的深拷贝代码
├── config/                                # Kubernetes 配置文件
│   ├── crd/                              # 自定义资源定义
│   │   ├── bases/                        # CRD 基础定义
│   │   │   └── apps.example.com_services.yaml
│   │   ├── kustomization.yaml            # CRD Kustomize 配置
│   │   └── kustomizeconfig.yaml          # Kustomize 配置选项
│   ├── default/                          # 默认部署配置
│   │   ├── kustomization.yaml            # 默认 Kustomize 配置
│   │   └── manager_auth_proxy_patch.yaml # 认证代理补丁
│   ├── manager/                          # Operator 管理器配置
│   │   ├── controller_manager_config.yaml # 控制器管理器配置
│   │   ├── kustomization.yaml            # 管理器 Kustomize 配置
│   │   ├── manager.yaml                  # Operator 部署配置
│   │   └── metrics_service.yaml          # 指标服务配置
│   ├── prometheus/                       # Prometheus 监控配置
│   │   ├── kustomization.yaml            # Prometheus Kustomize 配置
│   │   └── monitor.yaml                  # ServiceMonitor 配置
│   └── rbac/                            # RBAC 权限配置
│       ├── kustomization.yaml            # RBAC Kustomize 配置
│       ├── leader_election_role.yaml     # Leader 选举角色
│       ├── leader_election_role_binding.yaml # Leader 选举角色绑定
│       ├── role.yaml                     # 集群角色
│       ├── role_binding.yaml             # 角色绑定
│       └── service_account.yaml          # 服务账户
├── controllers/                          # 控制器实现
│   ├── service_controller.go             # 主控制器逻辑
│   ├── service_controller_test.go        # 控制器测试
│   └── suite_test.go                     # 测试套件
├── deploy/                               # 部署配置
│   └── helm/                            # Helm Chart
│       └── service-operator/            # Helm Chart 目录
│           ├── Chart.yaml               # Chart 元数据
│           ├── values.yaml              # 默认配置值
│           ├── crds/                    # CRD 定义
│           │   └── apps.example.com_services.yaml
│           └── templates/               # Helm 模板
│               ├── _helpers.tpl         # 模板助手函数
│               ├── deployment.yaml      # Deployment 模板
│               ├── rbac.yaml           # RBAC 模板
│               ├── service.yaml        # Service 模板
│               ├── serviceaccount.yaml # ServiceAccount 模板
│               └── servicemonitor.yaml # ServiceMonitor 模板
├── docs/                                # 项目文档
│   ├── api-reference.md                 # API 参考文档
│   ├── architecture.md                  # 架构设计文档
│   ├── deployment.md                    # 部署指南
│   ├── development.md                   # 开发指南
│   ├── getting-started.md               # 快速开始指南
│   └── troubleshooting.md               # 故障排除指南
├── examples/                            # 使用示例
│   ├── complete-application.yaml        # 完整应用示例
│   ├── database-service.yaml            # 数据库服务示例
│   ├── microservice-with-ingress.yaml   # 带 Ingress 的微服务示例
│   ├── namespace.yaml                   # 命名空间示例
│   └── sample-service.yaml              # 基础服务示例
├── hack/                                # 构建和开发脚本
│   └── boilerplate.go.txt               # Go 文件头部模板
├── scripts/                             # 自动化脚本
│   ├── dev-setup.sh                     # 开发环境设置脚本
│   ├── helm-install.sh                  # Helm 安装脚本
│   ├── install.sh                       # 安装脚本
│   └── uninstall.sh                     # 卸载脚本
├── .gitignore                           # Git 忽略文件配置
├── CHANGELOG.md                         # 变更日志
├── CONTRIBUTING.md                      # 贡献指南
├── Dockerfile                           # Docker 镜像构建文件
├── LICENSE                              # 开源许可证
├── Makefile                             # 构建和开发任务
├── PROJECT                              # Kubebuilder 项目配置
├── PROJECT_STRUCTURE.md                 # 项目结构说明（本文件）
├── README.md                            # 项目说明文档
├── go.mod                               # Go 模块定义
├── go.sum                               # Go 模块校验和
└── main.go                              # 程序入口点
```

## 核心文件说明

### API 定义 (`api/v1/`)

- **`service_types.go`**: 定义了 Service CRD 的 Go 结构体，包括 ServiceSpec、ServiceStatus 等
- **`groupversion_info.go`**: 定义 API 组和版本信息
- **`zz_generated.deepcopy.go`**: 由 controller-gen 自动生成的深拷贝方法

### 控制器 (`controllers/`)

- **`service_controller.go`**: 核心控制器逻辑，实现 Reconcile 方法
- **`service_controller_test.go`**: 控制器的单元测试和集成测试
- **`suite_test.go`**: 测试套件设置，使用 envtest 框架

### 配置文件 (`config/`)

- **`crd/`**: 包含 CRD 定义和相关配置
- **`rbac/`**: RBAC 权限配置，定义 Operator 所需的权限
- **`manager/`**: Operator 部署配置，包括 Deployment 和 Service
- **`default/`**: 默认的 Kustomize 配置，组合所有组件
- **`prometheus/`**: Prometheus 监控配置

### Helm Chart (`deploy/helm/`)

- **`Chart.yaml`**: Helm Chart 元数据
- **`values.yaml`**: 可配置的默认值
- **`templates/`**: Kubernetes 资源模板
- **`crds/`**: CRD 定义（用于 Helm 安装）

### 文档 (`docs/`)

- **`getting-started.md`**: 新用户快速上手指南
- **`api-reference.md`**: 完整的 API 文档
- **`development.md`**: 开发环境设置和贡献指南
- **`deployment.md`**: 生产环境部署指南
- **`architecture.md`**: 系统架构和设计文档
- **`troubleshooting.md`**: 常见问题和解决方案

### 示例 (`examples/`)

- **`sample-service.yaml`**: 基础的 Web 服务示例
- **`database-service.yaml`**: 数据库服务配置示例
- **`microservice-with-ingress.yaml`**: 带 Ingress 的微服务示例
- **`complete-application.yaml`**: 完整的多服务应用示例

### 脚本 (`scripts/`)

- **`install.sh`**: 一键安装脚本，支持多种安装选项
- **`uninstall.sh`**: 卸载脚本，清理所有资源
- **`helm-install.sh`**: Helm 安装脚本
- **`dev-setup.sh`**: 开发环境自动设置脚本

### 构建文件

- **`Makefile`**: 定义了构建、测试、部署等常用任务
- **`Dockerfile`**: 多阶段 Docker 构建文件
- **`go.mod`** 和 **`go.sum`**: Go 模块依赖管理

### 项目配置

- **`PROJECT`**: Kubebuilder 项目配置文件
- **`.gitignore`**: Git 版本控制忽略规则
- **`LICENSE`**: Apache 2.0 开源许可证

### 文档文件

- **`README.md`**: 项目主要说明文档
- **`CONTRIBUTING.md`**: 贡献者指南
- **`CHANGELOG.md`**: 版本变更记录

## 开发工作流

### 1. 修改 API

如果需要修改 API 定义：

1. 编辑 `api/v1/service_types.go`
2. 运行 `make generate` 生成深拷贝代码
3. 运行 `make manifests` 生成 CRD
4. 更新相关文档和示例

### 2. 修改控制器逻辑

如果需要修改控制器：

1. 编辑 `controllers/service_controller.go`
2. 添加或修改测试 `controllers/service_controller_test.go`
3. 运行 `make test` 验证更改
4. 更新文档

### 3. 添加新功能

添加新功能的典型流程：

1. 在 `api/v1/service_types.go` 中添加新字段
2. 在 `controllers/service_controller.go` 中实现逻辑
3. 添加测试用例
4. 更新 CRD 和 RBAC 配置
5. 添加示例和文档
6. 更新 Helm Chart

### 4. 发布新版本

发布流程：

1. 更新 `CHANGELOG.md`
2. 更新版本号
3. 创建 Git 标签
4. 构建和推送 Docker 镜像
5. 更新 Helm Chart 版本
6. 创建 GitHub Release

## 文件生成规则

### 自动生成的文件

这些文件由工具自动生成，不应手动编辑：

- `api/v1/zz_generated.deepcopy.go` - 由 `controller-gen` 生成
- `config/crd/bases/*.yaml` - 由 `controller-gen` 生成
- `config/rbac/role.yaml` - 由 `controller-gen` 生成

### 模板文件

这些文件作为模板使用：

- `hack/boilerplate.go.txt` - Go 文件头部模板
- `deploy/helm/service-operator/templates/*.yaml` - Helm 模板

## 依赖管理

### Go 依赖

主要依赖包括：

- `sigs.k8s.io/controller-runtime` - Controller Runtime 框架
- `k8s.io/api` - Kubernetes API 类型
- `k8s.io/apimachinery` - Kubernetes API 机制
- `k8s.io/client-go` - Kubernetes 客户端

### 开发工具

项目使用的开发工具：

- `controller-gen` - 代码和配置生成
- `kustomize` - Kubernetes 配置管理
- `envtest` - 集成测试框架

这个项目结构遵循了 Kubernetes Operator 的最佳实践，提供了完整的开发、测试、部署和文档支持。