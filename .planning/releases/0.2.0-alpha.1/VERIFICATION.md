# v0.2.0-alpha.1 release verification

**Current status: published GitHub prerelease; not Latest; repository remains private.** Earlier pending statements below describe preparation checkpoints.

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


## Publication completed

- Release: https://github.com/ndx-au/MoBlend/releases/tag/v0.2.0-alpha.1
- GitHub published timestamp: 2026-11-04T19:22:36Z (2026-11-05 Australia/Sydney).
- Remote metadata confirmed: isDraft=false, isPrerelease=true; repository isPrivate=true; no latest stable release endpoint points to this alpha.
- Exactly two uploaded assets: blendctl_0.2.0-alpha.1_linux_amd64.tar.gz and SHA256SUMS. Server SHA-256 digests match the verified local files; downloaded draft bytes and extracted binary already verified before publishing.
- Annotated remote tag object 48fa76984f9b5efc6cbb5c426f140fa838dd785a peels to immutable source 1ccc0038c9e9313ec59a709672fccd3c6669adfa. No tag or release asset was overwritten.
- Tag CI PASS: https://github.com/ndx-au/MoBlend/actions/runs/41894852743
- Test-only follow-up commit 0e8d42ac1b6052e5269466691049ac4226df8801 pushed on main; CI PASS: https://github.com/ndx-au/MoBlend/actions/runs/41896087094. Application executable remained byte-identical to the tagged binary. The older parallel branch timing failures remain visible in history and are disclosed in release notes.
- Publication command: `gh release edit v0.2.0-alpha.1 --draft=false --prerelease --latest=false --verify-tag` (via rtk), exit 0. Post-publication GitHub metadata/assets/digests independently queried and asserted.

The first alpha release is complete within its Linux/beta-runtime scope. Stable Blender, external model/onboarding qualification, broader visual approval and additional platforms remain follow-ups, not hidden completed gates. No license, repository visibility or supported-platform expansion was made.
