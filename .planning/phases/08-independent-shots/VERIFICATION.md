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


## 2026-11-05 Gauntlet repair addendum

SHOT-01 all-shots behavior now has a real CLI-path regression attributing the broken bad shot. SHOT-02 compile outputs use unique shot-nested jobs and hashes come from load-time bytes, including the manifest. CLI signal/process-group cancellation and exact render frame sets are covered. See [repair verification](../../fixes/gauntlet-round-1/VERIFICATION.md) for exact commands, fresh results, and remaining medium/low gaps. Full suite: `go test ./... -count=1 -timeout 10m` passed 103 tests/subtests with real Blender; project/CLI race checks passed 42. Historical evidence above describes the earlier implementation.


## Follow-up fixes — 2026-11-05

Concurrent shot previews now include an intentional failed shot, with independent success/failure markers, shot identity and dependency attribution. Version probes share tested process-group cancellation and bounded logs. Full Go/Blender suite: 134 tests/subtests passed; scoped race suite: 40 passed. See [follow-up verification](../../fixes/gauntlet-followup/VERIFICATION.md) for commands, complete finding mapping and remaining release/evaluation caveats.
