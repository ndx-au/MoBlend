package validate_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ndx-au/MoBlend/internal/diag"
	"github.com/ndx-au/MoBlend/internal/validate"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(wd, "../.."))
}

func TestPositiveExamples(t *testing.T) {
	root := repoRoot(t)
	for _, name := range []string{"title.json", "pedestal.json", "chart.json"} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(root, "examples", name))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := validate.ParseAndValidate(data); err != nil {
				t.Fatalf("validate: %v", err)
			}
		})
	}
}

func TestNegativeFixtures(t *testing.T) {
	root := repoRoot(t)
	dir := filepath.Join(root, "testdata/spec/negative")
	metaBytes, err := os.ReadFile(filepath.Join(dir, "meta.json"))
	if err != nil {
		t.Fatal(err)
	}
	var meta map[string]struct {
		Code    string `json:"code"`
		Pointer string `json:"pointer"`
	}
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		t.Fatal(err)
	}
	for name, want := range meta {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				t.Fatal(err)
			}
			_, err = validate.ParseAndValidate(data)
			if err == nil {
				t.Fatal("expected error")
			}
			elist, ok := err.(diag.ErrorList)
			if !ok || len(elist) == 0 {
				t.Fatalf("want ErrorList, got %T %v", err, err)
			}
			found := false
			for _, d := range elist {
				if d.Code == want.Code {
					found = true
					if want.Pointer != "" && d.Pointer != want.Pointer {
						t.Fatalf("pointer=%q want %q (diag=%+v)", d.Pointer, want.Pointer, d)
					}
					break
				}
			}
			if !found {
				t.Fatalf("missing code %s in %v", want.Code, elist)
			}
		})
	}
}

func TestSchemaModelParityRequired(t *testing.T) {
	// Structural parity: Go Declaration JSON tags cover schema required keys.
	required := []string{"schema_version", "asset_pack_version", "scene", "camera", "objects"}
	data, err := os.ReadFile(filepath.Join(repoRoot(t), "internal/resources/schema/v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	req, _ := schema["required"].([]any)
	got := map[string]bool{}
	for _, r := range req {
		got[r.(string)] = true
	}
	for _, k := range required {
		if !got[k] {
			t.Fatalf("schema missing required %s", k)
		}
	}
	sample := []byte(`{
	  "schema_version":"1.0.0",
	  "asset_pack_version":"0.1.0",
	  "scene":{"frame_start":1,"frame_end":1,"fps":24,"resolution":[64,64],"samples":4},
	  "camera":{"id":"cam","location":[0,-2,1],"rotation_euler_deg":[80,0,0]},
	  "objects":[{"id":"t","type":"text","text":"hi"}]
	}`)
	if _, err := validate.ParseAndValidate(sample); err != nil {
		t.Fatal(err)
	}
}
