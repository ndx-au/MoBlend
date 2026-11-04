# MoBlend test and agent evaluation playbook

Status: executable test battery. Scenarios are oracles, not automatic results. Written 2026-11-04 against the Linux RC; Module (M) and Shot (S) cases are mandatory on the Linux amd64 pin. Latest recorded run: `docs/TEST-PLAYBOOK-RUN-2026-11-04.md`.

Audience: agents and maintainers finding, reproducing and fixing transpiler failures. After reading, select a supported suite, run its cases, preserve evidence, minimize failures and add meaningful regression coverage before fixing code.

## What this battery proves

Measure five separate outcomes: declaration correctness, deterministic compilation, correct Blender state, usable rendered output and agent brief adherence. A PNG existing is not evidence of the other four. Existing prepared JSON fixtures are compiler/render tests, not natural-language authoring evaluations.

Current baseline: scene schema 1.0.0, pack 0.1.0, Blender 5.3 family, Linux amd64. Record exact versions and hashes for every run. Module and shot tests are active. Record BLOCKED for unavailable model/non-Linux evidence; never count skips as passes. Runtime remains local and requires no model; model access is only an optional evaluation dependency.

## Run levels and budgets

| Level | Run when | Cases | Budget |
|---|---|---|---|
| Quick | Each relevant compiler/resolver fix | A + relevant D/M cases | No Blender; bounded generated inputs |
| Integration | Before closing a runtime/composition task | B/C and relevant M/S cases | Small previews and short frame ranges on pinned CPU configuration |
| Visual | Asset, animation or colour behavior changes | B + prompt output contact sheets | Preview samples first; approved reference quality only for selected cases |
| Agent | Guide/schema changes and release qualification | P01–P12 + unsupported set; M/S prompts when enabled | Two repair revisions per brief; never unlimited retries |
| Extended | Before a release candidate or after a demonstrated resource bug | Boundary/property tests and concurrency | Cap processes, memory, wall time and output size; avoid rendering maximum-size scenes |

Start with existing checks and inspect verbose output for skipped integration tests:

```bash
scripts/check.sh
go test ./internal/runner -count=1 -timeout 20m -v
```

These are existing entry points, not a guarantee that every scenario below is implemented. Confirm Blender availability with `blendctl doctor`. Existing test helpers may skip unavailable Blender; a skipped test does not satisfy integration evidence. Prefer targeted test names while repairing one failure. Use a disposable output parent for every run; never destroy existing work to test cleanup.

Freeze compiler revision, embedded resources and input files for a measurement. Worktree changes from another agent make a run non-comparable: record the diff snapshot or rerun after the relevant files stop changing. Never revert another agent's work. Evaluate supported features against the contract, not current accidental behavior.

## Case catalogue

Every variant in a row is a separate execution record. Exact diagnostic codes/pointers should be pinned when fixtures are added; use the current documented contract. A rejection must be an actionable structured error, not a panic or a Blender traceback for an input-validation failure.

### A — Strict decoding, references and limits (no Blender)

| ID | Stimulus | Expected oracle |
|---|---|---|
| A01 | Smallest valid title, pedestal and chart fixtures | All validate/compile; defaults match explicit equivalent values |
| A02 | Truncated JSON, trailing second document, unknown nested field, duplicate nested key | Each rejected at decode/schema stage with meaningful location; no subprocess |
| A03 | Null required values, wrong scalar types, short/long transform vectors, fractional frame index | Reject wrong shape/type, no silent coercion or missing required data |
| A04 | Wrong scene/pack version and unknown asset/property | Explicit unsupported version/capability error; no fallback |
| A05 | Duplicate IDs across camera/light/object, missing animation target, forbidden parameter | Correct entity/reference error; no last-writer-wins behavior |
| A06 | Conflicting tracks, repeated keyframe time, frame outside range, reversed range | Reject according to timeline contract; keyframes never silently dropped |
| A07 | Each documented limit at lower/upper valid edge and one outside | Inside accepted, outside rejected; test object/light/text/frame/track/key/byte limits separately |
| A08 | Invalid colour, out-of-range roughness/chart value, exponent overflow, NaN/Infinity tokens | Valid values preserved; invalid/non-finite values rejected rather than clamped |
| A09 | Unicode text, quotes, backslashes, newline escapes and Python-like text in a text field | Valid string data survives unchanged into Blender text; no execution or unintended interpolation |
| A10 | Unsupported raw Python, bpy path, expressions, external asset path or invented field | Rejected locally with a repair hint; no external reads/downloads |

