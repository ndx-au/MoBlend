package runner_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/runner"
	"github.com/ndx-au/MoBlend/internal/spec"
)

func TestMixedAnimationInterpolationAndOrder(t *testing.T) {
	blender := blenderAvailable(t)
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata/compile/linear_midpoint.json"))
	if err != nil {
		t.Fatal(err)
	}
	_, v, err := compile.CompileJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	decl := v.Decl
	decl.Scene.FrameEnd = 9
	makeTrack := func(prop, interp string, from, to any) spec.Track {
		return spec.Track{TargetID: "cube", Property: prop, Interpolation: interp, Keys: []spec.Key{{Frame: 1, Value: from}, {Frame: 9, Value: to}}}
	}
	tracks := []spec.Track{
		makeTrack("location", "step", []float64{0, 0, 0}, []float64{8, 0, 0}),
		makeTrack("scale", "linear", []float64{1, 1, 1}, []float64{2, 2, 2}),
		makeTrack("rotation_euler_deg", "ease_in_out", []float64{0, 0, 0}, []float64{0, 0, 80}),
		makeTrack("material.roughness", "step", 0.0, 1.0),
		makeTrack("material.base_color", "linear", "#000000", "#FFFFFF"),
		{TargetID: "cam", Property: "location", Interpolation: "linear", Keys: []spec.Key{{Frame: 1, Value: []float64{0, -4, 1}}, {Frame: 9, Value: []float64{0, -8, 1}}}},
	}
	for _, reverse := range []bool{false, true} {
		t.Run(map[bool]string{false: "forward", true: "reversed"}[reverse], func(t *testing.T) {
			decl.Animation = append([]spec.Track(nil), tracks...)
			if reverse {
				for i, j := 0, len(decl.Animation)-1; i < j; i, j = i+1, j-1 {
					decl.Animation[i], decl.Animation[j] = decl.Animation[j], decl.Animation[i]
				}
			}
			art, err := compile.CompileDeclaration(decl)
			if err != nil {
				t.Fatal(err)
			}
			res, err := runner.RunApply(context.Background(), runner.ApplyOptions{Blender: blender, Parent: t.TempDir(), Artifacts: art, Decl: decl, RenderFrame: 3, Timeout: 3 * time.Minute})
			if err != nil {
				t.Fatal(err)
			}
			probe := filepath.Join(res.JobDir, "probe_interpolation.py")
			src := `
import bpy, math
s = bpy.context.scene
obj = bpy.data.objects["cube"]
bsdf = next(n for n in obj.data.materials[0].node_tree.nodes if n.type == "BSDF_PRINCIPLED")
def check(actual, expected, label):
    assert abs(actual - expected) < 0.002, (label, actual, expected)
# Endpoints plus asymmetric interior samples distinguish easing from linear.
for frame, ratio, eased in [(1, 0, 0), (3, .25, .15625), (7, .75, .84375), (9, 1, 1)]:
    s.frame_set(frame)
    held = 1 if frame == 9 else 0
    check(obj.location.x, 8 * held, "location step")
    check(obj.scale.x, 1 + ratio, "scale linear")
    check(obj.rotation_euler.z, math.radians(80) * eased, "rotation ease")
    check(bsdf.inputs["Roughness"].default_value, held, "material step")
    for channel in bsdf.inputs["Base Color"].default_value[:3]:
        check(channel, ratio, "material linear")
    check(s.camera.location.y, -4 - 4 * ratio, "camera linear")
print("MIXED_INTERPOLATION_OK")
`
			if err = os.WriteFile(probe, []byte(src), 0600); err != nil {
				t.Fatal(err)
			}
			result, err := runner.Run(context.Background(), runner.Options{Blender: blender, BlendFile: res.Blend, Python: probe, Timeout: 2 * time.Minute})
			if err != nil {
				t.Fatalf("probe %v: %s", err, result.Stdout)
			}
		})
	}
}
