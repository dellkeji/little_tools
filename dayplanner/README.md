# 日报工具 (Daily Report Tool)

[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](VERSION)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://golang.org/)
[![Platform](https://img.shields.io/badge/platform-Windows-0078D6.svg)](https://www.microsoft.com/windows)

一个基于 Golang 开发的 Windows 桌面应用程序，用于帮助用户记录和管理每日工作任务。

## 功能特性

- 📅 日历视图：以月历形式查看和管理每日任务
- ✍️ Markdown 编辑：支持 Markdown 语法编写任务内容
- 👁️ 实时预览：实时渲染 Markdown 内容
- 🔔 自动提醒：通过企业微信 Webhook 发送每日提醒
- 💾 本地存储：任务数据安全存储在本地文件系统

## 项目结构

```
daily-report-tool/
├── cmd/
│   └── daily-report/
│       └── main.go                 # 应用程序入口
├── internal/
│   ├── ui/                         # 表示层 - UI 组件
│   │   ├── calendar.go            # 日历视图组件
│   │   ├── editor.go              # Markdown 编辑器组件
│   │   ├── preview.go             # Markdown 预览组件
│   │   ├── settings.go            # 设置界面组件
│   │   └── mainwindow.go          # 主窗口
│   ├── service/                    # 业务逻辑层
│   │   ├── task_service.go        # 任务管理服务
│   │   ├── reminder_service.go    # 提醒服务
│   │   └── config_service.go      # 配置管理服务
│   ├── repository/                 # 数据访问层
│   │   ├── task_repository.go     # 任务数据仓库
│   │   └── config_repository.go   # 配置数据仓库
│   ├── model/                      # 数据模型
│   │   ├── task.go                # 任务模型
│   │   └── config.go              # 配置模型
│   └── util/                       # 工具函数
│       ├── markdown.go            # Markdown 处理工具
│       └── webhook.go             # Webhook 调用工具
├── config/
│   └── config.json                # 配置文件
├── data/
│   └── tasks/                     # 任务数据目录
│       └── YYYY-MM-DD.json        # 按日期存储的任务文件
├── go.mod
├── go.sum
└── README.md
```

## 技术栈

- **语言**: Go 1.21+
- **GUI 框架**: [Fyne v2](https://fyne.io/) - 跨平台 Go GUI 工具包
- **Markdown 渲染**: [goldmark](https://github.com/yuin/goldmark) - CommonMark 兼容的 Markdown 解析器
- **数据存储**: JSON 文件
- **HTTP 客户端**: Go 标准库 net/http

## 环境要求

- Go 1.21 或更高版本
- Windows 10 或更高版本
- GCC 编译器（用于 Fyne 的 CGO 依赖）

### Windows 开发环境设置

1. 安装 Go: https://golang.org/dl/
2. 安装 GCC (推荐使用 TDM-GCC): https://jmeubank.github.io/tdm-gcc/

## 快速开始

### 使用发布包（推荐）

1. 从 [Releases](../../releases) 下载最新版本
2. 解压到任意目录
3. 双击 `daily-report.exe` 启动应用
4. 查看 `快速开始.txt` 了解基本使用方法

详细安装说明请查看 [INSTALL.md](INSTALL.md)

## 构建方法

### 开发构建

```bash
# 克隆或下载项目
cd daily-report-tool

# 安装依赖
go mod download

# 运行应用程序
go run ./cmd/daily-report
```

### 生产构建

```bash
# 方式一：使用构建脚本（推荐）
.\build.bat              # Windows
./build.sh               # Linux/Mac

# 方式二：带图标构建（推荐用于发布）
.\build-with-icon.bat    # Windows
./build-with-icon.sh     # Linux/Mac

# 方式三：手动构建
go build -ldflags="-H windowsgui" -o daily-report.exe ./cmd/daily-report

# 方式四：使用 Fyne 打包工具（包含图标和资源）
.\package.bat            # Windows
./package.sh             # Linux/Mac
```

### 设置应用程序图标

如果要为应用程序添加自定义图标：

1. **转换图标格式**（三选一）：
   ```bash
   # 方法 A: 使用 Python 脚本（需要安装 Pillow）
   python convert_icon.py
   
   # 方法 B: 使用在线工具
   # 访问 https://convertio.co/zh/png-ico/
   # 上传 icon.png，下载 icon.ico
   
   # 方法 C: 使用 ImageMagick
   magick convert icon.png -define icon:auto-resize=256,128,64,48,32,16 icon.ico
   ```

2. **构建带图标的应用**：
   ```bash
   .\build-with-icon.bat    # Windows
   ./build-with-icon.sh     # Linux/Mac
   ```

详细说明请查看 [ICON_SETUP.md](ICON_SETUP.md) 和 [CONVERT_ICON.md](CONVERT_ICON.md)

### 创建发布包

```bash
# 创建完整的发布包（包含文档和示例配置）
.\release.bat            # Windows
./release.sh             # Linux/Mac

# 发布包将创建在 release/ 目录下
```

### 构建参数说明

- `-ldflags="-H windowsgui"`: 隐藏控制台窗口，创建纯 GUI 应用
- `-o daily-report.exe`: 指定输出文件名

## 配置说明

首次运行时，应用程序会自动创建配置文件 `config/config.json`：

```json
{
  "webhook_url": "",
  "reminder_time": "10:00",
  "reminder_enabled": false,
  "data_path": "./data/tasks"
}
```

### 配置项说明

- `webhook_url`: 企业微信 Webhook 地址（用于发送提醒）
- `reminder_time`: 每日提醒时间（24小时格式，如 "10:00"）
- `reminder_enabled`: 是否启用自动提醒功能
- `data_path`: 任务数据存储路径

### 获取企业微信 Webhook

1. 登录企业微信管理后台
2. 进入"应用管理" -> "群机器人"
3. 创建新的群机器人
4. 复制 Webhook 地址到配置文件

## 使用说明

1. **启动应用**: 双击 `daily-report.exe` 启动应用程序
2. **选择日期**: 在日历视图中点击日期
3. **编辑任务**: 在编辑器中使用 Markdown 语法编写任务内容
4. **实时预览**: 右侧预览区域会实时显示渲染后的内容
5. **自动保存**: 停止输入 2 秒后自动保存
6. **配置提醒**: 点击菜单栏的"设置"配置企业微信提醒

## 开发指南

### 架构设计

项目采用分层架构：

- **表示层 (UI Layer)**: 使用 Fyne 框架构建的用户界面组件
- **业务逻辑层 (Service Layer)**: 封装业务规则和流程
- **数据访问层 (Repository Layer)**: 处理数据持久化

### 添加新功能

1. 在 `internal/model` 中定义数据模型
2. 在 `internal/repository` 中实现数据访问接口
3. 在 `internal/service` 中实现业务逻辑
4. 在 `internal/ui` 中创建 UI 组件
5. 在 `cmd/daily-report/main.go` 中集成新功能

### 运行测试

```bash
# 运行所有单元测试
go test ./...

# 运行集成测试
go test ./internal/integration/... -v

# 运行端到端验证
go run scripts/e2e_verify.go

# 运行特定包的测试
go test ./internal/service

# 运行测试并显示覆盖率
go test -cover ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 测试覆盖

项目包含三个层次的测试：

1. **单元测试**: 测试各个组件的独立功能
   - Repository 层测试：数据访问和文件操作
   - Service 层测试：业务逻辑和服务交互
   - UI 组件测试：界面组件功能

2. **集成测试**: 测试组件间的协作（`internal/integration/`）
   - 任务创建、编辑和保存流程
   - 日历导航和日期查询
   - 配置管理和验证
   - 提醒服务集成
   - 错误处理和数据持久化

3. **端到端测试**: 验证完整的应用流程（`scripts/e2e_verify.go`）
   - 完整的用户工作流程
   - 数据持久化验证
   - 文件系统结构验证

4. **手动测试**: UI 交互和系统集成测试
   - 查看 [TESTING_GUIDE.md](TESTING_GUIDE.md) 了解详细的手动测试指南
   - 包括 UI 交互、提醒定时触发、Windows 兼容性等测试场景

## 数据存储

任务数据以 JSON 格式存储在 `data/tasks/` 目录下，每个文件对应一天的任务：

```
data/tasks/
├── 2025-11-10.json
├── 2025-11-11.json
└── 2025-11-12.json
```

每个任务文件的格式：

```json
{
  "date": "2025-11-10T00:00:00Z",
  "content": "# 今日工作\n\n- 完成需求文档\n- 代码评审",
  "created_at": "2025-11-10T09:00:00Z",
  "updated_at": "2025-11-10T15:30:00Z"
}
```

## 故障排查

### 应用无法启动

- 检查是否安装了 GCC 编译器
- 确认 Go 版本是否满足要求
- 查看 `logs/app.log` 日志文件

### 提醒功能不工作

- 检查配置文件中的 `webhook_url` 是否正确
- 确认 `reminder_enabled` 设置为 `true`
- 检查网络连接是否正常
- 查看日志文件中的错误信息

### 数据丢失

- 任务数据存储在 `data/tasks/` 目录
- 建议定期备份该目录
- 可以手动复制 JSON 文件进行恢复

## 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本更新历史和未来计划。

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 贡献

欢迎提交 Issue 和 Pull Request！

### 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 联系方式

如有问题或建议，请通过 Issue 反馈。

## 致谢

- [Fyne](https://fyne.io/) - 优秀的 Go GUI 框架
- [goldmark](https://github.com/yuin/goldmark) - 强大的 Markdown 解析器
- 所有贡献者和用户的支持
