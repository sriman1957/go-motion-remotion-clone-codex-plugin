package runtime

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestResolveTools_ReportsMissingRuntime(t *testing.T) {
	_, err := ResolveTools("C:/tmp/does-not-exist")
	if err == nil {
		t.Fatal("expected missing runtime error")
	}
	if !strings.Contains(err.Error(), "required runtimes not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBrowserCandidatePaths_IncludeWindowsInstallLocations(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("windows-only path expectation")
	}

	paths := browserCandidatePaths(filepath.Clean("C:/plugin"))
	joined := strings.Join(paths, "\n")
	if !strings.Contains(joined, `C:\Program Files\Google\Chrome\Application\chrome.exe`) {
		t.Fatalf("expected Chrome install path in candidates, got %s", joined)
	}
	if !strings.Contains(joined, `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`) {
		t.Fatalf("expected Edge install path in candidates, got %s", joined)
	}
}
