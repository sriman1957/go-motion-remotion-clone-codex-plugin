# Go Motion

Go Motion is a Codex plugin that generates prompt-to-video compositions using a Go MCP server and a browser-native HTML/CSS/vanilla JS runtime.

## Status

This repository currently contains the initial plugin scaffold and render pipeline architecture. The code is designed for bundled Chromium-compatible and FFmpeg runtimes, but those binaries are not checked into source control.

## Runtime requirements

Go Motion expects bundled Chromium-compatible and FFmpeg runtimes under `runtime/<platform>/`.

During development, the server may optionally fall back to compatible system-installed tools when plugin-local runtimes are unavailable.

## Current capabilities

- prompt planning into a structured composition
- HTML/CSS/vanilla JS composition package generation
- MCP stdio tool surface with `generate_video` and `list_styles`
- frame-by-frame browser screenshot rendering
- MP4 encoding orchestration through FFmpeg

## Current limitation

Real MP4 output requires both a browser runtime and FFmpeg. On machines where only the browser is available, the plugin will still generate the composition package and render plan, but it will report the missing FFmpeg runtime instead of silently failing.
