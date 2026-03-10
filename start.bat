@echo off
chcp 65001 >nul
echo ========================================
echo   TTS 项目快速启动
echo ========================================
echo.

echo [1/4] 启动 MySQL 和 Redis...
docker start mysql redis >nul 2>&1
echo       启动 DBHub (Database MCP, 端口 8080)...
docker compose up -d dbhub >nul 2>&1
echo       DBHub 已启动
echo       MySQL 和 Redis 已启动

echo [2/4] 启动 user-api (端口 8081)...
start "user-api" cmd /k "set GOTELEMETRY=off&& set GOCACHE=%~dp0gocache&& set GOMODCACHE=%~dp0gocache\\gomod&& cd /d %~dp0tts-backend\\user-api && go run . -f etc\\user-api.yaml"

echo [3/4] 启动 voice-api (端口 8082)...
start "voice-api" cmd /k "set GOTELEMETRY=off&& set GOCACHE=%~dp0gocache&& set GOMODCACHE=%~dp0gocache\\gomod&& cd /d %~dp0tts-backend\\voice-api && go run . -f etc\\voice-api.yaml"

echo [4/4] 启动 tts-api (端口 8083)...
start "tts-api" cmd /k "set GOTELEMETRY=off&& set GOCACHE=%~dp0gocache&& set GOMODCACHE=%~dp0gocache\\gomod&& cd /d %~dp0tts-backend\\tts-api && go run . -f etc\\tts-api.yaml"

start "tts-worker" cmd /k "set GOTELEMETRY=off&& set GOCACHE=%~dp0gocache&& set GOMODCACHE=%~dp0gocache\\gomod&& cd /d %~dp0tts-backend\\tts-worker && go run . -f etc\\tts-worker.yaml"

echo.
echo 等待后端服务启动...
timeout /t 5 /nobreak >nul

echo [5/5] 启动前端 (端口 3000)...
start "frontend" cmd /k "cd /d %~dp0tts-front && npm run dev -- --host 0.0.0.0"

echo.
echo ========================================
echo   启动完成!
echo ========================================
echo   前端: http://localhost:3000
echo   Frontend(LAN):   http://^<本机IP^>:3000
echo   user-api:  http://localhost:8081
echo   voice-api: http://localhost:8082
echo   tts-api:   http://localhost:8083
echo ========================================
echo.
pause
