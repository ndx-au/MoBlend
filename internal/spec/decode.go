package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/ndx-au/MoBlend/internal/diag"
)

// Parse reads a declaration from JSON bytes with strict rules:
// size limit, no duplicate keys, no trailing data, known fields only,
// finite numbers.
func Parse(data []byte) (Declaration, error) {
	if len(data) > MaxDeclarationBytes {
		return Declaration{}, diag.ErrorList{{
			Code:    "LIM_DECLARATION_BYTES",
			Stage:   "parse",
			Pointer: "",
			Message: fmt.Sprintf("declaration exceeds %d bytes", MaxDeclarationBytes),
			Hint:    "split scenes or reduce embedded data",
		}}
	}
	if err := rejectDuplicateKeys(data); err != nil {
		return Declaration{}, err
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	dec.UseNumber()

	var raw map[string]json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return Declaration{}, diag.ErrorList{{
			Code:    "SPEC_JSON",
			Stage:   "parse",
			Message: err.Error(),
		}}
	}
	if err := rejectTrailing(dec); err != nil {
		return Declaration{}, err
	}

	// Re-encode canonical object map then decode into typed struct so
	// DisallowUnknownFields applies at nested levels via a second pass.
	normalized, err := json.Marshal(raw)
	if err != nil {
		return Declaration{}, err
	}
	if err := rejectNonFinite(normalized); err != nil {
		return Declaration{}, err
	}

	dec2 := json.NewDecoder(bytes.NewReader(normalized))
	dec2.DisallowUnknownFields()
	var d Declaration
	if err := dec2.Decode(&d); err != nil {
		code := "SPEC_SCHEMA"
		msg := err.Error()
		if strings.Contains(msg, "unknown field") {
			code = "SPEC_UNKNOWN_FIELD"
		}
		return Declaration{}, diag.ErrorList{{
			Code:    code,
			Stage:   "schema",
			Message: msg,
		}}
	}
	return d, nil
}

func rejectTrailing(dec *json.Decoder) error {
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return diag.ErrorList{{
				Code:    "SPEC_TRAILING_DATA",
				Stage:   "parse",
				Message: "trailing data after JSON value",
			}}
		}
		if err != io.EOF {
			return diag.ErrorList{{
				Code:    "SPEC_TRAILING_DATA",
				Stage:   "parse",
				Message: err.Error(),
			}}
		}
	}
	return nil
}

func rejectNonFinite(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var walk func(path string) error
	walk = func(path string) error {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		switch v := tok.(type) {
		case json.Delim:
			switch v {
			case '{':
				for dec.More() {
					keyTok, err := dec.Token()
					if err != nil {
						return err
					}
					key, _ := keyTok.(string)
					child := path + "/" + key
					if err := walk(child); err != nil {
						return err
					}
				}
				_, err = dec.Token() // }
				return err
			case '[':
				i := 0
				for dec.More() {
					child := fmt.Sprintf("%s/%d", path, i)
					if err := walk(child); err != nil {
						return err
					}
					i++
				}
				_, err = dec.Token() // ]
				return err
			}
		case json.Number:
			f, err := v.Float64()
			if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
				return diag.ErrorList{{
					Code:    "SPEC_NONFINITE_NUMBER",
					Stage:   "parse",
					Pointer: pointerOrRoot(path),
					Message: "non-finite or invalid number",
				}}
			}
		}
		return nil
	}
	if err := walk(""); err != nil {
		return err
	}
	return nil
}

func pointerOrRoot(path string) string {
	if path == "" {
		return ""
	}
	return path
}

// rejectDuplicateKeys scans object keys at every nesting level.
func rejectDuplicateKeys(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	var walk func(path string) error
	walk = func(path string) error {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		switch v := tok.(type) {
		case json.Delim:
			switch v {
			case '{':
				seen := map[string]struct{}{}
				for dec.More() {
					keyTok, err := dec.Token()
					if err != nil {
						return err
					}
					key, ok := keyTok.(string)
					if !ok {
						return fmt.Errorf("object key type")
					}
					child := path + "/" + escapePointer(key)
					if _, dup := seen[key]; dup {
						return diag.ErrorList{{
							Code:    "SPEC_DUPLICATE_KEY",
							Stage:   "parse",
							Pointer: child,
							Message: fmt.Sprintf("duplicate key %q", key),
						}}
					}
					seen[key] = struct{}{}
					if err := walk(child); err != nil {
						return err
					}
				}
				_, err = dec.Token()
				return err
			case '[':
				i := 0
				for dec.More() {
					child := fmt.Sprintf("%s/%d", path, i)
					if err := walk(child); err != nil {
						return err
					}
					i++
				}
				_, err = dec.Token()
				return err
			}
		}
		return nil
	}
	return walk("")
}

func escapePointer(s string) string {
	s = strings.ReplaceAll(s, "~", "~0")
	s = strings.ReplaceAll(s, "/", "~1")
	return s
}

// CheckJSON rejects duplicate keys, trailing values, and non-finite numbers for
// all versioned input documents before their own typed schema decoding.
func CheckJSON(data []byte) error {
	if err := rejectDuplicateKeys(data); err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	var value any
	if err := dec.Decode(&value); err != nil {
		return err
	}
	if err := rejectTrailing(dec); err != nil {
		return err
	}
	return rejectNonFinite(data)
}