Boundary-test large counts/frames/resolutions at validation only. Use tiny bounded fixtures when checking execution.

### B — Scene, animation and visual truth (real Blender)

| ID | Stimulus | Expected oracle |
|---|---|---|
| B01 | Text/form/camera/light/material from a reference fixture | Reopen saved scene; probe IDs/types, transforms, active camera, shader values and explicit settings |
| B02 | Two chart instances with different counts, values, width and colour | Probe actual modifier inputs and evaluated geometry/material assignment; changing one leaves the other unchanged |
| B03 | Linear location from [0,0,0] at frame 1 to [2,0,0] at frame 21 | Frame 11 evaluates to [1,0,0] within documented tolerance; endpoints exact within tolerance |
| B04 | Same keys with step interpolation | Sample immediately before frame 21 and at 21; previous value held until the key |
| B05 | Same keys with ease_in_out | Endpoints and documented handle behavior correct; interior samples match the declared curve, not merely a different-looking PNG |
| B06 | Non-default focal length, colour, resolution, FPS and frame range | Saved Blender values and manifest match input; output dimensions and inclusive frame count correct |
| B07 | Title at first/middle/last frame, extreme supported text length and framing | Text legible at intended visible frames; no unintended clipping/occlusion; missing glyphs identified explicitly |
| B08 | Repeat identical construction and reopen artifacts | Equivalent owned scene/datablocks; no duplicates or stale state; no raw .blend hash equality claim |

Structural probes and render assertions must be independent of the emitter where practical. A test that only repeats the generated assignment cannot detect a wrong Blender binding. For colours, inspect shader values and output colour settings; displayed pixels are affected by lighting and colour management and need not equal the input hex.

### C — Runner failures and recovery

| ID | Stimulus | Expected oracle |
|---|---|---|
| C01 | Missing executable and unsupported Blender version | Actionable environment error, nonzero exit; no success manifest |
| C02 | Trusted fixture deliberately raises Python exception | Nonzero exit, bounded useful logs, failed job identified; never complete |
| C03 | Timeout and explicit cancellation during a bounded long-running fixture | Parent returns promptly, tested child process tree exits; incomplete artifacts never advertised as complete |
| C04 | Unwritable destination, path containing spaces/Unicode, existing unrelated sentinel file | Valid paths work; failures explain cause; sentinel never altered or deleted |
| C05 | Two previews of different scenes at once; repeat with same output parent | Unique jobs/outputs, correct scene for each result, no mutable resource cross-contamination |
| C06 | Success followed by failure and failure followed by success | Completion belongs only to the correct job; stale outputs cannot make a failed run look successful |
| C07 | Corrupted/missing pack, font or template in a controlled fixture | Resource integrity/binding failure is explicit; does not pass based on a PNG alone |
| C08 | Validate with JSON output while diagnostics/logging occur | Stdout parses as documented machine output, logs on stderr, correct failure exit |

Use controlled subprocess/filesystem fixtures to exercise failures. Do not kill unrelated Blender sessions, fill the real disk or modify the installed pack to inject a failure.

### D — Determinism, generated cases and change isolation

| ID | Transformation | Required invariant |
|---|---|---|
| D01 | Repeat compile; reorder JSON object keys | Equal generated script/data bytes; no timestamps or random IDs enter identity |
| D02 | Change only title text | Intended object data changes; other transforms/camera/assets remain equivalent |
| D03 | Omit documented default versus specify that default | Equal normalized meaning and emitted artifacts where the canonical contract promises it |
| D04 | Generate seeded valid bounded declarations | No panic; validate→compile→decode canonical data preserves supported meaning |
| D05 | Seeded mutation of one field into a known-invalid value | Expected rejection, not crash or accidentally accepted data |
| D06 | Same pinned CPU renderer/environment repeated | Structural equivalence and documented image tolerance; record image differences instead of demanding cross-GPU identity |

Record seed and minimized input for every generated failure. Freeze image thresholds before a run from repeatability measurements; do not tune the baseline to make a failed change pass. Preserve intentional array ordering.

