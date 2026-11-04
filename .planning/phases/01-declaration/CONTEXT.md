# Phase 1 — context and decisions

## Versions

| Field | v1 value |
|---|---|
| `schema_version` | `1.0.0` |
| `asset_pack_version` accepted | `0.0.1-phase0` (Phase 0 pack; chart asset ID reserved, may fail semantic bind until Phase 3) |

## Units and colour

- Distances: metres (Blender space, Z-up)
- Angles at JSON boundary: degrees; compiler converts to radians
- Colours: opaque `#RRGGBB` sRGB hex only (no alpha in v1 materials)
- Frames: inclusive integer `frame_start`/`frame_end`; integer `fps`

## Limits (reject, do not clamp)

| Limit | Value |
|---|---|
| Max declaration bytes | 256 KiB |
| Max objects | 32 |
| Max lights | 8 |
| Max text length | 256 runes |
| Max frames (`end-start+1`) | 600 |
| Max resolution edge | 3840 |
| Max samples | 4096 |
| Max animation tracks | 64 |
| Max keys per track | 64 |

## Diagnostic shape

```json
{
  "code": "SPEC_DUPLICATE_KEY",
  "stage": "parse|schema|semantic",
  "pointer": "/objects/0/id",
  "message": "human readable",
  "hint": "optional remediation"
}
```

Stable codes use `SPEC_*`, `SEM_*`, `CAT_*`, `LIM_*` prefixes.

## Catalogue seed

| asset_id | type | Phase availability |
|---|---|---|
| `base_cube` | mesh | Phase 0 pack object binding |
| `column_chart` | geometry_nodes | In catalogue for Phase 1 declaration fixtures; pack binding / availability enforced at compile in Phase 3 |

## Deferred fields (explicit)

- Parenting / collections hierarchy
- Custom shaders beyond Principled base colour + roughness
- Video/FFmpeg outputs
- Multiple cameras
- World/HDRI controls beyond a single world colour
- Compositor looks
- Easing beyond declaration of `interpolation` enum values needed by examples (evaluation Phase 3)

## JSON Schema dependency

Validate with a Go JSON Schema library **or** hand-written structural checks kept in lockstep with `schema/v1.json` via a parity test that walks required properties. Prefer embedding `schema/v1.json` and using `santhosh-tekuri/jsonschema/v6` if module fetch works; otherwise hand validation with parity checklist in tests.
