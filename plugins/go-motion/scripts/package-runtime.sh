#!/usr/bin/env bash
set -euo pipefail

PLUGIN_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ "$#" -eq 0 ]; then
  TARGETS=(
    "windows/amd64"
    "windows/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
  )
else
  TARGETS=("$@")
fi

pushd "$PLUGIN_ROOT" >/dev/null
for target in "${TARGETS[@]}"; do
  IFS="/" read -r goos goarch <<<"$target"
  if [ -z "${goos:-}" ] || [ -z "${goarch:-}" ]; then
    echo "invalid target: $target" >&2
    exit 1
  fi

  platform_key="${goos}-${goarch}"
  bin_dir="$PLUGIN_ROOT/runtime/$platform_key/bin"
  mkdir -p "$bin_dir"

  binary_name="go-motiond"
  if [ "$goos" = "windows" ]; then
    binary_name="go-motiond.exe"
  fi

  output="$bin_dir/$binary_name"
  echo "Building $platform_key -> $output"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -o "$output" ./cmd/go-motiond
done
popd >/dev/null
