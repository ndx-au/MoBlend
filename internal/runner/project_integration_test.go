package runner_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/project"
	"github.com/ndx-au/MoBlend/internal/runner"
)

func TestModularFlatPreviewEquivalence(t *testing.T) {
	blender := blenderAvailable(t)
	root := filepath.Join("..", "..")
	mod, err := project.LoadAndResolve(filepath.Join(root, "examples/projects/title_modular/project.json"))
	if err != nil {
		t.Fatal(err)
	}
	flat, err := project.LoadSceneOrProject(filepath.Join(root, "testdata/project/positive/title_flat.json"), "")
	if err != nil {
		t.Fatal(err)
	}
	am, err := compile.CompileDeclaration(mod.Decl)
	if err != nil {
		t.Fatal(err)
	}
	af, err := compile.CompileDeclaration(flat.Decl)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(am.Data, af.Data) || !bytes.Equal(am.Script, af.Script) {
		t.Fatal("compile artifacts diverge before Blender")
	}
	parent := t.TempDir()
	rm, err := runner.RunApply(context.Background(), runner.ApplyOptions{
		Blender: blender, Parent: parent, Artifacts: am, Decl: mod.Decl,
		Mode: "preview", RenderFrame: 1, Timeout: 3 * time.Minute,
	})
	if err != nil {
		t.Fatalf("modular preview: %v", err)
	}
	rf, err := runner.RunApply(context.Background(), runner.ApplyOptions{
		Blender: blender, Parent: parent, Artifacts: af, Decl: flat.Decl,
		Mode: "preview", RenderFrame: 1, Timeout: 3 * time.Minute,
	})
	if err != nil {
		t.Fatalf("flat preview: %v", err)
	}
	for _, p := range []string{rm.PNG, rf.PNG, rm.Blend, rf.Blend} {
		if _, err := os.Stat(p); err != nil {
			t.Fatal(err)
		}
	}
	// Reopen each independently saved scene and compare actual Blender state.
	probe := filepath.Join(parent, "probe_equivalence.py")
	src := `
import bpy, json, sys
s = bpy.context.scene
s.frame_set(1)
title = bpy.data.objects.get("title__text")
assert title is not None and title.type == "FONT"
assert title.data.body == "MoBlend"
assert abs(title.location.z - 0.5) < 1e-5
assert s.camera is not None and s.camera.name == "cam"
assert list(s.camera.location) == [0.0, -4.0, 1.0]
assert len(title.data.materials) == 1
objects = []
for o in sorted(s.objects, key=lambda o: o.name):
    mats = []
    for mat in getattr(o.data, "materials", []):
        bsdf = mat.node_tree.nodes.get("Principled BSDF")
        mats.append({"color": list(bsdf.inputs["Base Color"].default_value), "roughness": bsdf.inputs["Roughness"].default_value})
    objects.append({"id": o.name, "type": o.type, "location": list(o.location), "rotation": list(o.rotation_euler), "scale": list(o.scale), "materials": mats})
state = {"objects": objects, "camera": s.camera.name, "lens": s.camera.data.lens, "resolution": [s.render.resolution_x, s.render.resolution_y], "frames": [s.frame_start, s.frame_end], "fps": s.render.fps}
with open(sys.argv[sys.argv.index("--")+1], "w") as f:
    json.dump(state, f, sort_keys=True)
`
	if err := os.WriteFile(probe, []byte(src), 0600); err != nil {
		t.Fatal(err)
	}
	var states []any
	for _, result := range []runner.ApplyResult{rm, rf} {
		stateFile := filepath.Join(result.JobDir, "scene_probe.json")
		if _, err := runner.Run(context.Background(), runner.Options{Blender: blender, BlendFile: result.Blend, Python: probe, ExtraArgs: []string{stateFile}, Timeout: 2 * time.Minute}); err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(stateFile)
		if err != nil {
			t.Fatal(err)
		}
		var state any
		if err := json.Unmarshal(raw, &state); err != nil {
			t.Fatal(err)
		}
		states = append(states, state)
	}
	if !reflect.DeepEqual(states[0], states[1]) {
		t.Fatalf("saved Blender states differ: %v / %v", states[0], states[1])
	}
	// Decode and compare pixels under this same CPU Blender pin, allowing a
	// small channel tolerance. Require spatial variation to reject empty renders.
	images := []string{rm.PNG, rf.PNG}
	var pixels [][][4]uint32
	for _, path := range images {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		img, err := png.Decode(f)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		bounds := img.Bounds()
		if bounds.Dx() != 64 || bounds.Dy() != 64 {
			t.Fatalf("unexpected PNG bounds %v", bounds)
		}
		values := [][4]uint32{}
		distinct := map[[4]uint32]bool{}
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				p := [4]uint32{r, g, b, a}
				values = append(values, p)
				distinct[p] = true
			}
		}
		if len(distinct) < 8 {
			t.Fatalf("render lacks content: %d colors", len(distinct))
		}
		pixels = append(pixels, values)
	}
	for i, a := range pixels[0] {
		for c, v := range a {
			delta := int64(v) - int64(pixels[1][i][c])
			if delta < -257 || delta > 257 {
				t.Fatalf("PNG mismatch pixel %d channel %d: %d", i, c, delta)
			}
		}
	}

}

