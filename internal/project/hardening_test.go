package project_test

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ndx-au/MoBlend/internal/diag"
	"github.com/ndx-au/MoBlend/internal/project"
)

func requireCode(t *testing.T, err error, code string) {
	t.Helper()
	el, ok := err.(diag.ErrorList)
	if !ok {
		t.Fatalf("want %s, got %T: %v", code, err, err)
	}
	for _, d := range el {
		if d.Code == code {
			return
		}
	}
	t.Fatalf("want %s, got %v", code, err)
}

func writeJSON(t *testing.T, path string, doc any) {
	t.Helper()
	b, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestGraphBounds(t *testing.T) {
	for _, tc := range []struct {
		name, code      string
		count, padding  int
		chain, absolute bool
	}{
		{name: "depth", code: "LIM_IMPORT_DEPTH", count: 10, chain: true},
		{name: "file bytes", code: "LIM_FILE_BYTES", count: 1, padding: project.MaxFileBytes + 1},
		{name: "project bytes", code: "LIM_PROJECT_BYTES", count: 20, padding: project.MaxFileBytes - 1024},
		{name: "module count", code: "LIM_MODULES", count: 33},
		{name: "absolute path", code: "PROJ_PATH_ABS", count: 1, absolute: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			mods := []project.ManifestMod{}
			for i := 0; i < tc.count; i++ {
				id := fmt.Sprintf("m%d", i)
				path := id + ".json"
				doc := project.ModuleFile{ModuleVersion: project.ModuleVersion, Kind: "module", Namespace: id}
				if tc.chain && i+1 < tc.count {
					doc.Imports = []project.Import{{Alias: "next", Module: fmt.Sprintf("m%d", i+1)}}
				}
				raw, err := json.Marshal(doc)
				if err != nil {
					t.Fatal(err)
				}
				raw = append(raw, []byte(strings.Repeat(" ", tc.padding))...)
				if err := os.WriteFile(filepath.Join(dir, path), raw, 0600); err != nil {
					t.Fatal(err)
				}
				if tc.absolute {
					path = filepath.Join(dir, path)
				}
				mods = append(mods, project.ManifestMod{ID: id, Path: path})
			}
			man := project.Manifest{Kind: "project", ProjectVersion: project.ProjectVersion, Entry: "m0", Modules: mods}
			path := filepath.Join(dir, "project.json")
			writeJSON(t, path, man)
			_, err := project.LoadAndResolve(path)
			requireCode(t, err, tc.code)
		})
	}
}

func TestUnreachableListedGraphRejected(t *testing.T) {
	for _, tc := range []struct{ name, code string }{{"missing", "PROJ_FILE_MISSING"}, {"version", "MOD_VERSION"}, {"cycle", "PROJ_CYCLE"}} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			copyTree(t, filepath.Join(repoRoot(t), "examples/projects/title_modular"), dir)
			path := filepath.Join(dir, "project.json")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var man project.Manifest
			if err := json.Unmarshal(raw, &man); err != nil {
				t.Fatal(err)
			}
			man.Modules = append(man.Modules, project.ManifestMod{ID: "unused", Path: "unused.json"})
			if tc.name != "missing" {
				mod := project.ModuleFile{Kind: "module", ModuleVersion: project.ModuleVersion, Namespace: "unused"}
				if tc.name == "version" {
					mod.ModuleVersion = "999"
				} else {
					mod.Imports = []project.Import{{Alias: "self", Module: "unused"}}
				}
				writeJSON(t, filepath.Join(dir, "unused.json"), mod)
			}
			writeJSON(t, path, man)
			_, err = project.LoadAndResolve(path)
			requireCode(t, err, tc.code)
		})
	}
}

func TestDependencyManifestUsesLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	copyTree(t, filepath.Join(repoRoot(t), "examples/projects/three_shots"), dir)
	path := filepath.Join(dir, "project.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	r, err := project.LoadAndResolveShot(path, "intro")
	if err != nil {
		t.Fatal(err)
	}
	before, err := project.BuildDependencyManifest(r)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range r.Files {
		if err := os.Remove(filepath.Join(dir, f)); err != nil {
			t.Fatal(err)
		}
	}
	after, err := project.BuildDependencyManifest(r)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before.Files, after.Files) || before.ResolvedDataSHA != after.ResolvedDataSHA {
		t.Fatal("provenance re-read disk")
	}
	expected := fmt.Sprintf("%x", sha256.Sum256(raw))
	found := false
	for _, f := range after.Files {
		if f.Path == "project.json" {
			found = true
			if f.SHA256 != expected {
				t.Fatalf("manifest digest = %s", f.SHA256)
			}
		}
	}
	if !found {
		t.Fatal("manifest snapshot omitted")
	}
}

func TestUnsafeShotIDsRejectedGlobally(t *testing.T) {
	for _, id := range []string{"../escape", "/tmp/escape", "a/b", "a\\b", "", ".", "..", "a-b"} {
		t.Run(id, func(t *testing.T) {
			dir := t.TempDir()
			copyTree(t, filepath.Join(repoRoot(t), "examples/projects/three_shots"), dir)
			path := filepath.Join(dir, "project.json")
			raw, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var man project.Manifest
			if err := json.Unmarshal(raw, &man); err != nil {
				t.Fatal(err)
			}
			man.Shots[2].ID = id
			writeJSON(t, path, man)
			_, err = project.LoadAndResolveShot(path, "intro")
			requireCode(t, err, "SHOT_ID")
		})
	}
}

func TestDetectKindRequiresKind(t *testing.T) {
	k, err := project.DetectKind([]byte(`{"schema_version":"1.0.0","project_version":"1.0.0"}`))
	if err != nil || k != project.KindScene {
		t.Fatalf("kind=%v err=%v", k, err)
	}
}
