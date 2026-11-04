package runner_test

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ndx-au/MoBlend/internal/runner"
)

func blenderAvailable(t *testing.T) string {
	t.Helper()
	path, err := runner.ResolveBlender()
	if err != nil {
		t.Skipf("blender unavailable: %v", err)
	}
	return path
}

func TestSmokeRender(t *testing.T) {
	blender := blenderAvailable(t)
	ctx := context.Background()
	res, err := runner.RunSmoke(ctx, runner.SmokeOptions{
		Blender: blender,
		Parent:  t.TempDir(),
		Width:   64,
		Height:  64,
		Scale:   1.5,
		Timeout: 3 * time.Minute,
	})
	if err != nil {
		t.Fatalf("smoke: %v\nargs=%v", err, res.Args)
	}
	if res.ExitCode != 0 {
		t.Fatalf("exit=%d", res.ExitCode)
	}
	st, err := os.Stat(res.PNG)
	if err != nil {
		t.Fatal(err)
	}
	if st.Size() < 64 {
		t.Fatalf("png too small: %d", st.Size())
	}
	w, h := pngDims(t, res.PNG)
	if w != 64 || h != 64 {
		t.Fatalf("dims %dx%d", w, h)
	}
	if _, err := os.Stat(res.Blend); err != nil {
		t.Fatal(err)
	}

	// Reopen saved scene under Blender.
	probe := filepath.Join(res.JobDir, "probe.py")
	probeSrc := `
import bpy, sys
obj = bpy.data.objects.get("MoBlend_BaseCube")
if obj is None:
    raise SystemExit("missing cube")
sx = abs(obj.scale[0] - 1.5)
if sx > 1e-4:
    raise SystemExit(f"scale={tuple(obj.scale)}")
print("REOPEN_OK")
`
	if err := os.WriteFile(probe, []byte(probeSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	reopen, err := runner.Run(ctx, runner.Options{
		Blender:   blender,
		BlendFile: res.Blend,
		Python:    probe,
		Timeout:   2 * time.Minute,
	})
	if err != nil {
		t.Fatalf("reopen: %v\n%s", err, reopen.Stdout)
	}
}

func TestSmokePythonFailure(t *testing.T) {
	blender := blenderAvailable(t)
	_, err := runner.RunSmoke(context.Background(), runner.SmokeOptions{
		Blender: blender,
		Parent:  t.TempDir(),
		Fail:    true,
		Timeout: 2 * time.Minute,
	})
	if err == nil {
		t.Fatal("expected failure from --fail")
	}
}

func TestSmokeCancel(t *testing.T) {
	blender := blenderAvailable(t)
	ctx, cancel := context.WithCancel(context.Background())
	// Cancel almost immediately after start; render should not finish cleanly.
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()
	_, err := runner.RunSmoke(ctx, runner.SmokeOptions{
		Blender: blender,
		Parent:  t.TempDir(),
		Width:   256,
		Height:  256,
		Timeout: 3 * time.Minute,
	})
	if err == nil {
		t.Fatal("expected cancellation error")
	}
}

func pngDims(t *testing.T, path string) (int, int) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 24 || string(data[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Fatalf("not png")
	}
	w := binary.BigEndian.Uint32(data[16:20])
	h := binary.BigEndian.Uint32(data[20:24])
	return int(w), int(h)
}
