# MoBlend quick start (Linux amd64)

## Prerequisites

1. Blender **5.3** family on `PATH`, or set `MOBLEND_BLENDER` to the binary  
   (`BUILD-INFO.txt` in the package records the tested beta build; source
   checkouts also have `.planning/DEVELOPMENT-ENVIRONMENT.md`).
2. Extract the native release archive and put `blendctl` on PATH, or use
   `./blendctl` from its package directory. Go is needed only for source builds.

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

## Build a package from source

```bash
scripts/release-linux-amd64.sh --verify-render
```

This creates a labelled snapshot under `dist/releases/`, with a versioned tar.gz,
SHA256SUMS and an unpacked package. The script needs Go/Bash and standard Linux
archive tools; no Node/npm. It never tags or publishes. See [release strategy](RELEASING.md)
for clean annotated-tag builds, maturity gates and publication steps.
