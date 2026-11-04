package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/project"
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
	Resolved    *project.Resolved // optional load-time provenance written before success
}

// ApplyResult paths from a successful apply/render.
type ApplyResult struct {
	Shot     string   `json:"shot,omitempty"`
	JobDir   string   `json:"job_dir"`
	Blend    string   `json:"blend"`
	PNG      string   `json:"png"`
	Frames   []string `json:"frames,omitempty"`
	Manifest string   `json:"manifest,omitempty"`
	Args     []string `json:"args,omitempty"`
}

// RunApply extracts pack+font, writes compile artifacts, launches Blender.
func RunApply(ctx context.Context, opt ApplyOptions) (result ApplyResult, runErr error) {
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
	shot := ""
	if opt.Resolved != nil {
		shot = opt.Resolved.Shot
	}
	result.JobDir = jobDir
	result.Shot = shot
	defer func() {
		if runErr != nil {
			_ = os.MkdirAll(outDir, 0o755)
			if _, err := os.Stat(filepath.Join(outDir, "manifest.failed.json")); os.IsNotExist(err) {
				_ = WriteFailedMarker(outDir, "RUN_FAILED", runErr.Error(), shot)
			}
			result.JobDir = jobDir
			result.Shot = shot
		}
	}()
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

	doc, err := Doctor(ctx, blender)
	if err != nil {
		return ApplyResult{}, err
	}
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
		Shot:   shot,
		Blend:  filepath.Join(outDir, "scene.blend"),
		PNG:    filepath.Join(outDir, "preview.png"),
		Args:   res.Args,
	}
	if err != nil {
		_ = WriteFailedMarker(outDir, "RUN_BLENDER", err.Error(), shot)
		return out, err
	}
	if mode == "render" {
		framesDir := filepath.Join(outDir, "frames")
		entries, rdErr := os.ReadDir(framesDir)
		if rdErr != nil {
			_ = WriteFailedMarker(outDir, "RUN_OUTPUT", rdErr.Error(), shot)
			return out, fmt.Errorf("missing frames dir: %w", rdErr)
		}
		expected := opt.Decl.Scene.FrameEnd - opt.Decl.Scene.FrameStart + 1
		if len(entries) != expected {
			return out, fmt.Errorf("incomplete frame set: got %d want %d", len(entries), expected)
		}
		for frame := opt.Decl.Scene.FrameStart; frame <= opt.Decl.Scene.FrameEnd; frame++ {
			path := filepath.Join(framesDir, fmt.Sprintf("frame_%04d.png", frame))
			if err := validatePNG(path, opt.Decl.Scene.Resolution); err != nil {
				return out, fmt.Errorf("invalid frame %d: %w", frame, err)
			}
			out.Frames = append(out.Frames, path)
		}
	} else if err := validatePNG(out.PNG, opt.Decl.Scene.Resolution); err != nil {
		return out, fmt.Errorf("invalid preview png: %w", err)
	}

	if st, err := os.Stat(out.Blend); err != nil || !st.Mode().IsRegular() || st.Size() == 0 {
		return out, fmt.Errorf("missing or invalid scene.blend")
	}
	if opt.Resolved != nil {
		if _, err := project.WriteDependencyManifest(outDir, *opt.Resolved); err != nil {
			return out, err
		}
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
	if opt.Resolved != nil {
		outputs["shot"] = shot
		outputs["dependencies"] = "dependencies.json"
	}
	if err := WriteSuccessManifest(outDir, mode, blender, doc.VersionLine, opt.Decl, outputs); err != nil {
		return out, err
	}
	out.Manifest = filepath.Join(outDir, "manifest.json")
	return out, nil
}
