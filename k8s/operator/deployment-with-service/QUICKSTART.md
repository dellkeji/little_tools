# Service Operator 快速开始

这是一个 5 分钟快速开始指南，帮助您立即体验 Service Operator。

## 🚀 一键体验

### 前置条件
- Kubernetes 集群 (v1.20+)
- kubectl 已配置

### 1. 安装 Service Operator

```bash
# 方式一：使用安装脚本（推荐）
curl -sSL https://raw.githubusercontent.com/example/service-operator/main/scripts/install.sh | bash

# 方式二：手动安装
kubectl apply -f https://github.com/example/service-operator/releases/latest/download/install.yaml
```

### 2. 验证安装

```bash
# 检查 Operator 状态
kubectl get pods -n service-operator-system

# 检查 CRD
kubectl get crd services.apps.example.com
```

### 3. 创建第一个服务

```bash
# 创建示例服务
cat <<EOF | kubectl apply -f -
apiVersion: apps.example.com/v1
kind: Service
metadata:
  name: hello-world
  namespace: default
spec:
  image: nginx:1.21
  replicas: 2
  port: 80
  configData:
    index.html: |
      <!DOCTYPE html>
      <html>
      <head><title>Hello Service Operator!</title></head>
      <body>
        <h1>🎉 Service Operator 工作正常！</h1>
        <p>这是通过 Service Operator 部署的服务</p>
      </body>
      </html>
EOF
```

### 4. 查看结果

```bash
# 查看服务状态
kubectl get services.apps.example.com

# 查看创建的资源
kubectl get deployments,services,configmaps -l app=hello-world

# 测试服务
kubectl port-forward service/hello-world 8080:80
# 在浏览器中访问 http://localhost:8080
```

### 5. 清理资源

```bash
# 删除服务
kubectl delete services.apps.example.com hello-world

# 卸载 Operator（可选）
curl -sSL https://raw.githubusercontent.com/example/service-operator/main/scripts/uninstall.sh | bash
```

## 🎯 下一步

- 查看 [完整文档](docs/getting-started.md)
- 浏览 [更多示例](examples/)
- 了解 [API 参考](docs/api-reference.md)

## 🆘 遇到问题？

- 查看 [故障排除指南](docs/troubleshooting.md)
- 查看 [常见问题](docs/faq.md)
- 提交 [GitHub Issue](../../issues)

---

**恭喜！** 您已经成功体验了 Service Operator！🎉