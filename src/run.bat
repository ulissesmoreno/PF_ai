@echo off
setlocal EnableDelayedExpansion
cd /d "%~dp0"

title PF Agent

set CGO_ENABLED=0
set GOCACHE=%cd%\.gocache
set BINARY=.bin\pf-agent_%RANDOM%.exe

echo.
echo  ================================================
echo   PF Agent
echo   Hexagonal Architecture Active ^| GSD Loaded
echo  ================================================
echo.

echo [1/6] Verificando Go...
where go >nul 2>&1
if errorlevel 1 (
    echo.
    echo  ERRO: Go nao encontrado no PATH.
    echo  Instale em: https://go.dev/dl/
    echo.
    pause
    exit /b 1
)
for /f "tokens=3" %%v in ('go version') do echo        Go %%v encontrado
echo.

echo [2/6] Verificando .env...
if not exist ".env" (
    echo.
    echo  ERRO: .env nao encontrado em %cd%
    echo  Copie .env.example para .env e configure as variaveis.
    echo.
    pause
    exit /b 1
)
echo        .env encontrado
echo.

echo [3/6] Criando estrutura de pastas...
for %%d in (
    ..\.agent_handoff\raw
    ..\.agent_handoff\processing
    ..\.agent_handoff\pending_python
    ..\.agent_handoff\extracted
    ..\.agent_handoff\processed_python
    ..\.agent_handoff\success
    ..\.agent_handoff\failed
    ..\AGENTS
    .bin
    .runtime
    ..\output
    ..\data
    ..\logs
) do (
    if not exist "%%d" (
        mkdir "%%d"
        echo        Criada: %%d
    )
)
echo.

echo [4/6] Verificando Ollama...
set OLLAMA_URL=http://localhost:11434
for /f "usebackq tokens=1,* delims==" %%a in (".env") do (
    if /i "%%a"=="OLLAMA_URL" set OLLAMA_URL=%%b
)

curl -sf --max-time 3 "%OLLAMA_URL%/api/tags" >nul 2>&1
if errorlevel 1 (
    echo.
    echo  Ollama nao respondeu em %OLLAMA_URL%
    echo  Tentando iniciar Ollama...

    where ollama >nul 2>&1
    if errorlevel 1 (
        echo.
        echo  AVISO: ollama nao encontrado no PATH.
        echo  Instale em: https://ollama.com/download
        echo.
        choice /c SN /m "Continuar mesmo assim? (S/N)"
        if errorlevel 2 exit /b 1
    ) else (
        start "Ollama" /min ollama serve
        set OLLAMA_READY=0
        for /l %%i in (1,1,20) do (
            if "!OLLAMA_READY!"=="0" (
                curl -sf --max-time 2 "%OLLAMA_URL%/api/tags" >nul 2>&1
                if errorlevel 1 (
                    timeout /t 1 /nobreak >nul
                ) else (
                    set OLLAMA_READY=1
                )
            )
        )

        if "!OLLAMA_READY!"=="0" (
            echo.
            echo  AVISO: Ollama foi iniciado, mas ainda nao respondeu em %OLLAMA_URL%
            echo.
            choice /c SN /m "Continuar mesmo assim? (S/N)"
            if errorlevel 2 exit /b 1
        ) else (
            echo        Ollama iniciado e acessivel em %OLLAMA_URL%
        )
    )
) else (
    echo        Ollama acessivel em %OLLAMA_URL%
)
echo.

if "%DB_PATH%"=="" set DB_PATH=%cd%\.runtime\agent_runtime.db
if "%AGENT_TIMEOUT_SECONDS%"=="" set AGENT_TIMEOUT_SECONDS=1800
echo        DB runtime: %DB_PATH%
echo.

echo [5/6] Compilando...
go build -ldflags="-s -w" -o %BINARY% .
if errorlevel 1 (
    echo.
    echo  ERRO: Falha na compilacao. Verifique os erros acima.
    echo.
    pause
    exit /b 1
)
echo        Compilado: %BINARY%
echo.

echo [6/6] Iniciando PF Agent...
echo ================================================================
echo  PF Agent iniciado. Pressione Ctrl+C para encerrar.
echo ================================================================
echo.

%BINARY%

echo.
echo  PF Agent encerrado.
pause
