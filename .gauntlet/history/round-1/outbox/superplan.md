# Composite leftover plan: Gauntlet review: MoBlend .planning/ROADMAP.md

Round 1/3. Verdict: **revise**.

## Coverage

- correctness: 47 (0 critical, 3 high, 10 findings)
- tests: 47 (0 critical, 3 high, 10 findings)
- failure: 40 (0 critical, 3 high, 8 findings)
- concurrency: 29 (0 critical, 2 high, 6 findings)
- security: 22 (0 critical, 1 high, 5 findings)

## Must fix

### aggregate-limit-assertion-soft (tests, high)

- Location:
```
internal/project/resolve_test.go:TestAggregateObjectLimit
```
- Bug:
```
After expecting any error, acceptance is nearly tautological: nonempty Message plus 'object'/32/ /objects marks found=true, so many non-LIM diagnostics satisfy the check; VERIFICATION claims LIM_OBJECTS specifically.
```
- Predicted:
```
Broken or removed aggregate object enforcement can still pass if any objects-related error is returned.
```
- Task:
```
Require Code == LIM_OBJECTS (or the single documented code) and fail on other codes.
```

### all-shots-untested (tests, high)

- Location:
```
internal/project/resolve_test.go:TestShotSelectionAndAllShots; cmd/blendctl.cmdValidateAllShots
```
- Bug:
```
Test name and Phase 8 VERIFICATION credit --all-shots surfacing the broken unselected shot, but the test only calls LoadAndResolveShot; no automated assertion exercises ValidateAllShots/CLI --all-shots (manual one-liner only).
```
- Predicted:
```
Regression in project-wide shot aggregation can ship green while selected-shot paths still pass.
```
- Task:
```
Add a Go test that validates broken_unselected with all-shots mode and asserts a diagnostic for shot bad while documenting selected-shot success separately.
```

### backlog-phases78-pending-lie (correctness, high)

- Location:
```
.planning/BACKLOG.md (Phases 7–8 pending line) vs STATE.md / REQUIREMENTS MOD/SHOT
```
- Bug:
```
BACKLOG still says Phases 7–8 “implementation and acceptance evidence are pending” while STATE/REQUIREMENTS/ROADMAP phase bodies mark MOD-01–07 and SHOT-01–03 met with VERIFICATION dates.
```
- Predicted:
```
Follow-on agents re-implement or reopen closed gates, or treat shipped modular/shot behavior as unfinished and skip regression maintenance.
```
- Task:
```
Rewrite the BACKLOG sentence to “shipped; parameterized modules / build cache remain deferred” and leave only true leftovers.
```

### cli-no-signal-ctx (failure, high)

- Location:
```
cmd/blendctl/main.go:cmdApply
```
- Bug:
```
preview/render/smoke use context.Background() with no os/signal wiring; Ctrl+C kills blendctl without cooperative cancel of Run/RunApply.
```
- Predicted:
```
Interactive cancel leaves Blender running against the job tree; subsequent agents may see late PNGs or a half-written scene.blend with no failed marker update.
```
- Task:
```
Propagate a signal-notified context into RunApply/RunSmoke so cancel takes the same path as Timeout.
```

### compile-out-clobber (concurrency, high)

- Location:
```
cmd/blendctl/main.go:cmdCompile vs cmdApply shot nesting; Phase 8 CONTEXT Outputs
```
- Bug:
```
Phase 8 required --out DIR --shot ID to nest under DIR/<shot>/ and always use distinct job dirs. preview/render do both via shot join + JobDir/MkdirTemp; compile writes straight into --out with no shot nest and no unique job dir.
```
- Predicted:
```
Concurrent blendctl compile … --shot intro|outro --out SAME_DIR (multi-agent / parallel CI) overwrite each other's script/data/source_map/dependencies.json; SHOT-02 isolation holds only for preview/render paths covered by tests.
```
- Task:
```
Make compile apply the same filepath.Join(out, shot) nesting as cmdApply, or always write under a MkdirTemp job dir like RunApply.
```

### deps-hash-toctou (concurrency, high)

- Location:
```
internal/project/provenance.go:BuildDependencyManifest
```
- Bug:
```
SHOT-02 / comment claim hashes from the load snapshot, but BuildDependencyManifest re-ReadFiles disk paths in r.Files after resolve (and after Blender in cmdApply). Loaded module bytes are not retained on Resolved.
```
- Predicted:
```
While concurrent shot jobs run, an edit to a shared module yields dependencies.json digests that disagree with Decl/ResolvedDataSHA actually compiled and rendered—false provenance under concurrent agents.
```
- Task:
```
Hash the bytes retained at load time (store digests or raw snapshots on Resolved) and never re-read project inputs when writing dependencies.json.
```

