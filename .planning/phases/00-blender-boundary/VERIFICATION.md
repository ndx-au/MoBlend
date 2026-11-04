# Phase 0 — verification evidence

**Date:** 2026-11-04  
**Host:** gandalf, Ubuntu Linux amd64  
**Status:** gate met (Linux amd64 / pinned beta Blender)

## Supported environment table

| Component | Version / identity | Notes |
|---|---|---|
| OS | Linux x86_64 (Ubuntu 26.04.1 LTS) | First integration target |
| Go | `go1.25.6` (`go` directive 1.25) | `go version` |
| Blender | 5.3.0 Beta, hash `7e7c173d45c8` | Via `MOBLEND_BLENDER` / PATH `blender` → beta install |
| Module | `github.com/ndx-au/MoBlend` | |
| Pack | `0.0.1-phase0` | Embedded `internal/resources/pack/base.blend` |
| Pack SHA-256 | `882acdcadfaa1d5e180355a7fe2b1ab1d09865d95e2d9a8d531a88e97da083f5` | `blendctl pack-sha256` / `base.blend.sha256` |
| Fonts | N/A for Phase 0 pack | No text objects in base pack |

Release Blender baseline remains open (beta is a local development candidate only).

## Commands run

```bash
scripts/build-base-pack.sh
gofmt -l -w .
go vet ./...
go test ./... -count=1 -timeout 10m
go build -o /tmp/blendctl ./cmd/blendctl
/tmp/blendctl version
/tmp/blendctl pack-sha256
/tmp/blendctl smoke --out /tmp/moblend-smoke-out
```

## Gate checks

| Gate | Result | Evidence |
|---|---|---|
| Real PNG, expected dims, nonempty | PASS | 64×64 PNG, 5280 bytes under job `out/smoke.png` |
| Saved scene reopens | PASS | `TestSmokeRender` probe loads `smoke.blend`, cube scale 1.5 |
| Deliberate Python exception → failure | PASS | `TestSmokePythonFailure` (`--fail`) |
| Offline with Blender installed | PASS | No network used; local Blender + embedded pack |
| Argv order | PASS | `TestBuildArgsOrder` |
| CPU Cycles + isolation flags | PASS | Smoke script sets Cycles CPU; runner always passes `--background --factory-startup --disable-autoexec --python-exit-code 1` |
| Cancellation | PASS | `TestSmokeCancel` |

## Example smoke JSON (abridged)

Job under `/tmp/moblend-smoke-out/moblend-job-*` produced `smoke.png`, `smoke.blend`, `result.json` with pack version `0.0.1-phase0`, engine Cycles, device CPU.

## Known constraints

- Blender 5.3 warns: `Material.use_nodes` expected removed in 6.0 — pack authoring still works; revisit before Blender 6.
- CI (`ubuntu-latest`) runs unit tests only; Blender integration Skips when binary missing. Local/dev must run `go test ./internal/runner` with Blender.
- Beta pin is not a release compatibility claim.

## Invocation contract (recorded)

```
blender --background --factory-startup --disable-autoexec --python-exit-code 1 \
  <extracted>/base.blend --python <job>/smoke.py -- \
  --out-dir <job>/out --scale 1.5 --width 64 --height 64
```
