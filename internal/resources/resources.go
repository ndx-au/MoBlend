// Package resources embeds the first-party asset pack and extracts it for jobs.
package resources

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// PackVersion is the identity string stored on the Phase 0 base scene.
const PackVersion = "0.1.0"

// BaseBlendRel is the relative path of the embedded pack file inside the FS.
const BaseBlendRel = "pack/base.blend"

// ExtractBasePack writes the embedded base.blend into destDir and returns its path.
// destDir must already exist. The file is written with mode 0o644.
func ExtractBasePack(destDir string) (string, error) {
	data, err := packFS.ReadFile(BaseBlendRel)
	if err != nil {
		return "", fmt.Errorf("read embedded pack: %w", err)
	}
	out := filepath.Join(destDir, "base.blend")
	if err := os.WriteFile(out, data, 0o644); err != nil {
		return "", fmt.Errorf("write pack: %w", err)
	}
	return out, nil
}

// BasePackSHA256 returns the hex SHA-256 of the embedded base.blend bytes.
func BasePackSHA256() (string, error) {
	data, err := packFS.ReadFile(BaseBlendRel)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// PackFS exposes the embedded filesystem for tests.
func PackFS() fs.FS {
	return packFS
}
