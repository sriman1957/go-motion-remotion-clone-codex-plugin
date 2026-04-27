# Contributing

Thanks for wanting to contribute to Go Motion.

## What this project is

This repository contains a Codex plugin that aims to be a Go-powered, browser-native alternative to a Remotion-style video workflow:

- Go MCP server
- HTML/CSS/vanilla JS compositions
- bundled Chromium for rendering
- bundled FFmpeg for MP4 encoding

## Good first contribution areas

- improve prompt planning and scene generation
- add better motion systems and transitions
- improve rendering speed
- add Linux and macOS runtime packaging
- add more tests around render fidelity and runtime detection
- improve docs and onboarding

## Local development

From the plugin directory:

```powershell
cd plugins/go-motion
go test ./...
powershell -ExecutionPolicy Bypass -File .\scripts\package-windows.ps1
```

## Repo structure

- `plugins/go-motion/`: main plugin code and bundled runtimes
- `docs/`: design and usage documentation

## Contribution rules

- keep changes focused
- add or update tests when behavior changes
- prefer small, reviewable pull requests
- document user-visible behavior changes in the README or usage docs

## Runtime note

The Windows runtime bundle is currently committed so the plugin can work without requiring end users to install Node.js or Go.

