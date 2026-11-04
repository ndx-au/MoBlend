# Phase 2 — Deliver the first useful composition

**Status:** complete (gate met 2026-11-04; see VERIFICATION.md)  
**Depends on:** Phase 1  
**Requirements:** CORE-03, CORE-04, SCENE-01 (title path)

## Outcome

JSON title scene → deterministic compile artifacts → editable Blender scene → still image. Compilation is byte-stable and requires no Blender runtime.

## Tasks

### T1 — Typed ops + emitter
- **Owns:** `internal/compile/`
- **Does:** Lower normalized declaration to ops; emit trusted Python entry + canonical JSON data (user strings only in data file)
- **Check:** Repeat compile equal bytes; malicious/Unicode text round-trip in data file

### T2 — Bindings
- **Owns:** templates under `internal/resources/templates/`
- **Does:** Text, forms, materials, transforms, camera, lights via data API; avoid selection dependence
- **Check:** Probe asserts object IDs, transforms, materials, camera on saved scene

### T3 — `blendctl compile` + pristine pack jobs
- **Owns:** CLI + runner integration for compile-only and execute
- **Does:** Write readable artifacts; every execution starts from extracted pristine pack
- **Check:** Title example renders; text/colour/layout changes visible

### T4 — Fixtures / goldens
- **Owns:** `testdata/compile/`
- **Does:** Golden script/data stability; missing bindings; saved-scene property checks
- **Check:** Unit + Blender integration

## Gate

Title example renders; changed text/colour/layout visible; probe verifies IDs/transforms/materials/camera; compile byte-stable without Blender.

## Font decision

Ship or document a bundled font path for text objects (Phase 2 blocker for text). Prefer a redistributable font under `internal/resources/fonts/` with NOTICE; record in CONTEXT.
