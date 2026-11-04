package runner_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
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
}

func TestConcurrentShotPreviewsIsolated(t *testing.T) {
	blender := blenderAvailable(t)
	path := filepath.Join("..", "..", "examples/projects/three_shots/project.json")
	parent := t.TempDir()
	shots := []string{"intro", "product", "outro"}
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
			shotParent := filepath.Join(parent, shot)
			_ = os.MkdirAll(shotParent, 0o755)
			res, err := runner.RunApply(context.Background(), runner.ApplyOptions{
				Blender: blender, Parent: shotParent, Artifacts: art, Decl: r.Decl,
				Mode: "preview", RenderFrame: 1, Timeout: 3 * time.Minute,
			})
			errs[i] = err
			if err == nil {
				outs[i] = res.JobDir
				_, _ = project.WriteDependencyManifest(filepath.Join(res.JobDir, "out"), r)
			}
		}(i, shot)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("shot %s: %v", shots[i], err)
		}
		if outs[i] == "" {
			t.Fatalf("missing out for %s", shots[i])
		}
		if _, err := os.Stat(filepath.Join(outs[i], "out", "preview.png")); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(outs[i], "out", "dependencies.json")); err != nil {
			t.Fatal(err)
		}
	}
	if outs[0] == outs[1] || outs[1] == outs[2] {
		t.Fatal("shot jobs must use distinct directories")
	}
}
