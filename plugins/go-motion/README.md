# Go Motion Plugin

This directory contains the actual Codex plugin bundle for Go Motion.

## Purpose

The plugin exposes a Go MCP server that accepts prompt-based video generation requests and renders MP4 output through a browser-native pipeline:

- scene planning
- HTML, CSS, and vanilla JS composition generation
- headless Chromium frame capture
- FFmpeg encoding

## Bundle structure

- [`.codex-plugin/plugin.json`](./.codex-plugin/plugin.json): plugin metadata
- [`.mcp.json`](./.mcp.json): MCP server launcher
- [`cmd/go-motiond`](./cmd/go-motiond): server entrypoint
- [`internal/mcp`](./internal/mcp): MCP protocol and tool handlers
- [`internal/planner`](./internal/planner): prompt-to-composition planning
- [`internal/templates`](./internal/templates): composition templates
- [`internal/render`](./internal/render): browser rendering and encoding orchestration
- [`internal/runtime`](./internal/runtime): bundled runtime discovery
- [`runtime/windows-amd64`](./runtime/windows-amd64): packaged Windows runtime

## Cross-platform foundation

The runtime layer now treats platform as an explicit OS and CPU target instead of assuming a single Windows bundle. The intended release matrix is:

- `windows-amd64`
- `windows-arm64`
- `darwin-amd64`
- `darwin-arm64`
- `linux-amd64`
- `linux-arm64`

Each platform release is expected to provide the same layout under `runtime/<platform-key>/`:

- `bin/go-motiond` or `bin/go-motiond.exe`
- `chromium/...`
- `ffmpeg/...`

## Tooling model

Go Motion is intentionally Node-free on the user machine.

- No Node.js install required
- No npm dependency required
- No Go install required for Windows end users of the packaged bundle

The shipped runtime provides:

- `go-motiond.exe`
- Chromium
- FFmpeg

## Current capabilities

- `initialize`, `tools/list`, and `tools/call` MCP flow
- `generate_video` for prompt-to-video rendering
- `list_styles` for available composition styles
- local composition packaging into HTML, CSS, and JS files
- frame-by-frame rendering into MP4 output

## Current constraints

- packaged runtime currently targets Windows `amd64`
- rendering is local and can take time for larger jobs
- the composition model is browser-native rather than React-based

## Packaging scripts

- [`scripts/package-runtime.ps1`](./scripts/package-runtime.ps1): build one or many platform server binaries on Windows
- [`scripts/package-runtime.sh`](./scripts/package-runtime.sh): build one or many platform server binaries on Unix-like systems
- [`scripts/package-windows.ps1`](./scripts/package-windows.ps1): quick Windows `amd64` wrapper

## Related docs

- [`../../README.md`](../../README.md)
- [`../../docs/USAGE.md`](../../docs/USAGE.md)
