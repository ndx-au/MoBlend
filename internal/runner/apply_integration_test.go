package runner_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/runner"
)

func TestApplyTitleSmall(t *testing.T) {
	blender := blenderAvailable(t)
	root := filepath.Clean(filepath.Join("..", ".."))
	raw, err := os.ReadFile(filepath.Join(root, "testdata/compile/title_small.json"))
	if err != nil {
		t.Fatal(err)
	}
	art, validated, err := compile.CompileJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	res, err := runner.RunApply(context.Background(), runner.ApplyOptions{
		Blender:     blender,
		Parent:      t.TempDir(),
		Artifacts:   art,
		Decl:        validated.Decl,
		RenderFrame: 24,
		Timeout:     5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("apply: %v\nargs=%v", err, res.Args)
	}
	if st, err := os.Stat(res.PNG); err != nil || st.Size() < 64 {
		t.Fatalf("png missing/small: %v", err)
	}
	if _, err := os.Stat(res.Blend); err != nil {
		t.Fatal(err)
	}

	probe := filepath.Join(res.JobDir, "probe_title.py")
	src := `
import bpy
title = bpy.data.objects.get("title")
cam = bpy.data.objects.get("cam")
if title is None or cam is None:
    raise SystemExit("missing objects")
if title.type != "FONT":
    raise SystemExit(f"type={title.type}")
if "MoBlend" not in title.data.body:
    raise SystemExit(f"body={title.data.body!r}")
mat = title.data.materials[0] if title.data.materials else None
if mat is None:
    raise SystemExit("no material")
print("TITLE_PROBE_OK", list(title.location), list(cam.location))
`
	if err := os.WriteFile(probe, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	reopen, err := runner.Run(context.Background(), runner.Options{
		Blender:   blender,
		BlendFile: res.Blend,
		Python:    probe,
		Timeout:   2 * time.Minute,
	})
	if err != nil {
		t.Fatalf("probe: %v\n%s", err, reopen.Stdout)
	}
}
