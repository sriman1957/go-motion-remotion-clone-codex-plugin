package runtime

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestResolveTools_ReportsMissingRuntime(t *testing.T) {
	_, err := ResolveToolsForPlatform("C:/tmp/does-not-exist", Platform{OS: "linux", Arch: "amd64"})
	if err == nil {
		t.Fatal("expected missing runtime error")
	}
	if !strings.Contains(err.Error(), "required runtimes not found for linux-amd64") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "supported release targets") {
		t.Fatalf("expected supported targets in error, got %v", err)
	}
}

func TestBrowserCandidatePaths_IncludeWindowsInstallLocations(t *testing.T) {
	paths := browserCandidatePaths(filepath.Clean("C:/plugin"), Platform{OS: "windows", Arch: "amd64"})
	joined := strings.Join(paths, "\n")
	if !strings.Contains(joined, `C:\Program Files\Google\Chrome\Application\chrome.exe`) {
		t.Fatalf("expected Chrome install path in candidates, got %s", joined)
	}
	if !strings.Contains(joined, `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`) {
		t.Fatalf("expected Edge install path in candidates, got %s", joined)
	}
}

func TestBrowserCandidatePaths_IncludeMacBundleLocations(t *testing.T) {
	paths := browserCandidatePaths("/plugin", Platform{OS: "darwin", Arch: "arm64"})
	joined := strings.Join(paths, "\n")
	if !strings.Contains(joined, `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`) {
		t.Fatalf("expected macOS Chrome path in candidates, got %s", joined)
	}
}

func TestLauncherPath_UsesPlatformBinaryName(t *testing.T) {
	got := LauncherPath("/plugin", Platform{OS: "windows", Arch: "amd64"})
	want := filepath.Join("/plugin", "runtime", "windows-amd64", "bin", "go-motiond.exe")
	if got != want {
		t.Fatalf("LauncherPath() = %q, want %q", got, want)
	}
}

func TestCurrentPlatform_IsSupportedTarget(t *testing.T) {
	current := CurrentPlatform().Key()
	if !strings.Contains(strings.Join(SupportedPlatformKeys(), ","), current) && runtime.GOOS != "freebsd" {
		t.Fatalf("current platform %q should be part of supported matrix during tests", current)
	}
}
