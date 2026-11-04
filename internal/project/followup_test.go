package project_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/diag"
	"github.com/ndx-au/MoBlend/internal/project"
	"github.com/ndx-au/MoBlend/internal/spec"
)

func modularCopy(t *testing.T) (string, *project.Manifest, *project.ModuleFile, *project.ModuleFile) {
	t.Helper()
	dir := t.TempDir()
	copyTree(t, filepath.Join(repoRoot(t), "examples/projects/title_modular"), dir)
	var man project.Manifest
	var entry, title project.ModuleFile
	for name, value := range map[string]any{"project.json": &man, "main.entry.json": &entry, "title.module.json": &title} {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(raw, value); err != nil {
			t.Fatal(err)
		}
	}
	return dir, &man, &entry, &title
}

func TestRejectDuplicateImportAliases(t *testing.T) {
	dir, _, entry, _ := modularCopy(t)
	entry.Imports = append(entry.Imports, project.Import{Alias: "title", Module: "main"})
	writeJSON(t, filepath.Join(dir, "main.entry.json"), entry)
	_, err := project.LoadAndResolve(filepath.Join(dir, "project.json"))
	requireCode(t, err, "MOD_DUP_IMPORT_ALIAS")
	el := err.(diag.ErrorList)
	if el[0].Source != "main.entry.json" || el[0].Pointer != "/imports/1/alias" {
		t.Fatal(el)
	}
}

func TestImportPermutationAndOrphanIdentity(t *testing.T) {
	dir, man, entry, title := modularCopy(t)
	extra := project.ModuleFile{ModuleVersion: project.ModuleVersion, Kind: "module", Namespace: "other", Objects: []spec.Object{{ID: "text", Type: "text", Text: "Other"}}}
	orphan := project.ModuleFile{ModuleVersion: project.ModuleVersion, Kind: "module", Namespace: "title"}
	man.Modules = append(man.Modules, project.ManifestMod{ID: "other", Path: "other.json"}, project.ManifestMod{ID: "orphan", Path: "orphan.json"})
	entry.Imports = append(entry.Imports, project.Import{Alias: "other", Module: "other"})
	title.Animation[0].TargetID = "title:text"
	entry.Animation = title.Animation
	title.Animation = nil
	entry.Animation = append(entry.Animation, spec.Track{TargetID: "other:text", Property: "scale", Interpolation: "linear", Keys: []spec.Key{{Frame: 1, Value: []float64{1, 1, 1}}, {Frame: 24, Value: []float64{2, 2, 2}}}})
	for name, doc := range map[string]any{"project.json": man, "main.entry.json": entry, "title.module.json": title, "other.json": extra, "orphan.json": orphan} {
		writeJSON(t, filepath.Join(dir, name), doc)
	}
	var baseline compile.Artifacts
	for i := 0; i < 50; i++ {
		if i == 1 {
			entry.Imports[0], entry.Imports[1] = entry.Imports[1], entry.Imports[0]
			for a, b := 0, len(man.Modules)-1; a < b; a, b = a+1, b-1 {
				man.Modules[a], man.Modules[b] = man.Modules[b], man.Modules[a]
			}
			writeJSON(t, filepath.Join(dir, "main.entry.json"), entry)
			writeJSON(t, filepath.Join(dir, "project.json"), man)
		}
		r, err := project.LoadAndResolve(filepath.Join(dir, "project.json"))
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		a, err := compile.CompileDeclaration(r.Decl)
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			baseline = a
		} else if string(a.Data) != string(baseline.Data) || string(a.Script) != string(baseline.Script) {
			t.Fatal("import/module order changed compilation")
		}
	}
}

func TestModularCameraAnimationMatchesFlat(t *testing.T) {
	dir, _, entry, _ := modularCopy(t)
	entry.Animation = []spec.Track{{TargetID: entry.Camera.ID, Property: "location", Interpolation: "linear", Keys: []spec.Key{{Frame: 1, Value: []float64{0, -4, 1}}, {Frame: 24, Value: []float64{0, -5, 2}}}}}
	writeJSON(t, filepath.Join(dir, "main.entry.json"), entry)
	r, err := project.LoadAndResolve(filepath.Join(dir, "project.json"))
	if err != nil {
		t.Fatal(err)
	}
	flat, err := project.LoadSceneOrProject(filepath.Join(repoRoot(t), "testdata/project/positive/title_flat.json"), "")
	if err != nil {
		t.Fatal(err)
	}
	// The title track is identical in both fixtures; add the same camera track.
	flat.Decl.Animation = append(entry.Animation, flat.Decl.Animation...)
	a, err := compile.CompileDeclaration(r.Decl)
	if err != nil {
		t.Fatal(err)
	}
	b, err := compile.CompileDeclaration(flat.Decl)
	if err != nil {
		t.Fatal(err)
	}
	if string(a.Data) != string(b.Data) || string(a.Script) != string(b.Script) {
		t.Fatal("modular camera tracks differ from flat")
	}
}

