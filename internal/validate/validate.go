package validate

import (
	"fmt"
	"math"
	"regexp"
	"unicode/utf8"

	"github.com/ndx-au/MoBlend/internal/diag"
	"github.com/ndx-au/MoBlend/internal/resources"
	"github.com/ndx-au/MoBlend/internal/spec"
)

var (
	idRe  = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
	hexRe = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
)

// Result is a normalized declaration plus any non-fatal notes (none in v1).
type Result struct {
	Decl spec.Declaration
}

// ParseAndValidate parses JSON and validates against the embedded catalogue.
// Does not launch Blender.
func ParseAndValidate(data []byte) (Result, error) {
	d, err := spec.Parse(data)
	if err != nil {
		return Result{}, err
	}
	cat, err := resources.LoadCatalogue()
	if err != nil {
		return Result{}, err
	}
	return Validate(d, cat)
}

// Validate checks versions, structure semantics, limits, and catalogue refs.
func Validate(d spec.Declaration, cat resources.Catalogue) (Result, error) {
	var errs diag.ErrorList

	if d.SchemaVersion != spec.SchemaVersion {
		errs = append(errs, diag.Diagnostic{
			Code:    "SPEC_SCHEMA_VERSION",
			Stage:   "schema",
			Pointer: "/schema_version",
			Message: fmt.Sprintf("unsupported schema_version %q", d.SchemaVersion),
			Hint:    fmt.Sprintf("use %q", spec.SchemaVersion),
		})
	}
	if d.AssetPackVersion != spec.AssetPackVersion {
		errs = append(errs, diag.Diagnostic{
			Code:    "SPEC_ASSET_PACK_VERSION",
			Stage:   "schema",
			Pointer: "/asset_pack_version",
			Message: fmt.Sprintf("unsupported asset_pack_version %q", d.AssetPackVersion),
			Hint:    fmt.Sprintf("use %q", spec.AssetPackVersion),
		})
	}

	errs = append(errs, checkScene(d)...)
	errs = append(errs, checkCamera(d)...)
	errs = append(errs, checkLights(d)...)
	errs = append(errs, checkObjects(d, cat)...)
	errs = append(errs, checkAnimation(d)...)

	if len(errs) > 0 {
		return Result{}, errs
	}
	return Result{Decl: normalize(d)}, nil
}

func checkScene(d spec.Declaration) diag.ErrorList {
	var errs diag.ErrorList
	s := d.Scene
	if s.FrameEnd < s.FrameStart {
		errs = append(errs, diag.Diagnostic{
			Code: "SEM_FRAME_RANGE", Stage: "semantic", Pointer: "/scene/frame_end",
			Message: "frame_end must be >= frame_start",
		})
	}
	span := s.FrameEnd - s.FrameStart + 1
	if span > spec.MaxFrameSpan {
		errs = append(errs, diag.Diagnostic{
			Code: "LIM_FRAME_SPAN", Stage: "semantic", Pointer: "/scene/frame_end",
			Message: fmt.Sprintf("frame span %d exceeds %d", span, spec.MaxFrameSpan),
		})
	}
	if s.FPS < 1 || s.FPS > 120 {
		errs = append(errs, diag.Diagnostic{
			Code: "SEM_FPS", Stage: "semantic", Pointer: "/scene/fps",
			Message: "fps out of range",
		})
	}
	for i, edge := range s.Resolution {
		if edge < 16 || edge > spec.MaxResolutionEdge {
			errs = append(errs, diag.Diagnostic{
				Code: "LIM_RESOLUTION", Stage: "semantic",
				Pointer: fmt.Sprintf("/scene/resolution/%d", i),
				Message: "resolution edge out of limits",
			})
		}
	}
	if s.Samples < 1 || s.Samples > spec.MaxSamples {
		errs = append(errs, diag.Diagnostic{
			Code: "LIM_SAMPLES", Stage: "semantic", Pointer: "/scene/samples",
			Message: "samples out of limits",
		})
	}
	if s.WorldColor != "" && !hexRe.MatchString(s.WorldColor) {
		errs = append(errs, diag.Diagnostic{
			Code: "SEM_COLOR", Stage: "semantic", Pointer: "/scene/world_color",
			Message: "world_color must be #RRGGBB",
		})
	}
	return errs
}

func checkCamera(d spec.Declaration) diag.ErrorList {
	var errs diag.ErrorList
	if !idRe.MatchString(d.Camera.ID) {
		errs = append(errs, diag.Diagnostic{
			Code: "SEM_ID", Stage: "semantic", Pointer: "/camera/id",
			Message: "invalid id",
		})
	}
	if d.Camera.FocalLengthMM != nil {
		f := *d.Camera.FocalLengthMM
		if f <= 0 || f > 500 || math.IsNaN(f) || math.IsInf(f, 0) {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_CAMERA_FOCAL", Stage: "semantic", Pointer: "/camera/focal_length_mm",
				Message: "focal_length_mm out of range",
			})
		}
	}
	return errs
}