### mod-bound-limits-untested (tests, high)

- Location:
```
internal/project/resolve.go:loadOne/readBounded; playbook M06; MOD-02
```
- Bug:
```
Roadmap/MOD-02 require offline rejection of import-depth, per-file bytes, total project bytes, module-count, and absolute paths; resolve.go emits LIM_IMPORT_DEPTH/LIM_FILE_BYTES/LIM_PROJECT_BYTES/LIM_MODULES/PROJ_PATH_ABS but no *_test.go covers those codes (only cycle/missing/escape/symlink/aggregate objects).
```
- Predicted:
```
Depth/byte/count/absolute-path bypasses can land without failing CI despite containment claims.
```
- Task:
```
Add negative fixtures and assert exact codes for depth>8, oversized file, project-byte overrun, >32 modules, and absolute module path.
```

### mod02-orphan-listed-modules (correctness, high)

- Location:
```
internal/project/resolve.go:resolveManifest loadOne closure / MOD-02
```
- Bug:
```
Resolver loads only the entry/shot import closure. Modules listed in project.modules with missing paths, bad versions, or cycles that are not reachable from the selected entry never fail validate/compile.
```
- Predicted:
```
Agents get OK on projects whose module table is broken; inspect may show stubs while CI green; MOD-02 “missing files … rejected offline” is overstated for the full explicit graph.
```
- Task:
```
On non-shot validate (and --all-shots), either resolve every listed module path for existence/version/containment or document and test that only closure integrity is guaranteed.
```

### mod06-blender-equiv-existence-only (correctness, high)

- Location:
```
internal/runner/project_integration_test.go:TestModularFlatPreviewEquivalence / ROADMAP Phase 7 gate / MOD-06
```
- Bug:
```
Phase 7 gate and MOD-06 require modular≡flat Blender scene properties and render equivalence; the recorded Blender test only asserts preview.png/.blend exist for both sides after compile-byte equality—no scene probe or PNG compare.
```
- Predicted:
```
A resolver/compiler change that keeps compile bytes identical but diverges applied Blender state (or silently empty/wrong pixels) still passes VERIFICATION while MOD-06 stays marked met.
```
- Task:
```
Extend TestModularFlatPreviewEquivalence with the same object/camera/material probes used for title scenes and a bounded PNG structural check; demote MOD-06 until that lands.
```

### partial-frames-ok (failure, high)

- Location:
```
internal/runner/apply.go:RunApply
```
- Bug:
```
Render success only requires len(frames)>0; never checks inclusive decl frame_end-frame_start+1 or expected frame_NNNN.png names. Cancel mid-loop already leaves partial frames + scene.blend (saved before the render loop in apply_scene.py).
```
- Predicted:
```
Any exit-0 with a non-empty subset (or agents that trust PNG/.blend presence over manifests) treats incomplete shot output as finished; C03/C06 poison remains in the job out/ tree.
```
- Task:
```
Require exact expected frame set before WriteSuccessManifest; on failure keep failed marker and document that PNG/.blend alone are not success.
```

### proc-tree-orphan (failure, high)

- Location:
```
internal/runner/runner.go:Run
```
- Bug:
```
Comment and ARCHITECTURE/C03 claim cancellation kills the Blender process tree; implementation is bare exec.CommandContext with no Setpgid/process-group kill. TestSmokeCancel only asserts parent error, not child-tree exit.
```
- Predicted:
```
Timeout/cancel (and CLI SIGINT) can leave orphan Blender/render children holding GPU/CPU and writing into a job dir after the parent already returned failure.
```
- Task:
```
Start Blender in its own process group and kill the group on ctx Done; assert no live descendants in timeout+cancel tests.
```

### shot-id-outdir-escape (security, high)

- Location:
```
cmd/blendctl/main.go:cmdApply
```
- Bug:
```
Phase 8 nests --out DIR --shot ID via filepath.Join(outParent, shot) with no shot-id allowlist; project_v1.json and LoadAndResolveShot accept any nonempty shot id. Absolute ids discard DIR; ../ ids escape the intended nest.
```
- Predicted:
```
A malicious or agent-authored project can redirect job dirs (and Blender writes under them) outside the caller’s --out tree while still passing resolve/validate.
```
- Task:
```
Constrain shot ids to ^[A-Za-z][A-Za-z0-9_]*$ in schema+loader and reject Join results that leave the resolved --out root.
```

