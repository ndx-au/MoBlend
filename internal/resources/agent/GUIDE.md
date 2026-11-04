# MoBlend agent guide (v1)

Local-only authoring. There is **no RAG**, embeddings, or remote retrieval. Use exact asset IDs from `blendctl assets list`.

## Product

Write a JSON scene declaration → `validate` → `preview`/`render`. The Go CLI compiles deterministically and drives Blender.

## Exact commands

```bash
blendctl doctor
blendctl schema
blendctl assets list
blendctl assets describe base_cube
blendctl assets describe column_chart
blendctl validate scene.json --json
blendctl compile scene.json --out build/
blendctl preview scene.json --frame 24 --out preview/
blendctl render scene.json --out render/
```

Environment: set `MOBLEND_BLENDER` if Blender is not on `PATH`. Supported family: Blender **5.3**.

## Declaration essentials

- Required: `schema_version` (`1.0.0`), `asset_pack_version` (`0.1.0`), `scene`, `camera`, `objects`
- Units: metres, degrees at JSON boundary, colours `#RRGGBB`
- Object types: `text`, `form` (`cube|plane|cylinder`), `asset` (`base_cube`, `column_chart`)
- Animation: tracks with `step|linear|ease_in_out` on approved properties
- Limits: see schema / Phase 1 CONTEXT (max 32 objects, 600 frames, etc.)

## Worked examples

- `examples/title.json` — text title reveal
- `examples/pedestal.json` — form + `base_cube` product
- `examples/chart.json` — `column_chart` with `params.values` (0..1) and `bar_count`

## Common diagnostics

| Code | Meaning | Fix |
|---|---|---|
| `SPEC_DUPLICATE_KEY` | Duplicate JSON key | Remove duplicate |
| `SPEC_UNKNOWN_FIELD` | Extra field | Remove unknown keys |
| `SPEC_SCHEMA_VERSION` | Wrong schema | Use `1.0.0` |
| `CAT_UNKNOWN_ASSET` | Bad `asset_id` | `assets list` |
| `SEM_DUPLICATE_ID` | ID clash | Unique ids across camera/lights/objects |
| `ENV_BLENDER_*` | Doctor/preflight | Install/pin Blender 5.3 |

Machine diagnostics go to stdout with `--json`. Progress/logs go to stderr.

## Unsupported in v1 (do not invent fallbacks)

- Arbitrary Python / bpy paths / user `.blend` loading
- RAG / fuzzy asset search / marketplace downloads
- Physics, rigs, Grease Pencil authoring, custom node graphs
- Video/FFmpeg encoding, MCP server, GUI
- Multiple cameras, parenting hierarchies, HDRI worlds beyond `world_color`

## Repair loop

1. `validate --json` and fix pointers/codes
2. `preview --frame N` and inspect PNG + `scene.blend`
3. Adjust text/colours/layout/params; re-validate
4. `render` for full frame range; check `manifest.json` status `complete`
