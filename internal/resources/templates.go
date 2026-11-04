package resources

import (
	"fmt"
	"os"
	"path/filepath"
)

// ApplyScript returns the trusted apply_scene.py template bytes.
func ApplyScript() ([]byte, error) {
	return contentFS.ReadFile("templates/apply_scene.py")
}

// ExtractFont writes DejaVu Sans into destDir/fonts and returns the ttf path.
func ExtractFont(destDir string) (string, error) {
	data, err := contentFS.ReadFile("fonts/DejaVuSans.ttf")
	if err != nil {
		return "", err
	}
	dir := filepath.Join(destDir, "fonts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	out := filepath.Join(dir, "DejaVuSans.ttf")
	if err := os.WriteFile(out, data, 0o644); err != nil {
		return "", err
	}
	if notice, err := contentFS.ReadFile("fonts/DejaVuSans.LICENSE"); err == nil {
		_ = os.WriteFile(filepath.Join(dir, "DejaVuSans.LICENSE"), notice, 0o644)
	}
	return out, nil
}

// MustApplyScript is for tests.
func MustApplyScript() []byte {
	b, err := ApplyScript()
	if err != nil {
		panic(fmt.Sprintf("apply script: %v", err))
	}
	return b
}

// GuideMarkdown returns the embedded agent guide.
func GuideMarkdown() ([]byte, error) {
	return contentFS.ReadFile("agent/GUIDE.md")
}
