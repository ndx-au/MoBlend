# Gauntlet round 1 repair plan

Scope: twelve high-priority findings from `.gauntlet/outbox/superplan.md`, plus directly related failure/provenance/doc corrections. Existing Phase 0–8 plans remain authoritative; no new product scope.

| Order | Owner/files | Tasks | Requirements | Checks |
|---|---|---|---|---|
| 1 | primary: internal/project, schema | Exact diagnostics, bounded graph tests, safe shot IDs, load-time digests; validate full listed graph for non-shot/all-project operations | MOD-01–06, SHOT-01–03 | Go project + CLI regression tests |
| 2 | primary: cmd/blendctl | Signal context, unique nested compile outputs, all-shots regression | RUN-01/03, SHOT-01/02 | CLI tests incl concurrent compilation |
| 3 | primary: internal/runner | Process groups, exact frame sets, terminal failure markers, provenance before success | RUN-01–03, SHOT-02 | Artificial failure/cancel/timeout tests and real Blender |
| 4 | primary: runner integration + planning/docs | Reopen scene probes and PNG comparison; correct stale completion claims | MOD-06, QA-01, AGENT-01 | Pinned Blender integration, fmt/test/vet |

Completion: all twelve high items have executable regression evidence; record exact commands and any unavailable verification. Preserve remaining medium/low risks explicitly. Run judge compose if the harness exists.
