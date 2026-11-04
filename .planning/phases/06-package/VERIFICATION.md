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
