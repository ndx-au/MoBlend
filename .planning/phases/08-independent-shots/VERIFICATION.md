# Phase 8 — verification evidence

**Date:** 2026-11-04  
**Status:** gate met (SHOT-01–03)

## Commands

```bash
go test ./internal/project -count=1 -run 'TestShot|TestThreeShot|TestDependency' -v
go test ./internal/runner -count=1 -timeout 20m -run TestConcurrentShotPreviewsIsolated -v
go run ./cmd/blendctl validate examples/projects/three_shots/project.json
# → SHOT_REQUIRED
go run ./cmd/blendctl validate examples/projects/three_shots/project.json --shot intro
go run ./cmd/blendctl validate testdata/project/shots/broken_unselected/project.json --shot good
go run ./cmd/blendctl validate testdata/project/shots/broken_unselected/project.json --all-shots
# → fails on bad shot missing module
go run ./cmd/blendctl compile examples/projects/three_shots/project.json --shot intro --out /tmp/moblend-p8-intro
```

## Results

| Check | Result |
|---|---|
| SHOT-01 selection + closure | `TestShotSelectionAndAllShots`; good shot omits `bad.entry.json`; multi-shot needs `--shot`; `--all-shots` surfaces broken shot |
| SHOT-02 isolation + provenance | `dependencies.json` via `TestDependencyManifest`; `TestConcurrentShotPreviewsIsolated` distinct job dirs |
| SHOT-03 change isolation | `TestThreeShotChangeIsolation` — outro edit leaves intro/product hashes; shared edit changes dependents |
| Example | `examples/projects/three_shots/` (intro/product/outro + shared) |

## Caveats

- No incremental cache, scheduler, or video assembly (deferred).
- Model evaluation and non-Linux platforms remain unclaimed.
