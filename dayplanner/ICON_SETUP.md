# 应用程序图标设置指南

## 概述

本文档说明如何为日报工具设置 Windows 应用程序图标。

## 步骤

### 1. 准备图标文件

项目中已有 `icon.png` 文件。需要将其转换为 `.ico` 格式。

#### 方法 A: 使用在线工具（推荐）

1. 访问在线转换工具（例如）：
   - https://convertio.co/zh/png-ico/
   - https://www.icoconverter.com/
   - https://cloudconvert.com/png-to-ico

2. 上传 `icon.png` 文件

3. 选择输出尺寸（推荐包含多个尺寸）：
   - 16x16 像素
   - 32x32 像素
   - 48x48 像素
   - 256x256 像素

4. 下载生成的 `icon.ico` 文件

5. 将 `icon.ico` 放到项目根目录

#### 方法 B: 使用 ImageMagick（命令行）

如果已安装 ImageMagick：

```bash
# Windows
magick convert icon.png -define icon:auto-resize=256,128,64,48,32,16 icon.ico

# Linux/Mac
convert icon.png -define icon:auto-resize=256,128,64,48,32,16 icon.ico
```

#### 方法 C: 使用 GIMP（图形界面）

1. 用 GIMP 打开 `icon.png`
2. 调整图像大小为 256x256（图像 -> 缩放图像）
3. 导出为 ICO 格式（文件 -> 导出为 -> 选择 .ico 格式）
4. 在导出对话框中选择多个尺寸

### 2. 复制图标到构建目录

将生成的 `icon.ico` 复制到 `cmd/daily-report/` 目录：

```bash
# Windows
copy icon.ico cmd\daily-report\icon.ico

# Linux/Mac
cp icon.ico cmd/daily-report/icon.ico
```

### 3. 编译资源文件

项目已创建资源文件 `cmd/daily-report/icon.rc`。

#### 安装 windres（如果尚未安装）

windres 通常随 MinGW 或 TDM-GCC 一起安装。

检查是否已安装：
```bash
windres --version
```

如果未安装，可以：
- 安装 MinGW-w64: https://www.mingw-w64.org/
- 或安装 TDM-GCC: https://jmeubank.github.io/tdm-gcc/

#### 编译资源文件

```bash
# 进入目录
cd cmd/daily-report

# 编译资源文件
windres -o icon.syso icon.rc
```

这会生成 `icon.syso` 文件，Go 编译器会自动将其链接到可执行文件中。

### 4. 构建应用程序

使用现有的构建脚本：

```bash
# Windows
.\build.bat

# Linux/Mac
./build.sh
```

或手动构建：

```bash
go build -ldflags="-H windowsgui" -o daily-report.exe ./cmd/daily-report
```

### 5. 验证图标

构建完成后，检查 `daily-report.exe`：
- 在文件资源管理器中查看图标
- 右键点击 -> 属性 -> 查看图标
- 运行程序，检查任务栏图标

## 自动化脚本

为了简化流程，可以创建一个自动化脚本：

### build-with-icon.bat (Windows)

```batch
@echo off
echo 正在构建带图标的应用程序...

REM 检查图标文件
if not exist "icon.ico" (
    echo 错误: 找不到 icon.ico 文件
    echo 请先将 icon.png 转换为 icon.ico
    exit /b 1
)

REM 复制图标到构建目录
copy /Y icon.ico cmd\daily-report\icon.ico

REM 编译资源文件
cd cmd\daily-report
windres -o icon.syso icon.rc
if errorlevel 1 (
    echo 错误: 资源文件编译失败
    echo 请确保已安装 windres (MinGW)
    cd ..\..
    exit /b 1
)
cd ..\..

REM 构建应用程序
go build -ldflags="-H windowsgui" -o daily-report.exe ./cmd/daily-report
if errorlevel 1 (
    echo 错误: 应用程序构建失败
    exit /b 1
)

echo 构建成功！
echo 可执行文件: daily-report.exe
```

### build-with-icon.sh (Linux/Mac)

```bash
#!/bin/bash
set -e

echo "正在构建带图标的应用程序..."

# 检查图标文件
if [ ! -f "icon.ico" ]; then
    echo "错误: 找不到 icon.ico 文件"
    echo "请先将 icon.png 转换为 icon.ico"
    exit 1
fi

# 复制图标到构建目录
cp icon.ico cmd/daily-report/icon.ico

# 编译资源文件
cd cmd/daily-report
windres -o icon.syso icon.rc || {
    echo "错误: 资源文件编译失败"
    echo "请确保已安装 windres (MinGW)"
    cd ../..
    exit 1
}
cd ../..

# 构建应用程序
go build -ldflags="-H windowsgui" -o daily-report.exe ./cmd/daily-report

echo "构建成功！"
echo "可执行文件: daily-report.exe"
```

## 故障排查

### 问题 1: windres 命令未找到

**解决方案**: 安装 MinGW-w64 或 TDM-GCC

### 问题 2: 图标未显示

**可能原因**:
1. `icon.syso` 文件未生成或未在正确位置
2. 图标文件格式不正确
3. Windows 图标缓存问题

**解决方案**:
1. 确认 `cmd/daily-report/icon.syso` 存在
2. 重新生成 ICO 文件，确保包含多个尺寸
3. 清除 Windows 图标缓存：
   ```batch
   ie4uinit.exe -show
   ```

### 问题 3: 构建后图标消失

**原因**: 每次构建前需要确保 `icon.syso` 存在

**解决方案**: 使用自动化构建脚本，或将 `icon.syso` 添加到版本控制

## 注意事项

1. **icon.syso 文件**: 这是编译后的资源文件，Go 编译器会自动识别并链接
2. **位置很重要**: `icon.syso` 必须与 `main.go` 在同一目录
3. **跨平台**: 资源文件仅在 Windows 上有效，不影响其他平台的构建
4. **版本控制**: 可以将 `icon.syso` 添加到 Git，避免每次都重新编译

## 更新现有构建脚本

建议更新 `build.bat` 和 `release.bat` 以自动处理图标：

在构建命令前添加：
```batch
REM 确保资源文件存在
if exist "cmd\daily-report\icon.syso" (
    echo 使用现有的图标资源文件
) else (
    echo 警告: 未找到图标资源文件，应用程序将使用默认图标
)
```

## 参考资源

- [Go Windows 图标设置](https://github.com/akavel/rsrc)
- [windres 文档](https://sourceware.org/binutils/docs/binutils/windres.html)
- [Windows 图标格式规范](https://docs.microsoft.com/en-us/windows/win32/menurc/icon-resource)
