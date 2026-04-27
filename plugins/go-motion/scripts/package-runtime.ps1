param(
  [string[]]$Targets = @(
    "windows/amd64",
    "windows/arm64",
    "darwin/amd64",
    "darwin/arm64",
    "linux/amd64",
    "linux/arm64"
  )
)

$ErrorActionPreference = "Stop"

$pluginRoot = Split-Path -Parent $PSScriptRoot
Push-Location $pluginRoot
try {
  foreach ($target in $Targets) {
    $parts = $target.Split("/")
    if ($parts.Length -ne 2) {
      throw "Invalid target '$target'. Expected format os/arch."
    }

    $goos = $parts[0]
    $goarch = $parts[1]
    $platformKey = "$goos-$goarch"
    $binDir = Join-Path $pluginRoot "runtime\$platformKey\bin"
    New-Item -ItemType Directory -Force -Path $binDir | Out-Null

    $binaryName = "go-motiond"
    if ($goos -eq "windows") {
      $binaryName = "go-motiond.exe"
    }

    $output = Join-Path $binDir $binaryName
    Write-Output "Building $platformKey -> $output"

    $env:GOOS = $goos
    $env:GOARCH = $goarch
    $env:CGO_ENABLED = "0"

    go build -o $output .\cmd\go-motiond
  }
} finally {
  Remove-Item Env:GOOS -ErrorAction SilentlyContinue
  Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
  Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
  Pop-Location
}
