@echo off
setlocal EnableDelayedExpansion

title PF Agent

:: ── Configuração ─────────────────────────────────────────────────
set GOPATH=%USERPROFILE%\go
set GOMODCACHE=%USERPROFILE%\go\pkg\mod
set CGO_ENABLED=0
set BINARY=pf-agent.exe

:: ── Banner ───────────────────────────────────────────────────────
echo.
echo  ██████╗ ███████╗    █████╗  ██████╗ ███████╗███╗   ██╗████████╗
echo  ██╔══██╗██╔════╝   ██╔══██╗██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝
echo  ██████╔╝█████╗     ███████║██║  ███╗█████╗  ██╔██╗ ██║   ██║   
echo  ██╔═══╝ ██╔══╝     ██╔══██║██║   ██║██╔══╝  ██║╚██╗██║   ██║   
echo  ██║     ██║        ██║  ██║╚██████╔╝███████╗██║ ╚████║   ██║   
echo  ╚═╝     ╚═╝        ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   
echo.
echo  [STATUS] Hexagonal Architecture Active | GSD Framework Loaded
echo.

:: ── Verificar Go instalado ────────────────────────────────────────
echo [1/5] Verificando Go...
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

:: ── Carregar .env.example ─────────────────────────────────────────
echo [2/5] Verificando .env...
if not exist "..\.env" (
    echo.
    echo  ERRO: .env.example nao encontrado em %cd%\..
    echo.
    pause
    exit /b 1
)
echo        .env.example encontrado
echo.

:: ── Criar pastas necessarias ──────────────────────────────────────
echo [3/5] Criando estrutura de pastas...
for %%d in (
    ..\.agent_handoff\raw
    ..\.agent_handoff\processing
    ..\.agent_handoff\pending_python
    ..\.agent_handoff\extracted
    ..\.agent_handoff\processed_python
    ..\.agent_handoff\success
    ..\.agent_handoff\failed
    vault
    data
    logs
) do (
    if not exist "%%d" (
        mkdir "%%d"
        echo        Criada: %%d
    )
)
echo.

:: ── Compilar ──────────────────────────────────────────────────────
echo [4/5] Compilando...
if exist "%BINARY%" del /f /q "%BINARY%"

go build -o %BINARY% .
if errorlevel 1 (
    echo.
    echo  ERRO: Falha na compilacao. Verifique os erros acima.
    echo.
    pause
    exit /b 1
)
echo        Compilado: %BINARY%
echo.

:: ── Rodar ─────────────────────────────────────────────────────────
echo [5/5] Iniciando PF Agent...
echo ══════════════════════════════════════════════════════════════
echo  PF Agent iniciado. Pressione Ctrl+C para encerrar.
echo ══════════════════════════════════════════════════════════════
echo.

%BINARY%

:: ── Encerrado ─────────────────────────────────────────────────────
echo.
echo  PF Agent encerrado.
pause