// Package diag defines stable machine-readable diagnostics.
package diag

import (
	"fmt"
	"sort"
)

// Diagnostic is emitted on stdout for --json validate failures.
type Diagnostic struct {
	Code    string `json:"code"`
	Stage   string `json:"stage"`
	Pointer string `json:"pointer"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
	Source  string `json:"source,omitempty"` // project-relative file; omitted for single-file
}

func (d Diagnostic) Error() string {
	prefix := d.Code
	if d.Source != "" {
		prefix = fmt.Sprintf("%s %s", d.Code, d.Source)
	}
	if d.Pointer != "" {
		return fmt.Sprintf("%s %s: %s", prefix, d.Pointer, d.Message)
	}
	return fmt.Sprintf("%s: %s", prefix, d.Message)
}

// ErrorList is one or more diagnostics.
type ErrorList []Diagnostic

func (e ErrorList) Error() string {
	if len(e) == 0 {
		return "validation failed"
	}
	return e[0].Error()
}

// Sort orders diagnostics by source, then pointer, then code for stable output.
func Sort(e ErrorList) {
	sort.SliceStable(e, func(i, j int) bool {
		if e[i].Source != e[j].Source {
			return e[i].Source < e[j].Source
		}
		if e[i].Pointer != e[j].Pointer {
			return e[i].Pointer < e[j].Pointer
		}
		return e[i].Code < e[j].Code
	})
}
