@echo off
chcp 65001 >nul
echo 正在编译 xiaohongshu-mcp 项目...

echo.
echo [1/2] 编译 MCP 服务器...
go build -o xiaohongshu-mcp-windows-amd64.exe .
if %errorlevel% neq 0 (
    echo 编译 MCP 服务器失败！
    pause
    exit /b 1
)
echo ✓ MCP 服务器编译成功

echo.
echo [2/2] 编译登录工具...
go build -o xiaohongshu-login-windows-amd64.exe ./cmd/login
if %errorlevel% neq 0 (
    echo 编译登录工具失败！
    pause
    exit /b 1
)
echo ✓ 登录工具编译成功

echo.
echo ========================================
echo 编译完成！
echo ========================================
echo.
echo 生成的文件：
echo   - xiaohongshu-mcp-windows-amd64.exe (MCP服务器)
echo   - xiaohongshu-login-windows-amd64.exe (登录工具)
echo.
echo 使用步骤：
echo   1. 先运行登录工具: xiaohongshu-login-windows-amd64.exe
echo   2. 然后启动MCP服务器: xiaohongshu-mcp-windows-amd64.exe
echo.
pause
