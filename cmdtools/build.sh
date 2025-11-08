#!/bin/bash

# 在 macOS 上交叉编译为 Windows 可执行文件

echo "正在编译 Windows 可执行文件..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o opencmd.exe

if [ $? -eq 0 ]; then
    echo "编译成功！生成文件: opencmd.exe"
    ls -lh opencmd.exe
else
    echo "编译失败"
    exit 1
fi
