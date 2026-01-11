# Cross-Platform Agent

一个用Rust编写的跨平台Agent，支持Windows、Linux和macOS，提供远程控制和监控功能。

## 功能特性

### 控制平面 (Control Plane)
- 远程命令执行
- 文件部署
- 配置管理
- 进程控制
- 心跳监控

### 数据平面 (Data Plane)
- 系统指标收集 (CPU、内存、磁盘)
- 日志收集
- 自定义指标收集
- 多种导出器支持 (HTTP、文件)
- 类似OpenTelemetry Collector的架构

## 快速开始

### 1. 编译

```bash
cargo build --release
```

### 2. 生成默认配置

```bash
./target/release/agent config -o config.toml
```

### 3. 启动Agent

```bash
./target/release/agent start -c config.toml
```

## 配置说明

### 基本配置

```toml
[agent]
id = "550e8400-e29b-41d4-a716-446655440000"  # 可选，自动生成
name = "default-agent"
tags = ["production"]
heartbeat_interval = 30  # 秒
command_timeout = 300    # 秒

[control_plane]
enabled = true
server_url = "http://localhost:8080"
api_key = "your-api-key"  # 可选
poll_interval = 10        # 秒
max_concurrent_commands = 5

[data_plane]
enabled = true
buffer_size = 1000
flush_interval = 60  # 秒
```

### 收集器配置

```toml
[[data_plane.collectors]]
name = "system_metrics"
collector_type = "system"
enabled = true

[data_plane.collectors.config]
interval = 30
metrics = ["cpu", "memory", "disk"]

[[data_plane.collectors]]
name = "custom_script"
collector_type = "custom"
enabled = true

[data_plane.collectors.config]
command = "python /path/to/script.py"
interval = 60
```

### 导出器配置

```toml
[[data_plane.exporters]]
name = "http_exporter"
exporter_type = "http"
endpoint = "http://localhost:8081/metrics"
batch_size = 100
enabled = true

[data_plane.exporters.headers]
"Content-Type" = "application/json"
"Authorization" = "Bearer your-token"

[[data_plane.exporters]]
name = "file_exporter"
exporter_type = "file"
endpoint = "/var/log/agent-metrics.json"
batch_size = 50
enabled = true
```

## API接口

### 控制平面API

#### 注册Agent
```
POST /api/agents/register
Content-Type: application/json

{
  "id": "uuid",
  "hostname": "hostname",
  "platform": "linux",
  "arch": "x86_64",
  "version": "0.1.0",
  "last_heartbeat": "2023-12-26T10:00:00Z"
}
```

#### 心跳
```
POST /api/agents/{agent_id}/heartbeat
```

#### 获取命令
```
GET /api/agents/{agent_id}/commands
```

#### 提交命令结果
```
POST /api/commands/{command_id}/result
Content-Type: application/json

{
  "command_id": "uuid",
  "success": true,
  "output": "command output",
  "error": null,
  "execution_time": 1500
}
```

### 命令类型

#### 执行命令
```json
{
  "id": "uuid",
  "command_type": "Execute",
  "payload": {
    "command": "ls",
    "args": ["-la", "/tmp"]
  },
  "timeout": 30
}
```

#### 部署文件
```json
{
  "id": "uuid",
  "command_type": "Deploy",
  "payload": {
    "source": "/path/to/source",
    "destination": "/path/to/destination"
  }
}
```

#### 配置管理
```json
{
  "id": "uuid",
  "command_type": "Configure",
  "payload": {
    "path": "/etc/app/config.json",
    "config": {
      "key": "value"
    }
  }
}
```

#### 停止进程
```json
{
  "id": "uuid",
  "command_type": "Stop",
  "payload": {
    "process": "nginx"
  }
}
```

## 数据格式

### 指标数据
```json
{
  "name": "system_cpu_usage",
  "value": 45.2,
  "labels": {
    "collector": "system_metrics",
    "host": "server-01"
  },
  "timestamp": "2023-12-26T10:00:00Z"
}
```

### 日志数据
```json
{
  "level": "INFO",
  "message": "Application started",
  "source": "app.log",
  "timestamp": "2023-12-26T10:00:00Z",
  "labels": {
    "service": "web-server"
  }
}
```

## 跨平台支持

### Windows
- 使用WMI和PowerShell命令获取系统信息
- 支持Windows服务管理
- 使用taskkill进行进程管理

### Linux
- 使用/proc文件系统获取系统信息
- 支持systemd服务管理
- 使用pkill进行进程管理

### macOS
- 使用system_profiler获取系统信息
- 支持launchd服务管理
- 使用pkill进行进程管理

## 安全考虑

1. **API密钥认证**: 支持Bearer token认证
2. **命令白名单**: 可配置允许执行的命令
3. **文件路径限制**: 限制文件操作的路径范围
4. **超时控制**: 防止长时间运行的命令
5. **日志审计**: 记录所有执行的命令和操作

## 扩展开发

### 自定义收集器

```rust
use crate::data_plane::Collector;
use async_trait::async_trait;

pub struct CustomCollector {
    // 自定义字段
}

#[async_trait]
impl Collector for CustomCollector {
    async fn collect(&self) -> Result<Vec<MetricData>> {
        // 实现自定义收集逻辑
        Ok(vec![])
    }
}
```

### 自定义导出器

```rust
use crate::data_plane::Exporter;
use async_trait::async_trait;

pub struct CustomExporter {
    // 自定义字段
}

#[async_trait]
impl Exporter for CustomExporter {
    async fn export(&self, data: &[MetricData]) -> Result<()> {
        // 实现自定义导出逻辑
        Ok(())
    }
}
```

## 部署建议

### 系统服务部署

#### Linux (systemd)
```ini
[Unit]
Description=Cross Platform Agent
After=network.target

[Service]
Type=simple
User=agent
ExecStart=/usr/local/bin/agent start -c /etc/agent/config.toml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

#### Windows (服务)
使用NSSM或类似工具将Agent注册为Windows服务。

#### macOS (launchd)
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.company.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/agent</string>
        <string>start</string>
        <string>-c</string>
        <string>/etc/agent/config.toml</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

## 许可证

MIT License