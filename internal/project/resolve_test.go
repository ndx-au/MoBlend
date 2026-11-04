package project_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ndx-au/MoBlend/internal/compile"
	"github.com/ndx-au/MoBlend/internal/diag"
	"github.com/ndx-au/MoBlend/internal/project"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..")
}

func TestModularMatchesFlat(t *testing.T) {
	root := repoRoot(t)
	mod, err := project.LoadAndResolve(filepath.Join(root, "examples/projects/title_modular/project.json"))
	if err != nil {
		t.Fatal(err)
	}
	flatPath := filepath.Join(root, "testdata/project/positive/title_flat.json")
	flatLR, err := project.LoadSceneOrProject(flatPath, "")
	if err != nil {
		t.Fatal(err)
	}
	a1, err := compile.CompileDeclaration(mod.Decl)
	if err != nil {
		t.Fatal(err)
	}
	a2, err := compile.CompileDeclaration(flatLR.Decl)
	if err != nil {
		t.Fatal(err)
	}
	if string(a1.Data) != string(a2.Data) {
		t.Fatalf("data mismatch\nmod=%s\nflat=%s", a1.Data, a2.Data)
	}
	if string(a1.Script) != string(a2.Script) {
		t.Fatal("script mismatch")
	}
}

func TestDeterminismSerializationAndRelocate(t *testing.T) {
	root := repoRoot(t)
	src := filepath.Join(root, "examples/projects/title_modular")
	a := t.TempDir()
	b := t.TempDir()
	copyTree(t, src, a)
	copyTree(t, src, b)

	// Re-serialize copy B to change formatting and key order.
	entryPath := filepath.Join(b, "main.entry.json")
	raw, err := os.ReadFile(entryPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	rewritten, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(entryPath, append(rewritten, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}

	ra, err := project.LoadAndResolve(filepath.Join(a, "project.json"))
	if err != nil {
		t.Fatal(err)
	}
	rb, err := project.LoadAndResolve(filepath.Join(b, "project.json"))
	if err != nil {
		t.Fatal(err)
	}
	ca, err := compile.CompileDeclaration(ra.Decl)
	if err != nil {
		t.Fatal(err)
	}
	cb, err := compile.CompileDeclaration(rb.Decl)
	if err != nil {
		t.Fatal(err)
	}
	if string(ca.Data) != string(cb.Data) || string(ca.Script) != string(cb.Script) {
		t.Fatal("relocated/reordered project must compile identically")
	}
}

func TestRejectCycle(t *testing.T) {
	_, err := project.LoadAndResolve(filepath.Join(repoRoot(t), "testdata/project/negative/cycle/project.json"))
	if err == nil {
		t.Fatal("expected cycle error")
	}
	el, ok := err.(diag.ErrorList)
	if !ok || el[0].Code != "PROJ_CYCLE" {
		t.Fatalf("got %v", err)
	}
}

func TestRejectMissing(t *testing.T) {
	_, err := project.LoadAndResolve(filepath.Join(repoRoot(t), "testdata/project/negative/missing_file/project.json"))
	if err == nil {
		t.Fatal("expected missing")
	}
	requireCode(t, err, "PROJ_FILE_MISSING")
}

func TestRejectEscape(t *testing.T) {
	_, err := project.LoadAndResolve(filepath.Join(repoRoot(t), "testdata/project/negative/escape/project.json"))
	if err == nil {
		t.Fatal("expected escape")
	}
	el, ok := err.(diag.ErrorList)
	if !ok {
		t.Fatalf("%T %v", err, err)
	}
	if el[0].Code != "PROJ_PATH_ESCAPE" {
		t.Fatalf("code %s", el[0].Code)
	}
}

func TestRejectSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	outsideFile := filepath.Join(outside, "evil.entry.json")
	if err := os.WriteFile(outsideFile, []byte(`{
  "module_version":"1.0.0","kind":"entry","namespace":"main",
  "schema_version":"1.0.0","asset_pack_version":"0.1.0",
  "scene":{"frame_start":1,"frame_end":1,"fps":24,"resolution":[64,64],"samples":4},
  "camera":{"id":"cam","location":[0,-4,1],"rotation_euler_deg":[90,0,0]},
  "objects":[{"id":"o","type":"form","form":"cube","location":[0,0,0],
    "material":{"base_color":"#FFFFFF","roughness":0.5}}]
}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outsideFile, filepath.Join(root, "main.entry.json")); err != nil {
		t.Fatal(err)
	}
	man := `{
  "project_version":"1.0.0","kind":"project","entry":"main",
  "modules":[{"id":"main","path":"main.entry.json"}]
}`
	if err := os.WriteFile(filepath.Join(root, "project.json"), []byte(man), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := project.LoadAndResolve(filepath.Join(root, "project.json"))
	if err == nil {
		t.Fatal("expected symlink escape rejection")
	}
	el, ok := err.(diag.ErrorList)
	if !ok || el[0].Code != "PROJ_PATH_ESCAPE" {
		t.Fatalf("got %v", err)
	}
}

func TestAggregateObjectLimit(t *testing.T) {
	root := t.TempDir()
	mods := []any{}
	imports := []any{}
	for i := 0; i < 8; i++ {
		id := string(rune('a' + i))
		ns := "ns" + id
		objs := make([]any, 0, 5)
		for j := 0; j < 5; j++ {
			objs = append(objs, map[string]any{
				"id": "o" + string(rune('0'+j)), "type": "form", "form": "cube",
				"location": []any{0, 0, 0},
				"material": map[string]any{"base_color": "#FFFFFF", "roughness": 0.5},
			})
		}
		path := id + ".module.json"
		body, _ := json.Marshal(map[string]any{
			"module_version": "1.0.0", "kind": "module", "namespace": ns,
			"imports": []any{}, "objects": objs, "animation": []any{},
		})
		if err := os.WriteFile(filepath.Join(root, path), body, 0o644); err != nil {
			t.Fatal(err)
		}
		mods = append(mods, map[string]any{"id": id, "path": path})
		imports = append(imports, map[string]any{"alias": ns, "module": id})
	}
	entry, _ := json.Marshal(map[string]any{
		"module_version": "1.0.0", "kind": "entry", "namespace": "main",
		"schema_version": "1.0.0", "asset_pack_version": "0.1.0",
		"imports": imports,
		"scene": map[string]any{
			"frame_start": 1, "frame_end": 1, "fps": 24,
			"resolution": []any{64, 64}, "samples": 4,
		},
		"camera": map[string]any{
			"id": "cam", "location": []any{0, -4, 1}, "rotation_euler_deg": []any{90, 0, 0},
		},
		"objects": []any{
			map[string]any{
				"id": "extra", "type": "form", "form": "cube", "location": []any{0, 0, 0},
				"material": map[string]any{"base_color": "#FFFFFF", "roughness": 0.5},
			},
		},
	})
	if err := os.WriteFile(filepath.Join(root, "main.entry.json"), entry, 0o644); err != nil {
		t.Fatal(err)
	}
	mods = append(mods, map[string]any{"id": "main", "path": "main.entry.json"})
	man, _ := json.Marshal(map[string]any{
		"project_version": "1.0.0", "kind": "project", "entry": "main", "modules": mods,
	})
	if err := os.WriteFile(filepath.Join(root, "project.json"), man, 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := project.LoadAndResolve(filepath.Join(root, "project.json"))
	if err == nil {
		t.Fatal("expected aggregate object limit")
	}
	el, ok := err.(diag.ErrorList)
	if !ok {
		t.Fatalf("%v", err)
	}
	for _, d := range el {
		if d.Code == "LIM_OBJECTS" {
			return
		}
	}
	t.Fatalf("expected LIM_OBJECTS, got %v", el)
}

func TestSourceAwareDiagnostic(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("project.json", `{
  "project_version":"1.0.0","kind":"project","entry":"main",
  "modules":[{"id":"main","path":"main.entry.json"},{"id":"prop","path":"prop.module.json"}]
}`)
	write("main.entry.json", `{
  "module_version":"1.0.0","kind":"entry","namespace":"main",
  "schema_version":"1.0.0","asset_pack_version":"0.1.0",
  "imports":[{"alias":"prop","module":"prop"}],
  "scene":{"frame_start":1,"frame_end":1,"fps":24,"resolution":[64,64],"samples":4},
  "camera":{"id":"cam","location":[0,-4,1],"rotation_euler_deg":[90,0,0]},
  "objects":[]
}`)
	write("prop.module.json", `{
  "module_version":"1.0.0","kind":"module","namespace":"prop",
  "imports":[],
  "objects":[{"id":"x","type":"asset","asset_id":"does_not_exist","location":[0,0,0]}],
  "animation":[]
}`)
	_, err := project.LoadAndResolve(filepath.Join(root, "project.json"))
	if err == nil {
		t.Fatal("expected unknown asset")
	}
	el, ok := err.(diag.ErrorList)
	if !ok {
		t.Fatal(err)
	}
	if el[0].Source == "" {
		t.Fatalf("expected source file on diagnostic: %+v", el[0])
	}
	if el[0].Source != "prop.module.json" {
		t.Fatalf("source=%q", el[0].Source)
	}
}

func TestInspect(t *testing.T) {
	rep, err := project.Inspect(filepath.Join(repoRoot(t), "examples/projects/title_modular/project.json"))
	if err != nil {
		t.Fatal(err)
	}
	if rep.Entry != "main" || len(rep.Modules) != 2 {
		t.Fatalf("%+v", rep)
	}
}

func TestDiamondDedup(t *testing.T) {
	r, err := project.LoadAndResolve(filepath.Join(repoRoot(t), "testdata/project/positive/diamond/project.json"))
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]int{}
	for _, o := range r.Decl.Objects {
		ids[o.ID]++
	}
	if ids["base__mark"] != 1 {
		t.Fatalf("diamond must include base once: %v", ids)
	}
	if len(r.Decl.Objects) != 4 {
		t.Fatalf("want 4 objects, got %d (%v)", len(r.Decl.Objects), ids)
	}
}

func TestRejectDuplicateNamespace(t *testing.T) {
	root := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("project.json", `{"project_version":"1.0.0","kind":"project","entry":"a","modules":[{"id":"a","path":"a.json"},{"id":"b","path":"b.json"}]}`)
	write("a.json", `{"module_version":"1.0.0","kind":"entry","namespace":"same","schema_version":"1.0.0","asset_pack_version":"0.1.0","imports":[{"alias":"b","module":"b"}],"scene":{"frame_start":1,"frame_end":1,"fps":24,"resolution":[64,64],"samples":4},"camera":{"id":"cam","location":[0,-4,1],"rotation_euler_deg":[90,0,0]},"objects":[{"id":"t","type":"text","text":"a"}]}`)
	write("b.json", `{"module_version":"1.0.0","kind":"module","namespace":"same","imports":[],"objects":[{"id":"t","type":"text","text":"b"}],"animation":[]}`)
	_, err := project.LoadAndResolve(filepath.Join(root, "project.json"))
	if err == nil {
		t.Fatal("expected duplicate namespace")
	}
	el, ok := err.(diag.ErrorList)
	if !ok || el[0].Code != "MOD_DUP_NAMESPACE" {
		t.Fatalf("got %v", err)
	}
}

func TestShotSelectionAndAllShots(t *testing.T) {
	path := filepath.Join(repoRoot(t), "testdata/project/shots/broken_unselected/project.json")
	good, err := project.LoadAndResolveShot(path, "good")
	if err != nil {
		t.Fatal(err)
	}
	if len(good.Files) < 2 {
		t.Fatalf("expected good closure files, got %v", good.Files)
	}
	for _, f := range good.Files {
		if f == "bad.entry.json" {
			t.Fatal("good shot must not load bad entry")
		}
	}
	_, err = project.LoadAndResolveShot(path, "bad")
	if err == nil {
		t.Fatal("bad shot should fail")
	}
	_, err = project.LoadAndResolveShot(path, "")
	if err == nil {
		t.Fatal("multi-shot without selection must fail")
	}
	el, ok := err.(diag.ErrorList)
	if !ok || el[0].Code != "SHOT_REQUIRED" {
		t.Fatalf("got %v", err)
	}
}

func TestThreeShotChangeIsolation(t *testing.T) {
	src := filepath.Join(repoRoot(t), "examples/projects/three_shots")
	dir := t.TempDir()
	copyTree(t, src, dir)
	man := filepath.Join(dir, "project.json")

	hashShot := func(id string) string {
		t.Helper()
		r, err := project.LoadAndResolveShot(man, id)
		if err != nil {
			t.Fatal(err)
		}
		h, err := project.ResolvedDataSHA(r.Decl)
		if err != nil {
			t.Fatal(err)
		}
		return h
	}

	intro0 := hashShot("intro")
	product0 := hashShot("product")
	outro0 := hashShot("outro")
	if intro0 == product0 || product0 == outro0 {
		t.Fatal("shots should resolve differently")
	}

	// Edit outro only.
	outroPath := filepath.Join(dir, "outro.entry.json")
	raw, err := os.ReadFile(outroPath)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	objs := doc["objects"].([]any)
	o0 := objs[0].(map[string]any)
	o0["text"] = "OutroEdited"
	rewritten, _ := json.MarshalIndent(doc, "", "  ")
	if err := os.WriteFile(outroPath, append(rewritten, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if hashShot("intro") != intro0 || hashShot("product") != product0 {
		t.Fatal("unrelated shots changed after outro edit")
	}
	if hashShot("outro") == outro0 {
		t.Fatal("outro hash should change")
	}

	// Edit shared module — dependent shots change.
	intro1 := hashShot("intro")
	sharedPath := filepath.Join(dir, "shared.module.json")
	sraw, err := os.ReadFile(sharedPath)
	if err != nil {
		t.Fatal(err)
	}
	var sdoc map[string]any
	if err := json.Unmarshal(sraw, &sdoc); err != nil {
		t.Fatal(err)
	}
	sobjs := sdoc["objects"].([]any)
	so0 := sobjs[0].(map[string]any)
	mat := so0["material"].(map[string]any)
	mat["base_color"] = "#FF0000"
	srewritten, _ := json.MarshalIndent(sdoc, "", "  ")
	if err := os.WriteFile(sharedPath, append(srewritten, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if hashShot("intro") == intro1 || hashShot("product") == product0 {
		t.Fatal("shared edit must change dependent shot hashes")
	}
}

func TestDependencyManifest(t *testing.T) {
	path := filepath.Join(repoRoot(t), "examples/projects/three_shots/project.json")
	r, err := project.LoadAndResolveShot(path, "intro")
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	p, err := project.WriteDependencyManifest(out, r)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	var m project.DependencyManifest
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if m.Shot != "intro" || len(m.Files) == 0 || m.ResolvedDataSHA == "" {
		t.Fatalf("%+v", m)
	}
}

func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}