### M — Modular projects (Phase 7 capability gate)

| ID | Scenario | Expected oracle |
|---|---|---|
| M01 | Split paired monolithic scene into explicit modules | Equal normalized scene/script and equivalent scene probe/render; old single-file input still works |
| M02 | Reused local IDs in distinct namespaces; duplicate namespace; unresolved alias | Valid namespaced case works; ambiguous/unknown references rejected |
| M03 | Nested imports, diamond dependency, missing file, direct and indirect cycles | Contract-defined deduplication/identity; missing/cyclic dependencies fail deterministically |
| M04 | Import reordering and relocation to a different checkout root | Same generated script/data; only separate source provenance may differ |
| M05 | Parse/semantic error in deeply imported module | Diagnostic identifies original relative file and JSON pointer after reference rewriting |
| M06 | Traversal, absolute import, escaping symlink and dependency/byte/depth overrun | Reject before Blender; symlink spelling cannot bypass root containment |
| M07 | Modules each under local limits but resolved scene over its limit | Reject aggregate scene limit; no bypass by file splitting |
| M08 | One module edited by an agent while another remains unchanged | Scope-limited diff; expected resolved data changes only; shared contracts explicit |

### S — Independent shots (Phase 8 capability gate)

| ID | Scenario | Expected oracle |
|---|---|---|
| S01 | Three shots with one selected; omit selection in multi-shot project | Only selected closure built; ambiguous selection fails unless explicit all-shot mode requested |
| S02 | Broken import in an unselected shot | Selected valid shot succeeds; explicit project-wide validation reports broken shot |
| S03 | Change a shared module then a shot-only module | Shared change alters dependent content hashes; unrelated shot-only edit leaves selected scene unchanged |
| S04 | Concurrent shot builds plus one failed shot | Isolated outputs/manifests; success and failure recorded for the right shots |

## Natural-language authoring evaluation

Keep the evaluator's scoring checklist separate from the authoring agent. The agent receives the shipped guide, schema, catalogue, allowed examples and one brief. Do not supply the expected answer JSON or a previous candidate. Freeze prompts and checks before comparing model or compiler revisions. Existing supported fixtures can inform scoring, but copying them into a successful render is not an authoring measurement.

Use the following prefix, then one P brief:

> Use the installed MoBlend CLI and its local guide/schema/catalogue to author the requested scene. Stay inside supported capabilities. Write JSON, validate, compile and preview. Keep user text as data. Explain unsupported parts precisely; do not invent fields or fetch assets. Save artifacts and command output. Do not modify the compiler, schemas, guide or tests. Use a new output directory. Start with a small preview. Do not run final animation rendering until the brief requests it or the evaluator asks.

Supported briefs are P01–P12. Numeric render settings are intentional; camera/layout choices may vary while meeting the checklist. All timing uses integer frames. Use palette and font capabilities actually available; inability to display a requested glyph is recorded, not hidden.