func checkLights(d spec.Declaration) diag.ErrorList {
	var errs diag.ErrorList
	if len(d.Lights) > spec.MaxLights {
		errs = append(errs, diag.Diagnostic{
			Code: "LIM_LIGHTS", Stage: "semantic", Pointer: "/lights",
			Message: fmt.Sprintf("more than %d lights", spec.MaxLights),
		})
	}
	seen := map[string]struct{}{d.Camera.ID: {}}
	for i, l := range d.Lights {
		p := fmt.Sprintf("/lights/%d", i)
		if !idRe.MatchString(l.ID) {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_ID", Stage: "semantic", Pointer: p + "/id", Message: "invalid id",
			})
		}
		if _, dup := seen[l.ID]; dup {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_DUPLICATE_ID", Stage: "semantic", Pointer: p + "/id",
				Message: fmt.Sprintf("duplicate id %q", l.ID),
			})
		}
		seen[l.ID] = struct{}{}
		switch l.Type {
		case "AREA", "SUN", "POINT":
		default:
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_LIGHT_TYPE", Stage: "semantic", Pointer: p + "/type",
				Message: "unsupported light type",
			})
		}
		if l.Energy < 0 {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_LIGHT_ENERGY", Stage: "semantic", Pointer: p + "/energy",
				Message: "energy must be >= 0",
			})
		}
		if l.Color != "" && !hexRe.MatchString(l.Color) {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_COLOR", Stage: "semantic", Pointer: p + "/color",
				Message: "color must be #RRGGBB",
			})
		}
	}
	return errs
}

func checkObjects(d spec.Declaration, cat resources.Catalogue) diag.ErrorList {
	var errs diag.ErrorList
	if len(d.Objects) == 0 {
		errs = append(errs, diag.Diagnostic{
			Code: "SEM_OBJECTS_EMPTY", Stage: "semantic", Pointer: "/objects",
			Message: "at least one object is required",
		})
	}
	if len(d.Objects) > spec.MaxObjects {
		errs = append(errs, diag.Diagnostic{
			Code: "LIM_OBJECTS", Stage: "semantic", Pointer: "/objects",
			Message: fmt.Sprintf("more than %d objects", spec.MaxObjects),
		})
	}
	seen := map[string]struct{}{d.Camera.ID: {}}
	for _, l := range d.Lights {
		seen[l.ID] = struct{}{}
	}
	for i, o := range d.Objects {
		p := fmt.Sprintf("/objects/%d", i)
		if !idRe.MatchString(o.ID) {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_ID", Stage: "semantic", Pointer: p + "/id", Message: "invalid id",
			})
		}
		if _, dup := seen[o.ID]; dup {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_DUPLICATE_ID", Stage: "semantic", Pointer: p + "/id",
				Message: fmt.Sprintf("duplicate id %q", o.ID),
			})
		}
		seen[o.ID] = struct{}{}

		switch o.Type {
		case "text":
			if o.Text == "" {
				errs = append(errs, diag.Diagnostic{
					Code: "SEM_TEXT_REQUIRED", Stage: "semantic", Pointer: p + "/text",
					Message: "text objects require text",
				})
			}
			if utf8.RuneCountInString(o.Text) > spec.MaxTextRunes {
				errs = append(errs, diag.Diagnostic{
					Code: "LIM_TEXT_LENGTH", Stage: "semantic", Pointer: p + "/text",
					Message: "text exceeds max length",
				})
			}
		case "form":
			switch o.Form {
			case "cube", "plane", "cylinder":
			default:
				errs = append(errs, diag.Diagnostic{
					Code: "SEM_FORM", Stage: "semantic", Pointer: p + "/form",
					Message: "unsupported form",
				})
			}
		case "asset":
			if o.AssetID == "" {
				errs = append(errs, diag.Diagnostic{
					Code: "SEM_ASSET_REQUIRED", Stage: "semantic", Pointer: p + "/asset_id",
					Message: "asset objects require asset_id",
				})
				break
			}
			if _, ok := cat.ByID[o.AssetID]; !ok {
				errs = append(errs, diag.Diagnostic{
					Code: "CAT_UNKNOWN_ASSET", Stage: "semantic", Pointer: p + "/asset_id",
					Message: fmt.Sprintf("unknown asset_id %q", o.AssetID),
					Hint:    "use blendctl assets list",
				})
			}
			// reserved_phase3 assets (e.g. column_chart) may appear in declarations;
			// pack binding is enforced at compile/render time in later phases.
		default:
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_OBJECT_TYPE", Stage: "semantic", Pointer: p + "/type",
				Message: "unsupported object type",
			})
		}
		if o.Material != nil && !hexRe.MatchString(o.Material.BaseColor) {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_COLOR", Stage: "semantic", Pointer: p + "/material/base_color",
				Message: "base_color must be #RRGGBB",
			})
		}
		if o.Material != nil && o.Material.Roughness != nil {
			r := *o.Material.Roughness
			if r < 0 || r > 1 {
				errs = append(errs, diag.Diagnostic{
					Code: "SEM_ROUGHNESS", Stage: "semantic", Pointer: p + "/material/roughness",
					Message: "roughness must be in [0,1]",
				})
			}
		}
	}
	return errs
}

