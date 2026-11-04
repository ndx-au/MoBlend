// Package compile lowers validated declarations to deterministic Blender job artifacts.
package compile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/ndx-au/MoBlend/internal/resources"
	"github.com/ndx-au/MoBlend/internal/spec"
	"github.com/ndx-au/MoBlend/internal/validate"
)

// Artifacts are compile outputs (no Blender required).
type Artifacts struct {
	Data       []byte
	Script     []byte
	DataName   string
	ScriptName string
}

// Compile validates (if raw JSON) or accepts a normalized declaration and emits artifacts.
func CompileDeclaration(d spec.Declaration) (Artifacts, error) {
	data, err := MarshalCanonical(d)
	if err != nil {
		return Artifacts{}, err
	}
	script, err := resources.ApplyScript()
	if err != nil {
		return Artifacts{}, err
	}
	return Artifacts{
		Data:       data,
		Script:     script,
		DataName:   "scene_data.json",
		ScriptName: "apply_scene.py",
	}, nil
}

// CompileJSON parses, validates, and compiles.
func CompileJSON(raw []byte) (Artifacts, validate.Result, error) {
	res, err := validate.ParseAndValidate(raw)
	if err != nil {
		return Artifacts{}, res, err
	}
	art, err := CompileDeclaration(res.Decl)
	return art, res, err
}

// WriteToDir writes script + data into dir.
func WriteToDir(dir string, art Artifacts) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, art.ScriptName), art.Script, 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, art.DataName), art.Data, 0o644)
}

// MarshalCanonical encodes with sorted object keys for byte-stable output.
func MarshalCanonical(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var anyv any
	if err := json.Unmarshal(raw, &anyv); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := writeCanonical(&buf, anyv, 0); err != nil {
		return nil, err
	}
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

func writeCanonical(buf *bytes.Buffer, v any, indent int) error {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		buf.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			kb, _ := json.Marshal(k)
			buf.Write(kb)
			buf.WriteByte(':')
			if err := writeCanonical(buf, t[k], indent+1); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
	case []any:
		buf.WriteByte('[')
		for i, el := range t {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := writeCanonical(buf, el, indent+1); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
	case float64, bool, string, nil:
		b, err := json.Marshal(t)
		if err != nil {
			return err
		}
		buf.Write(b)
	default:
		return fmt.Errorf("unsupported canonical type %T", v)
	}
	return nil
}
