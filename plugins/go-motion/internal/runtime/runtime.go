package runtime

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Tools struct {
	BrowserPath string
	FFmpegPath  string
}

func ResolveTools(pluginRoot string) (Tools, error) {
	platformDir := platformKey()
	browserCandidates := browserCandidatePaths(pluginRoot)
	ffmpegCandidates := ffmpegCandidatePaths(pluginRoot)

	tools := Tools{
		BrowserPath: firstExisting(browserCandidates),
		FFmpegPath:  firstExisting(ffmpegCandidates),
	}

	if tools.BrowserPath == "" {
		tools.BrowserPath = lookPathAny(browserFallbacks())
	}
	if tools.FFmpegPath == "" {
		tools.FFmpegPath = lookPathAny([]string{"ffmpeg"})
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
		"required runtimes not found under %s (missing: %s)",
		filepath.Join(pluginRoot, "runtime", platformDir),
		strings.Join(missing, ", "),
	)
}

func browserCandidatePaths(pluginRoot string) []string {
	paths := []string{
		filepath.Join(pluginRoot, "runtime", platformKey(), "chromium", browserExecutableName()),
		filepath.Join(pluginRoot, "runtime", platformKey(), "chrome", browserExecutableName()),
		filepath.Join(pluginRoot, "runtime", platformKey(), "edge", edgeExecutableName()),
	}

	if runtime.GOOS == "windows" {
		paths = append(paths,
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
		)
	}

	return paths
}

func ffmpegCandidatePaths(pluginRoot string) []string {
	paths := []string{
		filepath.Join(pluginRoot, "runtime", platformKey(), "ffmpeg", ffmpegExecutableName()),
	}
	if runtime.GOOS == "windows" {
		paths = append(paths,
			`C:\ffmpeg\bin\ffmpeg.exe`,
			`C:\Program Files\ffmpeg\bin\ffmpeg.exe`,
			`C:\Program Files (x86)\ffmpeg\bin\ffmpeg.exe`,
		)
	}
	return paths
}

func platformKey() string {
	return runtime.GOOS + "-" + runtime.GOARCH
}

func browserExecutableName() string {
	if runtime.GOOS == "windows" {
		return "chrome.exe"
	}
	return "chrome"
}

func edgeExecutableName() string {
	if runtime.GOOS == "windows" {
		return "msedge.exe"
	}
	return "msedge"
}

func ffmpegExecutableName() string {
	if runtime.GOOS == "windows" {
		return "ffmpeg.exe"
	}
	return "ffmpeg"
}

func browserFallbacks() []string {
	if runtime.GOOS == "windows" {
		return []string{"chrome", "msedge", "chromium"}
	}
	if runtime.GOOS == "darwin" {
		return []string{"Google Chrome", "Chromium", "chrome", "chromium"}
	}
	return []string{"google-chrome", "chromium", "chrome"}
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