func checkAnimation(d spec.Declaration) diag.ErrorList {
	var errs diag.ErrorList
	if len(d.Animation) > spec.MaxTracks {
		errs = append(errs, diag.Diagnostic{
			Code: "LIM_TRACKS", Stage: "semantic", Pointer: "/animation",
			Message: "too many animation tracks",
		})
	}
	ids := map[string]struct{}{d.Camera.ID: {}}
	for _, l := range d.Lights {
		ids[l.ID] = struct{}{}
	}
	for _, o := range d.Objects {
		ids[o.ID] = struct{}{}
	}
	type trackKey struct{ target, prop, param string }
	seenTracks := map[trackKey]struct{}{}
	for i, t := range d.Animation {
		p := fmt.Sprintf("/animation/%d", i)
		if _, ok := ids[t.TargetID]; !ok {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_TRACK_TARGET", Stage: "semantic", Pointer: p + "/target_id",
				Message: fmt.Sprintf("unknown target_id %q", t.TargetID),
			})
		}
		switch t.Property {
		case "location", "rotation_euler_deg", "scale", "material.base_color", "material.roughness", "param":
		default:
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_TRACK_PROPERTY", Stage: "semantic", Pointer: p + "/property",
				Message: "unsupported animation property",
			})
		}
		if t.Property == "param" && t.ParamName == "" {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_TRACK_PARAM", Stage: "semantic", Pointer: p + "/param_name",
				Message: "param tracks require param_name",
			})
		}
		interp := t.Interpolation
		if interp == "" {
			interp = "linear"
		}
		switch interp {
		case "step", "linear", "ease_in_out":
		default:
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_INTERPOLATION", Stage: "semantic", Pointer: p + "/interpolation",
				Message: "unsupported interpolation",
			})
		}
		if len(t.Keys) == 0 || len(t.Keys) > spec.MaxKeysPerTrack {
			errs = append(errs, diag.Diagnostic{
				Code: "LIM_KEYS", Stage: "semantic", Pointer: p + "/keys",
				Message: "key count out of limits",
			})
		}
		frames := map[int]struct{}{}
		for j, k := range t.Keys {
			kp := fmt.Sprintf("%s/keys/%d", p, j)
			if k.Frame < d.Scene.FrameStart || k.Frame > d.Scene.FrameEnd {
				errs = append(errs, diag.Diagnostic{
					Code: "SEM_KEY_FRAME", Stage: "semantic", Pointer: kp + "/frame",
					Message: "key frame outside scene range",
				})
			}
			if _, dup := frames[k.Frame]; dup {
				errs = append(errs, diag.Diagnostic{
					Code: "SEM_DUPLICATE_KEYFRAME", Stage: "semantic", Pointer: kp + "/frame",
					Message: "duplicate keyframe time on track",
				})
			}
			frames[k.Frame] = struct{}{}
		}
		tk := trackKey{t.TargetID, t.Property, t.ParamName}
		if _, dup := seenTracks[tk]; dup {
			errs = append(errs, diag.Diagnostic{
				Code: "SEM_CONFLICTING_TRACK", Stage: "semantic", Pointer: p,
				Message: "conflicting track for same target/property",
			})
		}
		seenTracks[tk] = struct{}{}
	}
	return errs
}

func normalize(d spec.Declaration) spec.Declaration {
	if d.Scene.WorldColor == "" {
		d.Scene.WorldColor = "#1A1A1A"
	}
	if d.Camera.FocalLengthMM == nil {
		f := 50.0
		d.Camera.FocalLengthMM = &f
	}
	for i := range d.Lights {
		if d.Lights[i].Color == "" {
			d.Lights[i].Color = "#FFFFFF"
		}
	}
	for i := range d.Objects {
		o := &d.Objects[i]
		if o.Location == nil {
			v := [3]float64{0, 0, 0}
			o.Location = &v
		}
		if o.RotationEulerDeg == nil {
			v := [3]float64{0, 0, 0}
			o.RotationEulerDeg = &v
		}
		if o.Scale == nil {
			v := [3]float64{1, 1, 1}
			o.Scale = &v
		}
		if o.Material == nil {
			r := 0.4
			o.Material = &spec.Material{BaseColor: "#4A90D9", Roughness: &r}
		} else if o.Material.Roughness == nil {
			r := 0.4
			o.Material.Roughness = &r
		}
	}
	for i := range d.Animation {
		if d.Animation[i].Interpolation == "" {
			d.Animation[i].Interpolation = "linear"
		}
	}
	return d
}
