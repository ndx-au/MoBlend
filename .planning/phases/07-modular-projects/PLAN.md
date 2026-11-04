# Phase 7 — implementation plan

Status: complete.

## Goal and prerequisites

Deliver explicit multi-file composition on top of the working single-file RC. Read CONTEXT and the roadmap; inspect current code and regression commands before edits. Completed Phase 0–6 work is the baseline. Do not confuse existing phase 5 fixture results with model evaluation.

## Tasks and ownership boundaries

| Order | Responsibility / proposed files | Work | Requirements |
|---|---|---|---|
| 1 | Project schema/model: new `internal/project/` and embedded project/module schemas | Freeze manifest shape, input discriminator, namespaces, import aliases, entry semantics, limits; add positive/negative fixtures | MOD-01, MOD-02 |
| 2 | Resolver: `internal/project/` | Bounded immutable input snapshot, local graph/cycle detection, path containment, deterministic traversal and assembly; rewrite typed references, validate aggregate limits | MOD-02, MOD-03, MOD-05 |
| 3 | Diagnostics: `internal/diag/`, resolver provenance, validator integration | Add source context and map assembled errors back to original files/pointers; preserve old diagnostic contracts | MOD-04 |
| 4 | Integration: CLI entry point, compiler/runner adapter | Shared file-aware input loading for validate/compile/preview/render; retain raw single-file APIs; dependency inspection and separate source-map artifact | MOD-01, MOD-05, MOD-06 |
| 5 | Examples, embedded guide and schemas, integration fixtures | Modular/flat paired example, safe agent ownership guidance, regression and real Blender checks | MOD-06, MOD-07 |

These are file responsibility boundaries, not a request to spawn agents. Keep resolver → CLI integration sequential until contracts are stable. Other ongoing work must not be reverted.

## Meaningful verification

- Strict project/module fixtures: unknown/duplicate keys, unsupported versions, empty/missing imports, duplicate local/qualified IDs, missing references, cycles and conflicts.
- Local paths: traversal, absolute paths, escaping symlinks, repeated dependency aliases and controlled nested imports. Verify byte/depth/count limits fail before Blender starts.
- Aggregate scene limits: each module individually small but combined declaration too large must fail.
- Determinism: key/import reordering and copied project root yield equal emitted data/script; paired flat/modular fixture equality; source paths appear only in separate provenance artifacts.
- Diagnostics: nested module parse and semantic failures identify original relative file/pointer for humans and JSON consumers; existing single-file output remains compatible.
- Offline validate/compile with no Blender available; dependency inspection returns machine-readable IDs/files/imports.
- Real Blender: paired modular/flat scene probes and render comparison in pinned environment; project preview and full short render save editable scene and manifests.
- Existing Go tests/vet/format checks plus the repository's established Blender integration suite. Document exact actual commands/results in VERIFICATION.md; no broad support claim from compile-only tests.

## Completion and handoff

Mark MOD-01–07 met only after evidence. Update STATE and append an immutable M007 progress entry. Retain historical phase verification; do not edit committed progress entries. Queue Phase 8 next, while leaving model authoring evaluation and platform/publication caveats explicit.
