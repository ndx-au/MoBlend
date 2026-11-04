package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ndx-au/MoBlend/internal/diag"
	"github.com/ndx-au/MoBlend/internal/spec"
)

// InputKind classifies a JSON root for CLI dispatch.
type InputKind int

const (
	KindScene InputKind = iota
	KindProject
)

// DetectKind peeks at JSON bytes without full validation.
func DetectKind(data []byte) (InputKind, error) {
	var peek struct {
		Kind           string `json:"kind"`
		ProjectVersion string `json:"project_version"`
		SchemaVersion  string `json:"schema_version"`
	}
	if err := json.Unmarshal(data, &peek); err != nil {
		return KindScene, diag.ErrorList{{Code: "SPEC_JSON", Stage: "parse", Message: err.Error()}}
	}
	if peek.Kind == "project" {
		return KindProject, nil
	}
	return KindScene, nil
}

// LoadAndResolve loads a project manifest path and resolves to a Declaration.
func LoadAndResolve(manifestPath string) (Resolved, error) {
	return loadProject(manifestPath, "", true)
}

// LoadAndResolveShot resolves only the dependency closure of shotID (Phase 8).
// If shotID is empty and the project has no shots, behaves like LoadAndResolve.
func LoadAndResolveShot(manifestPath, shotID string) (Resolved, error) {
	return loadProject(manifestPath, shotID, false)
}

// LoadAndResolveAllListed validates the complete listed graph while assembling one shot.
func LoadAndResolveAllListed(manifestPath, shotID string) (Resolved, error) {
	return loadProject(manifestPath, shotID, true)
}

func loadProject(manifestPath, shotID string, allListed bool) (Resolved, error) {
	abs, err := filepath.Abs(manifestPath)
	if err != nil {
		return Resolved{}, err
	}
	data, err := readBounded(abs, MaxFileBytes)
	if err != nil {
		return Resolved{}, err
	}
	return resolveProjectSnapshot(abs, data, shotID, allListed)
}

func resolveProjectSnapshot(abs string, data []byte, shotID string, allListed bool) (Resolved, error) {
	root := filepath.Dir(abs)
	var man Manifest
	if err := strictDecode(data, &man); err != nil {
		return Resolved{}, withSource(err, rel(root, abs))
	}
	if man.Kind != "project" {
		return Resolved{}, diag.ErrorList{{Code: "PROJ_KIND", Stage: "schema", Source: rel(root, abs), Pointer: "/kind", Message: "expected kind project"}}
	}
	if man.ProjectVersion != ProjectVersion {
		return Resolved{}, diag.ErrorList{{Code: "PROJ_VERSION", Stage: "schema", Source: rel(root, abs), Pointer: "/project_version", Message: "unsupported project_version", Hint: "use " + ProjectVersion}}
	}
	if err := validateShots(man, rel(root, abs)); err != nil {
		return Resolved{}, err
	}
	if len(man.Shots) == 0 {
		if shotID != "" {
			return Resolved{}, diag.ErrorList{{
				Code: "SHOT_UNKNOWN", Stage: "schema", Pointer: "/shots", Source: rel(root, abs),
				Message: "project has no shots",
			}}
		}
		return resolveManifest(root, rel(root, abs), data, man, "", true)
	}
	if shotID == "" {
		if len(man.Shots) == 1 {
			shotID = man.Shots[0].ID
		} else {
			return Resolved{}, diag.ErrorList{{
				Code: "SHOT_REQUIRED", Stage: "schema", Pointer: "/shots", Source: rel(root, abs),
				Message: "multi-shot project requires --shot selection",
				Hint:    "pass --shot <id> or --all-shots for project-wide validation",
			}}
		}
	}
	entry := ""
	for i, s := range man.Shots {
		if s.ID == shotID {
			if s.Entry == "" {
				return Resolved{}, diag.ErrorList{{
					Code: "SHOT_ENTRY", Stage: "schema",
					Pointer: fmt.Sprintf("/shots/%d/entry", i), Source: rel(root, abs),
					Message: "shot missing entry module id",
				}}
			}
			entry = s.Entry
			break
		}
	}
	if entry == "" {
		return Resolved{}, diag.ErrorList{{
			Code: "SHOT_UNKNOWN", Stage: "schema", Pointer: "/shots", Source: rel(root, abs),
			Message: fmt.Sprintf("unknown shot id %q", shotID),
		}}
	}
	man.Entry = entry
	return resolveManifest(root, rel(root, abs), data, man, shotID, allListed)
}

