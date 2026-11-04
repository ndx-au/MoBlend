# Test playbook run — 2026-11-04

Compiler: worktree after Phases 7–8 (uncommitted). Platform: Linux amd64. Blender: 5.3.0 Beta (`~/.local/bin/blender`). Pack SHA-256: `cc4290cb23bc07d0981b793bffe342a964e5e23789db25d51dda1550f5792088`. Render device: Cycles CPU.

Bulky previews/renders kept under `/tmp/moblend-playbook/` (not committed).

## Baseline

| Check | Result |
|---|---|
| `scripts/check.sh` | PASS |
| `go test ./internal/runner -count=1 -timeout 20m -v` | PASS (no skips) |
| `blendctl doctor` | OK Blender 5.3.0 Beta |

## Suite coverage (this run)

| Suite | Status | Notes |
|---|---|---|
| A decode/limits/refs | PASS | Extended negative fixtures: conflicting track, duplicate keyframe, reversed range, bad colour |
| B scene/animation probes | PASS | Added `TestLinearLocationMidpoint` (B03), `TestStepInterpolationHold` (B04); pedestal/chart probes retained |
| C runner failures | PASS | Existing doctor/smoke failure/cancel/concurrent; C04 path-with-spaces CLI PASS |
| D determinism | PASS | Key-reorder compile byte equality; modular relocate/reorder |
| M modular | PASS | Flat≡modular, diamond dedup, symlink escape, aggregate limits, source diagnostics |
| S shots | PASS | Selection/closure, all-shots broken unselected, change isolation, concurrent shot previews |
| P01–P12 authoring | PASS (agent-authored, attempt 0 after framing repair) | See below |
| U unsupported fixtures | PASS reject | Fixture set U01–U05; prompt U01–U04 refused without invented fallbacks |
| Visual aesthetic | PASS with caveats | Chart GN tint still weak (known debt); titles/pedestal repaired |

## P briefs (agent authoring)

First pass had cropped titles (missing text `rotation_euler_deg: [90,0,0]` / too-close camera) and intersecting pedestal. Repaired within the two-revision budget; guide updated with framing/pedestal/chart notes.

| ID | Validity | Preview/render | Visual/brief | Notes |
|---|---|---|---|---|
| P01 | PASS | preview f24 | PASS | Dark blue plate, readable MoBlend |
| P02 | PASS | preview | PASS | Exact string including quotes/`\` survives compile + readable |
| P03 | PASS | preview f11 | PASS | Covered by linear midpoint Blender probe |
| P04 | PASS | preview f20 | PASS | Covered by step-hold Blender probe |
| P05 | PASS | preview f24 | PASS | Launch readable at end |
| P06 | PASS | preview | PASS | Green cube on neutral plinth after camera pull-back |
| P07 | PASS | preview | PASS structural | Heights 0.2/0.9/0.4 visible; colour tint weak (known) |
| P08 | PASS | preview | PASS structural | Two charts, 3+2 bars isolated |
| P09 | PASS | preview | PASS | 50 mm + 320×240 product |
| P10 | PASS | render | PASS | Frames 5–8 + complete manifest |
| P11 | PASS | preview | PASS | Broken dup-id rejected; repaired ids unique |
| P12 | PASS | preview | PASS | Left values → [0.4,0.4,0.4]; right/camera unchanged |
| M-P01 | PASS | inspect+preview | PASS | Modular title |
| S-P01 | PASS | shot preview | PASS | Product-only edit/build |

**Score:** 12/12 supported briefs eventual PASS after framing repairs. First-pass visual failures on P01/P02/P05/P06 were authoring/guide gaps, not compiler regressions. Chart colour remains documented product debt, not a fail of structural oracles.

## Regressions added this run

- `testdata/spec/negative/{conflicting_track,duplicate_keyframe,reversed_range,bad_color}.json`
- `testdata/project/positive/diamond/`
- `testdata/compile/{linear_midpoint,step_jump}.json`
- `TestDiamondDedup`, `TestRejectDuplicateNamespace`
- `TestLinearLocationMidpoint`, `TestStepInterpolationHold`
- Guide: title framing, pedestal height, chart tint caveat

## Remaining / blocked

- Multi-model LLM authoring battery (frozen 3× comparative runs): not executed; this run used one agent authoring pass.
- Cross-GPU pixel identity: not claimed.
- Chart material tint on GN instances: backlog polish.
