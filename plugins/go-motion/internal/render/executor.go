package render

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	gruntime "go-motion/internal/runtime"
)

type CommandInvocation struct {
	Name string
	Args []string
}

type RunResult struct {
	FramesDir   string
	OutputPath  string
	FrameCount  int
	BrowserPath string
	FFmpegPath  string
}

type Executor struct {
	run func(context.Context, CommandInvocation) error
}

func NewExecutor() Executor {
	return Executor{run: runCommand}
}

func (e Executor) Run(ctx context.Context, tools gruntime.Tools, jobDir string, compositionDir string, durationFrames int, fps int, width int, height int) (RunResult, error) {
	if e.run == nil {
		e.run = runCommand
	}

	framesDir := filepath.Join(jobDir, "frames")
	if err := os.MkdirAll(framesDir, 0o755); err != nil {
		return RunResult{}, err
	}

	htmlPath := filepath.Join(compositionDir, "index.html")
	for frame := range durationFrames {
		framePath := filepath.Join(framesDir, fmt.Sprintf("frame-%06d.png", frame))
		args := []string{
			"--headless=new",
			"--no-sandbox",
			"--disable-gpu",
			"--hide-scrollbars",
			"--window-size=" + strconv.Itoa(width) + "," + strconv.Itoa(height),
			"--virtual-time-budget=3000",
			"--run-all-compositor-stages-before-draw",
			"--screenshot=" + framePath,
			fileURL(htmlPath, frame),
		}
		if err := e.run(ctx, CommandInvocation{Name: tools.BrowserPath, Args: args}); err != nil {
			return RunResult{}, fmt.Errorf("capture frame %d: %w", frame, err)
		}
	}

	outputPath := filepath.Join(jobDir, "output.mp4")
	ffmpegArgs := []string{
		"-y",
		"-framerate", strconv.Itoa(fps),
		"-i", filepath.Join(framesDir, "frame-%06d.png"),
		"-pix_fmt", "yuv420p",
		outputPath,
	}
	if err := e.run(ctx, CommandInvocation{Name: tools.FFmpegPath, Args: ffmpegArgs}); err != nil {
		return RunResult{}, fmt.Errorf("encode mp4: %w", err)
	}

	return RunResult{
		FramesDir:   framesDir,
		OutputPath:  outputPath,
		FrameCount:  durationFrames,
		BrowserPath: tools.BrowserPath,
		FFmpegPath:  tools.FFmpegPath,
	}, nil
}

func runCommand(ctx context.Context, inv CommandInvocation) error {
	cmd := exec.CommandContext(ctx, inv.Name, inv.Args...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func fileURL(path string, frame int) string {
	return "file:///" + filepath.ToSlash(path) + "?frame=" + strconv.Itoa(frame)
}
