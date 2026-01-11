# 常见问题 (FAQ)

本文档回答了关于 Service Operator 的常见问题。

## 一般问题

### Q: Service Operator 是什么？

A: Service Operator 是一个 Kubernetes Operator，用于简化和自动化完整服务的部署和管理。通过单个自定义资源 (CRD)，您可以管理 Deployment、ConfigMap、Service 和 Ingress 等多个 Kubernetes 资源。

### Q: 为什么需要 Service Operator？

A: 传统的 Kubernetes 部署需要管理多个相关的资源文件，Service Operator 将这些资源统一管理，提供：
- 简化的部署流程
- 自动化的资源管理
- 一致的配置体验
- 减少人为错误

### Q: Service Operator 与 Helm 有什么区别？

A: 
- **Helm**: 模板化的包管理工具，主要用于应用打包和分发
- **Service Operator**: 运行时的资源管理器，提供持续的调谐和状态管理

两者可以结合使用：用 Helm 安装 Service Operator，用 Service Operator 管理应用。

## 安装和部署

### Q: 支持哪些 Kubernetes 版本？

A: Service Operator 支持 Kubernetes v1.20 及以上版本。推荐使用 v1.24+ 以获得最佳体验。

### Q: 如何选择安装方式？

A: 
- **脚本安装**: 适合快速试用和开发环境
- **Helm 安装**: 适合生产环境，支持自定义配置
- **Kustomize 安装**: 适合需要深度定制的场景

### Q: 可以在多个命名空间中安装 Service Operator 吗？

A: Service Operator 设计为集群级别的 Operator，通常只需要安装一个实例。它可以管理所有命名空间中的 Service 资源。

### Q: 如何升级 Service Operator？

A: 
```bash
# 使用脚本升级
./scripts/install.sh -i service-operator:v0.2.0

# 使用 Helm 升级
helm upgrade service-operator deploy/helm/service-operator --set image.tag=v0.2.0
```

## 使用问题

### Q: Service 资源创建后，相关的 Kubernetes 资源没有出现怎么办？

A: 请检查：
1. Operator 是否正常运行：`kubectl get pods -n service-operator-system`
2. 查看 Operator 日志：`kubectl logs -n service-operator-system deployment/service-operator-controller-manager`
3. 检查 Service 资源状态：`kubectl describe services.apps.example.com <service-name>`

### Q: 如何更新已部署的服务？

A: 直接修改 Service 资源即可，Operator 会自动处理更新：
```bash
kubectl patch services.apps.example.com my-service -p '{"spec":{"replicas":5}}'
```

### Q: 支持哪些类型的配置文件？

A: Service Operator 支持任意类型的配置文件，通过 `configData` 字段指定：
```yaml
spec:
  configData:
    application.yml: |
      server:
        port: 8080
    nginx.conf: |
      server {
        listen 80;
      }
```

### Q: 如何设置资源限制？

A: 使用 `resources` 字段：
```yaml
spec:
  resources:
    requests:
      cpu: "100m"
      memory: "128Mi"
    limits:
      cpu: "500m"
      memory: "512Mi"
```

### Q: Ingress 不工作怎么办？

A: 请检查：
1. 集群中是否安装了 Ingress Controller
2. Ingress 配置是否正确
3. DNS 解析是否正确
4. 防火墙和网络策略设置

## 故障排除

### Q: Pod 一直处于 Pending 状态

A: 可能的原因：
- 资源不足：检查节点资源
- 镜像拉取失败：检查镜像名称和权限
- 调度约束：检查节点选择器和污点容忍

### Q: ConfigMap 更新后应用没有重新加载配置

A: 这是 Kubernetes 的正常行为。可以：
1. 重启 Deployment：`kubectl rollout restart deployment <service-name>`
2. 使用支持配置热重载的应用
3. 实现配置变更检测机制

### Q: 服务无法访问

A: 检查步骤：
1. Pod 是否正常运行
2. Service 端点是否正确
3. 网络策略是否阻止访问
4. 防火墙设置

