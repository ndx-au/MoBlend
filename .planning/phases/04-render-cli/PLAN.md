# Phase 4 — Make rendering dependable

**Status:** complete (gate met 2026-11-04; see VERIFICATION.md)  
**Depends on:** Phase 3  
**Requirements:** ENV-01, RUN-01, RUN-02, RUN-03, UX-01

## Outcome

User-facing local CLI handles success and failure predictably: `doctor`, `preview`, `render` with structured diagnostics, job dirs, manifests, timeouts/cancellation, and honest failure.

## Tasks

1. `doctor` — Blender discovery/version gate; actionable missing/unsupported errors
2. Harden `preview`/`render` — frame range, effective settings in manifest, stdout/stderr discipline
3. Job hygiene — unique dirs, integrity checks, no pack mutation, no overwrite without flag
4. Failure matrix — missing Blender, bad version, unwritable dest, Python fail, cancel, concurrent jobs
5. Reopen delivered `.blend` files in tests

## Gate

Automated success/failure matrix passes; incomplete output never looks complete; jobs do not mutate embedded pack.
