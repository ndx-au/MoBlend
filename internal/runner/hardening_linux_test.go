//go:build linux

package runner_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/runner"
)

func fakeExecutable(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "blender")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0700); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunTerminatesDescendants(t *testing.T) {
	for _, timeout := range []bool{false, true} {
		t.Run(fmt.Sprintf("timeout=%v", timeout), func(t *testing.T) {
			dir := t.TempDir()
			fake := fakeExecutable(t, `sleep 60 &
child=$!
printf '%s %s' "$$" "$child" > pids
wait
`)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			duration := time.Duration(0)
			if timeout {
				duration = time.Second
			}
			done := make(chan error, 1)
			go func() {
				_, err := runner.Run(ctx, runner.Options{Blender: fake, WorkDir: dir, Timeout: duration})
				done <- err
			}()
			var raw []byte
			deadline := time.Now().Add(3 * time.Second)
			for time.Now().Before(deadline) {
				raw, _ = os.ReadFile(filepath.Join(dir, "pids"))
				if len(strings.Fields(string(raw))) == 2 {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			if len(strings.Fields(string(raw))) != 2 {
				t.Fatal("fake process never became ready")
			}
			if !timeout {
				cancel()
			}
			select {
			case err := <-done:
				want := context.Canceled
				if timeout {
					want = context.DeadlineExceeded
				}
				if !errors.Is(err, want) {
					t.Fatalf("want %v, got %v", want, err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("cancellation hung")
			}
			for _, pid := range strings.Fields(string(raw)) {
				// A reparented zombie is dead; Linux PID 1 may defer reaping it.
				alive := func() bool {
					stat, err := os.ReadFile("/proc/" + pid + "/stat")
					if os.IsNotExist(err) || errors.Is(err, syscall.ESRCH) {
						return false
					}
					if err != nil {
						t.Fatal(err)
					}
					fields := strings.Fields(string(stat)[strings.LastIndex(string(stat), ")")+1:])
					return fields[0] != "Z"
				}
				until := time.Now().Add(time.Second)
				for alive() && time.Now().Before(until) {
					time.Sleep(10 * time.Millisecond)
				}
				if alive() {
					t.Errorf("descendant %s still running", pid)
				}
			}
		})
	}
}

func TestRunEnvironmentIsolated(t *testing.T) {
	t.Setenv("MOBLEND_TEST_SECRET", "secret")
	fake := fakeExecutable(t, `test -z "$MOBLEND_TEST_SECRET" || exit 9
test -n "$HOME" && test "$HOME" = "$XDG_CACHE_HOME" || exit 8
printf '%s' "$HOME"
`)
	dir := t.TempDir()
	result, err := runner.Run(context.Background(), runner.Options{Blender: fake, WorkDir: dir})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(result.Stdout, dir+"/") {
		t.Fatalf("home not private: %s", result.Stdout)
	}
}

func TestApplyRejectsIncompleteFrameSets(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata/compile/title_render_tiny.json"))
	if err != nil {
		t.Fatal(err)
	}
	art, v, err := compile.CompileJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	for _, names := range [][]string{{"frame_0001.png"}, {"frame_0001.png", "wrong.png"}, {"frame_0001.png", "frame_0002.png", "extra.png"}} {
		t.Run(strings.Join(names, ","), func(t *testing.T) {
			body := `if [ "$1" = "--version" ]; then echo 'Blender 5.3.0 Beta'; exit 0; fi
while [ "$1" != "--out-dir" ]; do shift; done
shift
out="$1"
mkdir -p "$out/frames"
printf scene > "$out/scene.blend"
`
			for _, name := range names {
				body += "printf png > \"$out/frames/" + name + "\"\n"
			}
			fake := fakeExecutable(t, body)
			result, err := runner.RunApply(context.Background(), runner.ApplyOptions{Blender: fake, Parent: t.TempDir(), Artifacts: art, Decl: v.Decl, Mode: "render"})
			if err == nil {
				t.Fatal("partial/misnamed frame set accepted")
			}
			if _, err := os.Stat(filepath.Join(result.JobDir, "out", "manifest.json")); !os.IsNotExist(err) {
				t.Fatal("success manifest on incomplete output")
			}
			if _, err := os.Stat(filepath.Join(result.JobDir, "out", "manifest.failed.json")); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestApplyTimeoutHasFailedMarker(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata/compile/title_render_tiny.json"))
	if err != nil {
		t.Fatal(err)
	}
	art, v, err := compile.CompileJSON(raw)
	if err != nil {
		t.Fatal(err)
	}

	fake := fakeExecutable(t, `if [ "$1" = "--version" ]; then echo 'Blender 5.3.0 Beta'; exit 0; fi
sleep 60
`)
	result, err := runner.RunApply(context.Background(), runner.ApplyOptions{Blender: fake, Parent: t.TempDir(), Artifacts: art, Decl: v.Decl, Timeout: 50 * time.Millisecond})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got %v", err)
	}
	if _, err := os.Stat(filepath.Join(result.JobDir, "out", "manifest.failed.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(result.JobDir, "out", "manifest.json")); !os.IsNotExist(err) {
		t.Fatal("unexpected success")
	}
}

func TestDoctorTerminatesProbeDescendants(t *testing.T) {
	for _, timeout := range []bool{false, true} {
		t.Run(fmt.Sprintf("timeout=%v", timeout), func(t *testing.T) {
			dir := t.TempDir()
			pids := filepath.Join(dir, "pids")
			fake := fakeExecutable(t, "sleep 60 &\nprintf '%s %s' \"$$\" \"$!\" > '"+pids+"'\nwait\n")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if timeout {
				var stop context.CancelFunc
				ctx, stop = context.WithTimeout(ctx, time.Second)
				defer stop()
			}
			done := make(chan error, 1)
			go func() { _, err := runner.Doctor(ctx, fake); done <- err }()
			var raw []byte
			until := time.Now().Add(3 * time.Second)
			for time.Now().Before(until) {
				raw, _ = os.ReadFile(pids)
				if len(strings.Fields(string(raw))) == 2 {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			if len(strings.Fields(string(raw))) != 2 {
				t.Fatal("probe not ready")
			}
			if !timeout {
				cancel()
			}
			select {
			case err := <-done:
				want := context.Canceled
				if timeout {
					want = context.DeadlineExceeded
				}
				if !errors.Is(err, want) {
					t.Fatalf("want %v got %v", want, err)
				}
			case <-time.After(4 * time.Second):
				t.Fatal("probe cancellation hung")
			}
			for _, pid := range strings.Fields(string(raw)) {
				// Killing the process group schedules each member for termination; the
				// parent returning does not make every child terminal instantaneously.
				deadline := time.Now().Add(time.Second)
				for {
					stat, err := os.ReadFile("/proc/" + pid + "/stat")
					if os.IsNotExist(err) || errors.Is(err, syscall.ESRCH) {
						break
					}
					if err != nil {
						t.Fatal(err)
					}
					fields := strings.Fields(string(stat)[strings.LastIndex(string(stat), ")")+1:])
					if fields[0] == "Z" || fields[0] == "X" {
						break
					}
					if time.Now().After(deadline) {
						t.Fatalf("probe process %s still alive after cancellation deadline", pid)
					}
					time.Sleep(10 * time.Millisecond)
				}
			}
		})
	}
}

func TestSubprocessLogsBounded(t *testing.T) {
	fake := fakeExecutable(t, `if [ "$1" = "--version" ]; then
 printf 'Blender 5.3.0 Beta\n'
 head -c 4194304 /dev/zero | tr '\000' x
 exit 0
fi
printf 'BEGIN\n'
head -c 4194304 /dev/zero | tr '\000' x
printf '\nEND\n'
exit 1
`)
	rep, err := runner.Doctor(context.Background(), fake)
	if err != nil || !rep.OK {
		t.Fatalf("bounded probe lost version: %+v %v", rep, err)
	}
	result, err := runner.Run(context.Background(), runner.Options{Blender: fake, Timeout: 5 * time.Second})
	if err == nil {
		t.Fatal("expected failure")
	}
	if !result.LogTruncated || len(result.Stdout) > runner.MaxLogBytes+128 || len(err.Error()) > runner.MaxLogBytes+256 {
		t.Fatalf("logs unbounded: stdout=%d err=%d", len(result.Stdout), len(err.Error()))
	}
	if !strings.HasPrefix(result.Stdout, "BEGIN") || !strings.HasSuffix(result.Stdout, "END\n") || !strings.Contains(result.Stdout, "logs truncated") {
		t.Fatal("missing bounded head/tail/truncation metadata")
	}
}
