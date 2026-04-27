# Go Motion

Go Motion is a Codex plugin for generating prompt-to-video renders with:

- Go for the MCP server
- HTML, CSS, and vanilla JS for compositions
- bundled Chromium for browser rendering
- bundled FFmpeg for MP4 encoding

It is designed to work without requiring Node.js on the end user's machine, and the Windows packaging in this repo also avoids requiring Go to be installed by the user.

## Repository layout

- [`plugins/go-motion`](./plugins/go-motion): the actual Codex plugin
- [`docs/superpowers/specs`](./docs/superpowers/specs): design docs
- [`docs/superpowers/plans`](./docs/superpowers/plans): implementation plan

## Current status

The Windows plugin is working end to end:

- packaged server binary
- bundled Chromium runtime
- bundled FFmpeg runtime
- real prompt-to-video rendering to MP4 through the MCP tool surface

## Main plugin files

- [plugin manifest](./plugins/go-motion/.codex-plugin/plugin.json)
- [MCP launcher](./plugins/go-motion/.mcp.json)
- [plugin README](./plugins/go-motion/README.md)

## Example output

One successful rendered output from verification is located at:

- `plugins/go-motion/jobs/job-7f00da0f16c26c1c/output.mp4`

## Notes

- The current bundled runtime targets Windows `amd64`.
- The composition engine is intentionally browser-native rather than React-based.
- The default planner currently generates a concise 4-second style promo to keep local rendering practical.
