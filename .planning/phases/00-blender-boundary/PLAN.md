# Phase 0 — Prove the Blender boundary

**Status:** complete (gate met 2026-11-04; see VERIFICATION.md)  
**Depends on:** none  
**Requirements:** ENV-01 (partial), ASSET-01 (partial), RUN-03 (Python failure path seed)

## Outcome

A Go program drives a real headless Blender render using an embedded minimal asset pack. Gate evidence is recorded before any schema expansion.

## Tasks (dependency order)

### T1 — Module, toolchain, harness
- **Owns:** `go.mod`, `go.sum`, `.github/workflows/ci.yml` (or minimal `scripts/check.sh`), `README.md` status update
- **Does:** Initialize module `github.com/ndx-au/MoBlend`; pin Go from local evidence; add `go test ./...` + `gofmt`/`govet` entry points
- **Check:** `go test ./...` passes with no Blender required for pure unit tests

### T2 — First-party base pack authoring
- **Owns:** `assets/source/build_base_pack.py`, `scripts/build-base-pack.sh`, generated `internal/resources/pack/base.blend` (+ checksum note)
- **Does:** Camera, light, simple object, known material, explicit Cycles CPU render/colour settings; save small `.blend`
- **Check:** Script runs under pinned Blender; pack file exists and is reproducible enough for checksum recording

### T3 — Embed, extract, smoke job
- **Owns:** `internal/resources/`, `internal/runner/`, `cmd/blendctl` smoke path (or `blendctl smoke`)
- **Does:** `go:embed` pack; extract to unique job dir; argv-spawn Blender with trusted smoke script; change one property; save scene; render one PNG
- **Check:** Integration test produces PNG with expected dimensions and nonempty pixels; saved `.blend` reopens

### T4 — Boundary hardening
- **Owns:** runner flags + tests under `internal/runner` / `testdata/phase0`
- **Does:** Verify argument order, `--python-exit-code` failure, CPU Cycles, config isolation (`--factory-startup`, `--disable-autoexec`), cancellation on Linux amd64; note fonts if relevant
- **Check:** Deliberate Python exception → nonzero exit; cancel kills child; offline run with Blender installed

### T5 — Evidence and environment table
- **Owns:** `CONTEXT.md`, `VERIFICATION.md`, `.planning/STATE.md`, supported-environment table
- **Does:** Record asset checksum, Go/Blender versions, exact invocation, outputs, constraints
- **Check:** Phase gate criteria below all have command evidence; STATE next action points to Phase 1 only if gate met

## Completion conditions (gate)

1. Real PNG with expected dimensions and nonempty content
2. Saved scene reopens under Blender
3. Deliberate Python exception returns failure
4. Executable/tooling works offline with Blender installed
5. Versions, checksum, and supported-environment table recorded

## Out of scope

Schema/JSON language, catalogue, Geometry Nodes chart asset, CLI `validate`/`compile`/`render` product surfaces beyond smoke proof.