package runner_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/runner"
)

func TestDoctorOK(t *testing.T) {
	blender := blenderAvailable(t)
	rep, err := runner.Doctor(context.Background(), blender)
	if err != nil || !rep.OK {
		t.Fatalf("doctor: %+v err=%v", rep, err)
	}
	if rep.Family != runner.SupportedBlenderFamily {
		t.Fatalf("family %s", rep.Family)
	}
}

func TestDoctorMissing(t *testing.T) {
	_, err := runner.Doctor(context.Background(), "/nonexistent/blender-binary")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRenderWritesManifestAndFrames(t *testing.T) {
	blender := blenderAvailable(t)
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata/compile/title_render_tiny.json"))
	if err != nil {
		t.Fatal(err)
	}
	art, validated, err := compile.CompileJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	res, err := runner.RunApply(context.Background(), runner.ApplyOptions{
		Blender:   blender,
		Parent:    t.TempDir(),
		Artifacts: art,
		Decl:      validated.Decl,
		Mode:      "render",
		Timeout:   5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if len(res.Frames) != 2 {
		t.Fatalf("frames=%d want 2: %v", len(res.Frames), res.Frames)
	}
	data, err := os.ReadFile(res.Manifest)
	if err != nil {
		t.Fatal(err)
	}
	var man map[string]any
	if err := json.Unmarshal(data, &man); err != nil {
		t.Fatal(err)
	}
	if man["status"] != "complete" {
		t.Fatalf("manifest=%v", man)
	}
	// Failed marker must not exist alongside success.
	if _, err := os.Stat(filepath.Join(res.JobDir, "out", "manifest.failed.json")); err == nil {
		t.Fatal("unexpected failed manifest")
	}
}