| ID | Brief supplied to agent | Evaluator checklist (withheld from author) |
|---|---|---|
| P01 | “Make a title card reading MoBlend. Off-white text on a dark blue background. 384×160, 24 FPS, frames 1–24. Preview frame 24.” | Exact text, palette intent, dimensions/FPS/range, visible readable title |
| P02 | “Make a title card containing exactly: She said "Go!" — C:\temp\logo. Treat it as literal text. Preview at 384×160.” | Exact decoded string; punctuation/backslashes survive; no executable content |
| P03 | “A cube's X location must move linearly from 0 at frame 1 to 2 at frame 21; Y and Z remain 0. Use 24 FPS. Preview frame 11.” | Required keys/interpolation; Blender evaluated X=1 at frame 11, Y=Z=0 |
| P04 | “A cube stays at X=0 through frame 20, then jumps to X=2 at frame 21. Other coordinates stay 0. Preview frame 20.” | Step interpolation; probes at 20/21 return 0/2 |
| P05 | “Animate the title Launch upward with the supported ease-in-out preset over frames 1–24. Make its end position readable. Preview frames 1, 12 and 24.” | Correct target, preset and keys; three previews; intended upward motion and readable final framing |
| P06 | “Show the bundled base_cube product on a cylinder pedestal. Green product, neutral pedestal, warm key light. Preview at 320×240.” | Exact catalogue asset and valid form; independently assigned materials, visible product and pedestal |
| P07 | “Use column_chart for three bars with normalized heights 0.2, 0.9 and 0.4. Use blue bars and show all bars. Preview at 320×240.” | Values/count preserved in modifier and geometry; all bars visible; no fabricated labels/capabilities |
| P08 | “Show two separate column charts. Left: three bars [0.2,0.9,0.4], blue. Right: two bars [0.8,0.3], red. Preview both together.” | Counts/values/materials isolated; two independent instances; no overlap hiding bars |
| P09 | “Create a product composition using only supported assets. Explicitly use a 50 mm camera lens and a 320×240 output.” | Lens and dimensions match scene settings; valid visible product composition |
| P10 | “Create a simple cube scene at 24 FPS, frames 5–8 inclusive, 64×64, four samples. Render the complete PNG sequence.” | Exactly frames 5,6,7,8, correct dimensions, editable scene and success manifest |
| P11 | “Here is a title declaration with one duplicate object ID [evaluator supplies fixed invalid fixture]. Repair the identifier/reference error while preserving the title text, camera and colours.” | Minimal targeted repair; no unrelated changes; IDs unique and references valid |
| P12 | “Change only the left chart's values to [0.4,0.4,0.4] in this valid two-chart scene [evaluator supplies fixed fixture]. Keep the right chart and camera unchanged.” | Scope-limited diff and Blender proof right chart/camera unchanged; left updated |

For P02 supply the literal payload as a separate fenced text block so quoting in a prompt transport cannot change it. Expected actual characters include quotes around Go! and single backslashes in the Windows-style path. Keep that payload file with the run:

```text
She said "Go!" — C:\temp\logo
```

Unsupported prompt set (scored separately):

| ID | Brief | Expected behavior |
|---|---|---|
| U01 | “Load my arbitrary .blend and run this Python to construct custom nodes.” | Explicit capability refusal; no script execution or arbitrary blend load |
| U02 | “Find a realistic car online, download it and animate its rig.” | Identify unsupported search/download/rigging; no invented asset ID |
| U03 | “Render 10,000 frames with 10,000 objects.” | Explain existing bounded limits; no silent truncation or large render launch |
| U04 | “Output an MP4 with synchronized audio.” | Explain unsupported video/audio pipeline; PNG alternative may be offered explicitly, never claimed as the requested deliverable |

Include prompt M-P01: “Split this valid scene into two owned modules without changing the image. Give me the dependency listing and preview.” Score M01/M04/M05 and scope of edits. Include S-P01: “Change only shot two's title and render only shot two.” Score S01/S03/S04. They do not enter the P01–P12 denominator.

### Attempt protocol and scoring

1. Use a fresh agent context for each brief. Record model/provider/version and the exact guide/schema/catalogue context supplied. Local file inspection is allowed and is not RAG. Do not auto-inject unrelated source code or expected fixtures.
2. Attempt 0 ends when the agent submits its first candidate. The evaluator runs the fixed validator/probes/render checks; retain stdout, stderr, exit code, declaration and previews.
3. Return only diagnostics and observed requirement mismatches. Allow at most two candidate repair revisions. No compiler edits or changed brief/checklist during these attempts. Record first-pass and eventual outcomes separately.
4. Score each brief on validity, compile success, Blender structural assertions, requested outputs and visual/brief adherence. Overall pass requires every applicable check. A valid declaration with an unreadable title fails brief adherence.
5. Record unsupported handling separately, including whether the agent invented fields/assets or attempted prohibited fallback. Refusing an in-scope brief is a failed supported case.
6. If model access or visual inspection is unavailable, mark the affected result BLOCKED. Do not substitute hand-authored fixture success or an automated PNG-existence check.

Initial release evaluation target: at least 10/12 supported briefs satisfy all applicable checks within two revisions, and every unsupported brief is handled without unsafe or misleading fallback. Require all deterministic correctness/regression tests to pass; the model target is not permission for compiler failures. Report first-pass count, eventual count, failure causes, revision counts and unscored/blocked cases. No single run establishes universal model reliability; repeat the same frozen set at least three times before claiming comparative improvements, retaining every run.

## Visual review oracle

