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
