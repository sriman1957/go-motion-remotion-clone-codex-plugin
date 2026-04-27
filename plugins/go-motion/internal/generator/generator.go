package generator

import (
	"encoding/json"
	"os"
	"path/filepath"

	"go-motion/internal/spec"
	"go-motion/internal/templates"
)

type Package struct {
	Dir string
}

func GeneratePackage(baseDir string, comp spec.Composition) (Package, error) {
	jobDir := filepath.Join(baseDir, "composition")
	if err := os.MkdirAll(jobDir, 0o755); err != nil {
		return Package{}, err
	}

	files := map[string][]byte{
		"index.html": []byte(templates.HTML),
		"styles.css": []byte(templates.CSS),
		"runtime.js": []byte(templates.JS),
	}

	for name, contents := range files {
		if err := os.WriteFile(filepath.Join(jobDir, name), contents, 0o644); err != nil {
			return Package{}, err
		}
	}

	data, err := json.MarshalIndent(comp, "", "  ")
	if err != nil {
		return Package{}, err
	}
	js := []byte("window.__GO_MOTION_COMPOSITION__ = " + string(data) + ";\n")
	if err := os.WriteFile(filepath.Join(jobDir, "composition.js"), js, 0o644); err != nil {
		return Package{}, err
	}

	return Package{Dir: jobDir}, nil
}
