# v0.2.0-alpha.1 release verification

Date: 2026-11-05. Owner authorized publishing the alpha prerelease. Repository ndx-au/MoBlend is private; visibility remains unchanged.

## Source qualification before tag

- `gofmt -l .`: empty output.
- `go vet ./...`: exit 0, no issues.
- `go test ./... -count=1 -timeout 10m`: **134 tests/subtests passed across 8 packages**, with real Blender integrations.
- `bash -n scripts/release-linux-amd64.sh scripts/test-release-linux-amd64.sh`: exit 0.
- `bash scripts/test-release-linux-amd64.sh`: PASS, real builds/extraction, identical repeated archives, checksums and tag/dirty/version/output rejection boundaries.
- `git diff --check`: clean.
- Actual Blender version output: 5.3.0 Beta, build `7e7c173d45c8`. Go 1.25.6; Linux amd64.
- GitHub identity ndx-agent has WRITE access; no existing releases or product tags conflict with v0.2.0-alpha.1. Main matched origin/main before this release preparation.

## Tag/build/publication gates

Pending after the clean source commit: annotate v0.2.0-alpha.1; build using VERSION=0.2.0-alpha.1 and --release --verify-render; verify exact tag/build metadata, extracted archive, examples/projects, rendered still/sequence and reopened saved scene; push source/tag; create draft prerelease with archive/SHA256SUMS; download/check draft assets; publish without Latest; confirm remote metadata. Append exact results after completion; no pre-emptive publication/pass claim.

## Scope

Linux amd64 only; exact tested Blender beta/CPU baseline, no stable Blender claim. Model repair-loop/refusal evaluation, broader visual/HD approval and non-Linux platforms remain incomplete. Project license remains undeclared (repository stays private); existing font license/notice are included. No Node/npm dependency or new hosted service.


## Tagged package and draft verification

Annotated tag v0.2.0-alpha.1 points to source commit `1ccc0038c9e9313ec59a709672fccd3c6669adfa`. Clean `VERSION=0.2.0-alpha.1 bash scripts/release-linux-amd64.sh --release --verify-render` succeeded. Build info reports version 0.2.0-alpha.1, mode release, dirty=false, expected commit/schema/pack/runtime and render_check=passed.

Archive: 1,890,727 bytes, SHA256 `d0b74822ebd1d1a472c735f62ae2158e2781ec42f58421ab30942adf0399489a`. Binary: 4,120,760 bytes, SHA256 `23ce883307962d34929ca8ad4b59012db0271b147bb2aa269a7f8bf4f7006bd0`. Extracted examples title/pedestal/chart and modular/three-shot all-shots validation passed; real 384x160 title preview and full 32x32 two-frame sequence completed. Reopened saved scene asserted editable MoBlend text, camera lens/active ID, resolution, title endpoint and material roughness (ALPHA_REOPEN_OK). Smoke preview visually inspected: legible title, no new HD/aesthetic qualification claim.

Source/tag push succeeded. Draft GitHub prerelease has the native tar.gz and SHA256SUMS only. Both draft assets were downloaded into a fresh directory and compared byte-for-byte; extracted downloaded binary and archive checksums verify. No tag/asset replacement performed.

## CI observation follow-up before publication

Tag CI `41894852743` passed. Parallel main CI `41894853461` exposed immediate Linux child-state assertion timing (CLI signal test); retry exposed the same assumption in Doctor's probe test. Publication stayed draft while a test-only main follow-up added a one-second bounded terminal-state wait (no descendant-liveness assertion removed). Twenty repeated runs passed 160 tests/subtests; scoped race run passed 8; full Go/Blender suite again passed 134. Rebuilt CLI with identical release flags compared byte-identical to the official tagged binary. Main follow-up CI and actual publication are still pending; source-test limitation is disclosed in draft release notes. Tag and archive remain immutable.