func TestRejectUndeclaredQualifiedTarget(t *testing.T) {
	dir, man, entry, _ := modularCopy(t)
	helper := project.ModuleFile{ModuleVersion: project.ModuleVersion, Kind: "module", Namespace: "helper", Animation: []spec.Track{{TargetID: "title__text", Property: "scale", Keys: []spec.Key{{Frame: 1, Value: []float64{1, 1, 1}}}}}}
	man.Modules = append(man.Modules, project.ManifestMod{ID: "helper", Path: "helper.json"})
	entry.Imports = append(entry.Imports, project.Import{Alias: "helper", Module: "helper"})
	for name, doc := range map[string]any{"project.json": man, "main.entry.json": entry, "helper.json": helper} {
		writeJSON(t, filepath.Join(dir, name), doc)
	}
	_, err := project.LoadAndResolve(filepath.Join(dir, "project.json"))
	requireCode(t, err, "SEM_TRACK_TARGET")
}

func TestRejectModuleEntryOnlyFields(t *testing.T) {
	for _, key := range []string{"scene", "camera", "schema_version", "asset_pack_version"} {
		t.Run(key, func(t *testing.T) {
			dir, _, _, _ := modularCopy(t)
			p := filepath.Join(dir, "title.module.json")
			raw, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			var doc map[string]any
			if err = json.Unmarshal(raw, &doc); err != nil {
				t.Fatal(err)
			}
			doc[key] = nil
			writeJSON(t, p, doc)
			_, err = project.LoadAndResolve(filepath.Join(dir, "project.json"))
			requireCode(t, err, "MOD_ENTRY_ONLY")
			if el := err.(diag.ErrorList); el[0].Source != "title.module.json" || el[0].Pointer != "/"+key {
				t.Fatal(el)
			}
		})
	}
}

func TestOriginalSourcePointers(t *testing.T) {
	for _, kind := range []string{"object", "track", "scene", "camera"} {
		t.Run(kind, func(t *testing.T) {
			dir, _, entry, title := modularCopy(t)
			entry.Objects = []spec.Object{{ID: "first", Type: "text", Text: "first"}}
			expected := diag.Diagnostic{Source: "title.module.json"}
			switch kind {
			case "object":
				title.Objects[0].Material.BaseColor = "invalid"
				expected.Code = "SEM_COLOR"
				expected.Pointer = "/objects/0/material/base_color"
			case "track":
				entry.Animation = []spec.Track{{TargetID: entry.Camera.ID, Property: "location", Keys: []spec.Key{{Frame: 1, Value: []float64{0, -4, 1}}}}}
				title.Animation[0].Property = "bad"
				expected.Code = "SEM_TRACK_PROPERTY"
				expected.Pointer = "/animation/0/property"
			case "scene":
				entry.Scene.FrameEnd = 0
				expected.Source = "main.entry.json"
				expected.Code = "SEM_FRAME_RANGE"
				expected.Pointer = "/scene/frame_end"
			case "camera":
				entry.Camera.ID = "bad id"
				expected.Source = "main.entry.json"
				expected.Code = "SEM_ID"
				expected.Pointer = "/camera/id"
			}
			writeJSON(t, filepath.Join(dir, "main.entry.json"), entry)
			writeJSON(t, filepath.Join(dir, "title.module.json"), title)
			_, err := project.LoadAndResolve(filepath.Join(dir, "project.json"))
			if err == nil {
				t.Fatal("expected diagnostic")
			}
			el, ok := err.(diag.ErrorList)
			if !ok {
				t.Fatal(err)
			}
			for _, d := range el {
				if d.Code == expected.Code && d.Source == expected.Source && d.Pointer == expected.Pointer {
					return
				}
			}
			t.Fatalf("want source %s pointer %s; got %+v", expected.Source, expected.Pointer, el)
		})
	}
}

func TestPublicInputByteLimit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "large.json")
	if err := os.WriteFile(path, []byte(strings.Repeat(" ", project.MaxFileBytes+1)), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := project.LoadSceneOrProject(path, "")
	requireCode(t, err, "LIM_FILE_BYTES")
}

func TestLoadedSourceMapEntry(t *testing.T) {
	dir, _, _, _ := modularCopy(t)
	r, err := project.LoadAndResolve(filepath.Join(dir, "project.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(r.SourceMap.Entry, project.SourceRef{File: "main.entry.json", Module: "main"}) {
		t.Fatal(r.SourceMap.Entry)
	}
}
