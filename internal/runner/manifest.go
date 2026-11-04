package runner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/ndx-au/MoBlend/internal/resources"
	"github.com/ndx-au/MoBlend/internal/spec"
)

// Manifest is written only for successful jobs.
type Manifest struct {
	Status           string         `json:"status"`
	SchemaVersion    string         `json:"schema_version"`
	AssetPackVersion string         `json:"asset_pack_version"`
	PackSHA256       string         `json:"pack_sha256"`
	Mode             string         `json:"mode"`
	Blender          map[string]any `json:"blender"`
	Settings         map[string]any `json:"settings"`
	Outputs          map[string]any `json:"outputs"`
	CompletedAt      string         `json:"completed_at"`
}

// WriteSuccessManifest writes out/manifest.json for a completed job.
func WriteSuccessManifest(outDir string, mode string, blenderPath string, versionLine string, decl spec.Declaration, outputs map[string]any) error {
	sha, err := resources.BasePackSHA256()
	if err != nil {
		return err
	}
	m := Manifest{
		Status:           "complete",
		SchemaVersion:    decl.SchemaVersion,
		AssetPackVersion: decl.AssetPackVersion,
		PackSHA256:       sha,
		Mode:             mode,
		Blender: map[string]any{
			"path":    blenderPath,
			"version": versionLine,
		},
		Settings: map[string]any{
			"frame_start": decl.Scene.FrameStart,
			"frame_end":   decl.Scene.FrameEnd,
			"fps":         decl.Scene.FPS,
			"resolution":  decl.Scene.Resolution,
			"samples":     decl.Scene.Samples,
			"engine":      "CYCLES",
			"device":      "CPU",
		},
		Outputs:     outputs,
		CompletedAt: time.Now().UTC().Format(time.RFC3339),
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(outDir, "manifest.json"), data, 0o644)
}

// WriteFailedMarker writes manifest.failed.json without status complete.
func WriteFailedMarker(outDir string, code, message string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	payload := map[string]any{
		"status":  "failed",
		"code":    code,
		"message": message,
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(outDir, "manifest.failed.json"), data, 0o644)
}
