# Project state

Updated: 2026-11-04
Status: Roadmap phases 0–8 complete on Linux amd64 under recorded verification. RC `0.1.0-rc1` baseline retained; modularity (Phase 7) and independent shots (Phase 8) implemented with evidence.
Active phase: none (roadmap gates 0–8 met).
Next action: owner decisions on public publish, LLM authoring eval tooling, and any non-Linux platform claims. Track deferred ideas in BACKLOG (parameterized modules, build cache, video assembly).

## Completed

Roadmap phases 0–8 (see phase VERIFICATION files). README product refresh on `main` (`6066c99`). Phase 7–8 evidence: `.planning/phases/07-modular-projects/VERIFICATION.md`, `.planning/phases/08-independent-shots/VERIFICATION.md`. Progress: `M007.002`, `M008.001`.

## Post-RC follow-ups

| Item | Status |
|---|---|
| Concurrent jobs (`RUN-01`) | `TestConcurrentJobsIsolated` |
| Contact sheet + visual review (`QA-01`) | `testdata/visual/` + `phases/05-agent-usability/VISUAL-REVIEW.md` |
| Modular projects (`MOD-01`–`07`) | met — Phase 7 VERIFICATION |
| Independent shots (`SHOT-01`–`03`) | met — Phase 8 VERIFICATION |
| LLM authoring eval | still blocked / not claimed |
| macOS/Windows | not claimed |
| Public publish | needs owner authorization |

## Modularity decision — 2026-11-04

Owner approved Phases 7–8 after the initial RC. Implemented: explicit local module graph, namespaces, source diagnostics, inspect/CLI, paired fixtures; then shot selection, closure loading, dependency provenance, three-shot example, change isolation. Scene schema `1.0.0` and compiler backend unchanged. No RAG. Parameterized components and build caching remain deferred.
