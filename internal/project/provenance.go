package project

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ndx-au/MoBlend/internal/resources"
	"github.com/ndx-au/MoBlend/internal/spec"
)

// DependencyManifest records selected-shot input provenance (not compile identity).
type DependencyManifest struct {
	Shot             string         `json:"shot,omitempty"`
	SchemaVersion    string         `json:"schema_version"`
	AssetPackVersion string         `json:"asset_pack_version"`
	PackSHA256       string         `json:"pack_sha256"`
	Files            []FileDigest   `json:"files"`
	ResolvedDataSHA  string         `json:"resolved_data_sha256"`
	Effective        map[string]any `json:"effective_settings"`
	RecordedAt       string         `json:"recorded_at"`
}

// FileDigest is one project-relative input file hash from the load snapshot.
type FileDigest struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

// BuildDependencyManifest hashes loaded project files and resolved declaration bytes.
func BuildDependencyManifest(r Resolved) (DependencyManifest, error) {
	packSHA, err := resources.BasePackSHA256()
	if err != nil {
		return DependencyManifest{}, err
	}
	digests := append([]FileDigest(nil), r.Digests...)
	sort.Slice(digests, func(i, j int) bool { return digests[i].Path < digests[j].Path })
	canon, err := json.Marshal(r.Decl)
	if err != nil {
		return DependencyManifest{}, err
	}
	dataSum := sha256.Sum256(canon)
	return DependencyManifest{
		Shot:             r.Shot,
		SchemaVersion:    r.Decl.SchemaVersion,
		AssetPackVersion: r.Decl.AssetPackVersion,
		PackSHA256:       packSHA,
		Files:            digests,
		ResolvedDataSHA:  hex.EncodeToString(dataSum[:]),
		Effective: map[string]any{
			"frame_start": r.Decl.Scene.FrameStart,
			"frame_end":   r.Decl.Scene.FrameEnd,
			"fps":         r.Decl.Scene.FPS,
			"resolution":  r.Decl.Scene.Resolution,
			"samples":     r.Decl.Scene.Samples,
			"engine":      "CYCLES",
			"device":      "CPU",
		},
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// WriteDependencyManifest writes dependencies.json beside other job artifacts.
func WriteDependencyManifest(dir string, r Resolved) (string, error) {
	m, err := BuildDependencyManifest(r)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	data = append(data, '\n')
	path := filepath.Join(dir, "dependencies.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// ResolvedDataSHA returns a stable hash of the resolved declaration for change-isolation tests.
func ResolvedDataSHA(decl spec.Declaration) (string, error) {
	canon, err := json.Marshal(decl)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canon)
	return hex.EncodeToString(sum[:]), nil
}

func digest(path string, raw []byte) FileDigest {
	sum := sha256.Sum256(raw)
	return FileDigest{Path: path, SHA256: hex.EncodeToString(sum[:])}
}
