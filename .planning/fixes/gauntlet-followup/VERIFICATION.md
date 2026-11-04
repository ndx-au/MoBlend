# Gauntlet follow-up fixes: verification

Date: 2026-11-05. Linux amd64, Go 1.25.6, Blender 5.3.0 Beta build 7e7c173d45c8. This session repairs all four high and actionable medium/low findings from the fresh superplan. No release tag/publication or additional platform/model/visual claims.

## Commands and evidence

Commands below were executed via `rtk` after the corresponding final source changes:

| Command | Result |
|---|---|
| `gofmt -l .` | empty; formatting clean |
| `go vet ./...` | exit 0, no issues |
| `go test ./... -count=1 -timeout 10m` | exit 0; **134 tests/subtests passed across 8 packages**, including real Blender tests |
| `go test ./cmd/blendctl ./internal/project -count=1` | 61 passed, including exact CLI source pointer regression |
| `go test -race ./internal/project ./cmd/blendctl ./internal/runner -run 'TestOriginal\|TestImport\|TestReject\|TestCLISignals\|TestRun\|TestDoctorTerminates\|TestSubprocessLogs\|TestPreviewRejects\|TestApplyTimeout\|TestConcurrentShot\|TestValidateJSONOriginal' -count=1` | exit 0; **40 tests/subtests passed across 3 packages** |
| `git diff --check` | exit 0 |
| `VERSION=0.1.0-rc2-dev bash scripts/release-linux-amd64.sh` | exit 0; static Linux amd64 executable, 4,120,760 bytes; version and doctor succeed |
| `dist/linux-amd64/blendctl preview testdata/compile/title_render_tiny.json --frame 1 --out <temporary parent>` | exit 0 with only PATH and MOBLEND_BLENDER environment; decoded PNG, editable scene, complete manifest |

No Blender tests were skipped. Release build/preview used the existing Go/shell script, with no Node/npm command. The ignored dist artifact is a local development build, not a tagged/published release. Its checksum is `a380540747fb41f5c40c640efae795891f1094725b62448001437f4fd705a79e`; the VCS stamp predates the commits in this session.

Official pinned API docs read before animation changes: Action, ActionChannelbag (slot_handle), bpy_struct (path_from_id), all under https://docs.blender.org/api/5.3/ (HTTP 200). Curve selection uses the animated datablock's action slot and exact RNA data path.

## Finding-to-regression map

| Finding | Fix and evidence |
|---|---|
| duplicate-import-alias-overwrite | Reject duplicates during module load with original source/pointer; `TestRejectDuplicateImportAliases`. |
| entry-camera-track-unresolvable | Entry local map includes unqualified camera ID; flat/module compile equivalence in `TestModularCameraAnimationMatchesFlat`; real camera interior sampling in mixed animation test. |
| later-track-overwrites-interpolation | Apply mode only to current path within animated owner's slot. `TestMixedAnimationInterpolationAndOrder` reopens real saved scenes in forward/reversed track order, samples step location, linear scale, eased rotation and camera motion. |
| material-track-interpolation-never-applied | Use material node-tree action and socket.path_from_id; same real probe verifies roughness step and color linear at endpoints and asymmetric interior frames. |
| dispatch-unbounded-input-read | Read bounded regular root bytes once, pass same snapshot through dispatch/resolution. Modules also reject special files. `TestPublicInputByteLimit` and `TestRejectSpecialFileInputs`; nonblocking Unix open and post-open fstat. |
| doctor-cancel-pipe-hang | Doctor and Run share isolated environment, process group, cancellation, bounded pipe wait and log sink. `TestDoctorTerminatesProbeDescendants` covers cancel/timeout; CLI regression cancels while --version runs. Errors propagate before execution. |
| import-permutation-test-noop | Real two-import cross-module graph reverses imports and manifest module order, checks equal script/data over 50 loads; original single-import test renamed accurately to serialization/relocation. |
| logs-still-unbounded / unbounded-child-log-capture | Shared draining 64 KiB head/tail collector plus explicit marker/Result metadata; 4 MiB real subprocess output tested by `TestSubprocessLogsBounded`, including version parsing with noisy output. |
| mixed-animation-channels-unprobed | New reopened pinned-Blender test covers mixed interpolation and material action owner, in both track orders. |
| module-entry-fields-silently-ignored | Shared kind-dependent contract check rejects presence of forbidden module fields, even null; `TestRejectModuleEntryOnlyFields`. Inspect reuses contract validation. |
| orphan-namespace-map-order | Alias points directly to imported module ID; unused graph members cannot influence lookup. Same-namespace orphan included in 50 repeated/permuted loads in `TestImportPermutationAndOrphanIdentity`. |
| preview-empty-marked-complete | Decode PNG header and full pixel stream, require regular/nonempty output and exact declared size; same checks for frame sequence. `TestPreviewRejectsMalformedOutputs` rejects empty/directory/truncated/wrong-size outputs with failed markers and no complete manifest. |
| qualified-reference-import-bypass | Only local IDs or declared alias:local references accepted; `TestRejectUndeclaredQualifiedTarget`. |
| s04-failed-shot-missing | Real concurrent shots share one parent, including one deliberate Python failure. Test asserts unique jobs, shot-specific dependency/complete manifests and failed marker, no failed-job complete manifest. Result/markers include shot identity. |
| source-original-pointer-unasserted / source-pointer-still-flattened | SourceRef replaces assembled object/light/track prefix; entry fields map to entry source. Resolver asserts exact code/source/pointer for shifted object/track and entry scene/camera; CLI test asserts original JSON diagnostic. |
| ease-inout-probe-gap | Mixed test samples two-key AUTO_CLAMPED ease at frames 1,3,7,9, against independent 0,.15625,.84375,1 fractions (tolerance .002 radians/value). |

## Scope and remaining caveats

The prior independent Gauntlet artifacts remain historical at `.gauntlet/outbox/`; they were not rewritten to fabricate a post-fix ship verdict. This task requested implementation, commit/push and packaging advice, not another five-critic after-pass. All reported fresh-superplan tasks have mapped implementation/regression evidence above.

Per-file input snapshots are not transactional filesystem snapshots, and v1 does not sandbox a hostile filesystem or arbitrary Blender code. Linux is the tested runtime. Manual job retention remains documented. Phase 5 actual-model repair-loop/refusal evaluation, HD aesthetic approval, stable-release Blender baseline and non-Linux support remain unclaimed. Public release publication was not performed.

## Packaging recommendation researched

Prefer a native static blendctl executable in a versioned Linux amd64 tar.gz with guide, examples, notices and SHA256SUMS. Blender stays separately installed. The current Go/shell build already needs no Node. Optional later automation: native GoReleaser CLI Go builder, archives and checksums (official goreleaser.com docs), without an npm wrapper or JavaScript release script. Add macOS/Windows only after their real render/cancel tests pass; qualify a stable Blender build before a stable release.
