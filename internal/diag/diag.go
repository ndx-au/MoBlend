// Package diag defines stable machine-readable diagnostics.
package diag

import "fmt"

// Diagnostic is emitted on stdout for --json validate failures.
type Diagnostic struct {
	Code    string `json:"code"`
	Stage   string `json:"stage"`
	Pointer string `json:"pointer"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
}

func (d Diagnostic) Error() string {
	if d.Pointer != "" {
		return fmt.Sprintf("%s %s: %s", d.Code, d.Pointer, d.Message)
	}
	return fmt.Sprintf("%s: %s", d.Code, d.Message)
}

// ErrorList is one or more diagnostics.
type ErrorList []Diagnostic

func (e ErrorList) Error() string {
	if len(e) == 0 {
		return "validation failed"
	}
	return e[0].Error()
}
