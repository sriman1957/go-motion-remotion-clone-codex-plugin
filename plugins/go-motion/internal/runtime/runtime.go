package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	gruntime "runtime"
	"slices"
	"strings"
)

type Platform struct {
	OS   string
	Arch string
}

func (p Platform) Key() string {
	return p.OS + "-" + p.Arch
}

func (p Platform) BrowserExecutableName() string {
	switch p.OS {
	case "windows":
		return "chrome.exe"
	case "darwin":
		return "Chromium.app/Contents/MacOS/Chromium"
	default:
		return "chrome"
	}
}

func (p Platform) EdgeExecutableName() string {
	switch p.OS {
	case "windows":
		return "msedge.exe"
	case "darwin":
		return "Microsoft Edge.app/Contents/MacOS/Microsoft Edge"
	default:
		return "microsoft-edge"
	}
}

func (p Platform) FFmpegExecutableName() string {
	if p.OS == "windows" {
		return "ffmpeg.exe"
	}
	return "ffmpeg"
}

func (p Platform) ServerExecutableName() string {
	if p.OS == "windows" {
		return "go-motiond.exe"
	}
	return "go-motiond"
}

type Tools struct {
	BrowserPath string
	FFmpegPath  string
	Platform    Platform
	RuntimeDir  string
	ServerPath  string
}

func CurrentPlatform() Platform {
	return Platform{
		OS:   gruntime.GOOS,
		Arch: gruntime.GOARCH,
	}
}

func SupportedPlatforms() []Platform {
	return []Platform{
		{OS: "windows", Arch: "amd64"},
		{OS: "windows", Arch: "arm64"},
		{OS: "darwin", Arch: "amd64"},
		{OS: "darwin", Arch: "arm64"},
		{OS: "linux", Arch: "amd64"},
		{OS: "linux", Arch: "arm64"},
	}
}

func ResolveTools(pluginRoot string) (Tools, error) {
	return ResolveToolsForPlatform(pluginRoot, CurrentPlatform())
}

func ResolveToolsForPlatform(pluginRoot string, platform Platform) (Tools, error) {
	runtimeDir := RuntimeDir(pluginRoot, platform)
	tools := Tools{
		BrowserPath: firstExisting(browserCandidatePaths(pluginRoot, platform)),
		FFmpegPath:  firstExisting(ffmpegCandidatePaths(pluginRoot, platform)),
		Platform:    platform,
		RuntimeDir:  runtimeDir,
		ServerPath:  LauncherPath(pluginRoot, platform),
	}

	if platform == CurrentPlatform() {
		if tools.BrowserPath == "" {
			tools.BrowserPath = lookPathAny(browserFallbacks(platform))
		}
		if tools.FFmpegPath == "" {
			tools.FFmpegPath = lookPathAny(ffmpegFallbacks(platform))
		}
	}

	if tools.BrowserPath != "" && tools.FFmpegPath != "" {
		return tools, nil
	}

	var missing []string
	if tools.BrowserPath == "" {
		missing = append(missing, "browser")
	}
	if tools.FFmpegPath == "" {
		missing = append(missing, "ffmpeg")
	}

	return Tools{}, fmt.Errorf(
		"required runtimes not found for %s under %s (missing: %s; supported release targets: %s)",
		platform.Key(),
		runtimeDir,
		strings.Join(missing, ", "),
		strings.Join(SupportedPlatformKeys(), ", "),
	)
}

func SupportedPlatformKeys() []string {
	keys := make([]string, 0, len(SupportedPlatforms()))
	for _, platform := range SupportedPlatforms() {
		keys = append(keys, platform.Key())
	}
	return keys
}

func RuntimeDir(pluginRoot string, platform Platform) string {
	return filepath.Join(pluginRoot, "runtime", platform.Key())
}

func LauncherPath(pluginRoot string, platform Platform) string {
	return filepath.Join(RuntimeDir(pluginRoot, platform), "bin", platform.ServerExecutableName())
}

func browserCandidatePaths(pluginRoot string, platform Platform) []string {
	runtimeDir := RuntimeDir(pluginRoot, platform)
	paths := []string{
		filepath.Join(runtimeDir, "chromium", platform.BrowserExecutableName()),
		filepath.Join(runtimeDir, "chrome", platform.BrowserExecutableName()),
		filepath.Join(runtimeDir, "edge", platform.EdgeExecutableName()),
	}

	switch platform.OS {
	case "windows":
		paths = append(paths,
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
		)
	case "darwin":
		paths = append(paths,
			`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`,
			`/Applications/Chromium.app/Contents/MacOS/Chromium`,
			`/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge`,
		)
	case "linux":
		paths = append(paths,
			`/usr/bin/google-chrome`,
			`/usr/bin/google-chrome-stable`,
			`/usr/bin/chromium`,
			`/usr/bin/chromium-browser`,
			`/usr/bin/microsoft-edge`,
		)
	}

	return dedupe(paths)
}

func ffmpegCandidatePaths(pluginRoot string, platform Platform) []string {
	paths := []string{
		filepath.Join(RuntimeDir(pluginRoot, platform), "ffmpeg", platform.FFmpegExecutableName()),
	}

	switch platform.OS {
	case "windows":
		paths = append(paths,
			`C:\ffmpeg\bin\ffmpeg.exe`,
			`C:\Program Files\ffmpeg\bin\ffmpeg.exe`,
			`C:\Program Files (x86)\ffmpeg\bin\ffmpeg.exe`,
		)
	case "darwin":
		paths = append(paths,
			`/opt/homebrew/bin/ffmpeg`,
			`/usr/local/bin/ffmpeg`,
			`/usr/bin/ffmpeg`,
		)
	case "linux":
		paths = append(paths,
			`/usr/bin/ffmpeg`,
			`/usr/local/bin/ffmpeg`,
		)
	}

	return dedupe(paths)
}

func browserFallbacks(platform Platform) []string {
	switch platform.OS {
	case "windows":
		return []string{"chrome", "msedge", "chromium"}
	case "darwin":
		return []string{"Google Chrome", "Chromium", "chrome", "chromium"}
	default:
		return []string{"google-chrome", "google-chrome-stable", "chromium", "chrome", "microsoft-edge"}
	}
}

func ffmpegFallbacks(platform Platform) []string {
	switch platform.OS {
	case "windows":
		return []string{"ffmpeg.exe", "ffmpeg"}
	default:
		return []string{"ffmpeg"}
	}
}

func firstExisting(paths []string) string {
	for _, p := range paths {
		info, err := os.Stat(p)
		if err == nil && !info.IsDir() {
			return p
		}
	}
	return ""
}

func lookPathAny(names []string) string {
	for _, name := range names {
		path, err := exec.LookPath(name)
		if err == nil {
			return path
		}
	}
	return ""
}

func dedupe(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || slices.Contains(out, value) {
			continue
		}
		out = append(out, value)
	}
	return out
}
