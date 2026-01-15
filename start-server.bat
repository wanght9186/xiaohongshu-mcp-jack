@echo off
chcp 65001 >nul
echo 正在启动 xiaohongshu-mcp 服务器...
echo.
echo 服务器将在 http://localhost:18060/mcp 运行
echo 按 Ctrl+C 可以停止服务器
echo.
xiaohongshu-mcp-windows-amd64.exe
