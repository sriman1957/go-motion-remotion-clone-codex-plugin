package planner

import (
	"strings"

	"go-motion/internal/spec"
)

func PlanPrompt(prompt string) spec.Composition {
	trimmed := strings.TrimSpace(prompt)
	title := trimmed
	if title == "" {
		title = "Untitled video"
	}

	return spec.Composition{
		Title:          title,
		FPS:            24,
		Width:          1280,
		Height:         720,
		DurationFrames: 96,
		Scenes: []spec.Scene{
			{
				ID:             "scene-1",
				Kind:           "hero",
				Headline:       title,
				Body:           "A polished prompt-driven video composition.",
				DurationFrames: 32,
			},
			{
				ID:             "scene-2",
				Kind:           "feature-grid",
				Headline:       "Key message",
				Body:           "Highlight the strongest benefit with motion-led emphasis.",
				DurationFrames: 32,
			},
			{
				ID:             "scene-3",
				Kind:           "cta",
				Headline:       "Call to action",
				Body:           "Close with a concise brand-forward ending.",
				DurationFrames: 32,
			},
		},
	}
}
