# MoBlend quick start (Linux amd64)

## Prerequisites

1. Blender **5.3** family on `PATH`, or set `MOBLEND_BLENDER` to the binary  
   (see `.planning/DEVELOPMENT-ENVIRONMENT.md` for the local beta pin).
2. The `blendctl` binary from a release build (`scripts/release-linux-amd64.sh`) or `go run ./cmd/blendctl`.

No separate Python, Node, model key, or network is required for ordinary use.

## First render

```bash
blendctl doctor
blendctl validate examples/title.json
blendctl preview testdata/compile/title_small.json --frame 24 --out /tmp/moblend-out
```

Open `/tmp/moblend-out/moblend-job-*/out/preview.png` and `scene.blend`.

## Agent authoring

```bash
blendctl guide
blendctl assets list
blendctl assets describe column_chart
```

Write JSON using exact asset IDs only. Repair from `validate --json` diagnostics.
