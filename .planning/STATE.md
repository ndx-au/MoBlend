# Project state

Updated: 2026-11-05
Status: v0.2.0-alpha.1 published as a Linux amd64 prerelease in the private GitHub repository. Roadmap phases 0–8 CLI/composition work executed on Linux amd64. Phase 5 proposed model repair-loop gate remains unmet/blocked; fixture smoke and visual review are not model-authoring evidence. RC `0.1.0-rc1` baseline retained; modularity (Phase 7) and independent shots (Phase 8) implemented with evidence.
Active phase: none. Release v0.2.0-alpha.1 is published with verified native archive/checksums; audit: `.planning/releases/0.2.0-alpha.1/VERIFICATION.md`. The release tag remains at 1ccc0038c9e9313ec59a709672fccd3c6669adfa; main includes a test-only cancellation timing follow-up with green CI.
Next action: controlled alpha evaluation, stable Blender qualification, measured model-authoring/onboarding and broader visual review. Other platforms remain unclaimed. Repository remains private; this alpha is not Latest. The historical independent Gauntlet verdict predates the fixes.



## Completed

Roadmap phases 0–8 implementation (see phase VERIFICATION files; Phase 5 proposed authoring gate is unmet/blocked). README product refresh on `main` (`6066c99`). Phase 7–8 evidence: `.planning/phases/07-modular-projects/VERIFICATION.md`, `.planning/phases/08-independent-shots/VERIFICATION.md`. Progress: `M007.002`, `M008.001`.

## Post-RC follow-ups

| Item | Status |
|---|---|
| Concurrent jobs (`RUN-01`) | `TestConcurrentJobsIsolated` |
| Contact sheet + visual review (`QA-01`) | `testdata/visual/` + `phases/05-agent-usability/VISUAL-REVIEW.md` |
| Modular projects (`MOD-01`–`07`) | met — Phase 7 VERIFICATION |
| Independent shots (`SHOT-01`–`03`) | met — Phase 8 VERIFICATION |
| LLM authoring eval | still blocked / not claimed |
| macOS/Windows | not claimed |
| GitHub alpha prerelease | **published v0.2.0-alpha.1** — verified assets, not Latest; private repository |

## Modularity decision — 2026-11-04

Owner approved Phases 7–8 after the initial RC. Implemented: explicit local module graph, namespaces, source diagnostics, inspect/CLI, paired fixtures; then shot selection, closure loading, dependency provenance, three-shot example, change isolation. Scene schema `1.0.0` and compiler backend unchanged. No RAG. Parameterized components and build caching remain deferred.

## Gauntlet repair — 2026-11-05

All twelve high-priority leftovers repaired with regression coverage. Commands: `gofmt -l .` (empty), `go vet ./...` (exit 0), `go test ./... -count=1 -timeout 10m` (103 tests/subtests passed including real Blender), `go test -race ./internal/project ./cmd/blendctl -count=1` (42 passed), `git diff --check` (exit 0). Blender 5.3.0 Beta build 7e7c173d45c8, Linux amd64. No additional visual/model evaluation claimed.

New compose attempt rejected all five critic files as stale relative to the updated review; round-one outbox retained. Repair verification lists remaining medium/low risks. Dot Progress: `F000001.gauntlet-round-one-hardening.md`.

## Fresh Gauntlet review — 2026-11-05

Verdict **revise**, no criticals, four highs: duplicate import aliases; modular camera animation targets; cross-property interpolation overwrite; material interpolation owner. New Blender probes demonstrate wrong interior animation values despite 103 existing tests passing. SCENE-02 and module alias/reference/source-pointer acceptance are reopened; earlier met statements describe the prior tested baseline, not proof of the uncovered paths.

All five critics refreshed and judge compose succeeded. Full details: `.planning/fixes/gauntlet-fresh-review/VERIFICATION.md`; leftovers: `.gauntlet/outbox/superplan.md`. Round 1 artifacts preserved under `.gauntlet/history/round-1/`. The judge counter reset to round 1 for the changed review hash; this was review pass 2 in the conversation. No implementation changes made in this review session.

## Follow-up implementation — 2026-11-05

Repaired all four high and actionable medium/low fresh-superplan tasks: deterministic explicit aliases, camera references, original diagnostics, module contracts, mixed transform/material interpolation, real easing samples, bounded regular snapshot reads, version-probe cancellation, bounded logs, decoded PNG validation and concurrent failed-shot attribution. Exact finding-to-test map: `.planning/fixes/gauntlet-followup/VERIFICATION.md`.

Final checks: `gofmt -l .` empty; `go vet ./...` exit 0; `go test ./... -count=1 -timeout 10m` passed **134 tests/subtests across 8 packages**, with real Blender; targeted race suite passed **40 across 3 packages**; `git diff --check` clean; guide copies match. Current Linux Go/shell release build produced a **4,120,760-byte static executable** and its minimal-environment preview completed. No Node build/runtime dependency added. Commit/push authorized by the owner; release tagging/publication not requested. Dot Progress: `F000002.gauntlet-followup-fixes.md`.

## Native archive packaging — 2026-11-05

Implemented Go/Bash-only versioned archives/checksums with explicit snapshot and clean annotated-tag build modes; tags/releases are never published by the script. Included docs/examples/notices/build provenance; real extracted-binary render succeeded. Packaging acceptance proves repeated archive byte equality, checksums, extraction and negative tag/dirty/version/output guards. Bash syntax, Go formatting/vet and 134 Go/Blender tests passed. See `.planning/fixes/release-packaging/VERIFICATION.md` for exact commands and missing ShellCheck/other qualification evidence.

Recommended current maturity: alpha; suggested next public version `v0.2.0-alpha.1`, preserving historic `0.1.0-rc1`. This is release-strategy advice, not an actual version/tag/publication change. Progress: `M006.002.native-release-archives.md`.

## Alpha.1 release preparation — 2026-11-05

Owner authorized v0.2.0-alpha.1. Finalize pending native packaging/docs in a clean source commit, run fresh qualification, annotate the tag, build with --release --verify-render, verify extracted package, then publish a GitHub prerelease with tar.gz/SHA256SUMS. Repository is private and its visibility will be preserved. No stable/platform/model/HD claims or new project license are added. Plan/context: `.planning/releases/0.2.0-alpha.1/`.

## Alpha.1 published — 2026-11-05

Release: https://github.com/ndx-au/MoBlend/releases/tag/v0.2.0-alpha.1 . Exact tagged source 1ccc0038c9e9313ec59a709672fccd3c6669adfa; native Linux amd64 tar.gz and SHA256SUMS uploaded/downloaded/verified, preview + full tiny sequence + saved-scene reopen passed. GitHub confirms prerelease=true, draft=false, repository private, not Latest.

Tag CI passed. Parallel branch CI exposed two immediate child-state test observations; test-only main follow-up 0e8d42a uses bounded termination waits, passed 160 repeated tests/subtests, 8 race checks, full 134-test Go/Blender suite and hosted CI. Recompiled executable matched the immutable tagged binary; tag/assets were not replaced. Limitation/follow-up disclosed in notes. Publication audit: `.planning/releases/0.2.0-alpha.1/VERIFICATION.md`; Dot Progress M006.004.
