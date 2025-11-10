@echo off
REM Build script for Daily Report Tool
REM This script builds the Windows executable with console window hidden

echo Building Daily Report Tool...

REM Build the executable with windowsgui flag to hide console
go build -ldflags="-H windowsgui" -o daily-report.exe ./cmd/daily-report

if %ERRORLEVEL% EQU 0 (
    echo Build successful! Executable: daily-report.exe
) else (
    echo Build failed!
    exit /b 1
)
