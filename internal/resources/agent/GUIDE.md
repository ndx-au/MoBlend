# MoBlend agent guide (v1)

Local-only authoring. There is **no RAG**, embeddings, or remote retrieval. Use exact asset IDs from `blendctl assets list`.

## Product

Write a JSON scene declaration (or a multi-file project) → `validate` → `preview`/`render`. The Go CLI compiles deterministically and drives Blender.

## Exact commands

```bash
blendctl doctor
blendctl schema
blendctl schema project
blendctl schema module
blendctl assets list
blendctl assets describe base_cube
blendctl assets describe column_chart
blendctl validate scene.json --json
blendctl compile scene.json --out build/
blendctl preview scene.json --frame 24 --out preview/
blendctl render scene.json --out render/
blendctl inspect project.json
blendctl validate project.json --shot intro
blendctl validate project.json --all-shots
blendctl preview project.json --shot intro --out preview/
blendctl guide
```

Environment: set `MOBLEND_BLENDER` if Blender is not on `PATH`. Supported family: Blender **5.3**.

## Declaration essentials

- Required: `schema_version` (`1.0.0`), `asset_pack_version` (`0.1.0`), `scene`, `camera`, `objects`
- Units: metres, degrees at JSON boundary, colours `#RRGGBB`
- Object types: `text`, `form` (`cube|plane|cylinder`), `asset` (`base_cube`, `column_chart`)
- Animation: tracks with `step|linear|ease_in_out` on approved properties
- Limits: see schema / Phase 1 CONTEXT (max 32 objects, 600 frames, etc.)
- **Title framing:** for front-facing title cards use text `rotation_euler_deg: [90,0,0]`, camera looking along +Y (e.g. location `[0,-7.5,1]`, rotation `[90,0,0]`, focal ~40). Too-close cameras crop letters. Copy `testdata/compile/title_small.json` as a starting plate.
- **Pedestal:** keep `base_cube` above the cylinder (`location.z` ≈ `0.7` with plinth scale z `0.2`); pull the camera back so the product is fully in frame.
- **Charts:** bar heights/`bar_count` bind to Geometry Nodes; mesh tint from `material.base_color` on GN instances may not read as solid colour in preview (known limitation — treat values/count as the structural oracle).

## Projects and modules (multi-file)

Use when separate agents own different parts of a production. Detection is by `kind`/`project_version`, never by filename.

- **Project manifest** (`kind: project`, `project_version: 1.0.0`): lists `entry`, explicit `modules[{id,path}]`, optional `shots[{id,entry}]`.
- **Entry** (`kind: entry`): owns `schema_version`, `asset_pack_version`, `scene`, `camera`; may declare objects/lights/animation and `imports`.
- **Module** (`kind: module`): objects/lights/animation and `imports` only.
- **Namespaces**: each file has `namespace` (`^[A-Za-z][A-Za-z0-9]*$`). Resolved IDs are `{namespace}__{local}`.
- **Import aliases**: unique within each file; duplicate aliases are errors. The entry camera keeps its bare ID for camera animation.
- **Cross-module refs**: `{alias}:{local_id}` after an explicit import; same-file refs stay bare local IDs. Do not write resolved `namespace__id` strings to bypass imports.
- **Limits**: ≤32 modules, import depth ≤8, ≤256 KiB/file, ≤1 MiB project; resolved scene limits still apply.
- **Ownership**: one agent edits one module file; coordinate only on shared entry/camera and import contracts. Do not invent directory scans or implicit overrides.
- **Inspect** lists modules/imports without resolving the whole graph into a scene.
- **Compile** writes to a fresh `moblend-job-*` directory under `--out` (nested under the selected shot ID) and prints the artifact paths. It emits `source_map.json` and `dependencies.json` beside script/data (not part of compile identity).

Examples: `examples/projects/title_modular/`, `examples/projects/three_shots/`.

## Shots (independent targets)

