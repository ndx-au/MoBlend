package runner_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/runner"
)

func TestConcurrentJobsIsolated(t *testing.T) {
	blender := blenderAvailable(t)
	root := filepath.Join("..", "..")
	files := []string{
		"testdata/compile/title_small.json",
		"testdata/compile/pedestal_small.json",
		"testdata/compile/chart_small.json",
	}
	parent := t.TempDir()
	type result struct {
		name string
		res  runner.ApplyResult
		err  error
	}
	out := make([]result, len(files))
	var wg sync.WaitGroup
	for i, rel := range files {
		wg.Add(1)
		go func(i int, rel string) {
			defer wg.Done()
			raw, err := os.ReadFile(filepath.Join(root, rel))
			if err != nil {
				out[i] = result{name: rel, err: err}
				return
			}
			art, validated, err := compile.CompileJSON(raw)
			if err != nil {
				out[i] = result{name: rel, err: err}
				return
			}
			res, err := runner.RunApply(context.Background(), runner.ApplyOptions{
				Blender:   blender,
				Parent:    parent,
				Artifacts: art,
				Decl:      validated.Decl,
				Mode:      "preview",
				Timeout:   5 * time.Minute,
			})
			out[i] = result{name: rel, res: res, err: err}
		}(i, rel)
	}
	wg.Wait()

	seenJob := map[string]string{}
	seenPNG := map[string]string{}
	for _, r := range out {
		if r.err != nil {
			t.Fatalf("%s: %v", r.name, r.err)
		}
		if r.res.JobDir == "" {
			t.Fatalf("%s: empty job dir", r.name)
		}
		if prev, ok := seenJob[r.res.JobDir]; ok {
			t.Fatalf("job dir collision %s used by %s and %s", r.res.JobDir, prev, r.name)
		}
		seenJob[r.res.JobDir] = r.name
		if prev, ok := seenPNG[r.res.PNG]; ok {
			t.Fatalf("png collision %s used by %s and %s", r.res.PNG, prev, r.name)
		}
		seenPNG[r.res.PNG] = r.name
		st, err := os.Stat(r.res.PNG)
		if err != nil || st.Size() < 64 {
			t.Fatalf("%s: bad png: %v", r.name, err)
		}
		if _, err := os.Stat(r.res.Blend); err != nil {
			t.Fatalf("%s: missing blend: %v", r.name, err)
		}
		if _, err := os.Stat(r.res.Manifest); err != nil {
			t.Fatalf("%s: missing manifest: %v", r.name, err)
		}
	}
}
