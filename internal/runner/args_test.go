package runner_test

import (
	"testing"

	"github.com/ndx-au/MoBlend/internal/runner"
)

func TestBuildArgsOrder(t *testing.T) {
	args := runner.BuildArgs(runner.Options{
		BlendFile: "/tmp/base.blend",
		Python:    "/tmp/smoke.py",
		ExtraArgs: []string{"--out-dir", "/tmp/out", "--scale", "1.5"},
	})
	want := []string{
		"--background",
		"--factory-startup",
		"--disable-autoexec",
		"--python-exit-code",
		"1",
		"/tmp/base.blend",
		"--python",
		"/tmp/smoke.py",
		"--",
		"--out-dir",
		"/tmp/out",
		"--scale",
		"1.5",
	}
	if len(args) != len(want) {
		t.Fatalf("len=%d want %d\n%v", len(args), len(want), args)
	}
	for i := range want {
		if args[i] != want[i] {
			t.Fatalf("arg[%d]=%q want %q\nfull=%v", i, args[i], want[i], args)
		}
	}
}
