# Contributing to Service Operator

感谢您对 Service Operator 项目的关注！我们欢迎各种形式的贡献，包括但不限于：

- 报告 bug
- 提出功能请求
- 提交代码改进
- 改进文档
- 分享使用经验

## 开始之前

在开始贡献之前，请：

1. 阅读项目的 [README](README.md) 了解项目概况
2. 查看 [开发指南](docs/development.md) 了解开发环境设置
3. 浏览现有的 [Issues](../../issues) 和 [Pull Requests](../../pulls)
4. 加入我们的社区讨论

## 报告问题

### Bug 报告

如果您发现了 bug，请创建一个 Issue 并包含以下信息：

**Bug 报告模板：**

```markdown
## Bug 描述
简要描述遇到的问题

## 复现步骤
1. 执行 '...'
2. 点击 '....'
3. 滚动到 '....'
4. 看到错误

## 期望行为
描述您期望发生的情况

## 实际行为
描述实际发生的情况

## 环境信息
- Kubernetes 版本: [例如 v1.28.0]
- Service Operator 版本: [例如 v0.1.0]
- 操作系统: [例如 Ubuntu 20.04]
- 其他相关信息

## 附加信息
- 错误日志
- 配置文件
- 截图（如果适用）
```

### 功能请求

如果您有新功能的想法，请创建一个 Issue 并包含：

**功能请求模板：**

```markdown
## 功能描述
简要描述您希望添加的功能

## 使用场景
描述这个功能解决的问题或改进的场景

## 详细设计
如果有具体的设计想法，请详细描述

## 替代方案
是否考虑过其他解决方案？

## 附加信息
任何其他相关信息、链接或截图
```

## 代码贡献

### 开发流程

1. **Fork 项目**
   ```bash
   # 在 GitHub 上 fork 项目
   git clone https://github.com/your-username/service-operator.git
   cd service-operator
   ```

2. **设置开发环境**
   ```bash
   # 运行开发环境设置脚本
   ./scripts/dev-setup.sh
   ```

3. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/your-bug-fix
   ```

4. **进行开发**
   - 编写代码
   - 添加测试
   - 更新文档

5. **测试您的更改**
   ```bash
   # 运行测试
   make test
   
   # 运行代码检查
   make fmt vet
   
   # 本地测试
   make run
   ```

6. **提交更改**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

7. **推送分支**
   ```bash
   git push origin feature/your-feature-name
   ```

8. **创建 Pull Request**
   - 在 GitHub 上创建 Pull Request
   - 填写 PR 模板
   - 等待代码审查

### 代码规范

#### Go 代码风格

- 遵循标准的 Go 代码风格
- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码质量
- 为公共函数和类型添加注释

```go
// ServiceController reconciles a Service object
type ServiceController struct {
    client.Client
    Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop
func (r *ServiceController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // Implementation here
}
```

#### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型 (type):**
- `feat`: 新功能
- `fix`: bug 修复
- `docs`: 文档更新
- `style`: 代码格式化（不影响功能）
- `refactor`: 代码重构
- `test`: 添加或修改测试
- `chore`: 构建过程或辅助工具的变动

**示例:**
```
feat(controller): add support for custom annotations

Add ability to specify custom annotations for managed resources
through the Service spec.

Closes #123
```

#### 测试要求

- 为新功能添加单元测试
- 确保所有测试通过
- 测试覆盖率不应降低
- 添加集成测试（如果适用）

```go
func TestServiceController_Reconcile(t *testing.T) {
    // Test implementation
}
```

#### 文档要求

- 更新相关的 API 文档
- 添加或更新示例
- 更新 README（如果需要）
- 添加 CHANGELOG 条目

### Pull Request 流程

#### PR 模板

```markdown
## 变更描述
简要描述这个 PR 的变更内容

## 变更类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 文档更新
- [ ] 代码重构
- [ ] 性能改进
- [ ] 其他: ___________

## 测试
- [ ] 添加了新的测试
- [ ] 所有测试通过
- [ ] 手动测试通过

## 检查清单
- [ ] 代码遵循项目的代码规范
- [ ] 自我审查了代码变更
- [ ] 代码有适当的注释
- [ ] 更新了相关文档
- [ ] 变更不会产生新的警告
- [ ] 添加了测试证明修复有效或功能正常
- [ ] 新的和现有的单元测试都通过

## 相关 Issue
Closes #(issue number)

## 截图（如果适用）
添加截图来帮助解释您的变更

## 附加信息
任何其他相关信息
```

#### 代码审查

所有的 PR 都需要经过代码审查：

1. **自动检查**
   - CI/CD 流水线必须通过
   - 代码格式检查
   - 测试覆盖率检查

2. **人工审查**
   - 至少一个维护者的批准
   - 代码质量和设计审查
   - 文档完整性检查

3. **合并要求**
   - 所有检查通过
   - 冲突已解决
   - 获得必要的批准

## 文档贡献

### 文档类型

- **API 文档**: 描述 API 接口和使用方法
- **用户指南**: 面向最终用户的使用说明
- **开发文档**: 面向开发者的技术文档
- **示例**: 实际使用场景的示例代码

### 文档规范

- 使用清晰、简洁的语言
- 提供实际可运行的示例
- 包含必要的截图或图表
- 保持文档的时效性

### 文档结构

```
docs/
├── api-reference.md      # API 参考文档
├── development.md        # 开发指南
├── deployment.md         # 部署指南
├── architecture.md       # 架构设计
└── troubleshooting.md    # 故障排除
```

## 社区参与

### 行为准则

我们致力于为每个人提供友好、安全和欢迎的环境。请遵循以下原则：

- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表示同理心

### 沟通渠道

- **GitHub Issues**: 报告 bug 和功能请求
- **GitHub Discussions**: 一般讨论和问答
- **Pull Requests**: 代码审查和讨论

### 社区角色

#### 贡献者 (Contributors)
- 提交过至少一个被合并的 PR
- 参与 Issue 讨论和代码审查

#### 维护者 (Maintainers)
- 有项目写权限
- 负责代码审查和发布管理
- 指导项目方向

#### 核心团队 (Core Team)
- 项目的主要维护者
- 负责重大决策和项目治理

## 发布流程

### 版本管理

我们使用 [Semantic Versioning](https://semver.org/):

- **MAJOR**: 不兼容的 API 变更
- **MINOR**: 向后兼容的功能添加
- **PATCH**: 向后兼容的 bug 修复

### 发布步骤

1. 更新 CHANGELOG.md
2. 创建发布分支
3. 更新版本号
4. 创建 Git 标签
5. 构建和推送镜像
6. 创建 GitHub Release
7. 更新文档

## 获取帮助

如果您在贡献过程中遇到问题：

1. 查看现有的文档和 FAQ
2. 搜索相关的 Issues
3. 在 GitHub Discussions 中提问
4. 联系维护者

## 致谢

感谢所有为 Service Operator 项目做出贡献的人！

您的贡献将被记录在项目的贡献者列表中。

---

再次感谢您的贡献！🎉