# Gauntlet review: MoBlend `.planning/ROADMAP.md`

Round: 1 (after-pass). Plan under review: full implementation roadmap Phases 0–8.

## Original plan (summary)

Roadmap gates for a Go compiler/orchestrator that drives headless Blender from JSON declarations, with a bundled first-party asset pack and **no RAG in v1**.

| Phase | Planned outcome | Gate gist |
|---|---|---|
| 0 | Prove Blender boundary | Real PNG, reopenable scene, Python failure exit, offline with Blender |
| 1 | Smallest useful declaration | Schema/catalogue/types; positive+negative fixtures; local validate |
| 2 | First useful composition | Title JSON → editable scene → still; byte-stable compile |
| 3 | Bounded motion | Keyframes + chart GN asset; three reference animations |
| 4 | Dependable rendering | doctor/preview/render; manifests; failure honesty; isolation |
| 5 | Agent usability without RAG | Guide + ≥12 supported briefs + unsupported set; LLM eval if available |
| 6 | Package/hand off v1 | Linux amd64 RC; checksums; docs; no unauthorized publish |
| 7 | Explicit local modules | Project/module schemas; resolver; namespaces; source diagnostics; inspect; modular≡flat |
| 8 | Independent shots | Shot IDs; closure loading; isolated outputs; three-shot change isolation |

Fixed constraints: Go primary, JSON data language (no raw Python in declarations), Blender render, embedded pack, no hosted model service, no MCP/GUI/RAG in these phases. Parameterized components and build caching deferred.

Roadmap status header (as checked in): Phases 0–6 executed as Linux amd64 `0.1.0-rc1`; Phase 5 model-based authoring evaluation remains blocked; Phases 7–8 marked **complete (2026-11-04)** in their sections. Intro still says “Phase 7 is the next implementation phase” (stale relative to body).

## Execution evidence

### Git / working tree

- Repo clean at review time (`git status` empty).
- `git diff` vs HEAD: empty (no uncommitted execution delta).
- Commits (recent):
  - `ef13b56` Add modular projects, independent shots, and playbook evidence.
  - `bbcdd47` Harden post-RC isolation and reference visuals.
  - `6066c99` Refresh README as a product overview for the Linux amd64 RC.
  - `e1926df` Ship MoBlend Linux amd64 RC through roadmap phases 0–6.
  - `13e1b2c` Initial commit
- Scope from `13e1b2c..HEAD`: **180 files**, ~11.6k insertions (Go CLI, internal packages, schemas, examples, testdata, phase VERIFICATION, docs, RC dist path evidence).

### STATE / REQUIREMENTS claim

- `.planning/STATE.md`: phases 0–8 complete on Linux amd64; active phase none; next = owner decisions (publish, LLM eval tooling, non-Linux claims).
- `.planning/REQUIREMENTS.md`: ENV/CORE/ASSET/SCENE/RUN/AGENT/UX/QA/DIST marked met with caveats; MOD-01–07 and SHOT-01–03 marked met; LLM authoring **blocked**; DIST-02 Linux amd64 only.

### Phase verification (recorded, not re-run in this review)

| Phase | VERIFICATION status | Notable caveats |
|---|---|---|
| 0 | gate met | Blender **5.3.0 Beta**; “Release Blender baseline remains open” |
| 1 | gate met | Local validate; deferred fields in CONTEXT |
| 2 | gate met | Title probe; Unicode in data not script |
| 3 | gate met | Chart isolation; compositor look deferred |
| 4 | gate met | Manifests; exit categories; concurrent jobs later hardened |
| 5 | automated gate met; LLM **blocked** | 12/12 preview fixtures; proposed ≥10/12 repair-loop **not** model-proven |
| 6 | Linux amd64 RC ready | macOS/Windows not claimed; public publish not done |
| 7 | MOD-01–07 met | modular≡flat + Blender preview equivalence under pin |
| 8 | SHOT-01–03 met | no cache/scheduler/video assembly |

Key implementation surfaces claimed:

- `cmd/blendctl` — schema/assets/validate/compile/preview/render/doctor/guide/inspect + project/shot flags
- `internal/project/` — load/resolve/inspect/provenance (~500 LOC resolve + tests)
- `schema/project_v1.json`, `schema/module_v1.json` embedded
- Examples: `examples/projects/title_modular/`, `examples/projects/three_shots/`
- Negatives: cycle, missing, escape, symlink escape, broken unselected shot
- Playbook: `docs/TEST-PLAYBOOK.md` + run note `docs/TEST-PLAYBOOK-RUN-2026-11-04.md`

### Document inconsistencies / residual product risks visible without code dive

1. **ROADMAP intro vs phase bodies**: header still frames Phase 7 as “next”; Phases 7–8 sections say complete.
2. **BACKLOG stale**: still says Phases 7–8 “implementation and acceptance evidence are pending” while STATE/REQUIREMENTS claim met.
3. **Phase 5 proposed gate**: model repair-loop + visual HD review partially blocked / smoke-only; must not be over-claimed.
4. **Pinned Blender is beta (5.3.0)**: Phase 0 explicitly left release baseline open — stability/support risk for “supported environment.”
5. **Non-Linux / public publish**: explicitly unclaimed — fine if messaging stays honest.
6. **Deferred backlog items** still out of scope: parameterized modules, content-addressed cache, video assembly, MCP, RAG — should not silently land.

### Todo / plan status in ROADMAP

No Cursor todo checklist in ROADMAP. Phase PLAN.md files exist under `.planning/phases/00`–`08`. Status narrative: gates 0–8 marked complete with VERIFICATION companions; blocked items called out (LLM eval, multi-platform, publish).

## Critic charge

Compare the roadmap plan to this execution record. File leftovers where claims, gates, safety boundaries, failure honesty, concurrency isolation, or tests are incomplete, overstated, or hazardous — especially around project path containment, shot isolation, diagnostics, beta Blender pin, and blocked LLM/platform claims.


## Implementation repair after round 1 — 2026-11-05

This section supersedes the original clean-working-tree execution snapshot above. Original critic JSON and round-one superplan are retained unchanged pending independent refresh.

All twelve high-priority findings have implementation changes and regression tests. The review delta includes new regression files (not just tracked-file diffs): `.gauntlet/inbox/implementation.diff`.

Changes: exact LIM_OBJECTS assertion; real all-shots CLI regression; BACKLOG/ROADMAP/Phase 5 status corrections; signal context and Linux process-group cancellation; unique shot-nested compile jobs; load-time dependency digests including manifest; graph limit/absolute-path tests and listed orphan validation; saved Blender state and nonempty decoded-pixel equivalence; exact inclusive frame set before success; schema/loader/output shot containment.

Related medium/low repairs: default minimal private environment; provenance written before success with errors propagated; terminal markers on post-job failures; strict inspect errors and real root; duplicate/trailing JSON checks; exact escape/missing diagnostics; all-shots help; fixture-smoke naming and NL refusal caveats; pairwise same-parent concurrent preview isolation. Remaining medium/low gaps are listed in `.planning/fixes/gauntlet-round-1/VERIFICATION.md`; no finding has been silently deleted.

Verification: Go formatting/vet, complete Go + real Blender suite, race checks for project/CLI, and direct CLI signal tests. Exact final output and environment are recorded in the repair VERIFICATION. This is implementation evidence, not an independent new critic verdict.

Implementation diff SHA-256: `30ebea05ee9cdf23bcc56b1f3350664968dfb37b11c0fbaa78a6d10d4a002280`.
