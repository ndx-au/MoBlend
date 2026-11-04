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