## Residual risks

- **medium** (concurrency)
- Id:
```
blender-cancel-env
```
- Bug:
```
Concurrent isolation relies on unique job dirs and --factory-startup, but Run inherits process env (shared Blender/XDG caches) and uses CommandContext without Setpgid/process-group kill.
```
- Task:
```
Set a per-job HOME/XDG_* under JobDir and kill the Blender process group on ctx cancel (SysProcAttr.Setpgid + negative PID signal).
```

- **medium** (security)
- Id:
```
blender-env-inheritance
```
- Bug:
```
ARCHITECTURE asked for explicit job environment; Run leaves Options.Env nil so Blender inherits the full parent environ (secrets, tokens, proxy creds).
```
- Task:
```
Pass a minimal allowlisted env (PATH, MOBLEND_*, display/locale as needed) instead of inheriting by default.
```

- **medium** (security)
- Id:
```
compile-out-overwrite
```
- Bug:
```
Plan required refusing nonempty destinations unless explicit overwrite; compile.WriteToDir unconditionally overwrites apply_scene.py, scene_data.json, source_map.json, and dependencies.json in the caller-supplied --out directory.
```
- Task:
```
Require empty --out or an explicit --overwrite limited to MoBlend-owned artifact names before WriteToDir.
```

- **medium** (failure)
- Id:
```
compile-partial-out
```
- Bug:
```
compile writes script then data straight into user --out; ARCHITECTURE nonempty-refuse/overwrite flag never implemented. Crash mid-write leaves mixed generations in the destination.
```
- Task:
```
Stage to a temp dir and rename into --out, or refuse nonempty destinations without an explicit overwrite mode.
```

- **medium** (failure)
- Id:
```
deps-after-success-ignored
```
- Bug:
```
dependencies.json is written only after RunApply success with _, _ = WriteDependencyManifest(...) — errors discarded; crash between manifest.json and dependencies.json leaves SHOT-02 provenance missing on an otherwise complete job.
```
- Task:
```
Write dependencies inside RunApply before advertising success (or fail the job if provenance write fails).
```

- **medium** (correctness)
- Id:
```
detectkind-project-version-alone
```
- Bug:
```
CONTEXT requires kind:project and project_version; DetectKind classifies KindProject if project_version is non-empty even without kind, routing scene-like JSON into project load errors.
```
- Task:
```
Require kind=="project" (and matching project_version) for project detection; add a negative fixture for schema_version+stray project_version.
```

- **medium** (tests)
- Id:
```
escape-code-ambiguous
```
- Bug:
```
Escape rejection accepts PROJ_PATH_ESCAPE or PROJ_FILE_MISSING, so a containment miss that degrades to missing-file still passes the adversarial path case.
```
- Task:
```
Pin PROJ_PATH_ESCAPE for ../ traversal; keep FILE_MISSING only on the dedicated missing_file fixture with an explicit code assert.
```

- **medium** (tests)
- Id:
```
eval-png-not-brief
```
- Bug:
```
Named as supported-brief evaluation but only validate/compile/preview + PNG existence for hand-authored fixtures—no checklist assertions (text, lens, frame counts, chart values). Playbook warns fixtures ≠ NL authoring; STATE still blocks LLM repair-loop evidence.
```
- Task:
```
Rename/clarify as fixture smoke; add per-fixture structural oracles for brief checklists; keep LLM repair-loop explicitly blocked until model runs exist.
```

- **medium** (correctness)
- Id:
```
inspect-swallows-module-failures
```
- Bug:
```
Inspect silently skips unreadable/invalid module files (empty kind/namespace/imports) instead of emitting diagnostics; uses loose json.Unmarshal, not strictDecode.
```
- Task:
```
Fail or attach diagnostics when a listed module path escapes, is missing, or fails strict module schema decode.
```

- **medium** (tests)
- Id:
```
modular-preview-existence-only
```
- Bug:
```
Plan gate asks for equivalent Blender scene properties/render outputs; test proves compile-byte equality then only os.Stat on PNG/blend—no reopen probe comparing IDs/transforms/materials between modular and flat.
```
- Task:
```
Reuse title/scene probe helpers to assert equivalent probed properties (and optionally image tolerance) for both previews.
```

