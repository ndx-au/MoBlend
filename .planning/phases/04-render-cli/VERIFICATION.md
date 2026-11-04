# Phase 4 — verification evidence

**Date:** 2026-11-04  
**Status:** gate met (Linux amd64)

## Commands

```bash
go test -a ./... -count=1 -timeout 15m
go run ./cmd/blendctl doctor
go run ./cmd/blendctl doctor --json
go run ./cmd/blendctl render testdata/compile/title_render_tiny.json --out /tmp/moblend-render
```

## Gate checks

| Gate | Result | Evidence |
|---|---|---|
| doctor succeeds on pinned build | PASS | `TestDoctorOK` |
| missing Blender fails | PASS | `TestDoctorMissing` |
| preview/render outputs + manifest | PASS | `TestRenderWritesManifestAndFrames` (`status=complete`, 2 frames) |
| failure marker not success | PASS | success path has no `manifest.failed.json`; failures write `manifest.failed.json` |
| pack not mutated | PASS | extract-per-job; embedded bytes unchanged (`pack-sha256`) |

## Exit codes

- doctor environment failure → exit 3
- validate/render execution failure → exit 1
- usage → exit 2


## Follow-up fixes — 2026-11-05

Version probes and execution now share process-tree cancellation, isolated environment and bounded logs. Preview/frame PNGs decode completely with matching dimensions before success. Tests cover probe cancel/timeout descendants, noisy output, malformed previews, and concurrent shot success/failure attribution. `go test ./... -count=1 -timeout 10m` passed 134 tests/subtests with real Blender; scoped race suite passed 40. See [follow-up verification](../../fixes/gauntlet-followup/VERIFICATION.md) for exact commands and remaining caveats.
