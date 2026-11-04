//go:build !linux && !darwin

package project

import "os"

func openInput(path string) (*os.File, error) { return os.Open(path) }