- **medium** (concurrency)
- Id:
```
path-contain-toctou
```
- Bug:
```
Containment uses Lstat+EvalSymlinks on real but returns unverified abs ; readBounded later opens abs . No O_NOFOLLOW/openat of the checked real path.
```
- Task:
```
Open/read the EvalSymlinks-validated real path (or O_NOFOLLOW) and refuse if the opened file identity diverges from the containment check.
```

- **medium** (correctness)
- Id:
```
phase5-complete-overclaim
```
- Bug:
```
STATE reports roadmap gates 0–8 complete; Phase 5’s proposed gate (≥10/12 model repair-loop + visual HD) remains blocked and was never model-proven—only automated fixtures + smoke visual.
```
- Task:
```
Keep 0–8 CLI/composition complete, but state Phase 5 proposed authoring gate explicitly unmet/blocked beside AGENT-01.
```

- **medium** (correctness)
- Id:
```
roadmap-intro-next-phase-stale
```
- Bug:
```
Intro still says “Phase 7 is the next implementation phase; Phase 8 follows” and sequencing says “New implementation: 7 → 8”, while Phase 7–8 sections and STATE claim complete (2026-11-04).
```
- Task:
```
Align the header and sequencing to “0–8 executed; open: LLM eval, release Blender baseline, non-Linux/publish”.
```

- **medium** (tests)
- Id:
```
s04-failed-shot-missing
```
- Bug:
```
S04 requires concurrent shot builds plus one failed shot with success/failure attributed correctly; the integration test only runs three successful previews and checks distinct job dirs.
```
- Task:
```
Extend concurrent shot test with one deliberately failing shot and assert isolation of success manifests vs failure for the right IDs.
```

- **medium** (concurrency)
- Id:
```
shared-mod-torn-load
```
- Bug:
```
Plan/ARCHITECTURE treat a shared module as immutable input during a run; execution only documents agent ownership. loadOne sequentially ReadFiles each closure member with no FS snapshot/copy into the job tree.
```
- Task:
```
Copy or mmap a consistent closure snapshot into the job dir before resolve/compile, or fail the run if file mtimes/hashes change between first read and manifest write.
```

- **medium** (correctness)
- Id:
```
shot02-provenance-rehash-disk
```
- Bug:
```
dependencies.json file digests are computed by re-reading paths from disk after resolve, not from bytes captured at load; Effective also hardcodes engine/device CYCLES/CPU.
```
- Task:
```
Hash the load-time byte snapshot stored on Resolved, and copy effective settings from the declaration/runner config actually used.
```

- **medium** (failure)
- Id:
```
success-write-no-rollback
```
- Bug:
```
If WriteSuccessManifest fails after Blender produced frames/PNG, RunApply returns error without WriteFailedMarker; early extract/compile failures after JobDir create also leave unmarked orphan trees.
```
- Task:
```
On any post-jobDir failure path, always WriteFailedMarker (best-effort) before returning; never leave success-looking bytes without a terminal status file.
```

- **medium** (security)
- Id:
```
symlink-containment-toctou
```
- Bug:
```
Containment EvalSymlinks then returns the pre-symlink abs path; later load/hash (readBounded, BuildDependencyManifest) re-open that path and follow whatever the symlink is at read time. rootReal also silently falls back to unevaluated root if EvalSymlinks(root) fails.
```
- Task:
```
Open with O_NOFOLLOW (or re-EvalSymlinks+containment immediately before read) and fail closed if root cannot be realpath’d.
```

- **medium** (failure)
- Id:
```
timeout-untested
```
- Bug:
```
Phase 4 plan and RUN-03 claim timeout honesty; suite only has immediate cancel. Playbook run marks C03 PASS via cancel alone — timeout path unexercised.
```
- Task:
```
Add a short artificial Timeout fixture that asserts timed-out error, no manifest.json, and process-group gone.
```

- **medium** (tests)
- Id:
```
unsupported-fixture-vs-prompt
```
- Bug:
```
Automated suite only asserts ParseAndValidate ErrorList on schema/semantic poison JSON (unknown asset, raw python field, etc.); playbook unsupported prompts (arbitrary .blend/Python, download/rig, huge limits, MP4) have no automated agent-refusal oracle.
```
- Task:
```
Document fixture IDs as schema negatives only; add separate blocked/pending cases for NL unsupported refusals until an agent eval harness exists.
```

- **low** (correctness)
- Id:
```
cli-all-shots-undocumented
```
- Bug:
```
SHOT-01 all-shots mode is implemented and tested, but top-level usage/validate usage omit --all-shots (hint text mentions it only inside SHOT_REQUIRED).
```
- Task:
```
Document --all-shots in usage and GUIDE; add a playbook S-case that requires it.
```

