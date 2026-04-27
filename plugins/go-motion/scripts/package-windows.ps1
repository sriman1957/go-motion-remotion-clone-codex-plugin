$ErrorActionPreference = "Stop"

$pluginRoot = Split-Path -Parent $PSScriptRoot
$binDir = Join-Path $pluginRoot "runtime\windows-amd64\bin"
New-Item -ItemType Directory -Force -Path $binDir | Out-Null

$output = Join-Path $binDir "go-motiond.exe"
Push-Location $pluginRoot
try {
  go build -o $output .\cmd\go-motiond
  Write-Output $output
} finally {
  Pop-Location
}
