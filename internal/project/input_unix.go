//go:build linux || darwin

package project

import (
	"os"
	"syscall"
)

// Nonblocking open prevents a replacement FIFO from hanging before fstat.
func openInput(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
}
