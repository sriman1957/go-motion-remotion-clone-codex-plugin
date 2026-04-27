# Go Motion Design

**Problem**

Build a Codex plugin that can generate professional prompt-to-video output locally without requiring Node.js or Go to be installed on the end user's system.

**Goal**

Ship a self-contained plugin that:

- accepts a user prompt for a video
- plans a polished composition
- renders the composition with HTML, CSS, and vanilla JS in a real browser engine
- encodes the result to `mp4`
- is hosted by a Go MCP server

**Non-Goals for v1**

- React or Remotion compatibility
- image sequence export
- timeline editing UI
- cloud rendering dependency
- bundled binaries checked into this repo

**Key Constraints**

- End users should not need Node.js.
- End users should not need Go installed.
- Output should target professional-looking marketing and explainer videos.
- The plugin should be structured as a Codex plugin, not just a Go CLI.
- The rendering model should use browser-native HTML/CSS/JS for high visual fidelity.

## Product Shape

The plugin acts as a prompt-to-video studio for Codex. The user asks for a video, the plugin turns the prompt into a structured composition spec, materializes HTML/CSS/JS scenes, renders those scenes frame-by-frame in headless Chromium, and encodes the final video to `mp4` with FFmpeg.

The system is "similar to Remotion" in output ambition and composition quality, not in React API compatibility. The authoring model is browser-native and intentionally avoids Node.js.

## Architecture

The system has six major parts:

1. **Codex plugin shell**
   Contains `.codex-plugin/plugin.json`, `.mcp.json`, docs, and packaging metadata.

2. **Go MCP server**
   Exposes tools such as `generate_video`, `list_styles`, and `inspect_render`. Owns request validation, workspace setup, orchestration, and status reporting.

3. **Prompt planner**
   Converts natural-language prompts into a normalized composition spec with duration, pacing, scenes, typography direction, motion presets, and audio hints.

4. **Composition generator**
   Converts the spec into HTML, CSS, and vanilla JS files plus a manifest describing FPS, duration, viewport size, and scene timing.

5. **Renderer**
   Launches a Chromium-compatible browser in headless mode, opens the generated composition runtime, advances frame state deterministically, and captures screenshots for each frame.

6. **Encoder**
   Invokes FFmpeg to combine frames and optional audio into a final `mp4`.

## Plugin Runtime Strategy

The plugin is designed around bundled native runtimes:

- Chromium-compatible browser binary
- FFmpeg binary
- Go server executable

This repo will implement the runtime discovery and extraction model, but will not check large vendor binaries into source control. Instead, the code will expect platform-specific runtime directories under the plugin root such as:

- `runtime/win64/chromium/`
- `runtime/win64/ffmpeg/`
- `runtime/linux-x64/chromium/`
- `runtime/darwin-arm64/chromium/`

The server must support:

- plugin-local runtime discovery
- clear errors when a required runtime is missing
- optional fallback to system `chrome` / `ffmpeg` during development

## Composition Model

The internal composition contract is a JSON manifest:

- video metadata: title, width, height, fps, durationFrames
- theme tokens: colors, fonts, spacing, motion style
- scenes: ordered list of timed scenes
- layers: text, shape, image, video, background, caption
- animations: entry, emphasis, exit, transform, opacity, blur
- audio tracks: optional voice, music, SFX

Each generated composition includes:

- `composition.json`
- `index.html`
- `styles.css`
- `runtime.js`
- scene partials or data payloads

The runtime JS reads the frame number from the renderer and computes the visual state for that frame. This keeps rendering deterministic.

## Rendering Flow

1. User calls `generate_video` with a prompt and optional settings.
2. Server creates a render job directory.
3. Prompt planner emits a normalized composition spec.
4. Composition generator writes HTML/CSS/JS runtime files.
5. Renderer launches browser and captures frames.
6. Encoder produces `output.mp4`.
7. Tool returns output path plus job metadata.

## Error Handling

The system must fail clearly when:

- Chromium runtime cannot be found
- FFmpeg runtime cannot be found
- composition generation produces invalid files
- renderer page crashes
- frame capture stalls
- encoding fails

Errors should say what was missing, where the server looked, and how to fix it.

## Testing Strategy

V1 testing focuses on behavior that can be verified cheaply in repo:

- prompt normalization tests
- composition manifest generation tests
- runtime path resolution tests
- render command construction tests
- smoke integration test for generating a composition package

Full browser rendering tests should be integration-gated and skipped automatically when runtimes are absent.

## Repo Layout

Planned layout:

- `plugins/go-motion/`
- `plugins/go-motion/.codex-plugin/plugin.json`
- `plugins/go-motion/.mcp.json`
- `plugins/go-motion/README.md`
- `plugins/go-motion/go.mod`
- `plugins/go-motion/cmd/go-motiond/main.go`
- `plugins/go-motion/internal/...`
- `plugins/go-motion/testdata/...`

## Success Criteria for v1

- Codex plugin scaffolds cleanly in this repo.
- Go server builds and starts locally.
- MCP tool surface is implemented.
- A prompt can generate a composition package on disk.
- If browser and FFmpeg runtimes are present, the tool can render an `mp4`.
- Missing runtimes produce actionable errors instead of crashes.

