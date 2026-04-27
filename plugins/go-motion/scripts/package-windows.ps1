$ErrorActionPreference = "Stop"

$script = Join-Path $PSScriptRoot "package-runtime.ps1"
& $script -Targets @("windows/amd64")
