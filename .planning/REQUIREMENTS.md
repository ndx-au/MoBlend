# V1 requirements

Update status only with recorded verification evidence.

| ID | Requirement | Acceptance evidence | Phase | Status |
|---|---|---|---|---|
| ENV-01 | Supported Blender discovery and version gate | doctor succeeds on pinned build; missing/unsupported executable yields actionable errors | 0, 4 | **met** — `TestDoctorOK` / `TestDoctorMissing`; `blendctl doctor` |
| CORE-01 | Strict versioned JSON | Valid fixtures accepted; unknown/duplicate fields, wrong types and unsupported versions rejected | 1 | **met** — `TestNegativeFixtures` / positives |
| CORE-02 | Semantic checks and limits | Bad references, duplicate IDs, forbidden targets, conflicting tracks and resource overruns rejected before Blender launch | 1 | **met** — validate package + fixtures |
| CORE-03 | Deterministic compilation | Repeat compilation and irrelevant input-key reordering produce equal script/data bytes | 2 | **met** — `TestCompileByteStableAndSafeText` |
| CORE-04 | Data cannot become executable syntax | Quotes, backslashes, Unicode and Python-like text round-trip safely; invalid declarations never execute | 2 | **met** — compile test + validate rejects |
| ASSET-01 | Embedded pack with documented controls | Offline release extracts verified bytes; Blender probe validates every catalogue binding | 0, 3 | **met** — pack `0.1.0`, chart/base_cube bindings |
| SCENE-01 | Editable scene composition | Text, solid forms, camera, lights and materials inspected in saved scene and rendered | 2 | **met** — title probe + examples |
| SCENE-02 | Bounded animation | Step, linear and documented easing checked at start, intermediate and end frames | 3 | **met** — `TestMixedAnimationInterpolationAndOrder`: transform/material step/linear/ease and camera interior/endpoints, both track orders |
| SCENE-03 | Geometry Nodes asset | Two independently parameterized instances render correctly without shared-state leakage | 3 | **met** — `TestApplyChartIsolation` |
| RUN-01 | Fresh isolated jobs | Repeated/concurrent jobs produce correct independent scenes without duplicate objects or output collisions | 4 | **met** — unique job dirs; `TestConcurrentJobsIsolated` |
| RUN-02 | Preview and final outputs | Preview PNG, full PNG sequence, editable scene and manifest have expected dimensions/counts/content | 4 | **met** — preview + `TestRenderWritesManifestAndFrames` |
| RUN-03 | Honest failure and cancellation | Timeout, invalid output path, Python error and killed process return failure; no success manifest | 0, 4 | **met** — render/probe process-tree cancellation and timeout, bounded logs, invalid-output rejection; follow-up VERIFICATION |
| AGENT-01 | No-RAG authoring package | Offline guide, schema, exact-ID catalogue and three examples support documented agent authoring exercise | 5 | **met** (docs/CLI); LLM exercise **blocked** |
| UX-01 | Machine-readable interface | JSON diagnostics parse and identify input location; logs never corrupt machine stdout | 4 | **met** — validate/doctor `--json` |
| QA-01 | Three visual references | Title, pedestal and chart pass scene assertions, render checks and recorded contact-sheet review | 5 | **met** for smoke fixtures — `testdata/visual/` + `VISUAL-REVIEW.md`; HD examples optional follow-up |
| DIST-01 | Self-contained tool distribution | Clean-machine exercise needs only tool + supported Blender; checksums and notices included | 6 | **met** — `dist/linux-amd64/` + docs |
| DIST-02 | Verified support matrix | Every advertised OS/architecture has real smoke-render evidence, not just cross-compilation | 6 | **met for Linux amd64 only**; others not advertised |

## Release gate

Primary platform (Linux amd64) requirements met with explicit partials/blocks above. No network/RAG/model key required for ordinary use. Public publication requires owner authorization. Missing visual review and LLM eval are documented blockers for those *claims*, not for the Linux CLI RC.

## Planned modularity extension — not included in RC completion claims

| ID | Requirement | Acceptance evidence | Phase | Status |
|---|---|---|---|---|
| MOD-01 | Backward-compatible project manifest and module schemas | Existing single-file corpus unchanged; strict project/module version and field validation | 7 | met — follow-up VERIFICATION |
| MOD-02 | Explicit bounded dependency graph | Missing files, cycles, depth/count/byte overruns and escaping paths/symlinks rejected offline before Blender | 7 | met — follow-up VERIFICATION |
| MOD-03 | Namespaced typed composition | Distinct modules may reuse local IDs; duplicate namespaces and unresolved/conflicting references fail; import order has no semantic effect | 7 | met — follow-up VERIFICATION |
| MOD-04 | Source-aware validation | Nested module errors report originating relative filename and original JSON pointer in human/JSON output | 7 | met — follow-up VERIFICATION |
| MOD-05 | Compiler equivalence and project-wide limits | Equivalent modular/flat fixture has equal normalized data/script; aggregate scene limits enforced; relocation preserves bytes | 7 | met — follow-up VERIFICATION |
| MOD-06 | Usable project CLI and render path | Dependency inspection plus validate/compile/preview/render work on a modular example; real Blender scene/render equivalence | 7 | met — follow-up VERIFICATION |
| MOD-07 | Agent ownership guide | Shipped guide/schema/examples describe file ownership, imports, namespaces and safe concurrent editing without a model dependency | 7 | met — follow-up VERIFICATION |
| SHOT-01 | Explicit independent shot targets | Selected shot loads only its closure; ambiguous selection fails; explicit whole-project validation checks every shot | 8 | met |
| SHOT-02 | Isolated outputs and provenance | Per-shot scene/frames/manifests; concurrent jobs isolated; dependency/content hashes recorded from input snapshot | 8 | met |
| SHOT-03 | Change isolation | Three-shot example: one-shot edit preserves unrelated resolved outputs; shared edit changes hashes of dependent shots | 8 | met |

Phases 7–8 gates are met with VERIFICATION evidence. These do not invent LLM-authoring or multi-platform claims; the blocked model-based authoring claim remains blocked.

2026-11-05 follow-up evidence: `.planning/fixes/gauntlet-followup/VERIFICATION.md`. The fresh review reopened animation and module claims; concrete regressions and real Blender probes now cover the repaired paths. Earlier independent verdict remains a historical pre-fix review, not a post-fix ship verdict.
