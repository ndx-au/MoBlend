# Phase 2 — verification evidence

**Date:** 2026-11-04  
**Status:** gate met (Linux amd64 / Blender 5.3.0 Beta)

## Commands

```bash
gofmt -l -w .
go vet ./...
go test ./... -count=1 -timeout 10m
go build -o /tmp/blendctl ./cmd/blendctl
/tmp/blendctl compile examples/title.json --out /tmp/moblend-title-build
/tmp/blendctl preview testdata/compile/title_small.json --frame 24 --out /tmp/moblend-preview
```

## Gate checks

| Gate | Result | Evidence |
|---|---|---|
| Title renders | PASS | `TestApplyTitleSmall`; CLI preview PNG nonempty |
| Text/colour/layout in scene | PASS | Probe asserts FONT body `MoBlend`, camera/title IDs present |
| Compile byte-stable, no Blender | PASS | `TestCompileByteStableAndSafeText` |
| User text not in script | PASS | awkward Unicode/quotes only in `scene_data.json` |
| Pristine pack each job | PASS | `RunApply` extracts pack+font per job dir |

## Notes

- Bundled font: `internal/resources/fonts/DejaVuSans.ttf` + `DejaVuSans.LICENSE`
- Blender 5.3 layered Actions: interpolation via `action.layers[].strips[].channelbags[].fcurves` (see `iter_action_fcurves`)
- Full `examples/title.json` (1280×720, 32 samples) is slower; CI/dev gate uses `testdata/compile/title_small.json`
- Pedestal/chart compile-bind beyond title is Phase 3+
