// Package runner launches Blender with an argument vector for isolated jobs.
package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// DefaultPythonExitCode is the nonzero code Blender returns when a CLI Python
// script raises (via --python-exit-code).
const DefaultPythonExitCode = 1

// Options configures a single Blender invocation.
type Options struct {
	Blender   string
	BlendFile string
	Python    string
	ExtraArgs []string // passed after "--"
	WorkDir   string
	Timeout   time.Duration
	Env       []string // if nil, inherits process env
}

// Result captures process outcome.
type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Args     []string
}

// BuildArgs returns the ordered Blender argv (without the executable).
func BuildArgs(opt Options) []string {
	args := []string{
		"--background",
		"--factory-startup",
		"--disable-autoexec",
		"--python-exit-code",
		fmt.Sprintf("%d", DefaultPythonExitCode),
	}
	if opt.BlendFile != "" {
		args = append(args, opt.BlendFile)
	}
	if opt.Python != "" {
		args = append(args, "--python", opt.Python)
	}
	if len(opt.ExtraArgs) > 0 {
		args = append(args, "--")
		args = append(args, opt.ExtraArgs...)
	}
	return args
}

// Run executes Blender with opt. Context cancellation kills the process tree
// via CommandContext.
func Run(ctx context.Context, opt Options) (Result, error) {
	if opt.Blender == "" {
		return Result{}, errors.New("blender executable path is required")
	}
	if opt.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.Timeout)
		defer cancel()
	}

	args := BuildArgs(opt)
	cmd := exec.CommandContext(ctx, opt.Blender, args...)
	if opt.WorkDir != "" {
		cmd.Dir = opt.WorkDir
	}
	if opt.Env != nil {
		cmd.Env = opt.Env
	}

	out, err := cmd.CombinedOutput()
	res := Result{
		Args:   append([]string{opt.Blender}, args...),
		Stdout: string(out),
		Stderr: "",
	}
	if err == nil {
		res.ExitCode = 0
		return res, nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		res.ExitCode = ee.ExitCode()
		return res, fmt.Errorf("blender exit %d: %w\n%s", res.ExitCode, err, out)
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return res, fmt.Errorf("blender timed out: %w\n%s", ctx.Err(), out)
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return res, fmt.Errorf("blender canceled: %w\n%s", ctx.Err(), out)
	}
	return res, fmt.Errorf("blender run failed: %w\n%s", err, out)
}

// ResolveBlender picks MOBLEND_BLENDER, then BLENDER, then "blender" on PATH.
func ResolveBlender() (string, error) {
	for _, key := range []string{"MOBLEND_BLENDER", "BLENDER"} {
		if v := os.Getenv(key); v != "" {
			return v, nil
		}
	}
	path, err := exec.LookPath("blender")
	if err != nil {
		return "", fmt.Errorf("blender not found on PATH (set MOBLEND_BLENDER): %w", err)
	}
	return path, nil
}

// JobDir creates a unique empty job directory under parent (or os.TempDir).
func JobDir(parent string) (string, error) {
	if parent == "" {
		parent = os.TempDir()
	}
	return os.MkdirTemp(parent, "moblend-job-*")
}

// MustAbs returns an absolute path or panics — for tests only.
func MustAbs(path string) string {
	a, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return a
}
