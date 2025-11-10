# 应用程序图标快速设置

## 🎯 目标

为日报工具添加自定义图标，让应用程序在 Windows 文件资源管理器和任务栏中显示你的图标。

## 📋 前提条件

- ✅ 项目中已有 `icon.png` 文件
- ✅ 已安装 Go 开发环境
- ⚠️ 需要安装 MinGW（用于编译资源文件）

## 🚀 快速开始（3 步完成）

### 步骤 1: 转换图标格式

选择以下任一方法：

#### 方法 A: 使用 Python 脚本（推荐）

```bash
# 安装 Pillow（如果尚未安装）
pip install Pillow

# 运行转换脚本
python convert_icon.py
```

#### 方法 B: 使用在线工具

1. 访问 https://convertio.co/zh/png-ico/
2. 上传 `icon.png`
3. 下载 `icon.ico` 到项目根目录

### 步骤 2: 安装 MinGW（如果尚未安装）

检查是否已安装：
```bash
windres --version
```

如果未安装，下载并安装：
- **推荐**: TDM-GCC https://jmeubank.github.io/tdm-gcc/
- 或 MinGW-w64 https://www.mingw-w64.org/

### 步骤 3: 构建应用程序

```bash
# Windows
.\build-with-icon.bat

# Linux/Mac
./build-with-icon.sh
```

## ✅ 验证

构建完成后：
1. 在文件资源管理器中查看 `daily-report.exe`
2. 应该能看到你的自定义图标
3. 运行程序，任务栏也会显示该图标

## 🔧 故障排查

### 问题 1: "windres 命令未找到"

**解决方案**: 安装 MinGW（见步骤 2）

### 问题 2: 图标未显示

**可能原因**:
- Windows 图标缓存问题

**解决方案**:
```bash
# 清除图标缓存
ie4uinit.exe -show

# 或重启文件资源管理器
taskkill /f /im explorer.exe
start explorer.exe
```

### 问题 3: "找不到 icon.ico"

**解决方案**: 确保已完成步骤 1，在项目根目录生成了 `icon.ico` 文件

## 📁 生成的文件

构建过程会生成以下文件：

```
项目根目录/
├── icon.ico                        # 转换后的图标文件
└── cmd/daily-report/
    ├── icon.rc                     # 资源定义文件（已创建）
    └── icon.syso                   # 编译后的资源文件（自动生成）
```

## 🎨 图标设计建议

为了获得最佳显示效果：

- **原始 PNG 尺寸**: 至少 256x256 像素，推荐 512x512
- **ICO 包含尺寸**: 16x16, 32x32, 48x48, 256x256
- **背景**: 建议使用透明背景
- **风格**: 简洁明了，在小尺寸下也能识别

## 📚 更多信息

- 详细设置指南: [ICON_SETUP.md](ICON_SETUP.md)
- 图标转换说明: [CONVERT_ICON.md](CONVERT_ICON.md)
- 构建说明: [README.md](README.md)

## 💡 提示

1. **一次设置，永久使用**: `icon.syso` 文件可以提交到版本控制，之后就不需要每次都重新编译资源文件

2. **自动化构建**: 使用 `build-with-icon.bat/sh` 脚本可以自动处理所有步骤

3. **跨平台**: 资源文件仅在 Windows 上有效，不影响 Linux/Mac 构建

---

**需要帮助？** 查看 [ICON_SETUP.md](ICON_SETUP.md) 获取详细说明
