---
name: moblend-authoring
description: Create and revise Blender motion graphics with MoBlend's blendctl CLI, including JSON scenes, modular projects and independent shots. Use for MoBlend production authoring, preview inspection and diagnostic-driven repair, rather than development of the transpiler itself.
---

# MoBlend authoring

Turn a brief into supported declarative motion graphics, an editable Blender scene and verified rendered outputs. Use the host's ordinary file, command and image-inspection capabilities; no particular harness, plugin, agent API or model provider is required.

## Establish the available tool contract

Locate `blendctl` on PATH or use an executable supplied by the user. If only a MoBlend source checkout is available, run `go run ./cmd/blendctl` from that checkout, using absolute input/output paths when the production lives elsewhere. Do not assume this skill's directory is the production or compiler checkout. If no executable/source checkout is available, report the missing prerequisite rather than inventing a command or downloading software implicitly.

Read the installed tool's guidance before authoring:

```text
blendctl version
blendctl guide
blendctl schema
blendctl assets list
blendctl doctor
```

Read `assets describe <asset-id>` only for assets needed by the brief. For multi-file work also read `schema project` and `schema module`. The installed guide, schemas and catalogue own syntax, limits, versions and known visual limitations; do not copy a fixed schema or asset list into this skill. Use documented CLI help to resolve command differences. Match the tool's supported Blender environment; `MOBLEND_BLENDER` can select an existing executable without changing a machine-wide installation.

If Blender is unavailable, author/validate/compile when possible and report preview/render as blocked. Do not treat compilation as proof of a rendered result.

## Choose the production boundary

- Use a single scene declaration for a small standalone composition.
- Use a project with explicit local modules when parts have independent ownership or reuse. Inspect an existing project with `blendctl inspect project.json` before editing. Read the target module, its imports and the relevant entry/schema rather than the entire production by default.
- Use independent shot entries for separate scenes or timelines. Select a shot explicitly; do not resolve selection errors by rendering every shot.

Keep scene settings and camera in the entry. Follow the guide's namespace and import rules; cross-module references must be explicit. Modules compose scene data, not arbitrary Python. Do not substitute expressions, implicit merge order or filename guesses for the supported contracts.

When several agents are involved, agree on ownership of module/shot files and coordinate shared entry, namespace and import changes. Do not spawn agents merely because the project is modular. Preserve unrelated edits and avoid rewriting shared declarations for a local change. Existing resource limits still apply to the resolved scene; file splitting does not increase them.

## Author to the brief

Translate the requested text, palette, dimensions, timing, camera framing, assets and outputs into checkable criteria. Choose sensible reversible defaults for unspecified details. Ask only when a missing decision prevents a useful supported result.

Use an installed example as a starting point when available; take framing and coordinate conventions from the current guide. Preserve literal user text as JSON data. Use exact asset IDs and published parameters. Keep production files separate from compiler resources and place outputs in a new, explicit directory.

For unsupported requests, identify the unsupported part and explain what can be produced with the available tool. Do not invent fields/assets, run arbitrary bpy, fetch external assets, add RAG or change the transpiler to make an authoring request pass. Obtain direction if an alternative would materially change the requested deliverable. Do not label PNG frames as an encoded video.

## Validate, compile and preview

The following are command patterns; replace the input, frame and output paths with the production's values. Use fresh directories for separate attempts. For projects, carry the same `--shot ID` selection through validate, compile, preview and render.

```text
blendctl validate scene.json --json
blendctl compile scene.json --out build-attempt-01/
blendctl preview scene.json --frame 24 --out preview-attempt-01/

blendctl inspect project.json
blendctl validate project.json --shot intro --json
blendctl compile project.json --shot intro --out build-intro-attempt-01/
blendctl preview project.json --shot intro --frame 24 --out preview-intro-attempt-01/
```

Choose an in-range frame and a modest render workload. For animation, inspect first/middle/last and important transition frames, not just a pleasing final still. Validate before rendering; do not silently shorten timelines or change requested dimensions to fit a limit. If using a reduced-cost draft declaration, keep it separate and preserve the final settings.

Inspect the actual emitted PNGs using whatever image viewer the harness provides. Check legibility, missing glyphs, framing, clipping, overlap, material appearance and fidelity to the brief. For charts, verify values/count and instance isolation; apply documented Geometry Nodes material limitations honestly. For timing claims, use frame samples or scene probes where available. A successful process or a PNG existing does not establish visual correctness. If image inspection is unavailable, say visual review is blocked.

## Repair within scope

Read diagnostic code, pointer and project `source` file before changing data. Fix the cause in the owning file; preserve intended text, timing, camera and unrelated modules. Revalidate and preview the changed behavior.

Distinguish invalid authored input from environment, tool/asset and visual-quality failures. Do not keep trying random schema fields, weaken checks or modify compiler code during an authoring task. If the same failure persists after two focused repairs, stop repeating that approach: minimize the failing input, retain logs and report the blocker or choose a different supported approach that preserves the brief. Continue independent work when useful. For a confirmed tool bug, provide the smallest reproduction rather than claiming the production complete.

For an explicitly requested benchmark, use the repository's test playbook when available. Keep frozen prompt scoring and author repair attempts separate from compiler fixes; a hand-edited fixture is not evidence of model authoring success.

## Render and hand off

Render the requested selected shot or standalone scene after checking representative previews:

```text
blendctl render scene.json --out final-render/
blendctl render project.json --shot intro --out final-intro/
```

For a whole-production deliverable, list shots and render the requested ones explicitly. Use `validate project.json --all-shots` for project-wide validation when needed; it does not render the production. A broken unselected shot is separate from the selected shot's status.

Use reported job paths; outputs can be nested in generated job and shot directories. Verify command exit, completion manifest, editable `.blend`, requested frame count/range and image dimensions. Failed or partial jobs must not inherit a previous job's success. Preserve project/module sources and relevant dependency/source-map artifacts with the handoff.

Report what was delivered, actual source/output locations, shots/frames checked and any remaining limitations or blocked review. Distinguish authored, validated, compiled, rendered and visually inspected results. Do not claim cross-platform support or pixel identity that was not tested.