func readBounded(path string, max int) ([]byte, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, diag.ErrorList{{Code: "PROJ_IO", Stage: "parse", Message: err.Error()}}
	}
	if !st.Mode().IsRegular() {
		return nil, diag.ErrorList{{Code: "PROJ_FILE_TYPE", Stage: "parse", Message: "input must be a regular file"}}
	}
	if st.Size() > int64(max) {
		return nil, diag.ErrorList{{Code: "LIM_FILE_BYTES", Stage: "parse", Message: fmt.Sprintf("file exceeds %d bytes", max)}}
	}
	f, err := openInput(path)
	if err != nil {
		return nil, diag.ErrorList{{Code: "PROJ_IO", Stage: "parse", Message: err.Error()}}
	}
	defer f.Close()
	st, err = f.Stat()
	if err != nil {
		return nil, err
	}
	if !st.Mode().IsRegular() {
		return nil, diag.ErrorList{{Code: "PROJ_FILE_TYPE", Stage: "parse", Message: "input must be a regular file"}}
	}
	data, err := io.ReadAll(io.LimitReader(f, int64(max)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > max {
		return nil, diag.ErrorList{{Code: "LIM_FILE_BYTES", Stage: "parse", Message: fmt.Sprintf("file exceeds %d bytes", max)}}
	}
	return data, nil
}

func strictDecode(data []byte, v any) error {
	// Use the scene decoder's duplicate-key/trailing-document guard before typed decoding.
	if err := spec.CheckJSON(data); err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return diag.ErrorList{{Code: "SPEC_SCHEMA", Stage: "schema", Message: err.Error()}}
	}
	return nil
}

// ValidShotID is also the output directory component contract.
func ValidShotID(id string) bool { return localRe.MatchString(id) }

func validateShots(man Manifest, source string) error {
	if len(man.Shots) > 64 {
		return diag.ErrorList{{Code: "LIM_SHOTS", Stage: "schema", Source: source, Message: "at most 64 shots"}}
	}
	seen := map[string]bool{}
	for i, shot := range man.Shots {
		p := fmt.Sprintf("/shots/%d", i)
		if !ValidShotID(shot.ID) {
			return diag.ErrorList{{Code: "SHOT_ID", Stage: "schema", Source: source, Pointer: p + "/id", Message: "shot id must match ^[A-Za-z][A-Za-z0-9_]*$"}}
		}
		if seen[shot.ID] {
			return diag.ErrorList{{Code: "SHOT_DUP_ID", Stage: "schema", Source: source, Pointer: p + "/id", Message: "duplicate shot id"}}
		}
		seen[shot.ID] = true
		found := false
		for _, m := range man.Modules {
			if m.ID == shot.Entry {
				found = true
			}
		}
		if !found {
			return diag.ErrorList{{Code: "SHOT_ENTRY", Stage: "schema", Source: source, Pointer: p + "/entry", Message: "shot entry must name a listed module"}}
		}
	}
	return nil
}

func withSource(err error, source string) error {
	elist, ok := err.(diag.ErrorList)
	if !ok {
		return err
	}
	for i := range elist {
		if elist[i].Source == "" {
			elist[i].Source = source
		}
	}
	return elist
}

func rel(root, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(r)
}

// IsProjectPath reports whether path content looks like a project (for CLI).
func IsProjectPath(path string) (bool, error) {
	data, err := readBounded(path, MaxFileBytes)
	if err != nil {
		return false, err
	}
	k, err := DetectKind(data)
	return k == KindProject, err
}

// Ensure scene schema import stays referenced for MaxDeclarationBytes alignment.
var _ = spec.SchemaVersion

// ListShots validates the outer manifest without reading shot contents.
func ListShots(path string) ([]ManifestShot, error) {
	data, err := readBounded(path, MaxFileBytes)
	if err != nil {
		return nil, err
	}
	var man Manifest
	if err := strictDecode(data, &man); err != nil {
		return nil, err
	}
	if man.Kind != "project" || man.ProjectVersion != ProjectVersion {
		return nil, fmt.Errorf("expected supported project manifest")
	}
	if err := validateShots(man, filepath.Base(path)); err != nil {
		return nil, err
	}
	return man.Shots, nil
}
