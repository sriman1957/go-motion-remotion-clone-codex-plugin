package render

import (
	"path/filepath"
)

type Command struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type Plan struct {
	Commands []Command `json:"commands"`
}

func BuildRenderPlan(pluginRoot string, jobDir string) Plan {
	compositionDir := filepath.Join(jobDir, "composition")
	framesDir := filepath.Join(jobDir, "frames")
	outputPath := filepath.Join(jobDir, "output.mp4")

	return Plan{
		Commands: []Command{
			{
				Name: "capture-frames",
				Args: []string{
					"--composition-dir", compositionDir,
					"--frames-dir", framesDir,
					"--plugin-root", pluginRoot,
				},
			},
			{
				Name: "encode-mp4",
				Args: []string{
					"-i", filepath.Join(framesDir, "frame-%06d.png"),
					"-pix_fmt", "yuv420p",
					outputPath,
				},
			},
		},
	}
}
