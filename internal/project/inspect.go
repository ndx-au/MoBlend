package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// InspectReport is machine-readable project dependency info.
type InspectReport struct {
	ProjectVersion string          `json:"project_version"`
	Entry          string          `json:"entry"`
	Modules        []InspectModule `json:"modules"`
	Shots          []InspectShot   `json:"shots,omitempty"`
}

// InspectModule describes one listed module.
type InspectModule struct {
	ID        string   `json:"id"`
	Path      string   `json:"path"`
	Kind      string   `json:"kind,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
	Imports   []string `json:"imports,omitempty"`
}

// InspectShot lists a shot target.
type InspectShot struct {
	ID    string `json:"id"`
	Entry string `json:"entry"`
}

// Inspect reads a project manifest and returns dependency listing without full resolve.
func Inspect(manifestPath string) (InspectReport, error) {
	abs, err := filepath.Abs(manifestPath)
	if err != nil {
		return InspectReport{}, err
	}
	root := filepath.Dir(abs)
	data, err := readBounded(abs, MaxFileBytes)
	if err != nil {
		return InspectReport{}, err
	}
	var man Manifest
	if err := strictDecode(data, &man); err != nil {
		return InspectReport{}, withSource(err, rel(root, abs))
	}
	if man.Kind != "project" {
		return InspectReport{}, fmt.Errorf("not a project manifest")
	}
	rep := InspectReport{
		ProjectVersion: man.ProjectVersion,
		Entry:          man.Entry,
	}
	for _, s := range man.Shots {
		rep.Shots = append(rep.Shots, InspectShot{ID: s.ID, Entry: s.Entry})
	}
	for _, m := range man.Modules {
		im := InspectModule{ID: m.ID, Path: m.Path}
		absM, err := resolveContainedPath(root, root, m.Path)
		if err == nil {
			if raw, err := os.ReadFile(absM); err == nil {
				var mf ModuleFile
				if json.Unmarshal(raw, &mf) == nil {
					im.Kind = mf.Kind
					im.Namespace = mf.Namespace
					for _, imp := range mf.Imports {
						im.Imports = append(im.Imports, imp.Module)
					}
					sort.Strings(im.Imports)
				}
			}
		}
		rep.Modules = append(rep.Modules, im)
	}
	return rep, nil
}