- Multi-shot projects require `--shot <id>` (or sole-shot default). `--all-shots` validates every shot closure and the entire listed dependency graph.
- Selected shot loads **only** its dependency closure; a broken unselected shot does not block the selected one.
- Preview/render with `--out DIR --shot ID` nests under `DIR/<shot>/` for isolation. Concurrent jobs stay in separate job directories.

## Worked examples

- `examples/title.json` — text title reveal
- `examples/pedestal.json` — form + `base_cube` product
- `examples/chart.json` — `column_chart` with `params.values` (0..1) and `bar_count`
- `examples/projects/title_modular/` — modular title paired with flat fixture
- `examples/projects/three_shots/` — intro/product/outro sharing `shared.module.json`

## Common diagnostics

| Code | Meaning | Fix |
|---|---|---|
| `SPEC_DUPLICATE_KEY` | Duplicate JSON key | Remove duplicate |
| `SPEC_UNKNOWN_FIELD` | Extra field | Remove unknown keys |
| `SPEC_SCHEMA_VERSION` | Wrong schema | Use `1.0.0` |
| `CAT_UNKNOWN_ASSET` | Bad `asset_id` | `assets list` |
| `SEM_DUPLICATE_ID` | ID clash | Unique ids across camera/lights/objects |
| `PROJ_CYCLE` / `PROJ_PATH_ESCAPE` | Bad module graph/path | Fix imports; keep paths inside the project root |
| `MOD_DUP_IMPORT_ALIAS` | Alias reused in one file | Give each imported module a distinct alias |
| `MOD_ENTRY_ONLY` | Entry settings in a module | Move scene/camera/version fields to the entry |
| `PROJ_FILE_TYPE` / `LIM_FILE_BYTES` | Nonregular or oversized input | Use a bounded regular JSON file |
| `SHOT_REQUIRED` | Multi-shot, no selection | Pass `--shot` or `--all-shots` |
| `ENV_BLENDER_*` | Doctor/preflight | Install/pin Blender 5.3 |

Machine diagnostics go to stdout with `--json`. Project diagnostics include `source` (relative file). Progress/logs go to stderr.

## Unsupported in v1 (do not invent fallbacks)

- Arbitrary Python / bpy paths / user `.blend` loading
- RAG / fuzzy asset search / marketplace downloads
- Physics, rigs, Grease Pencil authoring, custom node graphs
- Video/FFmpeg encoding, MCP server, GUI
- Multiple cameras, parenting hierarchies, HDRI worlds beyond `world_color`
- Parameterized module instances, remote imports, incremental build cache

## Repair loop

1. `validate --json` and fix pointers/codes (and `source` for projects)
2. `preview --frame N` and inspect PNG + `scene.blend`
3. Adjust text/colours/layout/params; re-validate
4. `render` for full frame range; check `manifest.json` status `complete`

## Output and snapshot safety

Shot IDs must match `^[A-Za-z][A-Za-z0-9_]*$`. Compile, preview and render allocate unique jobs under the output parent, including when the sole shot is selected by default. An escaping shot-directory symlink is rejected.

Only `manifest.json` with status `complete` establishes rendering success. PNGs or `.blend` alone may be partial artifacts from a failed job; render success requires the exact inclusive expected frame set. Failures retain `manifest.failed.json`. Job directories are retained for diagnosis; remove unwanted jobs manually after inspection.

Dependency hashes describe bytes read during load, including the manifest, and never re-read project files after rendering. Coordinate edits before starting jobs: this is a per-file snapshot, not an atomic snapshot across concurrent filesystem edits. The loader follows the containment-checked real path; adversarial directory replacement races remain outside the v1 execution-hygiene guarantee.

Version probes and render jobs share process-group cancellation on tested Linux, private configuration directories, and 64 KiB retained log limits with explicit truncation. Successful preview/frame PNGs are fully decoded and checked against declared dimensions before the completion manifest is written.
