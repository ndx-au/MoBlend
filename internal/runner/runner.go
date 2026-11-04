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
	Env       []string // optional explicit environment; default is an isolated allowlist
}

// Result captures process outcome.
type Result struct {
	LogTruncated bool
	ExitCode     int
	Stdout       string
	Stderr       string
	Args         []string
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
// using a dedicated process group on Linux/macOS.
func Run(ctx context.Context, opt Options) (Result, error) {
	if opt.Blender == "" {
		return Result{}, errors.New("blender executable path is required")
	}
	if opt.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.Timeout)
		defer cancel()
	}

	return runCommand(ctx, opt.Blender, BuildArgs(opt), opt.WorkDir, opt.Env)
}

// runCommand is the common process/environment/log boundary for version probes
// and execution. Both must have identical cancellation behavior.
func runCommand(ctx context.Context, blender string, args []string, workDir string, env []string) (Result, error) {
	cmd := exec.CommandContext(ctx, blender, args...)
	cmd.Dir = workDir
	configureProcessTree(cmd)
	cmd.WaitDelay = 2 * time.Second
	cmd.Env = env
	if cmd.Env == nil {
		envDir, err := os.MkdirTemp(workDir, "moblend-env-*")
		if err != nil {
			return Result{}, fmt.Errorf("create isolated environment: %w", err)
		}
		defer os.RemoveAll(envDir)
		cmd.Env = isolatedEnv(envDir)
	}
	logs := &boundedLog{}
	cmd.Stdout = logs
	cmd.Stderr = logs
	err := cmd.Run()
	out, truncated := logs.snapshot()
	res := Result{Args: append([]string{blender}, args...), Stdout: out, LogTruncated: truncated}
	if ctx.Err() != nil {
		res.ExitCode = -1
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return res, fmt.Errorf("blender timed out: %w\n%s", ctx.Err(), out)
		}
		return res, fmt.Errorf("blender canceled: %w\n%s", ctx.Err(), out)
	}
	if err == nil {
		return res, nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		res.ExitCode = ee.ExitCode()
		return res, fmt.Errorf("blender exit %d: %w\n%s", res.ExitCode, err, out)
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

func isolatedEnv(dir string) []string {
	env := []string{}
	for _, key := range []string{"PATH", "LANG", "LC_ALL", "SYSTEMROOT", "WINDIR"} {
		if v, ok := os.LookupEnv(key); ok {
			env = append(env, key+"="+v)
		}
	}
	for _, key := range []string{"HOME", "XDG_CONFIG_HOME", "XDG_CACHE_HOME", "XDG_DATA_HOME", "TMPDIR"} {
		env = append(env, key+"="+dir)
	}
	// Do not inherit Python startup paths, credentials, or Blender user scripts.
	return env
}
