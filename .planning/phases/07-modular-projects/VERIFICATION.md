# Phase 7 — verification evidence

**Date:** 2026-11-04  
**Status:** gate met (MOD-01–07)

## Commands

```bash
go test ./internal/project -count=1 -v
go test ./internal/runner -count=1 -timeout 20m -run TestModularFlatPreviewEquivalence -v
go vet ./...
go test ./... -count=1 -timeout 20m
go run ./cmd/blendctl schema project
go run ./cmd/blendctl schema module
go run ./cmd/blendctl validate examples/projects/title_modular/project.json
go run ./cmd/blendctl inspect examples/projects/title_modular/project.json
go run ./cmd/blendctl compile examples/projects/title_modular/project.json --out /tmp/moblend-p7-compile
go run ./cmd/blendctl preview examples/projects/title_modular/project.json --frame 1 --out /tmp/moblend-p7-preview
go run ./cmd/blendctl guide | head
```

## Results

| Check | Result |
|---|---|
| MOD-01 project/module schemas + single-file unchanged | `schema/project_v1.json`, `schema/module_v1.json` embedded; scene schema `1.0.0` retained |
| MOD-02 cycles / missing / escape / symlink | `TestRejectCycle`, `TestRejectMissing`, `TestRejectEscape`, `TestRejectSymlinkEscape` |
| MOD-03 namespaces + import-order independence | qualification `__`; `TestDeterminismImportReorderAndRelocate` |
| MOD-04 source-aware diagnostics | `TestSourceAwareDiagnostic` → `source=prop.module.json` |
| MOD-05 modular≡flat + aggregate limits | `TestModularMatchesFlat`; `TestAggregateObjectLimit` → `LIM_OBJECTS` |
| MOD-06 CLI + Blender | validate/inspect/compile/preview OK; `TestModularFlatPreviewEquivalence` PASS |
| MOD-07 guide | embedded `agent/GUIDE.md` documents ownership/imports/namespaces |

## Caveats

- LLM authoring evaluation remains blocked / not claimed.
- macOS/Windows not claimed.
- Pixel identity across GPUs/OS/Blender builds is not promised; equivalence is compile bytes + successful preview artifacts under pinned Linux amd64 Blender 5.3.
