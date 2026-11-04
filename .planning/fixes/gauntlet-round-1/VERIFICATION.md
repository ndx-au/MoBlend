# Gauntlet round 1 repair verification

Date: 2026-11-05 (Australia/Sydney). Implementation: twelve high-priority findings repaired; independent re-judging pending refreshed critic reports.

## Environment and commands

Linux amd64; Blender **5.3.0 Beta**, build **7e7c173d45c8**, same local development pin as Phase 0. This does not establish a stable-release Blender baseline, non-Linux support or public publication. Official pinned Scene API documentation was fetched/read from `https://docs.blender.org/api/5.3/bpy.types.Scene.html` (HTTP 200), including the requirement to use `frame_set` to evaluate animation.

Fresh checks after the final Go changes (commands executed via `rtk`):

| Exact command | Result |
|---|---|
| `gofmt -l .` | exit 0, empty output |
| `go vet ./...` | exit 0, no issues |
| `go test ./... -count=1 -timeout 10m` | exit 0, **103 tests/subtests passed across 8 packages**, including real Blender integrations |
| `go test -race ./internal/project ./cmd/blendctl -count=1` | exit 0, **42 tests/subtests passed across 2 packages** |
| `git diff --check` | exit 0 |

Blender was available throughout; integration tests were executed, not skipped. The real Blender equivalence test passed its reopened scene assertions and image checks. No new human aesthetic/HD review was performed; no model-based authoring/refusal evaluation was performed.

## Twelve high-priority findings

| Judge order | Repair / regression evidence |
|---|---|
| 1 | `TestAggregateObjectLimit` requires exactly `LIM_OBJECTS`; removed message heuristics and alternate codes. |
| 2 | `TestValidateAllShotsReportsBrokenUnselected` invokes `cmdValidate --all-shots --json`, requires `PROJ_FILE_MISSING` attributed to `bad`, and separately verifies selected `good` succeeds. |
| 3 | BACKLOG now records shipped composition/shots and deferred parameterized modules/cache. |
| 4 | CLI preview/render/smoke use `signal.NotifyContext`; `TestCLISignalsCancelJobs` sends real SIGINT/SIGTERM to a built CLI, requires cancellation output, dead descendants and preview/render failed markers. |
| 5 | Compile creates unique shot-nested jobs. `TestCompileConcurrentShotOutputs` runs intro/outro/intro concurrently into the same parent and checks all artifact sets. |
| 6 | Resolver hashes each loaded byte buffer, including manifest; dependency manifest copies retained digests. `TestDependencyManifestUsesLoadSnapshot` removes source files after loading and requires unchanged provenance. |
| 7 | `TestGraphBounds` asserts exact depth>8, file bytes, project bytes, >32 modules and absolute-path diagnostics. Actual reads are byte-bounded. Depth also checks longest DAG path to prevent diamond shortcuts. |
| 8 | Non-shot resolution and whole-project validation check every listed module graph for missing files, versions, containment, cycles and graph limits, while only entry closure entities are merged. `TestUnreachableListedGraphRejected` tests missing/version/cycle orphans; selected-shot isolation remains tested. |
| 9 | `TestModularFlatPreviewEquivalence` reopens both Blender files, asserts expected title/camera state, compares object IDs/types/transforms/material colors/roughness, camera lens, resolution/FPS/range. Decoded 64×64 PNGs require >=8 distinct colors and <=1/255 channel differences under the same CPU pin. |
| 10 | Render success requires exactly the inclusive expected frame names with regular nonempty files, plus scene.blend. `TestApplyRejectsIncompleteFrameSets` rejects partial, misnamed and extra sets with failed markers and no success manifest. Real `TestRenderWritesManifestAndFrames` still passes. |
| 11 | Linux/macOS launcher configures a dedicated process group and kills it on context cancellation. `TestRunTerminatesDescendants` verifies both explicit cancellation and timeout kill parent/child; `TestApplyTimeoutHasFailedMarker` checks terminal status. Linux zombie states count as terminated, not running children. Other platforms remain unclaimed. |
| 12 | Shot ID regex is enforced in schema and loader globally. Output joining rechecks IDs and resolved containment, including escaping symlinks. `TestUnsafeShotIDsRejectedGlobally` and `TestOutputParentRejectsEscape` cover traversal/absolute/separator IDs and output symlinks. |

## Related medium/low items

Also repaired: default environment uses a minimal allowlist and private HOME/XDG/TMPDIR (`TestRunEnvironmentIsolated`); compile cannot overwrite previous CLI outputs; dependency writes occur before success and errors propagate; post-job failures attempt terminal failed markers; DetectKind requires explicit project kind; exact traversal/missing diagnostics; inspect surfaces path/read/unknown-field/version errors and canonicalizes its root; duplicated/trailing JSON is rejected for modules/projects; all-shots help; stale roadmap/Phase 5 claims; fixture-smoke naming; unsupported JSON versus natural-language evaluation distinction; same-parent pairwise concurrent preview isolation. Guide copies remain synchronized.

Remaining medium/low work from the original superplan is retained, not dropped:

- `eval-png-not-brief`: per-brief independent structural/checklist oracles and model repair-loop/unsupported prompt runs remain pending. Renaming fixture smoke does not satisfy them.
- `s04-failed-shot-missing`: concurrent real-shot integration still covers three successful shots. Synthetic runner/CLI failure tests cover failed jobs separately; the combined concurrent successful/failed-shot matrix remains open.
- `shared-mod-torn-load`: per-file snapshots are immutable after read; no atomic filesystem-wide snapshot while agents edit several files. Coordinate input edits before load; document this limitation.
- `path-contain-toctou` / `symlink-containment-toctou`: reads use the containment-checked real path and fail closed on root resolution. No openat/O_NOFOLLOW defense against concurrent replacement of directory components; execution hygiene is not a hostile-filesystem sandbox.
- `ease-inout-probe-gap`: a dedicated interior easing probe remains open.
- `jobdir-never-cleaned`: manual retention is documented; automatic cleanup/latest pointer is deferred.
- Inspect is strict about JSON fields and versions but does not replace full semantic validation of module contents; use validate for that gate.
- Compile's internal writer accepts caller-owned destinations; the public CLI guarantees fresh directories. Terminal marker writes are best-effort when storage itself fails. Logs remain captured in memory. The separate Doctor `--version` probe still inherits its parent environment; the new allowlist is applied by Run to execution jobs.
- `shot02-provenance-rehash-disk`: the re-read bug is fixed. CYCLES/CPU effective fields match the fixed template engine/device contract; configurable engine/device is not a current feature.

Phase 5 proposed >=10/12 diagnostic repair-loop gate remains **unmet/blocked**, with smoke-only visual approval preserved. Stable Blender baseline, non-Linux verification and public publication remain open as before.

## Gauntlet compose

The literal `node scripts/judge.mjs compose` failed because this repository has no judge script. Located the installed implementation and ran:

`node /home/terence/code/ndx-skills/.cursor/skills/ndx-gauntlet/scripts/judge.mjs compose`

Result: exit 1, `gauntlet fail-closed: critic files incomplete`; all five critics are **stale (older than review)**. Updated `.gauntlet/inbox/review.md` and `.gauntlet/inbox/implementation.diff` contain the repair context and tracked/new Go delta. Original critic JSON and round-one outbox remain unchanged; no new ship verdict is claimed. Refresh the critics against this diff before composing the next verdict. Findings were not erased, relabeled as passed, or replaced with stubs.
