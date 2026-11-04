//go:build linux

package project_test

import (
	"github.com/ndx-au/MoBlend/internal/project"
	"path/filepath"
	"syscall"
	"testing"
)

func TestRejectSpecialFileInputs(t *testing.T) {
	dir := t.TempDir()
	fifo := filepath.Join(dir, "input.json")
	if err := syscall.Mkfifo(fifo, 0600); err != nil {
		t.Fatal(err)
	}
	_, err := project.LoadSceneOrProject(fifo, "")
	requireCode(t, err, "PROJ_FILE_TYPE")
}
