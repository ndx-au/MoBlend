package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// SupportedBlenderFamily is the major.minor gate for Phase 4 doctor.
const SupportedBlenderFamily = "5.3"

var versionLineRe = regexp.MustCompile(`(?m)^Blender\s+(\d+\.\d+(?:\.\d+)?(?:\s+\w+)?)`)

// DoctorReport is machine-readable environment status.
type DoctorReport struct {
	OK          bool   `json:"ok"`
	BlenderPath string `json:"blender_path,omitempty"`
	VersionLine string `json:"version_line,omitempty"`
	Family      string `json:"family_required"`
	ErrorCode   string `json:"error_code,omitempty"`
	Message     string `json:"message,omitempty"`
	Hint        string `json:"hint,omitempty"`
}

// Doctor resolves Blender and checks the supported version family.
func Doctor(ctx context.Context, blender string) (DoctorReport, error) {
	rep := DoctorReport{Family: SupportedBlenderFamily}
	var err error
	if blender == "" {
		blender, err = ResolveBlender()
		if err != nil {
			rep.ErrorCode = "ENV_BLENDER_MISSING"
			rep.Message = err.Error()
			rep.Hint = "install Blender or set MOBLEND_BLENDER"
			return rep, fmt.Errorf("%s", rep.Message)
		}
	}
	rep.BlenderPath = blender

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, blender, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		rep.ErrorCode = "ENV_BLENDER_UNUSABLE"
		rep.Message = fmt.Sprintf("blender --version failed: %v", err)
		return rep, fmt.Errorf("%s", rep.Message)
	}
	text := string(out)
	rep.VersionLine = firstLine(text)
	family := parseFamily(text)
	if family == "" {
		rep.ErrorCode = "ENV_BLENDER_VERSION_PARSE"
		rep.Message = "could not parse Blender version"
		return rep, fmt.Errorf("%s", rep.Message)
	}
	if family != SupportedBlenderFamily {
		rep.ErrorCode = "ENV_BLENDER_UNSUPPORTED"
		rep.Message = fmt.Sprintf("Blender family %s not supported (need %s)", family, SupportedBlenderFamily)
		rep.Hint = "install the pinned Blender build documented in DEVELOPMENT-ENVIRONMENT.md"
		return rep, fmt.Errorf("%s", rep.Message)
	}
	rep.OK = true
	return rep, nil
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return s
}

func parseFamily(versionOutput string) string {
	m := versionLineRe.FindStringSubmatch(versionOutput)
	if m == nil {
		return ""
	}
	// m[1] like "5.3.0 Beta" or "5.3.0"
	fields := strings.Fields(m[1])
	if len(fields) == 0 {
		return ""
	}
	parts := strings.Split(fields[0], ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[0] + "." + parts[1]
}

// DoctorJSON encodes a report.
func DoctorJSON(rep DoctorReport) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	_ = enc.Encode(rep)
	return buf.Bytes()
}
