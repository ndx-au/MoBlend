package runner_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/diag"
	"github.com/ndx-au/MoBlend/internal/runner"
	"github.com/ndx-au/MoBlend/internal/validate"
)

// TestEvalPreparedFixturesPreview is a render smoke suite, not a model authoring evaluation.
func TestEvalPreparedFixturesPreview(t *testing.T) {
	blender := blenderAvailable(t)
	dir := filepath.Join("..", "..", "testdata/eval/supported")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
			files = append(files, e.Name())
		}
	}
	if len(files) < 12 {
		t.Fatalf("need >=12 supported briefs, got %d", len(files))
	}
	pass := 0
	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				t.Fatal(err)
			}
			art, validated, err := compile.CompileJSON(raw)
			if err != nil {
				t.Fatalf("validate/compile: %v", err)
			}
			res, err := runner.RunApply(context.Background(), runner.ApplyOptions{
				Blender:   blender,
				Parent:    t.TempDir(),
				Artifacts: art,
				Decl:      validated.Decl,
				Mode:      "preview",
				Timeout:   5 * time.Minute,
			})
			if err != nil {
				t.Fatalf("preview: %v", err)
			}
			if _, err := os.Stat(res.PNG); err != nil {
				t.Fatal(err)
			}
			pass++
		})
	}
	t.Logf("prepared fixtures previewed: %d/%d", pass, len(files))
}

// TestEvalUnsupportedSchemaFixturesRejected checks JSON rejection, not natural-language refusals.
func TestEvalUnsupportedSchemaFixturesRejected(t *testing.T) {
	dir := filepath.Join("..", "..", "testdata/eval/unsupported")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		n++
		t.Run(e.Name(), func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			_, err = validate.ParseAndValidate(raw)
			if err == nil {
				t.Fatal("expected rejection")
			}
			if _, ok := err.(diag.ErrorList); !ok {
				t.Fatalf("want ErrorList, got %T %v", err, err)
			}
		})
	}
	if n < 4 {
		t.Fatalf("need several unsupported cases, got %d", n)
	}
}
