# 开发指南

本文档介绍如何设置开发环境并参与 Service Operator 的开发。

## 开发环境要求

### 必需工具

- **Go 1.21+**: 用于编译和运行 Operator
- **Docker**: 用于构建容器镜像
- **kubectl**: 用于与 Kubernetes 集群交互
- **kind** (推荐): 用于本地 Kubernetes 集群

### 可选工具

- **kustomize**: 用于 Kubernetes 资源管理 (会自动安装)
- **controller-gen**: 用于代码生成 (会自动安装)
- **envtest**: 用于运行测试 (会自动安装)

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd service-operator
```

### 2. 自动设置开发环境

运行开发环境设置脚本：

```bash
chmod +x scripts/dev-setup.sh
./scripts/dev-setup.sh
```

这个脚本会：
- 检查必需的工具
- 创建 kind 集群
- 安装开发工具
- 下载依赖
- 生成代码
- 运行测试
- 安装 CRD

### 3. 手动设置 (可选)

如果你想手动设置环境：

```bash
# 下载依赖
go mod download

# 安装开发工具
make controller-gen kustomize envtest

# 生成代码
make generate manifests

# 运行测试
make test

# 创建 kind 集群
kind create cluster --name service-operator-dev

# 安装 CRD
make install
```

## 开发工作流

### 本地运行 Operator

```bash
# 在本地运行 Operator (连接到集群)
make run
```

这会启动 Operator 进程，连接到当前 kubectl 上下文的集群。

### 测试更改

在另一个终端中：

```bash
# 应用示例
kubectl apply -f examples/sample-service.yaml

# 查看状态
kubectl get services.apps.example.com

# 查看详细信息
kubectl describe services.apps.example.com sample-web-service

# 查看创建的资源
kubectl get deployments,services,configmaps,ingresses
```

### 代码更改后的步骤

1. **更新 API 定义** (如果修改了 `api/v1/service_types.go`):
   ```bash
   make generate manifests
   make install  # 更新 CRD
   ```

2. **更新控制器逻辑** (如果修改了 `controllers/service_controller.go`):
   ```bash
   # 重启 make run 进程即可
   ```

3. **运行测试**:
   ```bash
   make test
   ```

## 项目结构

```
service-operator/
├── api/v1/                     # API 定义
│   ├── service_types.go        # Service CRD 定义
│   ├── groupversion_info.go    # API 组版本信息
│   └── zz_generated.deepcopy.go # 自动生成的深拷贝代码
├── controllers/                # 控制器实现
│   ├── service_controller.go   # 主控制器逻辑
│   ├── service_controller_test.go # 控制器测试
│   └── suite_test.go          # 测试套件
├── config/                     # Kubernetes 配置
│   ├── crd/                   # CRD 定义
│   ├── rbac/                  # RBAC 配置
│   ├── manager/               # Operator 部署配置
│   └── default/               # 默认配置
├── examples/                   # 使用示例
├── scripts/                    # 脚本工具
├── docs/                      # 文档
├── main.go                    # 入口文件
├── Makefile                   # 构建脚本
├── Dockerfile                 # 容器镜像构建
└── README.md                  # 项目说明
```

## 常用开发命令

### 构建和测试

```bash
# 格式化代码
make fmt

# 静态检查
make vet

# 运行测试
make test

# 构建二进制文件
make build

# 构建 Docker 镜像
make docker-build IMG=service-operator:dev
```

### 代码生成

```bash
# 生成深拷贝代码
make generate

# 生成 CRD 和 RBAC 配置
make manifests
```

### 集群操作

```bash
# 安装 CRD
make install

# 卸载 CRD
make uninstall

# 部署 Operator
make deploy IMG=service-operator:dev

# 卸载 Operator
make undeploy
```

## 测试

### 单元测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./controllers/...

# 运行测试并查看覆盖率
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out
```

### 集成测试

集成测试使用 envtest 框架，会启动一个真实的 etcd 和 kube-apiserver：

```bash
# 运行集成测试
make test
```

### 手动测试

```bash
# 启动 Operator
make run

# 在另一个终端应用测试资源
kubectl apply -f examples/sample-service.yaml

# 验证资源创建
kubectl get services.apps.example.com
kubectl get deployments,services,configmaps

# 测试更新
kubectl patch services.apps.example.com sample-web-service -p '{"spec":{"replicas":3}}'

# 清理
kubectl delete -f examples/sample-service.yaml
```

## 调试

### 查看日志

```bash
# 本地运行时，日志会直接输出到终端

# 如果部署到集群，查看 Operator 日志
kubectl logs -n service-operator-system deployment/service-operator-controller-manager -f
```

### 使用调试器

可以使用 VS Code 或 GoLand 等 IDE 的调试功能：

1. 设置断点
2. 使用调试模式运行 `make run`
3. 应用测试资源触发断点

### 常见问题排查

1. **CRD 未安装**:
   ```bash
   make install
   ```

2. **权限问题**:
   检查 `config/rbac/role.yaml` 中的权限配置

3. **资源创建失败**:
   查看 Operator 日志和事件：
   ```bash
   kubectl describe services.apps.example.com <resource-name>
   kubectl get events --sort-by=.metadata.creationTimestamp
   ```

## 贡献代码

### 代码规范

1. **Go 代码风格**: 遵循标准 Go 代码风格
2. **注释**: 为公共 API 添加适当的注释
3. **测试**: 为新功能添加测试
4. **文档**: 更新相关文档

### 提交流程

1. Fork 项目
2. 创建功能分支: `git checkout -b feature/new-feature`
3. 提交更改: `git commit -am 'Add new feature'`
4. 推送分支: `git push origin feature/new-feature`
5. 创建 Pull Request

### 代码检查

提交前运行：

```bash
# 格式化和检查
make fmt vet

# 运行测试
make test

# 生成最新的代码和配置
make generate manifests
```

## 发布流程

### 构建发布版本

```bash
# 构建并推送镜像
make docker-build docker-push IMG=myregistry/service-operator:v1.0.0

# 生成发布配置
cd config/manager && kustomize edit set image controller=myregistry/service-operator:v1.0.0
kustomize build config/default > service-operator-v1.0.0.yaml
```

### 版本管理

1. 更新版本号
2. 创建 Git 标签
3. 构建和推送镜像
4. 创建 GitHub Release

## 高级开发主题

### 添加新的 API 字段

1. 修改 `api/v1/service_types.go`
2. 运行 `make generate manifests`
3. 更新控制器逻辑
4. 添加测试
5. 更新文档

### 添加新的控制器

1. 创建新的控制器文件
2. 在 `main.go` 中注册控制器
3. 添加 RBAC 权限
4. 添加测试

### Webhook 支持

如果需要添加 admission webhook：

1. 使用 kubebuilder 生成 webhook 代码
2. 实现验证和变更逻辑
3. 配置证书管理
4. 更新部署配置

这个开发指南应该能帮助你快速上手 Service Operator 的开发。如果有任何问题，请查看项目文档或提交 Issue。