# Go Motion Usage

## What this plugin does

Go Motion exposes an MCP tool surface for Codex and renders videos locally from prompts.

The main tool today is:

- `generate_video`: create a composition, render frames in Chromium, and encode an MP4 with FFmpeg

## Windows runtime layout

The current working package expects:

- `plugins/go-motion/runtime/windows-amd64/bin/go-motiond.exe`
- `plugins/go-motion/runtime/windows-amd64/chromium/chrome.exe`
- `plugins/go-motion/runtime/windows-amd64/ffmpeg/ffmpeg.exe`

These are already bundled in this workspace.

## Planned cross-platform runtime matrix

The repo is now structured so future releases can ship the same plugin for multiple platform targets:

- `windows-amd64`
- `windows-arm64`
- `darwin-amd64`
- `darwin-arm64`
- `linux-amd64`
- `linux-arm64`

Each target should use the same directory contract:

- `plugins/go-motion/runtime/<platform>/bin/go-motiond`
- `plugins/go-motion/runtime/<platform>/chromium/...`
- `plugins/go-motion/runtime/<platform>/ffmpeg/...`

For Windows targets, the server binary name is `go-motiond.exe`.

## Build scripts for runtime binaries

From `plugins/go-motion/`, contributors can build server binaries for one or more targets:

PowerShell:

- `./scripts/package-runtime.ps1`
- `./scripts/package-runtime.ps1 -Targets @("windows/amd64", "darwin/arm64")`

Shell:

- `./scripts/package-runtime.sh`
- `./scripts/package-runtime.sh windows/amd64 linux/amd64`

## How the render flow works

1. Codex invokes `generate_video`
2. The Go MCP server creates a job directory under `plugins/go-motion/jobs/`
3. A composition package is generated with:
   - `index.html`
   - `styles.css`
   - `composition.js`
   - `runtime.js`
4. Headless Chromium captures PNG frames
5. FFmpeg encodes the frames into `output.mp4`

## Example output location

Successful renders are written under:

- `plugins/go-motion/jobs/<job-id>/output.mp4`

## Current constraints

- The packaged runtime currently targets Windows `amd64`
- Rendering is frame-by-frame and can take time on larger compositions
- The default planner currently generates a short promo-style composition for practical local runtime
