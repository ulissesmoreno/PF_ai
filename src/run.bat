@echo off
setlocal EnableDelayedExpansion

title PF Agent

:: ── Configuração ─────────────────────────────────────────────────
set CGO_ENABLED=0
set BINARY=..\bin\pf-agent.exe

:: ── Banner ───────────────────────────────────────────────────────
echo.
echo  ██████╗ ███████╗    █████╗  ██████╗ ███████╗███╗   ██╗████████╗
echo  ██╔══██╗██╔════╝   ██╔══██╗██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝
echo  ██████╔╝█████╗     ███████║██║  ███╗█████╗  ██╔██╗ ██║   ██║
echo  ██╔═══╝ ██╔══╝     ██╔══██║██║   ██║██╔══╝  ██║╚██╗██║   ██║
echo  ██║     ██║        ██║  ██║╚██████╔╝███████╗██║ ╚████║   ██║
echo  ╚═╝     ╚═╝        ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝
echo.
echo  [STATUS] Hexagonal Architecture Active ^| GSD Framework Loaded
echo.

:: ── [1/6] Verificar Go instalado ───────────────────────────────────
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

:: ── [2/6] Verificar .env ───────────────────────────────────────────
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

:: ── [3/6] Criar estrutura de pastas ────────────────────────────────
echo [3/6] Criando estrutura de pastas...

:: Pastas do pipeline (relativas a src/, pois o binario roda de src/)
for %%d in (    
    ..\.agent_handoff\processing
    ..\.agent_handoff\pending_python
    ..\.agent_handoff\extracted
    ..\.agent_handoff\processed_python
    ..\.agent_handoff\success
    ..\.agent_handoff\failed
    .agent_handoff
    ..\AGENTS
    ..\bin    
    output
    data
    logs
) do (
    if not exist "%%d" (
        mkdir "%%d"
        echo        Criada: %%d
    )
)
echo.

:: ── [4/6] Verificar Ollama acessivel ───────────────────────────────
echo [4/6] Verificando Ollama...

:: Lê OLLAMA_URL do .env, usa http://localhost:11434 como fallback
set OLLAMA_URL=http://localhost:11434
for /f "usebackq tokens=1,* delims==" %%a in (".env") do (
    if /i "%%a"=="OLLAMA_URL" set OLLAMA_URL=%%b
)

:: Testa conectividade com a API do Ollama (curl disponivel no Windows 10+)
curl -sf --max-time 3 "%OLLAMA_URL%/api/tags" >nul 2>&1
if errorlevel 1 (
    echo.
    echo  AVISO: Ollama nao respondeu em %OLLAMA_URL%
    echo  Certifique-se de que o Ollama esta rodando antes de continuar.
    echo.
    choice /c SN /m "Continuar mesmo assim? (S/N)"
    if errorlevel 2 exit /b 1
) else (
    echo        Ollama acessivel em %OLLAMA_URL%
)
echo.

:: ── [5/6] Compilar ─────────────────────────────────────────────────
echo [5/6] Compilando...
if exist "%BINARY%" del /f /q "%BINARY%"

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

:: ── [6/6] Rodar ────────────────────────────────────────────────────
echo [6/6] Iniciando PF Agent...
echo ══════════════════════════════════════════════════════════════
echo  PF Agent iniciado. Pressione Ctrl+C para encerrar.
echo ══════════════════════════════════════════════════════════════
echo.

%BINARY%

:: ── Encerrado ──────────────────────────────────────────────────────
echo.
echo  PF Agent encerrado.
pause