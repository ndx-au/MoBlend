# Phase 6 — Package and hand off v1

**Status:** complete for Linux amd64 RC 2026-11-04 (see VERIFICATION.md; publish blocked on authorization)  
**Depends on:** Phase 5  
**Requirements:** DIST-01, DIST-02 (Linux amd64 only unless others tested)

## Outcome

Installable documented release candidate for Linux amd64 with embedded resources, checksums, attribution, and offline onboarding using documented prerequisites only.

## Tasks

1. Release build script + version reporting + checksums + font/pack notices
2. Clean-environment offline onboarding notes/test on Linux amd64
3. Re-check examples against schema; rerun integration baselines
4. Troubleshooting / limitations / compatibility / asset-build docs + quick start
5. Reconcile REQUIREMENTS statuses; prepare artifacts (no public publish without authorization)

## Gate

REQUIREMENTS release gate for claimed platforms; independent onboarding with documented prerequisites; report binary size and a sample render time with hardware context.
