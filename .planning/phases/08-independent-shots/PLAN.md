# Phase 8 — independent shots implementation plan

Status: complete.

## Goal and prerequisites

Add explicit shot targets on top of Phase 7 modular projects. Phase 7 resolver, namespaces, inspect, and single-scene project path are prerequisites. Do not add caching, schedulers, or video assembly.

## Tasks and ownership

| Order | Area | Work | Requirements |
|---|---|---|---|
| 1 | Manifest + CONTEXT | Frozen `shots[{id,entry}]` semantics, selection rules, provenance fields | SHOT-01 |
| 2 | Resolver / CLI | Shot closure resolve, `--shot` / `--all-shots`, sole-shot default | SHOT-01 |
| 3 | Outputs | Per-shot nesting, `dependencies.json` hashes, concurrent isolation | SHOT-02 |
| 4 | Example + tests | Three-shot example; change isolation; broken unselected shot | SHOT-03 |

## Meaningful verification

- Multi-shot without `--shot` fails; sole shot may default.
- Selected shot does not load unselected entries; `--all-shots` reports broken unselected shot.
- Concurrent shot previews write distinct job dirs + dependency manifests.
- Edit one shot file: unrelated resolved hashes unchanged; edit shared module: dependents change.
- Real Blender preview for three concurrent shots under pinned environment.

## Completion

Mark SHOT-01–03 only with VERIFICATION evidence. Append M008 progress. Keep model-eval and cross-platform caveats explicit.
