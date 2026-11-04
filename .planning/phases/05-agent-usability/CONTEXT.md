# Phase 5 — context

## Evaluation approach

Automated gate uses fixture briefs under `testdata/eval/`:

- `supported/` — ≥12 JSON briefs derived from title/pedestal/chart
- `unsupported/` — separate deliberate rejects (unknown asset, raw python field, bad schema, duplicate id, expression track)

`TestEvalSupportedBriefsPreview` runs validate+preview for each supported brief (counts as structural adherence without an LLM author).

## Blocked evidence (explicit)

- Live LLM authoring/repair loop with a named model: **blocked** (no model invoked this session)
- Human visual contact-sheet review of typography/framing: **blocked** pending owner review of preview PNGs

Do not fabricate model scores. Automated preview success is recorded; visual QA remains open.
