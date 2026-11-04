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
| SCENE-02 | Bounded animation | Step, linear and documented easing checked at start, intermediate and end frames | 3 | **met** — pedestal motion probe; step fixture in eval |
| SCENE-03 | Geometry Nodes asset | Two independently parameterized instances render correctly without shared-state leakage | 3 | **met** — `TestApplyChartIsolation` |
| RUN-01 | Fresh isolated jobs | Repeated/concurrent jobs produce correct independent scenes without duplicate objects or output collisions | 4 | **partial** — unique job dirs verified; concurrent stress not automated |
| RUN-02 | Preview and final outputs | Preview PNG, full PNG sequence, editable scene and manifest have expected dimensions/counts/content | 4 | **met** — preview + `TestRenderWritesManifestAndFrames` |
| RUN-03 | Honest failure and cancellation | Timeout, invalid output path, Python error and killed process return failure; no success manifest | 0, 4 | **met** — smoke fail/cancel; failed marker; success-only `manifest.json` |
| AGENT-01 | No-RAG authoring package | Offline guide, schema, exact-ID catalogue and three examples support documented agent authoring exercise | 5 | **met** (docs/CLI); LLM exercise **blocked** |
| UX-01 | Machine-readable interface | JSON diagnostics parse and identify input location; logs never corrupt machine stdout | 4 | **met** — validate/doctor `--json` |
| QA-01 | Three visual references | Title, pedestal and chart pass scene assertions, render checks and recorded contact-sheet review | 5 | **partial** — automated previews/probes pass; contact-sheet review **blocked** |
| DIST-01 | Self-contained tool distribution | Clean-machine exercise needs only tool + supported Blender; checksums and notices included | 6 | **met** — `dist/linux-amd64/` + docs |
| DIST-02 | Verified support matrix | Every advertised OS/architecture has real smoke-render evidence, not just cross-compilation | 6 | **met for Linux amd64 only**; others not advertised |

## Release gate

Primary platform (Linux amd64) requirements met with explicit partials/blocks above. No network/RAG/model key required for ordinary use. Public publication requires owner authorization. Missing visual review and LLM eval are documented blockers for those *claims*, not for the Linux CLI RC.
