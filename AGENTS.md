# Agent instructions

## Mission and source of truth

Build the product described in `.planning/PROJECT.md`. Read PROJECT, ARCHITECTURE, REQUIREMENTS, ROADMAP, and STATE before implementation. This is a new project; do not import requirements or tooling from unrelated repositories.

Fixed constraints: Go compiler/orchestrator, JSON declaration, Blender rendering, bundled first-party asset pack, **no RAG in v1**. Python is generated glue executed inside Blender, not the primary application language. Do not add a hosted model service or custom programming language.

## Dot Progress (audit trail)

Implement [`.progress/README.md`](.progress/README.md). That file is authoritative for progress logging if anything else conflicts on logging rules.

- Plan and gates live under `.planning/` (especially `.planning/ROADMAP.md` and phase plans).
- Record history by **appending** entries under `.progress/milestone/`, `.progress/general/`, or `.progress/fixes/` — never edit, rename, or delete committed progress files.
- Map roadmap phases to milestone tokens: Phase 0 → `M000`, Phase 1 → `M001`, … Phase 6 → `M006`, Phase 7 → `M007`, Phase 8 → `M008`. Frontmatter may use short forms such as `M0 — Prove the Blender boundary`.
- Append an entry for each meaningful session of work (including partial milestone progress, decisions, reversals, and fixes). Skip only trivial typo-only changes.
- Optional harness: `.claude/rules/dotprogress.md` and `.claude/skills/dotprogress/`.
- Keep `.progress/` tracked in git (do not gitignore it).

## Execution

- Work through roadmap gates in order. First prove a real headless render; do not build a large abstraction layer before that proof.
- Before each phase, create `.planning/phases/NN-name/PLAN.md` with tasks, file ownership, dependency order, requirement IDs, checks, and completion conditions. Record consequential choices in `CONTEXT.md`.
- Make reasonable reversible choices independently. Ask only for unresolved product scope changes, necessary credentials, or externally consequential actions without authorization.
- Record implementation progress and exact verification commands in STATE and a phase `VERIFICATION.md`. A skipped test is not a pass. Do not mark a phase complete with an unmet gate.
- Also append a Dot Progress entry for the same session when the work is meaningful.
- Keep commits focused when Git is initialized. Preserve unrelated work. Do not publish releases or upload user media without authorization.
- Capture new ideas in `.planning/BACKLOG.md`; do not silently expand v1.

## Engineering invariants

- Treat declarations as data. No raw Python, expressions, shell commands, arbitrary bpy paths, or external scripts in the language.
- Compile deterministically. Use stable IDs and ordering; keep paths, timestamps, and run IDs out of generated-code identity.
- Prefer Blender's data API over context-sensitive operators. Explicitly set state whenever an operator is unavoidable.
- Build each job from the pristine embedded asset pack in a new Blender process. V1 does not reconcile a user's existing scene.
- Validate schema, references, asset contracts, limits, and capabilities before launching a render.
- Keep the schema, Go model, examples, catalogue, and agent guide consistent through automated checks.
- Spawn Blender with an argument vector, never a constructed shell command. Keep user strings out of executable source syntax.
- Verify claims with unit tests and real Blender integration tests. Golden Python alone does not prove a valid scene.
- Do not promise identical pixels across GPUs, operating systems, or Blender versions.

## Completion evidence

Run Go formatting, tests, vet, and the phase-specific Blender checks. Include missing prerequisites and remaining visual review explicitly. Pin exact supported tool versions after Phase 0 evidence. Read official documentation for that pinned version before relying on bpy behavior.
