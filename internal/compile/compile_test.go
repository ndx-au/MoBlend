package compile_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ndx-au/MoBlend/internal/compile"
)

func TestCompileByteStableAndSafeText(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", ".."))
	raw, err := os.ReadFile(filepath.Join(root, "examples/title.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	objs := doc["objects"].([]any)
	obj0 := objs[0].(map[string]any)
	awkward := "He said \"hi\" \\ and café 🎨"
	obj0["text"] = awkward
	raw2, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	a1, _, err := compile.CompileJSON(raw2)
	if err != nil {
		t.Fatal(err)
	}
	a2, _, err := compile.CompileJSON(raw2)
	if err != nil {
		t.Fatal(err)
	}
	if string(a1.Data) != string(a2.Data) {
		t.Fatal("data not byte-stable")
	}
	if string(a1.Script) != string(a2.Script) {
		t.Fatal("script not byte-stable")
	}
	if !strings.Contains(string(a1.Data), "café") {
		t.Fatalf("expected unicode text in data, got %s", a1.Data)
	}
	if !strings.Contains(string(a1.Data), `\"hi\"`) && !strings.Contains(string(a1.Data), `\u0022hi\u0022`) {
		// encoding/json escapes quotes as \"
		if !strings.Contains(string(a1.Data), "He said") {
			t.Fatalf("expected awkward text in data, got %s", a1.Data)
		}
	}
	if strings.Contains(string(a1.Script), "café") || strings.Contains(string(a1.Script), awkward) {
		t.Fatal("user text leaked into script")
	}
}
