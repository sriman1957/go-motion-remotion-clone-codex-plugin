package generator

import (
	"os"
	"path/filepath"
	"testing"

	"go-motion/internal/spec"
)

func TestGeneratePackage_WritesRuntimeFiles(t *testing.T) {
	dir := t.TempDir()
	comp := spec.Composition{
		Title:          "Demo",
		FPS:            30,
		Width:          1280,
		Height:         720,
		DurationFrames: 90,
	}

	result, err := GeneratePackage(dir, comp)
	if err != nil {
		t.Fatalf("GeneratePackage() error = %v", err)
	}

	for _, name := range []string{"index.html", "styles.css", "runtime.js", "composition.js"} {
		if _, err := os.Stat(filepath.Join(result.Dir, name)); err != nil {
			t.Fatalf("expected %s: %v", name, err)
		}
	}
}
