# Implementation roadmap

Status: ready for Phase 0 planning. Phases are dependency gates, not calendar promises. Complete one vertical slice at a time; preserve evidence before expanding the language.

## Phase 0 — Prove the Blender boundary

**Outcome:** a Go program drives a real headless Blender render using an embedded minimal asset pack.

Tasks:
1. Initialize module, basic CI and test harness. Select and pin a maintained Go toolchain and an available supported Blender build using official documentation and local execution evidence.
2. Create a first-party base scene from a checked-in Blender Python authoring script: camera, light, simple object, known material, explicit render/colour settings. Save the small `.blend` pack.
3. Embed and extract the pack; launch Blender with a trusted smoke script, change one property, save a scene, render one PNG.
4. Verify command argument ordering, Python failure exit, CPU-compatible rendering, fonts, configuration isolation and cancellation on Linux amd64.
5. Record asset checksum, versions, invocation, output and observed constraints in the phase context. Create a supported-environment table.

**Gate:** real PNG with expected dimensions and nonempty content; saved scene reopens; deliberate Python exception returns failure; executable works offline with Blender installed. If the boundary fails, resolve it before schema expansion.

## Phase 1 — Define the smallest useful declaration

**Depends on:** 0. **Outcome:** an executable language contract.

Tasks:
1. Write schema v1, catalogue contract and Go input types. Freeze units, colour handling, frame semantics, defaults, limits and diagnostic shape.
2. Implement strict parsing, duplicate-key rejection, version checking, schema and semantic validation, and canonical normalization.
3. Define only the fields needed by title, pedestal and chart examples. Author positive and negative fixtures before implementing further capabilities.
4. Add `schema`, `assets list/describe`, and `validate` CLI surfaces. Ensure validation is local and does not launch Blender.

**Gate:** all positive fixtures validate; every negative fixture has the expected stable code and JSON pointer. Schema/model parity is checked. Record any deliberately deferred fields.

## Phase 2 — Deliver the first useful composition

**Depends on:** 1. **Outcome:** JSON title scene → editable Blender scene → still image.

Tasks:
1. Lower normalized declarations to typed operations; generate Python entry point and canonical data using trusted templates.
2. Implement text, basic forms, materials, transform, camera and lighting bindings. Avoid selection/context dependence.
3. Add `compile`, retaining readable artifacts for debugging. Build from the pristine pack for every execution.
4. Test malicious-looking text and Unicode round-trip, stable ordering, repeated compilation, missing bindings and saved-scene properties.

**Gate:** the title example renders in Blender; changed text/colour/layout is visible; a probe verifies object IDs, transforms, materials and camera. Compilation is byte-stable and requires no Blender runtime.

## Phase 3 — Add expressive, bounded motion

**Depends on:** 2. **Outcome:** all three reference compositions can be expressed and animated.

Tasks:
1. Implement typed keyframe tracks for approved transforms, selected material values and asset parameters. Add step/linear interpolation, then the single documented easing preset.
2. Author a column-chart Geometry Nodes asset with bounded values and published input contracts. Package any required font with redistribution notices.
3. Verify instance isolation and socket identifiers; reject unsupported animation targets and incompatible pack versions.
4. Finish title reveal, geometric product pedestal and chart declarations. Add one fixed compositor look only if Phase 0 verified it.

**Gate:** Blender probes evaluate representative frames and check expected values within defined tolerances. Two asset instances stay independent. All three short animations render with correct frame counts. No arbitrary node graphs, rigs or simulations enter this phase.

## Phase 4 — Make rendering dependable

**Depends on:** 3. **Outcome:** a user-facing local CLI handles success and failure predictably.

Tasks:
1. Complete `doctor`, `preview`, and `render`; add structured diagnostics, stable exit categories and stdout/stderr discipline.
2. Implement unique job directories, safe output handling, resource integrity checks, timeouts, cancellation, bounded logs and complete/failed manifests.
3. Explicitly configure render engine, resolution, frame range, FPS, samples, seed where applicable and colour pipeline. Store all effective settings.
4. Exercise missing Blender, incompatible version, unwritable destinations, process failures, concurrent jobs, cancellation and repeat runs.
5. Reopen delivered `.blend` files and verify editability and expected scene structure.

**Gate:** automated end-to-end success and failure matrix passes. Incomplete output cannot be mistaken for a finished render. Jobs do not mutate the embedded pack or overwrite unrelated files.

## Phase 5 — Prove agent usability without RAG

**Depends on:** 4. **Outcome:** an agent can author and repair supported scenes from the shipped local documentation.

Tasks:
1. Write a compact agent guide with supported concepts, exact commands, common diagnostics, worked examples and explicit unsupported requests.
2. Make `guide` print the embedded documentation. Keep schema and catalogue accessible on demand through exact-ID commands; no semantic index or retrieval infrastructure.
3. Run an evaluation set of at least 12 supported briefs spanning the three compositions: text changes, palettes, timing, camera framing, chart values. Add a separate set of deliberate unsupported requests.
4. Record initial validation, repair attempts, final render success, brief adherence and out-of-scope handling. Document model/version/context if a model is used; keep it an evaluation dependency only.
5. Generate contact sheets and review typography, framing, intersections, chart fidelity, lighting and animation readability. Fix asset or compiler problems revealed by review.

**Proposed gate:** at least 10 of 12 supported briefs render and meet their checklist within two diagnostic-driven revisions; all unsupported cases are explicitly identified without unsafe fallback. Keep unsupported cases in a separate set so they do not inflate the score. Core compiler/integration tests all pass. If no model access or visual reviewer is available, mark that evidence blocked rather than fabricate it.

## Phase 6 — Package and hand off v1

**Depends on:** 5. **Outcome:** an installable, documented release candidate.

Tasks:
1. Produce release binaries with embedded schema, docs, catalogue, templates and assets. Include checksums, version reporting, asset/font attribution and build instructions.
2. Run a clean-environment offline onboarding test on Linux amd64. Test macOS arm64 and Windows amd64 before advertising support; document untested targets as experimental or omit them.
3. Check all example declarations against the final schema and rerun integration/render baselines on the pinned environment.
4. Write troubleshooting, limitations, compatibility policy, reproducible asset-build instructions and a short human/agent quick start.
5. Reconcile requirement status, phase evidence and known issues. Prepare release artifacts; public publication requires owner authorization.

**Gate:** REQUIREMENTS release gate is met; an independent onboarding run produces the reference outputs using only the documented prerequisites. Report measured binary size and render time with hardware/context rather than inventing performance targets.

## Sequencing discipline

Dependency chain: 0 → 1 → 2 → 3 → 4 → 5 → 6. Phase 0 resolves Blender/API uncertainty; Phase 1 establishes contracts; subsequent work expands a running product. Do not build MCP, a GUI, a provider abstraction, or RAG ahead of these gates.
