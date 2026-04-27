package mcp

import (
	"context"
	"testing"

	"go-motion/internal/runtime"
)

func TestGenerateVideo_CreatesJobMetadata(t *testing.T) {
	srv := NewServer("C:/plugin", t.TempDir())
	result, err := srv.GenerateVideo(GenerateVideoRequest{Prompt: "Create a launch video"})
	if err != nil {
		t.Fatalf("GenerateVideo() error = %v", err)
	}
	if result.JobID == "" {
		t.Fatal("expected JobID")
	}
	if result.CompositionDir == "" {
		t.Fatal("expected CompositionDir")
	}
	if result.Status == "" {
		t.Fatal("expected Status")
	}
}

func TestGenerateVideo_RendersWhenRuntimeAndRendererAvailable(t *testing.T) {
	srv := NewServer("C:/plugin", t.TempDir())
	srv.resolveTools = func(string) (runtime.Tools, error) {
		return runtime.Tools{
			BrowserPath: "browser",
			FFmpegPath:  "ffmpeg",
		}, nil
	}
	srv.renderer = func(_ context.Context, tools runtime.Tools, jobDir string, compositionDir string, durationFrames int, fps int, width int, height int) (RenderResult, error) {
		return RenderResult{
			OutputPath: jobDir + "/output.mp4",
			FrameCount: durationFrames,
		}, nil
	}

	result, err := srv.GenerateVideo(GenerateVideoRequest{Prompt: "Create a launch video"})
	if err != nil {
		t.Fatalf("GenerateVideo() error = %v", err)
	}
	if result.Status != "rendered" {
		t.Fatalf("Status = %q, want %q", result.Status, "rendered")
	}
	if result.OutputPath == "" {
		t.Fatal("expected OutputPath")
	}
}