Use contact sheets labeled with case, frame, build and effective settings. Review: title legibility/missing glyphs, framing, object overlap, chart value/count fidelity, material separation, lighting, clipping and motion readability. Compare first/middle/last frames plus boundary frames around step changes. Passing individual images does not prove temporal behavior; probe intermediate values and inspect a short sequence for flicker where relevant.

Record reviewer identity/type, each criterion and notes. Automated structural/image assertions and human/agent aesthetic review are separate fields. Never label uninspected output visually approved. Keep reference images and image-comparison thresholds versioned with their justification.

## Evidence record

Store each run in its own evaluation artifact directory; prefer an ignored local artifact root for bulky renders and sanitized committed summaries for durable evidence. Do not commit machine secrets, private assets or model credentials. A run record contains:

```yaml
run_id: <unique-run>
case_id: <P03-or-other-case>
status: NOT_RUN # PASS | FAIL | BLOCKED | NOT_RUN
compiler_revision: <commit-and-relevant-diff-hash>
resource_hashes: <record-pack-template-schema>
blender_version_and_build: <exact>
platform_and_render_device: <exact>
model_and_context: <exact-or-not-applicable>
seed: <if-generated>
attempt: 0
input_and_prompt_hashes: <exact>
commands_and_exit_codes: <artifact-reference>
oracle_results: <validity-structure-output-visual-brief-fields>
artifacts: <input-logs-scene-previews-manifest>
wall_time_and_timeout: <measured>
failure_class: <none-or-class-below>
```

Preserve baseline and failing artifacts before repair; capture source snapshots for concurrent work. A summary links the case IDs to requirements and regression-test names. Module cases cover MOD-01–07; shot cases cover SHOT-01–03. Baseline A/D cover CORE; B covers SCENE/QA; C covers ENV/RUN/UX; prompt runs cover AGENT plus brief-specific scene checks.

## Find → reproduce → fix → verify

1. Classify the failure: malformed author output, unclear guide/schema, parser/validator, resolver/provenance, compiler, asset binding, Blender runtime, environment, runner, or visual quality. A valid expected input rejected by a documented constraint may expose a design mismatch; do not remove limits silently.
2. Minimize the declaration/project and command while retaining the failure. Pin versions and seed. Keep a fresh reproduction independent of previous output directories.
3. Add a regression that fails for the observed reason, with an independent oracle. For a render bug, include an evaluated-scene probe or justified image comparison; for a diagnostic bug, assert file/pointer/code. Do not mirror the implementation and call it coverage.
4. Apply the smallest correct change. Keep public semantics and unsupported boundaries intact; document any intentional contract change. Do not weaken the fixture, expected outcome or threshold just to get green.
5. Rerun the minimized case, relevant neighboring cases and existing required checks. Rerun the full applicable battery before release qualification; avoid repeating costly unrelated rendering during every edit.
6. Record exact commands/results in the owning phase verification and append a new Dot Progress fix entry. Preserve earlier committed progress and evaluation history. New failures stay visible until resolved or explicitly blocked.

An execution-agent instruction for this loop:

> Read this playbook, current STATE and the active phase contract. Inventory existing coverage. Select the highest-risk applicable cases; do not execute blocked future features as though implemented. Run a bounded baseline, preserve failures, minimize each reproducible issue, add a meaningful regression, fix it, and verify the affected behavior. Keep authoring evaluation and implementation repair as separate runs. Report pass/fail/blocked counts and exact remaining issues; never invent results or expand scope to satisfy unsupported prompts.

## Adoption order

First extend the current suite with A/B/C gaps found during the active phase; reuse its 12 fixture renders as baseline regression coverage. Run the mandatory M and S cases for composition changes. Run P01–P12 with an actual model when available to resolve the outstanding authoring-evaluation evidence gap. Add every confirmed minimized bug to the permanent regression corpus. This playbook adds verification guidance, not a new runtime feature or a completed milestone.

## Gauntlet round 1 clarifications

- S02: invoke `validate testdata/project/shots/broken_unselected/project.json --all-shots --json`; require a `PROJ_FILE_MISSING` diagnostic attributed to shot `bad`. Selected `good` must succeed independently.
- Prepared supported JSON previews are fixture smoke tests only. Unsupported JSON negatives check the schema/validator, not natural-language refusal behavior. P01–P12 and U01–U04 agent-authoring/refusal evidence remain blocked until actual model evaluation runs are recorded.
