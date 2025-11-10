@echo off
echo ========================================
echo 构建带图标的日报工具
echo ========================================
echo.

REM 检查图标文件
if not exist "icon.ico" (
    echo [错误] 找不到 icon.ico 文件
    echo.
    echo 请先将 icon.png 转换为 icon.ico 格式
    echo 可以使用在线工具: https://convertio.co/zh/png-ico/
    echo.
    pause
    exit /b 1
)

echo [1/4] 复制图标文件...
copy /Y icon.ico cmd\daily-report\icon.ico >nul
if errorlevel 1 (
    echo [错误] 复制图标文件失败
    pause
    exit /b 1
)
echo       完成

echo [2/4] 编译资源文件...
cd cmd\daily-report
windres -o icon.syso icon.rc 2>nul
@REM if errorlevel 1 (
@REM     echo [错误] 资源文件编译失败
@REM     echo.
@REM     echo 请确保已安装 windres (MinGW)
@REM     echo 下载地址: https://www.mingw-w64.org/
@REM     cd ..\..
@REM     pause
@REM     exit /b 1
@REM )
cd ..\..
echo       完成

echo [3/4] 构建应用程序...
go build -ldflags="-H windowsgui" -o daily-report.exe ./cmd/daily-report
if errorlevel 1 (
    echo [错误] 应用程序构建失败
    pause
    exit /b 1
)
echo       完成

echo [4/4] 清理临时文件...
del /Q cmd\daily-report\icon.ico 2>nul
echo       完成

echo.
echo ========================================
echo 构建成功！
echo ========================================
echo 可执行文件: daily-report.exe
echo.
echo 提示: 如果图标未显示，请尝试：
echo   1. 重启文件资源管理器
echo   2. 运行: ie4uinit.exe -show
echo.
pause
