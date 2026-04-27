package render

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	gruntime "go-motion/internal/runtime"
)

func TestExecutor_Run_InvokesBrowserAndFFmpeg(t *testing.T) {
	jobDir := t.TempDir()
	compositionDir := filepath.Join(jobDir, "composition")
	if err := os.MkdirAll(compositionDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	var calls []CommandInvocation
	exec := Executor{
		run: func(_ context.Context, inv CommandInvocation) error {
			calls = append(calls, inv)
			if inv.Name == "browser" {
				frameFile := ""
				for _, arg := range inv.Args {
					if len(arg) > len("--screenshot=") && arg[:13] == "--screenshot=" {
						frameFile = arg[len("--screenshot="):]
						break
					}
				}
				if frameFile == "" {
					t.Fatal("expected screenshot argument")
				}
				if err := os.WriteFile(frameFile, []byte("png"), 0o644); err != nil {
					return err
				}
			}
			if inv.Name == "ffmpeg" {
				out := inv.Args[len(inv.Args)-1]
				if err := os.WriteFile(out, []byte("mp4"), 0o644); err != nil {
					return err
				}
			}
			return nil
		},
	}

	result, err := exec.Run(context.Background(), gruntime.Tools{
		BrowserPath: "browser",
		FFmpegPath:  "ffmpeg",
	}, jobDir, compositionDir, 2, 24, 1280, 720)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(calls) != 3 {
		t.Fatalf("expected 3 process invocations, got %d", len(calls))
	}
	if _, err := os.Stat(result.OutputPath); err != nil {
		t.Fatalf("expected output file: %v", err)
	}
	if calls[2].Args[2] != "24" {
		t.Fatalf("expected ffmpeg framerate arg 24, got %q", calls[2].Args[2])
	}
}