func TestConcurrentShotPreviewsIsolated(t *testing.T) {
	blender := blenderAvailable(t)
	path := filepath.Join("..", "..", "examples/projects/three_shots/project.json")
	parent := t.TempDir()
	shots := []string{"intro", "product", "outro", "product"}
	var wg sync.WaitGroup
	errs := make([]error, len(shots))
	outs := make([]string, len(shots))
	for i, shot := range shots {
		wg.Add(1)
		go func(i int, shot string) {
			defer wg.Done()
			r, err := project.LoadAndResolveShot(path, shot)
			if err != nil {
				errs[i] = err
				return
			}
			art, err := compile.CompileDeclaration(r.Decl)
			if err != nil {
				errs[i] = err
				return
			}

			if i == 3 {
				art.Script = []byte("raise RuntimeError('intentional failed shot')\n")
			}
			res, err := runner.RunApply(context.Background(), runner.ApplyOptions{
				Blender: blender, Parent: parent, Artifacts: art, Decl: r.Decl,
				Mode: "preview", RenderFrame: 1, Timeout: 3 * time.Minute, Resolved: &r,
			})
			errs[i] = err
			outs[i] = res.JobDir
			if res.Shot != shot {
				errs[i] = fmt.Errorf("shot attribution %q want %q", res.Shot, shot)
			}
		}(i, shot)
	}
	wg.Wait()
	for i, err := range errs {
		if i == 3 {
			if err == nil {
				t.Fatal("intentional failed shot succeeded")
			}
			raw, e := os.ReadFile(filepath.Join(outs[i], "out", "manifest.failed.json"))
			if e != nil {
				t.Fatal(e)
			}
			var failed struct{ Status, Shot string }
			if e = json.Unmarshal(raw, &failed); e != nil {
				t.Fatal(e)
			}
			if failed.Status != "failed" || failed.Shot != shots[i] {
				t.Fatalf("failed shot attribution: %s", raw)
			}
			if _, e = os.Stat(filepath.Join(outs[i], "out", "manifest.json")); !os.IsNotExist(e) {
				t.Fatal("failed shot has success manifest")
			}
			continue
		}
		if err != nil {
			t.Fatalf("shot %s: %v", shots[i], err)
		}
		if outs[i] == "" {
			t.Fatalf("missing out for %s", shots[i])
		}
		if _, err := os.Stat(filepath.Join(outs[i], "out", "preview.png")); err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(filepath.Join(outs[i], "out", "dependencies.json"))
		if err != nil {
			t.Fatal(err)
		}
		var deps project.DependencyManifest
		if err = json.Unmarshal(raw, &deps); err != nil {
			t.Fatal(err)
		}
		if deps.Shot != shots[i] {
			t.Fatalf("dependency shot %s want %s", deps.Shot, shots[i])
		}
		raw, err = os.ReadFile(filepath.Join(outs[i], "out", "manifest.json"))
		if err != nil {
			t.Fatal(err)
		}
		var manifest struct {
			Status  string
			Outputs map[string]any
		}
		if err = json.Unmarshal(raw, &manifest); err != nil {
			t.Fatal(err)
		}
		if manifest.Status != "complete" || manifest.Outputs["shot"] != shots[i] {
			t.Fatalf("success shot attribution: %s", raw)
		}
	}
	for i := range outs {
		for j := 0; j < i; j++ {
			if outs[i] == outs[j] {
				t.Fatal("shot jobs must use distinct directories")
			}
		}
	}
}
