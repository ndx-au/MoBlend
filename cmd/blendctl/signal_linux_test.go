//go:build linux

package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestCLISignalsCancelJobs(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "blendctl")
	if out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v %s", err, out)
	}
	for _, mode := range []string{"preview", "render", "smoke", "preview-version"} {
		t.Run(mode, func(t *testing.T) {
			parent := t.TempDir()
			ready := filepath.Join(parent, "ready")
			fake := filepath.Join(parent, "blender")
			script := `#!/bin/sh
if [ "$1" = "--version" ]; then echo 'Blender 5.3.0 Beta'; exit 0; fi
sleep 60 &
printf '%s' "$!" > '` + ready + `'
wait
`
			if err := os.WriteFile(fake, []byte(script), 0700); err != nil {
				t.Fatal(err)
			}
			command := mode
			if mode == "preview-version" {
				command = "preview"
				lines := strings.SplitN(script, "\n", 3)
				script = lines[0] + "\n" + lines[2]
				if err := os.WriteFile(fake, []byte(script), 0700); err != nil {
					t.Fatal(err)
				}
			}

			args := []string{command}
			if mode != "smoke" {
				args = append(args, filepath.Join("..", "..", "testdata/compile/title_render_tiny.json"))
			}
			args = append(args, "--out", parent)
			cmd := exec.Command(binary, args...)
			cmd.Env = append(os.Environ(), "MOBLEND_BLENDER="+fake)
			var output bytes.Buffer
			cmd.Stdout = &output
			cmd.Stderr = &output
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if cmd.ProcessState == nil {
					_ = cmd.Process.Kill()
					_ = cmd.Wait()
				}
			}()
			deadline := time.Now().Add(5 * time.Second)
			var child []byte
			for time.Now().Before(deadline) {
				child, _ = os.ReadFile(ready)
				if len(child) > 0 {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			if len(child) == 0 {
				t.Fatal("fake Blender not ready")
			}
			sig := os.Signal(os.Interrupt)
			if mode == "render" {
				sig = syscall.SIGTERM
			}
			if err := cmd.Process.Signal(sig); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { done <- cmd.Wait() }()
			select {
			case err := <-done:
				if err == nil {
					t.Fatalf("%s succeeded after signal", mode)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("CLI did not return on signal")
			}
			// The CLI must handle the signal through context cancellation, rather than
			// default signal termination (which cannot emit the structured job result).
			if !strings.Contains(output.String(), "canceled") {
				t.Fatalf("signal bypassed cancellation path: %s", output.String())
			}
			if mode != "smoke" {
				failed, _ := filepath.Glob(filepath.Join(parent, "moblend-job-*", "out", "manifest.failed.json"))
				success, _ := filepath.Glob(filepath.Join(parent, "moblend-job-*", "out", "manifest.json"))
				if len(failed) != 1 || len(success) != 0 {
					t.Fatalf("failed=%v success=%v", failed, success)
				}
			}
			// SIGKILL delivery is asynchronous across the process group. Waiting
			// for the CLI reaps its parent, but does not synchronously reap children.
			// Require descendant termination within a bounded interval, not the same
			// scheduling instant as parent exit.
			deadline = time.Now().Add(time.Second)
			for {
				stat, err := os.ReadFile("/proc/" + string(child) + "/stat")
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
					t.Fatalf("child still alive after cancellation deadline: %s", child)
				}
				time.Sleep(10 * time.Millisecond)
			}
		})
	}
}
