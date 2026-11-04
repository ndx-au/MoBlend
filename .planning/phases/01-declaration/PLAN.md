# Phase 1 — Define the smallest useful declaration

**Status:** complete (gate met 2026-11-04; see VERIFICATION.md)  
**Depends on:** Phase 0  
**Requirements:** CORE-01, CORE-02 (declaration side); ASSET-01 catalogue contract seed; UX-01 diagnostic shape seed

## Outcome

An executable language contract: schema v1, Go types, strict parse/validate/normalize, catalogue contract, and local CLI surfaces (`schema`, `assets list/describe`, `validate`) that never launch Blender.

## Tasks (dependency order)

### T1 — Freeze semantics + schema/catalogue
- **Owns:** `internal/resources/schema/v1.json`, `internal/resources/catalogue/v1.json`, `CONTEXT.md` limits table
- **Does:** Units, sRGB hex colours, integer frames/FPS, defaults, limits, diagnostic shape; fields only for title / pedestal / chart examples
- **Check:** Schema documents versions `schema_version` + `asset_pack_version`

### T2 — Go model + strict decode
- **Owns:** `internal/spec/`
- **Does:** Input types; reject unknown fields, duplicate keys, trailing data, non-finite numbers; version gate
- **Check:** Unit tests for each rejection class

### T3 — Semantic validation + normalization
- **Owns:** `internal/validate/`, normalization in `internal/spec/`
- **Does:** Refs, duplicate IDs, limits, catalogue asset IDs, stable ordering / defaults
- **Check:** Positive fixtures validate; negatives return stable `code` + JSON pointer

### T4 — Fixtures before capability creep
- **Owns:** `testdata/spec/positive/`, `testdata/spec/negative/`, three example stubs under `examples/`
- **Does:** Author fixtures first; implement only what they require
- **Check:** Table-driven tests cover every fixture

### T5 — CLI surfaces
- **Owns:** `cmd/blendctl` subcommands
- **Does:** `schema`, `assets list|describe`, `validate [--json]` — no Blender spawn
- **Check:** Manual + tests; `validate` works offline without Blender binary

## Completion conditions (gate)

1. All positive fixtures validate
2. Every negative fixture yields expected stable code and JSON pointer
3. Schema/model parity check automated
4. Deferred fields listed in CONTEXT

## Out of scope

Python emission, Blender compile/render, Geometry Nodes chart authoring, animation evaluation (declare tracks only if examples need them; full motion is Phase 3).
