// Package project loads and resolves multi-file MoBlend projects into a single Declaration.
package project

import "github.com/ndx-au/MoBlend/internal/spec"

const (
	ProjectVersion = "1.0.0"
	ModuleVersion  = "1.0.0"

	MaxModules      = 32
	MaxImportDepth  = 8
	MaxFileBytes    = spec.MaxDeclarationBytes
	MaxProjectBytes = 1 << 20 // 1 MiB
)

// Manifest is the project root document.
type Manifest struct {
	ProjectVersion string         `json:"project_version"`
	Kind           string         `json:"kind"`
	Entry          string         `json:"entry"`
	Modules        []ManifestMod  `json:"modules"`
	Shots          []ManifestShot `json:"shots,omitempty"` // Phase 8
}

// ManifestMod lists one module file.
type ManifestMod struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

// ManifestShot is a Phase 8 shot entry (ignored by Phase 7 resolve-all-as-one when empty).
type ManifestShot struct {
	ID    string `json:"id"`
	Entry string `json:"entry"` // module id used as entry for this shot
}

// ModuleFile is a module or entry document before resolution.
type ModuleFile struct {
	ModuleVersion    string        `json:"module_version"`
	Kind             string        `json:"kind"` // module | entry
	Namespace        string        `json:"namespace"`
	SchemaVersion    string        `json:"schema_version,omitempty"`
	AssetPackVersion string        `json:"asset_pack_version,omitempty"`
	Imports          []Import      `json:"imports,omitempty"`
	Scene            *spec.Scene   `json:"scene,omitempty"`
	Camera           *spec.Camera  `json:"camera,omitempty"`
	Lights           []spec.Light  `json:"lights,omitempty"`
	Objects          []spec.Object `json:"objects,omitempty"`
	Animation        []spec.Track  `json:"animation,omitempty"`
}

// Import aliases another module id inside this file.
type Import struct {
	Alias  string `json:"alias"`
	Module string `json:"module"`
}

// SourceRef maps a resolved entity back to its origin.
type SourceRef struct {
	File    string `json:"file"`
	Pointer string `json:"pointer"`
	Module  string `json:"module"`
}

// SourceMap is emitted beside compile artifacts (not part of script/data identity).
type SourceMap struct {
	Entry   SourceRef            `json:"entry"`
	Objects map[string]SourceRef `json:"objects"`
	Lights  map[string]SourceRef `json:"lights"`
	Tracks  []SourceRef          `json:"tracks"`
}

// Resolved is a flattened declaration plus provenance.
type Resolved struct {
	Decl      spec.Declaration
	SourceMap SourceMap
	RootDir   string
	Files     []string     // project-relative paths loaded
	Digests   []FileDigest // immutable hashes captured at load time
	Shot      string       // selected shot id when resolving a multi-shot project
}
