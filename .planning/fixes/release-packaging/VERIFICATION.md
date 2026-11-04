# Native release packaging verification

Date: 2026-11-05. Go 1.25.6, Linux amd64; Blender 5.3.0 Beta 7e7c173d45c8 for real render checks. Script uses Bash/Go/Git/GNU tar/gzip/coreutils; no Node/npm or publishing tool.

## Implemented

- Versioned unpacked package + tar.gz + SHA256SUMS, with binary, guide, examples, smoke inputs, user docs, provenance/runtime info and font license/notice.
- Default snapshot identity includes source revision and dirty state. Explicit --release requires SemVer VERSION, clean checkout and an annotated vVERSION tag at the selected commit; build inputs come from that immutable Git archive.
- Private staging, new-destination-only output, deterministic archive metadata, optional real render check. No automatic tagging, pushes, uploads or publication.
- Recommended maturity and version/promotion strategy documented in docs/RELEASING.md. Alpha is a recommendation based on current evidence; existing historic labels/contracts were not renamed.

## Exact commands and results

| Command | Result |
|---|---|
| `bash -n scripts/release-linux-amd64.sh scripts/test-release-linux-amd64.sh` | exit 0 |
| `bash scripts/test-release-linux-amd64.sh` | PASS: real release/snapshot builds in disposable Git checkout, extract/execute/validate, byte-identical repeated archives, checksums; rejects invalid versions, numeric leading zeroes, missing/lightweight/wrong-commit tags, tracked/untracked dirt and existing destinations |
| `gofmt -l .` | empty output |
| `go vet ./...` | exit 0, no issues |
| `go test ./... -count=1 -timeout 10m` | exit 0, **134 tests/subtests passed across 8 packages**, real Blender checks included |
| `VERSION=0.2.0-alpha.1 bash scripts/release-linux-amd64.sh --verify-render --out /tmp/moblend-package-smoke-lmu923np/bundle` | exit 0; snapshot version `0.2.0-alpha.1+git.f4270bd48b13.dirty`; render_check=passed; archive and executable checksums verify |
| Fresh archive extraction, bundled title/pedestal/chart validation and three-shot --all-shots validation | exit 0 |
| Extracted binary preview of bundled title_render_tiny.json, with only PATH/MOBLEND_BLENDER environment | exit 0; decoded PNG and editable scene, manifest status complete |
| `git diff --check` | exit 0 |

Snapshot archive measured 1,890,788 bytes (environment/snapshot-specific, not a size promise). Extracted preview: `/tmp/moblend-package-smoke-lmu923np/moblend-job-4265975336/out/preview.png`. Build artefacts are ephemeral/local and not release publication.

ShellCheck was not installed; no ShellCheck pass is claimed. Bash syntax and executed acceptance/failure paths are the shell evidence. No new aesthetic review or independent model-authoring evaluation was performed. Other platforms/stable Blender qualification remain open. No project software license was invented; BUILD-INFO records its absence until the owner chooses terms.

The real render smoke preceded a final metadata-only change to derive the short revision and default timestamp from the same captured commit; the complete packaging acceptance suite was rerun after that change. Application Go source did not change in this session.

## Repository boundaries

Preserved unrelated README edits and G000008 present at session start. No tag or GitHub release was created in the product checkout; annotated tags used by acceptance tests existed only in the disposable fixture. Publication remains an explicit maintainer action.
