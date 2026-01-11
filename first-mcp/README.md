# MCP Server

一个使用 MCP Python SDK 构建的服务器，提供版本查询功能并支持便捷扩展。

## 功能特性

- 🔍 内置版本查询工具
- 🔧 易于扩展的架构
- 📝 完整的日志记录
- 🚀 异步处理支持

## 安装依赖

```bash
pip install -r requirements.txt
```

## 运行测试

### 快速测试
```bash
python run_tests.py
```

### 详细测试选项
```bash
# 运行所有测试
python run_tests.py --type all --verbose

# 只运行单元测试
python run_tests.py --type unit

# 只运行集成测试
python run_tests.py --type integration

# 运行测试并生成覆盖率报告
python run_tests.py --coverage
```

### 直接使用 pytest
```bash
# 运行所有测试
pytest

# 运行特定测试文件
pytest test_server.py -v

# 运行集成测试
pytest test_integration.py -m integration -v

# 生成覆盖率报告
pytest --cov=server --cov-report=html
```

## 运行服务器

```bash
python server.py
```

## 可用工具

### get_server_version
查询当前服务器版本

**参数**: 无

**返回**: 服务器版本信息

## 添加新工具

要添加新工具，请按照以下步骤：

1. 在 `handle_list_tools()` 函数中添加工具定义
2. 在 `handle_call_tool()` 函数中添加工具调用处理
3. 实现具体的工具处理函数

### 示例：添加一个新工具

```python
# 1. 在 handle_list_tools() 中添加工具定义
Tool(
    name="my_new_tool",
    description="我的新工具描述",
    inputSchema={
        "type": "object",
        "properties": {
            "input_param": {
                "type": "string",
                "description": "输入参数描述"
            }
        },
        "required": ["input_param"],
    },
),

# 2. 在 handle_call_tool() 中添加处理逻辑
elif name == "my_new_tool":
    return await my_new_tool_handler(arguments or {})

# 3. 实现工具处理函数
async def my_new_tool_handler(arguments: dict[str, Any]) -> list[TextContent]:
    input_param = arguments.get("input_param", "")
    
    # 你的工具逻辑
    result = f"处理结果: {input_param}"
    
    return [
        TextContent(
            type="text",
            text=result
        )
    ]
```

## 项目结构

```
.
├── server.py              # 主服务器文件
├── test_server.py         # 单元测试
├── test_integration.py    # 集成测试
├── run_tests.py          # 测试运行脚本
├── pytest.ini           # pytest 配置
├── requirements.txt      # Python依赖
├── example_config.json   # MCP 配置示例
└── README.md            # 项目文档
```

## 测试覆盖

项目包含完整的测试套件：

### 单元测试 (`test_server.py`)
- ✅ 服务器版本查询功能
- ✅ 工具列表功能
- ✅ 工具调用处理
- ✅ 错误处理
- ✅ 示例工具处理器

### 集成测试 (`test_integration.py`)
- ✅ 服务器初始化
- ✅ MCP 协议合规性
- ✅ 并发处理能力
- ✅ 完整工作流程

### 测试特性
- 异步测试支持
- 错误场景覆盖
- 性能测试
- 协议合规性验证

## 配置 MCP 客户端

要在 Kiro 中使用此服务器，请在 `.kiro/settings/mcp.json` 中添加配置：

```json
{
  "mcpServers": {
    "my-mcp-server": {
      "command": "python",
      "args": ["path/to/server.py"],
      "disabled": false,
      "autoApprove": ["get_server_version"]
    }
  }
}
```