# Phase 5 — Prove agent usability without RAG

**Status:** complete for automated gate 2026-11-04; LLM/visual evidence blocked (see VERIFICATION.md)  
**Depends on:** Phase 4  
**Requirements:** AGENT-01, QA-01

## Outcome

An agent can author/repair supported scenes from shipped local docs only (no RAG).

## Tasks

1. Write compact `agent/GUIDE.md` (concepts, commands, diagnostics, examples, unsupported)
2. `blendctl guide` prints embedded guide; schema/catalogue via exact-ID commands
3. Evaluation set ≥12 supported briefs + separate unsupported set
4. Record evaluation results; contact sheets / visual notes (block if no reviewer)
5. Fix compiler/asset issues found

## Gate

≥10/12 supported briefs succeed within two diagnostic-driven revisions; unsupported cases identified; core tests pass. Mark blocked evidence explicitly if no model/visual reviewer.
