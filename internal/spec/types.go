package spec

// SchemaVersion is the only accepted declaration schema version in v1.
const SchemaVersion = "1.0.0"

// AssetPackVersion is the pack version declarations must target in Phase 1.
const AssetPackVersion = "0.1.0"

// Limits mirror CONTEXT.md — reject, never clamp.
const (
	MaxDeclarationBytes = 256 * 1024
	MaxObjects          = 32
	MaxLights           = 8
	MaxTextRunes        = 256
	MaxFrameSpan        = 600
	MaxResolutionEdge   = 3840
	MaxSamples          = 4096
	MaxTracks           = 64
	MaxKeysPerTrack     = 64
)

// Declaration is the Go input model for a scene JSON document.
type Declaration struct {
	SchemaVersion    string   `json:"schema_version"`
	AssetPackVersion string   `json:"asset_pack_version"`
	Scene            Scene    `json:"scene"`
	Camera           Camera   `json:"camera"`
	Lights           []Light  `json:"lights,omitempty"`
	Objects          []Object `json:"objects"`
	Animation        []Track  `json:"animation,omitempty"`
}

// Scene holds render/timeline settings.
type Scene struct {
	FrameStart int    `json:"frame_start"`
	FrameEnd   int    `json:"frame_end"`
	FPS        int    `json:"fps"`
	Resolution [2]int `json:"resolution"`
	Samples    int    `json:"samples"`
	WorldColor string `json:"world_color,omitempty"`
}

// Camera is the single active camera.
type Camera struct {
	ID               string     `json:"id"`
	Location         [3]float64 `json:"location"`
	RotationEulerDeg [3]float64 `json:"rotation_euler_deg"`
	FocalLengthMM    *float64   `json:"focal_length_mm,omitempty"`
}

// Light is a scene light.
type Light struct {
	ID               string      `json:"id"`
	Type             string      `json:"type"`
	Location         [3]float64  `json:"location"`
	RotationEulerDeg *[3]float64 `json:"rotation_euler_deg,omitempty"`
	Energy           float64     `json:"energy"`
	Color            string      `json:"color,omitempty"`
	Size             *float64    `json:"size,omitempty"`
}

// Material is a Principled subset.
type Material struct {
	BaseColor string   `json:"base_color"`
	Roughness *float64 `json:"roughness,omitempty"`
}

// Object is a text, form, or catalogue asset instance.
type Object struct {
	ID               string         `json:"id"`
	Type             string         `json:"type"`
	Location         *[3]float64    `json:"location,omitempty"`
	RotationEulerDeg *[3]float64    `json:"rotation_euler_deg,omitempty"`
	Scale            *[3]float64    `json:"scale,omitempty"`
	Material         *Material      `json:"material,omitempty"`
	Text             string         `json:"text,omitempty"`
	Form             string         `json:"form,omitempty"`
	AssetID          string         `json:"asset_id,omitempty"`
	Params           map[string]any `json:"params,omitempty"`
}

// Track is a bounded keyframe track.
type Track struct {
	TargetID      string `json:"target_id"`
	Property      string `json:"property"`
	ParamName     string `json:"param_name,omitempty"`
	Interpolation string `json:"interpolation,omitempty"`
	Keys          []Key  `json:"keys"`
}

// Key is a single keyframe.
type Key struct {
	Frame int `json:"frame"`
	Value any `json:"value"`
}
