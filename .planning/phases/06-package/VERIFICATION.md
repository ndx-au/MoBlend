# Phase 6 — verification evidence

**Date:** 2026-11-04  
**Status:** Linux amd64 release candidate ready; public publish not authorized

## Commands

```bash
VERSION=0.1.0-rc1 scripts/release-linux-amd64.sh
./dist/linux-amd64/blendctl doctor
./dist/linux-amd64/blendctl preview testdata/compile/title_small.json --frame 24 --out /tmp/moblend-rc-preview
go test ./... -count=1 -timeout 20m
```

## Measurements (gandalf, Linux amd64)

| Item | Value |
|---|---|
| Binary | `dist/linux-amd64/blendctl` |
| Size | 3805368 bytes (~3.6 MiB) |
| SHA-256 | `9a62f469ce3e3a24a84de5b595a564b5888692ac9200ec722ed31be6b5288a09` |
| Sample preview wall time | ~1.54 s (`title_small` 64×64 / 4 samples, frame 24) |
| Blender | 5.3.0 Beta `7e7c173d45c8` |
| Included notices | `DejaVuSans.LICENSE`, `GUIDE.md` |

## Gate

| Item | Result |
|---|---|
| Offline tooling + Blender only | PASS (no network in preview path) |
| Examples validate on final schema | PASS (`go test` / examples) |
| Docs (quick start, limitations, troubleshooting) | PASS under `docs/` |
| macOS/Windows evidence | **not claimed** |
| Public GitHub release upload | **not done** (needs owner authorization) |


## Native archive follow-up — 2026-11-05

Extended the existing script to emit versioned native tar.gz packages, SHA256SUMS, guide/examples/docs/notices and build/runtime metadata. Snapshot and clean annotated-tag modes are separate; --verify-render is explicit. No Node dependency or automatic publication added. Disposable tagged-build tests prove repeatability, extraction and version/dirty/tag/output rejection; a real extracted-binary Blender preview completed. Full Go/Blender suite: 134 tests/subtests passed; fmt/vet clean. Exact commands and caveats: [packaging verification](../../fixes/release-packaging/VERIFICATION.md). Recommended current product maturity is alpha; historic RC evidence above is preserved, not promoted to product-stability evidence.


## Alpha.1 source qualification — 2026-11-05

Owner authorized the alpha.1 GitHub prerelease in the existing private repository. Fresh full Go/Blender suite passed 134 tests/subtests; fmt/vet/Bash syntax and executed packaging acceptance passed. Pending immutable-tag package and draft-download checks are tracked in [release verification](../../releases/0.2.0-alpha.1/VERIFICATION.md). Earlier statements about publication authorization describe earlier sessions, not a current blocker.


## Alpha.1 published — 2026-11-05

[v0.2.0-alpha.1](https://github.com/ndx-au/MoBlend/releases/tag/v0.2.0-alpha.1) published as a GitHub prerelease with native Linux amd64 tar.gz and SHA256SUMS. Repository stays private; alpha is not Latest. Tagged build was extracted, rendered (still and full tiny sequence), reopened, and its draft downloads checked byte-for-byte. Tag CI and a disclosed test-only main follow-up CI passed. No tag/asset overwrite or unsupported stable/platform/model claims. Full evidence and CI timing-test chronology: [release verification](../../releases/0.2.0-alpha.1/VERIFICATION.md).