## 高级用法

### Q: 如何实现蓝绿部署？

A: 当前版本不直接支持蓝绿部署，但可以通过以下方式实现：
1. 创建新版本的 Service 资源
2. 测试新版本
3. 更新 Ingress 指向新版本
4. 删除旧版本

### Q: 支持自动扩缩容吗？

A: 当前版本不直接支持 HPA，但创建的 Deployment 可以配置 HPA：
```bash
kubectl autoscale deployment my-service --cpu-percent=50 --min=1 --max=10
```

### Q: 如何监控 Service Operator？

A: Service Operator 提供 Prometheus 指标：
```bash
# 端口转发到指标端点
kubectl port-forward -n service-operator-system service/service-operator-controller-manager-metrics-service 8443:8443

# 访问指标
curl -k https://localhost:8443/metrics
```

### Q: 支持多集群部署吗？

A: 当前版本不支持多集群管理。每个集群需要独立安装 Service Operator。

## 开发问题

### Q: 如何贡献代码？

A: 请参考 [贡献指南](../CONTRIBUTING.md)：
1. Fork 项目
2. 创建功能分支
3. 提交 Pull Request

### Q: 如何添加新的资源类型支持？

A: 需要修改：
1. API 定义 (`api/v1/service_types.go`)
2. 控制器逻辑 (`controllers/service_controller.go`)
3. RBAC 权限配置
4. 测试和文档

### Q: 如何调试 Operator？

A: 
1. 本地运行：`make run`
2. 启用详细日志：`--zap-log-level=debug`
3. 使用 IDE 调试器
4. 查看 Kubernetes 事件

## 性能问题

### Q: Service Operator 的性能如何？

A: Service Operator 设计为轻量级，典型资源使用：
- CPU: 10-100m
- 内存: 64-128Mi
- 可管理数百个 Service 资源

### Q: 如何优化性能？

A: 
1. 调整并发设置：`--max-concurrent-reconciles`
2. 增加资源限制
3. 使用节点选择器将 Operator 调度到高性能节点

## 安全问题

### Q: Service Operator 安全吗？

A: Service Operator 遵循安全最佳实践：
- 最小权限 RBAC
- 非 root 用户运行
- TLS 加密通信
- 定期安全扫描

### Q: 如何报告安全漏洞？

A: 请参考 [安全策略](../SECURITY.md)，通过邮件报告安全问题。

## 许可证和支持

### Q: Service Operator 的许可证是什么？

A: Apache License 2.0，允许商业使用、修改和分发。

### Q: 如何获得支持？

A: 
- **社区支持**: GitHub Issues 和 Discussions
- **文档**: 查看项目文档
- **企业支持**: 联系维护团队

### Q: 有商业版本吗？

A: 当前只有开源版本。如需企业级功能和支持，请联系维护团队。

## 路线图

### Q: 未来会支持哪些功能？

A: 计划中的功能：
- 更多资源类型支持 (StatefulSet, DaemonSet)
- 高级部署策略 (蓝绿部署、金丝雀发布)
- 多集群支持
- Web UI 管理界面

### Q: 如何影响产品路线图？

A: 
1. 提交功能请求 Issue
2. 参与社区讨论
3. 贡献代码实现

## 其他问题

### Q: 可以在生产环境使用吗？

A: Service Operator 设计为生产就绪，但建议：
1. 在测试环境充分验证
2. 制定备份和恢复计划
3. 监控 Operator 运行状态

### Q: 与其他 Operator 兼容吗？

A: Service Operator 遵循 Kubernetes Operator 模式，与其他 Operator 兼容。但要注意资源冲突。

### Q: 如何迁移现有应用？

A: 
1. 分析现有资源配置
2. 创建对应的 Service 资源
3. 逐步迁移，确保服务不中断

---

如果您的问题没有在这里找到答案，请：
1. 查看 [故障排除指南](troubleshooting.md)
2. 搜索 [GitHub Issues](../../issues)
3. 创建新的 Issue 或 Discussion