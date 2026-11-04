//go:build !linux && !darwin

package runner

import "os/exec"

// Other platforms are not advertised until process-tree cancellation is tested.
func configureProcessTree(cmd *exec.Cmd) {}
