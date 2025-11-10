@echo off
REM Release script for Daily Report Tool
REM This script builds the application and creates a release package

echo ========================================
echo Daily Report Tool - Release Builder
echo ========================================
echo.

REM Get version from user or use default
set /p VERSION="Enter version (default: 1.0.0): "
if "%VERSION%"=="" set VERSION=1.0.0

set RELEASE_DIR=release\daily-report-v%VERSION%-windows
set RELEASE_ZIP=daily-report-v%VERSION%-windows.zip

echo.
echo Building version %VERSION%...
echo.

REM Clean previous release
if exist release rmdir /s /q release
mkdir release
mkdir "%RELEASE_DIR%"

REM Build the executable
echo [1/5] Building executable...
go build -ldflags="-H windowsgui" -o "%RELEASE_DIR%\daily-report.exe" ./cmd/daily-report

if %ERRORLEVEL% NEQ 0 (
    echo Build failed!
    exit /b 1
)
echo Build successful!

REM Create directory structure
echo [2/5] Creating directory structure...
mkdir "%RELEASE_DIR%\config"
mkdir "%RELEASE_DIR%\data"
mkdir "%RELEASE_DIR%\data\tasks"

REM Copy configuration example
echo [3/5] Copying configuration files...
copy config\config.json.example "%RELEASE_DIR%\config\config.json.example" >nul

REM Copy documentation
echo [4/5] Copying documentation...
copy README.md "%RELEASE_DIR%\README.md" >nul

REM Copy icon if exists
if exist icon.png (
    copy icon.png "%RELEASE_DIR%\icon.png" >nul
)

REM Create a quick start guide
echo [5/5] Creating quick start guide...
(
echo # 快速开始指南
echo.
echo ## 首次使用
echo.
echo 1. 双击 `daily-report.exe` 启动应用程序
echo 2. 应用会自动创建配置文件和数据目录
echo 3. 点击日历选择日期，开始编写任务
echo.
echo ## 配置企业微信提醒（可选）
echo.
echo 1. 在企业微信中创建群机器人，获取 Webhook 地址
echo 2. 点击应用菜单栏的"设置"
echo 3. 填入 Webhook 地址
echo 4. 设置提醒时间（默认 10:00）
echo 5. 启用提醒开关
echo 6. 点击"保存"
echo.
echo ## 使用技巧
echo.
echo - 支持 Markdown 语法：标题、列表、粗体、斜体、代码块等
echo - 编辑器会在停止输入 2 秒后自动保存
echo - 日历上有标记的日期表示已有任务记录
echo - 任务数据存储在 `data/tasks/` 目录，可以手动备份
echo.
echo ## 故障排查
echo.
echo - 如果应用无法启动，请检查是否有杀毒软件拦截
echo - 提醒功能需要正确配置企业微信 Webhook 地址
echo - 查看 `logs/app.log` 了解详细错误信息
echo.
echo ## 数据备份
echo.
echo 建议定期备份以下目录：
echo - `data/tasks/` - 任务数据
echo - `config/config.json` - 配置文件
echo.
echo 更多信息请查看 README.md
) > "%RELEASE_DIR%\快速开始.txt"

echo.
echo ========================================
echo Release package created successfully!
echo ========================================
echo.
echo Location: %RELEASE_DIR%
echo.
echo Contents:
dir /b "%RELEASE_DIR%"
echo.

REM Create ZIP archive if 7z or PowerShell is available
echo Creating ZIP archive...
where 7z >nul 2>nul
if %ERRORLEVEL% EQU 0 (
    7z a -tzip "release\%RELEASE_ZIP%" ".\%RELEASE_DIR%\*" >nul
    echo ZIP archive created: release\%RELEASE_ZIP%
) else (
    powershell -command "Compress-Archive -Path '%RELEASE_DIR%\*' -DestinationPath 'release\%RELEASE_ZIP%' -Force" 2>nul
    if %ERRORLEVEL% EQU 0 (
        echo ZIP archive created: release\%RELEASE_ZIP%
    ) else (
        echo Note: Could not create ZIP archive. Please zip manually.
    )
)

echo.
echo Release is ready for distribution!
pause
