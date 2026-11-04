package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func capture(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	old := os.Stdout
	f, err := os.CreateTemp(t.TempDir(), "stdout")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	os.Stdout = f
	defer func() { os.Stdout = old }()
	callErr := fn()
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	raw, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw), callErr
}

func TestValidateAllShotsReportsBrokenUnselected(t *testing.T) {
	path := filepath.Join("..", "..", "testdata/project/shots/broken_unselected/project.json")
	raw, err := capture(t, func() error { return cmdValidate([]string{path, "--shot", "good", "--json"}) })
	if err != nil {
		t.Fatalf("selected good: %s %v", raw, err)
	}
	raw, err = capture(t, func() error { return cmdValidate([]string{path, "--all-shots", "--json"}) })
	if err == nil {
		t.Fatal("all-shots accepted broken shot")
	}
	var report struct {
		OK          bool                          `json:"ok"`
		Diagnostics []struct{ Code, Hint string } `json:"diagnostics"`
	}
	if e := json.Unmarshal([]byte(raw), &report); e != nil {
		t.Fatal(e)
	}
	found := false
	for _, d := range report.Diagnostics {
		if d.Code == "PROJ_FILE_MISSING" && strings.Contains(d.Hint, "shot bad") {
			found = true
		}
	}
	if report.OK || !found {
		t.Fatalf("missing attributed bad-shot error: %s", raw)
	}
}

func TestCompileConcurrentShotOutputs(t *testing.T) {
	path := filepath.Join("..", "..", "examples/projects/three_shots/project.json")
	parent := t.TempDir()
	shots := []string{"intro", "outro", "intro"}
	errs := make([]error, len(shots))
	var wg sync.WaitGroup
	for i, shot := range shots {
		wg.Add(1)
		go func(i int, shot string) {
			defer wg.Done()
			errs[i] = cmdCompile([]string{path, "--shot", shot, "--out", parent})
		}(i, shot)
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	for shot, want := range map[string]int{"intro": 2, "outro": 1} {
		jobs, err := filepath.Glob(filepath.Join(parent, shot, "moblend-job-*"))
		if err != nil || len(jobs) != want {
			t.Fatalf("%s jobs=%v err=%v", shot, jobs, err)
		}
		for _, job := range jobs {
			for _, name := range []string{"apply_scene.py", "scene_data.json", "source_map.json", "dependencies.json"} {
				if _, err := os.Stat(filepath.Join(job, name)); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
}

func TestOutputParentRejectsEscape(t *testing.T) {
	for _, id := range []string{"../escape", "/tmp/escape", "a/b"} {
		if _, err := outputParent(t.TempDir(), id); err == nil {
			t.Fatalf("accepted %q", id)
		}
	}
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "intro")); err != nil {
		t.Fatal(err)
	}
	if _, err := outputParent(root, "intro"); err == nil {
		t.Fatal("accepted escaping symlink")
	}
}

func TestValidateJSONOriginalSourcePointer(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join("..", "..", "examples/projects/title_modular")
	for _, name := range []string{"project.json", "main.entry.json", "title.module.json"} {
		raw, err := os.ReadFile(filepath.Join(source, name))
		if err != nil {
			t.Fatal(err)
		}
		var doc map[string]any
		if err = json.Unmarshal(raw, &doc); err != nil {
			t.Fatal(err)
		}
		if name == "main.entry.json" {
			doc["objects"] = []any{map[string]any{"id": "first", "type": "text", "text": "first"}}
		}
		if name == "title.module.json" {
			doc["objects"].([]any)[0].(map[string]any)["material"].(map[string]any)["base_color"] = "bad"
		}
		raw, err = json.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(dir, name), raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	raw, err := capture(t, func() error { return cmdValidate([]string{filepath.Join(dir, "project.json"), "--json"}) })
	if err == nil {
		t.Fatal("invalid color accepted")
	}
	var report struct {
		OK          bool                                     `json:"ok"`
		Diagnostics []struct{ Code, Source, Pointer string } `json:"diagnostics"`
	}
	if err = json.Unmarshal([]byte(raw), &report); err != nil {
		t.Fatal(err)
	}
	for _, d := range report.Diagnostics {
		if d.Code == "SEM_COLOR" && d.Source == "title.module.json" && d.Pointer == "/objects/0/material/base_color" {
			return
		}
	}
	t.Fatalf("missing original diagnostic: %s", raw)
}
