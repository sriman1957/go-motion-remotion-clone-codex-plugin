package templates

const HTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Go Motion Composition</title>
  <link rel="stylesheet" href="styles.css">
</head>
<body>
  <main id="app" class="app">
    <section class="hero">
      <p class="eyebrow">Go Motion</p>
      <h1 id="headline">Loading composition...</h1>
      <p id="body" class="body"></p>
    </section>
  </main>
  <script src="composition.js"></script>
  <script src="runtime.js"></script>
</body>
</html>
`

const CSS = `:root {
  --bg: #07111f;
  --panel: #0f1d35;
  --text: #f5f7fb;
  --muted: #b7c1d8;
  --accent: #38bdf8;
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  min-height: 100vh;
  display: grid;
  place-items: center;
  background:
    radial-gradient(circle at top, rgba(56, 189, 248, 0.26), transparent 35%),
    linear-gradient(160deg, #040814, #0b1530 58%, #07111f);
  color: var(--text);
  font-family: "Segoe UI", Arial, sans-serif;
}

.app {
  width: 100%;
  max-width: 960px;
  padding: 72px;
}

.hero {
  padding: 56px;
  border: 1px solid rgba(183, 193, 216, 0.18);
  border-radius: 28px;
  background: rgba(15, 29, 53, 0.78);
  backdrop-filter: blur(12px);
  box-shadow: 0 30px 80px rgba(0, 0, 0, 0.28);
}

.eyebrow {
  margin: 0 0 20px;
  color: var(--accent);
  text-transform: uppercase;
  letter-spacing: 0.16em;
  font-size: 14px;
}

h1 {
  margin: 0;
  font-size: 62px;
  line-height: 1;
}

.body {
  margin-top: 20px;
  max-width: 720px;
  color: var(--muted);
  font-size: 22px;
  line-height: 1.5;
}
`

const JS = `function currentFrame() {
  const params = new URLSearchParams(window.location.search);
  const value = Number(params.get("frame") || "0");
  if (Number.isNaN(value) || value < 0) {
    return 0;
  }
  return Math.floor(value);
}

function sceneForFrame(composition, frame) {
  let cursor = 0;
  for (const scene of composition.scenes || []) {
    const start = cursor;
    const end = cursor + (scene.durationFrames || 0);
    if (frame >= start && frame < end) {
      return scene;
    }
    cursor = end;
  }
  return (composition.scenes || [])[0] || {};
}

function boot() {
  const composition = window.__GO_MOTION_COMPOSITION__ || { title: "Untitled", scenes: [] };
  const frame = currentFrame();
  const scene = sceneForFrame(composition, frame);
  document.title = composition.title;
  document.getElementById("headline").textContent = scene.headline || composition.title;
  document.getElementById("body").textContent = scene.body || "Prompt-driven composition ready for rendering.";
  document.body.dataset.frame = String(frame);
  document.body.dataset.scene = scene.id || "";
}

boot();
`
