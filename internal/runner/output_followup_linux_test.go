//go:build linux

package runner_test

import (
	"context"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/runner"
)

func TestPreviewRejectsMalformedOutputs(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata/compile/title_render_tiny.json"))
	if err != nil {
		t.Fatal(err)
	}
	art, v, err := compile.CompileJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	wrong := filepath.Join(t.TempDir(), "wrong.png")
	f, err := os.Create(wrong)
	if err != nil {
		t.Fatal(err)
	}
	if err = png.Encode(f, image.NewNRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, script string }{
		{"empty", `: > "$out/preview.png"`},
		{"directory", `mkdir "$out/preview.png"`},
		{"truncated", `printf '\211PNG\r\n\032\n' > "$out/preview.png"`},
		{"wrong dimensions", `cp '` + wrong + `' "$out/preview.png"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fake := fakeExecutable(t, `if [ "$1" = "--version" ]; then echo 'Blender 5.3.0 Beta'; exit 0; fi
while [ "$1" != "--out-dir" ]; do shift; done
shift
out="$1"
printf scene > "$out/scene.blend"
`+tc.script+"\n")
			result, err := runner.RunApply(context.Background(), runner.ApplyOptions{Blender: fake, Parent: t.TempDir(), Artifacts: art, Decl: v.Decl})
			if err == nil {
				t.Fatal("malformed preview marked complete")
			}
			if _, err := os.Stat(filepath.Join(result.JobDir, "out", "manifest.json")); !os.IsNotExist(err) {
				t.Fatal("unexpected success manifest")
			}
			if _, err := os.Stat(filepath.Join(result.JobDir, "out", "manifest.failed.json")); err != nil {
				t.Fatal(err)
			}
		})
	}
}
