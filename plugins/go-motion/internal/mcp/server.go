package mcp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"

	"go-motion/internal/generator"
	"go-motion/internal/planner"
	"go-motion/internal/render"
	gruntime "go-motion/internal/runtime"
)

type Server struct {
	pluginRoot   string
	workspace    string
	resolveTools func(string) (gruntime.Tools, error)
	renderer     func(context.Context, gruntime.Tools, string, string, int, int, int, int) (RenderResult, error)
}

type GenerateVideoRequest struct {
	Prompt string `json:"prompt"`
}

type GenerateVideoResult struct {
	JobID          string      `json:"jobId"`
	Status         string      `json:"status"`
	Prompt         string      `json:"prompt"`
	CompositionDir string      `json:"compositionDir"`
	RenderPlan     render.Plan `json:"renderPlan"`
	OutputPath     string      `json:"outputPath,omitempty"`
	FrameCount     int         `json:"frameCount,omitempty"`
	RuntimeError   string      `json:"runtimeError,omitempty"`
}

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

type RenderResult struct {
	OutputPath string
	FrameCount int
}

func NewServer(pluginRoot string, workspace string) *Server {
	return &Server{
		pluginRoot:   pluginRoot,
		workspace:    workspace,
		resolveTools: gruntime.ResolveTools,
		renderer: func(ctx context.Context, tools gruntime.Tools, jobDir string, compositionDir string, durationFrames int, fps int, width int, height int) (RenderResult, error) {
			exec := render.NewExecutor()
			result, err := exec.Run(ctx, tools, jobDir, compositionDir, durationFrames, fps, width, height)
			if err != nil {
				return RenderResult{}, err
			}
			return RenderResult{
				OutputPath: result.OutputPath,
				FrameCount: result.FrameCount,
			}, nil
		},
	}
}

func (s *Server) ListTools() []Tool {
	return []Tool{
		{
			Name:        "generate_video",
			Description: "Generate a prompt-driven video composition and render plan.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"prompt": map[string]any{
						"type":        "string",
						"description": "The video request in natural language.",
					},
				},
				"required": []string{"prompt"},
			},
		},
		{
			Name:        "list_styles",
			Description: "List built-in visual directions for prompt planning.",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
	}
}

func (s *Server) GenerateVideo(req GenerateVideoRequest) (GenerateVideoResult, error) {
	jobID, err := randomID()
	if err != nil {
		return GenerateVideoResult{}, err
	}

	jobDir := filepath.Join(s.workspace, "jobs", jobID)
	if err := os.MkdirAll(jobDir, 0o755); err != nil {
		return GenerateVideoResult{}, err
	}

	comp := planner.PlanPrompt(req.Prompt)
	pkg, err := generator.GeneratePackage(jobDir, comp)
	if err != nil {
		return GenerateVideoResult{}, err
	}

	plan := render.BuildRenderPlan(s.pluginRoot, jobDir)
	result := GenerateVideoResult{
		JobID:          jobID,
		Status:         "planned",
		Prompt:         req.Prompt,
		CompositionDir: pkg.Dir,
		RenderPlan:     plan,
	}

	tools, err := s.resolveTools(s.pluginRoot)
	if err != nil {
		result.RuntimeError = err.Error()
		result.Status = "runtime-missing"
		return result, nil
	}

	renderResult, err := s.renderer(context.Background(), tools, jobDir, pkg.Dir, comp.DurationFrames, comp.FPS, comp.Width, comp.Height)
	if err != nil {
		result.RuntimeError = err.Error()
		result.Status = "render-failed"
		return result, nil
	}

	result.Status = "rendered"
	result.OutputPath = renderResult.OutputPath
	result.FrameCount = renderResult.FrameCount

	return result, nil
}

func randomID() (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "job-" + hex.EncodeToString(buf), nil
}
