# MoBlend

Declarative motion graphics for Blender, aimed at coding agents.

You (or an agent) write a small versioned **JSON** scene. **`blendctl`** (Go) validates it, compiles deterministic Blender Python against a **bundled asset pack**, and produces an editable `.blend`, PNG frames, and a completion manifest. No RAG, no hosted model, no separate Python install for users.

**Status:** Linux amd64 release candidate `0.1.0-rc1` (roadmap phases 0–6). Blender **5.3** family required. macOS/Windows not verified yet.

## Quick start

```bash
# Prerequisites: Go 1.25+, Blender 5.3 on PATH (or MOBLEND_BLENDER=/path/to/blender)

go test ./...
go run ./cmd/blendctl doctor
go run ./cmd/blendctl preview testdata/compile/title_small.json --frame 24 --out /tmp/moblend-out
```

Open `/tmp/moblend-out/moblend-job-*/out/preview.png` and `scene.blend`.

Release binary:

```bash
VERSION=0.1.0-rc1 scripts/release-linux-amd64.sh
./dist/linux-amd64/blendctl doctor
```

## CLI

| Command | Purpose |
|---|---|
| `doctor` | Check Blender path/version |
| `guide` | Print the local agent guide |
| `schema` | Print JSON Schema v1 |
| `assets list` / `assets describe <id>` | Exact-ID catalogue |
| `validate <scene.json> [--json]` | Strict parse + semantic checks (no Blender) |
| `compile <scene.json> --out DIR` | Emit trusted script + canonical data (no Blender) |
| `preview <scene.json> [--frame N] [--out DIR]` | One frame + `.blend` + `manifest.json` |
| `render <scene.json> [--out DIR]` | Full frame range under `frames/` |
| `smoke` | Pack/runner boundary smoke |

Set `MOBLEND_BLENDER` if Blender is not on `PATH`.

## Examples

| File | Composition |
|---|---|
| `examples/title.json` | Text title reveal |
| `examples/pedestal.json` | Cylinder plinth + `base_cube` product |
| `examples/chart.json` | `column_chart` Geometry Nodes bars |

Catalogue assets: `base_cube`, `column_chart`. Schema `1.0.0`, asset pack `0.1.0`.

Agents should start with `blendctl guide` (also `agent/GUIDE.md`). Use exact asset IDs only; repair from `validate --json` diagnostics.

## What v1 does not do

Arbitrary Python/`bpy`, user `.blend` import, RAG/fuzzy search, physics/rigs/custom node graphs, video/FFmpeg, MCP, or a GUI. See `docs/LIMITATIONS.md`.

## Development

```bash
scripts/check.sh                          # gofmt, vet, unit tests
go test ./internal/runner -count=1        # Blender integration (needs Blender)
scripts/build-base-pack.sh                # regenerate pack base
scripts/build-column-chart.sh             # add/update column chart GN
scripts/make-contact-sheet.sh             # title/pedestal/chart contact sheet
```

Layout: `cmd/blendctl`, `internal/{spec,validate,compile,runner,resources}`, embedded pack/fonts/schema under `internal/resources/`.

## Docs and project history

| Doc | Role |
|---|---|
| `docs/QUICKSTART.md` | First render |
| `docs/TROUBLESHOOTING.md` | Common failures |
| `docs/LIMITATIONS.md` | Compatibility boundaries |
| `AGENTS.md` | Agent execution rules |
| `.planning/` | Roadmap, phase plans, verification |
| `.progress/` | Append-only Dot Progress log |

## License notes

Bundled DejaVu Sans: see `internal/resources/fonts/DejaVuSans.LICENSE`. Public release publication still requires owner authorization.
