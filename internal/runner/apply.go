package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/resources"
	"github.com/ndx-au/MoBlend/internal/spec"
)

// ApplyOptions runs a compiled declaration against a pristine pack extract.
type ApplyOptions struct {
	Blender     string
	Parent      string
	Artifacts   compile.Artifacts
	Decl        spec.Declaration
	Mode        string // preview (default) or render
	RenderFrame int
	Timeout     time.Duration
}

// ApplyResult paths from a successful apply/render.
type ApplyResult struct {
	JobDir   string   `json:"job_dir"`
	Blend    string   `json:"blend"`
	PNG      string   `json:"png"`
	Frames   []string `json:"frames,omitempty"`
	Manifest string   `json:"manifest,omitempty"`
	Args     []string `json:"args,omitempty"`
}

// RunApply extracts pack+font, writes compile artifacts, launches Blender.
func RunApply(ctx context.Context, opt ApplyOptions) (ApplyResult, error) {
	if opt.Timeout == 0 {
		opt.Timeout = 5 * time.Minute
	}
	blender := opt.Blender
	if blender == "" {
		var err error
		blender, err = ResolveBlender()
		if err != nil {
			return ApplyResult{}, err
		}
	}
	jobDir, err := JobDir(opt.Parent)
	if err != nil {
		return ApplyResult{}, err
	}
	outDir := filepath.Join(jobDir, "out")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return ApplyResult{}, err
	}
	packPath, err := resources.ExtractBasePack(jobDir)
	if err != nil {
		return ApplyResult{}, err
	}
	fontPath, err := resources.ExtractFont(jobDir)
	if err != nil {
		return ApplyResult{}, err
	}
	if err := compile.WriteToDir(outDir, opt.Artifacts); err != nil {
		return ApplyResult{}, err
	}
	script := filepath.Join(outDir, opt.Artifacts.ScriptName)
	data := filepath.Join(outDir, opt.Artifacts.DataName)

	mode := opt.Mode
	if mode == "" {
		mode = "preview"
	}
	extra := []string{
		"--out-dir", outDir,
		"--data", data,
		"--font", fontPath,
		"--pack", packPath,
		"--mode", mode,
	}
	if opt.RenderFrame != 0 {
		extra = append(extra, "--render-frame", fmt.Sprintf("%d", opt.RenderFrame))
	}

	doc, _ := Doctor(ctx, blender)
	res, err := Run(ctx, Options{
		Blender: blender,
		Python:  script,
		// No BlendFile: factory empty scene; script loads pack assets explicitly.
		ExtraArgs: extra,
		WorkDir:   jobDir,
		Timeout:   opt.Timeout,
	})
	out := ApplyResult{
		JobDir: jobDir,
		Blend:  filepath.Join(outDir, "scene.blend"),
		PNG:    filepath.Join(outDir, "preview.png"),
		Args:   res.Args,
	}
	if err != nil {
		_ = WriteFailedMarker(outDir, "RUN_BLENDER", err.Error())
		return out, err
	}
	if mode == "render" {
		framesDir := filepath.Join(outDir, "frames")
		entries, rdErr := os.ReadDir(framesDir)
		if rdErr != nil {
			_ = WriteFailedMarker(outDir, "RUN_OUTPUT", rdErr.Error())
			return out, fmt.Errorf("missing frames dir: %w", rdErr)
		}
		for _, e := range entries {
			if !e.IsDir() {
				out.Frames = append(out.Frames, filepath.Join(framesDir, e.Name()))
			}
		}
		if len(out.Frames) == 0 {
			_ = WriteFailedMarker(outDir, "RUN_OUTPUT", "no frames written")
			return out, fmt.Errorf("no frames written")
		}
	} else if _, statErr := os.Stat(out.PNG); statErr != nil {
		_ = WriteFailedMarker(outDir, "RUN_OUTPUT", statErr.Error())
		return out, fmt.Errorf("missing preview png: %w", statErr)
	}

	outputs := map[string]any{
		"blend": "scene.blend",
		"mode":  mode,
	}
	if mode == "render" {
		outputs["frame_count"] = len(out.Frames)
		outputs["frames_dir"] = "frames"
	} else {
		outputs["preview"] = "preview.png"
		outputs["frame"] = opt.RenderFrame
	}
	if err := WriteSuccessManifest(outDir, mode, blender, doc.VersionLine, opt.Decl, outputs); err != nil {
		return out, err
	}
	out.Manifest = filepath.Join(outDir, "manifest.json")
	return out, nil
}
