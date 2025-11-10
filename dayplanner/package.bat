@echo off
REM Package script for Daily Report Tool using Fyne
REM This script creates a packaged Windows application

echo Packaging Daily Report Tool with Fyne...

REM Check if fyne command is available
where fyne >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Fyne CLI not found. Installing...
    go install fyne.io/fyne/v2/cmd/fyne@latest
)

REM Package the application
if exist icon.png (
    fyne package -os windows -icon icon.png -name "Daily Report Tool" -appID com.dailyreport.tool
) else (
    echo Warning: icon.png not found, packaging without icon
    fyne package -os windows -name "Daily Report Tool" -appID com.dailyreport.tool
)

if %ERRORLEVEL% EQU 0 (
    echo Package successful!
) else (
    echo Package failed!
    exit /b 1
)