- **low** (tests)
- Id:
```
ease-inout-probe-gap
```
- Bug:
```
Linear midpoint and step-hold have Blender probes; ease_in_out (declared in schema/fixtures) has no interior-sample probe matching B05.
```
- Task:
```
Add a Blender frame probe for ease_in_out interior samples against the documented curve contract.
```

- **low** (security)
- Id:
```
inspect-unevaluated-root
```
- Bug:
```
Inspect calls resolveContainedPath(root, root, …) without EvalSymlinks(root), unlike resolveManifest; containment failures are swallowed and the module is still listed.
```
- Task:
```
Reuse the same rootReal computation as resolve and surface containment errors instead of omitting fields.
```

- **low** (failure)
- Id:
```
jobdir-never-cleaned
```
- Bug:
```
Unique moblend-job-* dirs under --out accumulate forever; no retention, current symlink, or failed-tree scrub. Plan deferred hostile delete correctly but shipped no bounded cleanup policy.
```
- Task:
```
Emit a stable latest pointer atomically and document retention; optional --keep-failed only.
```

- **low** (tests)
- Id:
```
missing-diag-code-weak
```
- Bug:
```
Only asserts err != nil; no ErrorList/code pin (unlike cycle/symlink).
```
- Task:
```
Assert diag.ErrorList with PROJ_MODULE_MISSING or PROJ_FILE_MISSING as documented.
```

- **low** (correctness)
- Id:
```
playbook-module-activation-stale
```
- Bug:
```
Playbook status says M/S cases are active after 7–8, but baseline prose still says module/shot tests “activate when Phase 7/8 is implemented,” inviting BLOCKED skips of live gates.
```
- Task:
```
Replace activation hedges with “M/S mandatory on Linux amd64 pin; record BLOCKED only for LLM/non-Linux.”
```

- **low** (concurrency)
- Id:
```
shot-concurrent-test-gap
```
- Bug:
```
Gate evidence for SHOT-02 concurrent isolation uses separate Parents ( parent/<shot> ) and only checks outs[0]==outs[1] || outs[1]==outs[2], omitting outs[0] vs outs[2]; does not exercise same-parent MkdirTemp like TestConcurrentJobsIsolated.
```
- Task:
```
Run all three shots with the same Parent, assert pairwise distinct JobDir/PNG/manifest paths, and mirror that matrix for compile once nesting is fixed.
```


## Sequence

1. [tests]
```
Require Code == LIM_OBJECTS (or the single documented code) and fail on other codes.
```
2. [tests]
```
Add a Go test that validates broken_unselected with all-shots mode and asserts a diagnostic for shot bad while documenting selected-shot success separately.
```
3. [correctness]
```
Rewrite the BACKLOG sentence to “shipped; parameterized modules / build cache remain deferred” and leave only true leftovers.
```
4. [failure]
```
Propagate a signal-notified context into RunApply/RunSmoke so cancel takes the same path as Timeout.
```
5. [concurrency]
```
Make compile apply the same filepath.Join(out, shot) nesting as cmdApply, or always write under a MkdirTemp job dir like RunApply.
```
6. [concurrency]
```
Hash the bytes retained at load time (store digests or raw snapshots on Resolved) and never re-read project inputs when writing dependencies.json.
```
7. [tests]
```
Add negative fixtures and assert exact codes for depth>8, oversized file, project-byte overrun, >32 modules, and absolute module path.
```
8. [correctness]
```
On non-shot validate (and --all-shots), either resolve every listed module path for existence/version/containment or document and test that only closure integrity is guaranteed.
```
9. [correctness]
```
Extend TestModularFlatPreviewEquivalence with the same object/camera/material probes used for title scenes and a bounded PNG structural check; demote MOD-06 until that lands.
```
10. [failure]
```
Require exact expected frame set before WriteSuccessManifest; on failure keep failed marker and document that PNG/.blend alone are not success.
```
11. [failure]
```
Start Blender in its own process group and kill the group on ctx Done; assert no live descendants in timeout+cancel tests.
```
12. [security]
```
Constrain shot ids to ^[A-Za-z][A-Za-z0-9_]*$ in schema+loader and reject Join results that leave the resolved --out root.
```

## Guardrails

- Do not drop a finding without saying why it is invalid.
- After the next impl round, run `node scripts/judge.mjs compose` on an updated review (plan + new diff).
- This is round 1. If compose still says revise, another after-pass is allowed (max 3 rounds).
