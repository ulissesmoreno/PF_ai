param(
  [string]$Port = "8081"
)

$ErrorActionPreference = "Stop"

if (-not $env:PF_AI_AUTH_SECRET) {
  Write-Error "PF_AI_AUTH_SECRET must be set in the current shell."
}

$workspace = Split-Path -Parent $PSScriptRoot
$env:PF_AI_HTTP_PORT = $Port
$env:GOCACHE = Join-Path $workspace ".gocache"

Set-Location $workspace
& "C:\Program Files\Go\bin\go.exe" run ./cmd/pf-ai-service
