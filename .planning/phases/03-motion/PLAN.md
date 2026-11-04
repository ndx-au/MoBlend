# Phase 3 — Add expressive, bounded motion

**Status:** complete (gate met 2026-11-04; see VERIFICATION.md)  
**Depends on:** Phase 2  
**Requirements:** SCENE-02, SCENE-03, ASSET-01 (catalogue bindings), QA seed for three refs

## Outcome

Title reveal, geometric pedestal, and column chart declarations animate; chart uses a first-party Geometry Nodes asset with published inputs; two instances stay independent.

## Tasks

### T1 — Animation tracks hardened
- Verify step/linear/ease_in_out on transforms and selected material values via Blender frame probes
- Reject unsupported targets / conflicting tracks (already in validate; add probe coverage)

### T2 — Column chart GN asset
- Author `assets/source/build_column_chart.py`; bake into pack (new pack version or additive object/node group)
- Publish socket identifiers in catalogue; bump `asset_pack_version` if pack bytes change
- Package font NOTICE already present

### T3 — Instance isolation
- Duplicate mutable datablocks; two chart instances with different params must not leak
- Probe socket identifiers on supported Blender 5.3

### T4 — Three reference compositions
- Finish title / pedestal / chart examples + small render fixtures
- Optional fixed compositor look only if cheap and Phase 0-compatible (else defer)

## Gate

Probes at representative frames within tolerances; two asset instances independent; three short animations render with correct frame counts; no arbitrary node graphs/rigs/simulations.
