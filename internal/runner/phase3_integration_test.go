package runner_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/runner"
)

func TestApplyChartIsolation(t *testing.T) {
	blender := blenderAvailable(t)
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata/compile/chart_small.json"))
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
		Timeout:   5 * time.Minute,
	})
	if err != nil {
		t.Fatalf("apply: %v\n%s", err, res.Args)
	}
	probe := filepath.Join(res.JobDir, "probe_chart.py")
	src := `
import bpy
a = bpy.data.objects.get("chartA")
b = bpy.data.objects.get("chartB")
if a is None or b is None:
    raise SystemExit("missing charts")
va = list(a.get("moblend_bar_values", []))
vb = list(b.get("moblend_bar_values", []))
if va == vb:
    raise SystemExit(f"values leaked {va} {vb}")
if abs(va[1] - 0.9) > 1e-5:
    raise SystemExit(f"chartA values {va}")
if abs(vb[0] - 0.8) > 1e-5:
    raise SystemExit(f"chartB values {vb}")
# Node groups must be distinct copies.
ma = next(m for m in a.modifiers if m.type == "NODES")
mb = next(m for m in b.modifiers if m.type == "NODES")
if ma.node_group == mb.node_group:
    raise SystemExit("shared node group")
if ma.node_group.name == mb.node_group.name:
    raise SystemExit("same node group name")
print("CHART_ISOLATION_OK", va, vb, ma.node_group.name, mb.node_group.name)
`
	if err := os.WriteFile(probe, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runner.Run(context.Background(), runner.Options{
		Blender:   blender,
		BlendFile: res.Blend,
		Python:    probe,
		Timeout:   2 * time.Minute,
	})
	if err != nil {
		t.Fatalf("probe: %v\n%s", err, out.Stdout)
	}
}

func TestLinearLocationMidpoint(t *testing.T) {
	blender := blenderAvailable(t)
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata/compile/linear_midpoint.json"))
	if err != nil {
		t.Fatal(err)
	}
	art, validated, err := compile.CompileJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	res, err := runner.RunApply(context.Background(), runner.ApplyOptions{
		Blender: blender, Parent: t.TempDir(), Artifacts: art, Decl: validated.Decl,
		RenderFrame: 11, Timeout: 3 * time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}
	probe := filepath.Join(res.JobDir, "probe_linear.py")
	src := `
import bpy
obj = bpy.data.objects.get("cube")
if obj is None:
    raise SystemExit("missing cube")
bpy.context.scene.frame_set(1)
x1,y1,z1 = obj.location
bpy.context.scene.frame_set(11)
x11,y11,z11 = obj.location
bpy.context.scene.frame_set(21)
x21,y21,z21 = obj.location
tol = 1e-3
if abs(x1) > tol or abs(y1) > tol or abs(z1) > tol:
    raise SystemExit(f"f1={obj.location}")
if abs(x11-1.0) > 0.05 or abs(y11) > tol or abs(z11) > tol:
    raise SystemExit(f"f11=({x11},{y11},{z11})")
if abs(x21-2.0) > tol or abs(y21) > tol or abs(z21) > tol:
    raise SystemExit(f"f21=({x21},{y21},{z21})")
print("LINEAR_MIDPOINT_OK", x1, x11, x21)
`
	if err := os.WriteFile(probe, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runner.Run(context.Background(), runner.Options{
		Blender: blender, BlendFile: res.Blend, Python: probe, Timeout: 2 * time.Minute,
	})
	if err != nil {
		t.Fatalf("probe: %v\n%s", err, out.Stdout)
	}
	if !strings.Contains(out.Stdout, "LINEAR_MIDPOINT_OK") {
		t.Fatalf("missing ok: %s", out.Stdout)
	}
}

func TestStepInterpolationHold(t *testing.T) {
	blender := blenderAvailable(t)
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata/compile/step_jump.json"))
	if err != nil {
		t.Fatal(err)
	}
	art, validated, err := compile.CompileJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	res, err := runner.RunApply(context.Background(), runner.ApplyOptions{
		Blender: blender, Parent: t.TempDir(), Artifacts: art, Decl: validated.Decl,
		RenderFrame: 20, Timeout: 3 * time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}
	probe := filepath.Join(res.JobDir, "probe_step.py")
	src := `
import bpy
obj = bpy.data.objects.get("cube")
if obj is None:
    raise SystemExit("missing cube")
bpy.context.scene.frame_set(20)
x20 = obj.location.x
bpy.context.scene.frame_set(21)
x21 = obj.location.x
if abs(x20) > 1e-3:
    raise SystemExit(f"frame20 x={x20} want ~0")
if abs(x21-2.0) > 1e-3:
    raise SystemExit(f"frame21 x={x21} want ~2")
print("STEP_HOLD_OK", x20, x21)
`
	if err := os.WriteFile(probe, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runner.Run(context.Background(), runner.Options{
		Blender: blender, BlendFile: res.Blend, Python: probe, Timeout: 2 * time.Minute,
	})
	if err != nil {
		t.Fatalf("probe: %v\n%s", err, out.Stdout)
	}
	if !strings.Contains(out.Stdout, "STEP_HOLD_OK") {
		t.Fatalf("missing ok: %s", out.Stdout)
	}
}

func TestApplyPedestalMotionProbe(t *testing.T) {
	blender := blenderAvailable(t)
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata/compile/pedestal_small.json"))
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
		t.Fatalf("apply: %v", err)
	}
	probe := filepath.Join(res.JobDir, "probe_pedestal.py")
	src := `
import bpy, math
obj = bpy.data.objects.get("product")
if obj is None:
    raise SystemExit("missing product")
bpy.context.scene.frame_set(1)
z1 = obj.rotation_euler.z
bpy.context.scene.frame_set(24)
z24 = obj.rotation_euler.z
# Expect ~90 degrees at frame 24.
if abs(z24 - math.radians(90)) > math.radians(2):
    raise SystemExit(f"z24={z24} z1={z1}")
print("PEDESTAL_MOTION_OK", z1, z24)
`
	if err := os.WriteFile(probe, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runner.Run(context.Background(), runner.Options{
		Blender:   blender,
		BlendFile: res.Blend,
		Python:    probe,
		Timeout:   2 * time.Minute,
	})
	if err != nil {
		t.Fatalf("probe: %v\n%s", err, out.Stdout)
	}
	if !strings.Contains(out.Stdout, "PEDESTAL_MOTION_OK") {
		t.Fatalf("missing ok marker: %s", out.Stdout)
	}
}
