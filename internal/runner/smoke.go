package runner

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ndx-au/MoBlend/internal/resources"
)

//go:embed testdata/smoke.py
var smokeScript []byte

// SmokeOptions configures the Phase 0 boundary smoke job.
type SmokeOptions struct {
	Blender string
	Parent  string // parent for job dir; empty → temp
	Scale   float64
	Width   int
	Height  int
	Fail    bool
	Timeout time.Duration
}

// SmokeResult is the machine-readable smoke outcome.
type SmokeResult struct {
	JobDir   string          `json:"job_dir"`
	PNG      string          `json:"png"`
	Blend    string          `json:"blend"`
	PackSHA  string          `json:"pack_sha256"`
	Blender  json.RawMessage `json:"blender_info,omitempty"`
	ExitCode int             `json:"exit_code"`
	Args     []string        `json:"args"`
}

// RunSmoke extracts the pack, runs the trusted smoke script, and returns paths.
func RunSmoke(ctx context.Context, opt SmokeOptions) (SmokeResult, error) {
	if opt.Scale == 0 {
		opt.Scale = 1.5
	}
	if opt.Width == 0 {
		opt.Width = 64
	}
	if opt.Height == 0 {
		opt.Height = 64
	}
	if opt.Timeout == 0 {
		opt.Timeout = 3 * time.Minute
	}
	blender := opt.Blender
	if blender == "" {
		var err error
		blender, err = ResolveBlender()
		if err != nil {
			return SmokeResult{}, err
		}
	}

	jobDir, err := JobDir(opt.Parent)
	if err != nil {
		return SmokeResult{}, err
	}
	outDir := filepath.Join(jobDir, "out")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return SmokeResult{}, err
	}

	blendPath, err := resources.ExtractBasePack(jobDir)
	if err != nil {
		return SmokeResult{}, err
	}
	packSHA, err := resources.BasePackSHA256()
	if err != nil {
		return SmokeResult{}, err
	}

	scriptPath := filepath.Join(jobDir, "smoke.py")
	if err := os.WriteFile(scriptPath, smokeScript, 0o644); err != nil {
		return SmokeResult{}, err
	}

	extra := []string{
		"--out-dir", outDir,
		"--scale", fmt.Sprintf("%g", opt.Scale),
		"--width", fmt.Sprintf("%d", opt.Width),
		"--height", fmt.Sprintf("%d", opt.Height),
	}
	if opt.Fail {
		extra = append(extra, "--fail")
	}

	res, err := Run(ctx, Options{
		Blender:   blender,
		BlendFile: blendPath,
		Python:    scriptPath,
		ExtraArgs: extra,
		WorkDir:   jobDir,
		Timeout:   opt.Timeout,
	})
	out := SmokeResult{
		JobDir:   jobDir,
		PackSHA:  packSHA,
		ExitCode: res.ExitCode,
		Args:     res.Args,
		PNG:      filepath.Join(outDir, "smoke.png"),
		Blend:    filepath.Join(outDir, "smoke.blend"),
	}
	if err != nil {
		return out, err
	}

	infoPath := filepath.Join(outDir, "result.json")
	if data, readErr := os.ReadFile(infoPath); readErr == nil {
		out.Blender = json.RawMessage(data)
	}
	return out, nil
}
