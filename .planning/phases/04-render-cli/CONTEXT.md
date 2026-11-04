# Phase 4 — context

## Supported Blender gate

- Minimum accepted: major.minor matching pinned baseline family `5.3` (beta builds allowed for now).
- `doctor` reports binary path, version string, build hash when available, and pack/schema versions from the tool.

## Manifest

Written only on success as `manifest.json` in the job `out/` directory:

```json
{
  "status": "complete",
  "schema_version": "1.0.0",
  "asset_pack_version": "0.1.0",
  "pack_sha256": "...",
  "blender": {"path": "...", "version": "..."},
  "settings": {"frame_start":1,"frame_end":24,"fps":24,"resolution":[64,64],"samples":4,"engine":"CYCLES","device":"CPU"},
  "outputs": {"blend":"scene.blend","frames":["preview.png"], "mode":"preview|render"}
}
```

Failed jobs must not write `status: complete`. Prefer `manifest.failed.json` or omit success manifest.

## Exit categories

| Code | Meaning |
|---|---|
| 0 | success |
| 1 | validation / execution failure |
| 2 | usage |
| 3 | environment (missing/unsupported Blender) — doctor & preflight |
