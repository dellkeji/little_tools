#!/bin/bash
# Release script for Daily Report Tool
# This script builds the application and creates a release package

echo "========================================"
echo "Daily Report Tool - Release Builder"
echo "========================================"
echo

# Get version from user or use default
read -p "Enter version (default: 1.0.0): " VERSION
VERSION=${VERSION:-1.0.0}

RELEASE_DIR="release/daily-report-v${VERSION}-windows"
RELEASE_ZIP="daily-report-v${VERSION}-windows.zip"

echo
echo "Building version ${VERSION}..."
echo

# Clean previous release
rm -rf release
mkdir -p release
mkdir -p "${RELEASE_DIR}"

# Build the executable
echo "[1/5] Building executable..."
go build -ldflags="-H windowsgui" -o "${RELEASE_DIR}/daily-report.exe" ./cmd/daily-report

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi
echo "Build successful!"

# Create directory structure
echo "[2/5] Creating directory structure..."
mkdir -p "${RELEASE_DIR}/config"
mkdir -p "${RELEASE_DIR}/data/tasks"

# Copy configuration example
echo "[3/5] Copying configuration files..."
cp config/config.json.example "${RELEASE_DIR}/config/config.json.example"

# Copy documentation
echo "[4/5] Copying documentation..."
cp README.md "${RELEASE_DIR}/README.md"

# Copy icon if exists
if [ -f icon.png ]; then
    cp icon.png "${RELEASE_DIR}/icon.png"
fi

# Create a quick start guide
echo "[5/5] Creating quick start guide..."
cat > "${RELEASE_DIR}/快速开始.txt" << 'EOF'
# 快速开始指南

## 首次使用

1. 双击 `daily-report.exe` 启动应用程序
2. 应用会自动创建配置文件和数据目录
3. 点击日历选择日期，开始编写任务

## 配置企业微信提醒（可选）

1. 在企业微信中创建群机器人，获取 Webhook 地址
2. 点击应用菜单栏的"设置"
3. 填入 Webhook 地址
4. 设置提醒时间（默认 10:00）
5. 启用提醒开关
6. 点击"保存"

## 使用技巧

- 支持 Markdown 语法：标题、列表、粗体、斜体、代码块等
- 编辑器会在停止输入 2 秒后自动保存
- 日历上有标记的日期表示已有任务记录
- 任务数据存储在 `data/tasks/` 目录，可以手动备份

## 故障排查

- 如果应用无法启动，请检查是否有杀毒软件拦截
- 提醒功能需要正确配置企业微信 Webhook 地址
- 查看 `logs/app.log` 了解详细错误信息

## 数据备份

建议定期备份以下目录：
- `data/tasks/` - 任务数据
- `config/config.json` - 配置文件

更多信息请查看 README.md
EOF

echo
echo "========================================"
echo "Release package created successfully!"
echo "========================================"
echo
echo "Location: ${RELEASE_DIR}"
echo
echo "Contents:"
ls -1 "${RELEASE_DIR}"
echo

# Create ZIP archive
echo "Creating ZIP archive..."
cd release
zip -r "${RELEASE_ZIP}" "daily-report-v${VERSION}-windows" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "ZIP archive created: release/${RELEASE_ZIP}"
else
    echo "Note: Could not create ZIP archive. Please zip manually."
fi
cd ..

echo
echo "Release is ready for distribution!"
