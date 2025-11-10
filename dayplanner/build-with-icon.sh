#!/bin/bash
set -e

echo "========================================"
echo "构建带图标的日报工具"
echo "========================================"
echo ""

# 检查图标文件
if [ ! -f "icon.ico" ]; then
    echo "[错误] 找不到 icon.ico 文件"
    echo ""
    echo "请先将 icon.png 转换为 icon.ico 格式"
    echo "可以使用在线工具: https://convertio.co/zh/png-ico/"
    echo ""
    exit 1
fi

echo "[1/4] 复制图标文件..."
cp icon.ico cmd/daily-report/icon.ico
echo "      完成"

echo "[2/4] 编译资源文件..."
cd cmd/daily-report
if ! windres -o icon.syso icon.rc 2>/dev/null; then
    echo "[错误] 资源文件编译失败"
    echo ""
    echo "请确保已安装 windres (MinGW)"
    echo "下载地址: https://www.mingw-w64.org/"
    cd ../..
    exit 1
fi
cd ../..
echo "      完成"

echo "[3/4] 构建应用程序..."
go build -ldflags="-H windowsgui" -o daily-report.exe ./cmd/daily-report
echo "      完成"

echo "[4/4] 清理临时文件..."
rm -f cmd/daily-report/icon.ico
echo "      完成"

echo ""
echo "========================================"
echo "构建成功！"
echo "========================================"
echo "可执行文件: daily-report.exe"
echo ""
echo "提示: 如果图标未显示，请尝试："
echo "  1. 重启文件资源管理器"
echo "  2. 运行: ie4uinit.exe -show"
echo ""
