# Phase 3 — verification evidence

**Date:** 2026-11-04  
**Status:** gate met (Linux amd64 / Blender 5.3.0 Beta)

## Commands

```bash
scripts/build-column-chart.sh   # bumps pack to 0.1.0 when needed
go test -a ./... -count=1 -timeout 15m
go run ./cmd/blendctl preview testdata/compile/title_small.json --frame 24 --out /tmp/p3-title
go run ./cmd/blendctl preview testdata/compile/pedestal_small.json --frame 24 --out /tmp/p3-ped
go run ./cmd/blendctl preview testdata/compile/chart_small.json --out /tmp/p3-chart
```

## Gate checks

| Gate | Result | Evidence |
|---|---|---|
| Frame probes / motion | PASS | `TestApplyPedestalMotionProbe` (~90° at frame 24) |
| Two chart instances independent | PASS | `TestApplyChartIsolation` distinct values + node group copies |
| Three compositions render | PASS | title / pedestal / chart small fixtures via `RunApply` |
| Pack/catalogue sockets | PASS | `MoBlend_ColumnChart`, inputs `Socket_1` Geometry, `Socket_2` Bar Width |
| Pack version | `0.1.0` | SHA-256 `cc4290cb23bc07d0981b793bffe342a964e5e23789db25d51dda1550f5792088` |

## Notes

- Blender 5.3 GN modifier inputs: `mod.properties.inputs.Socket_2.value` (not IDProperty map)
- Bar heights via mesh point attribute `bar_value`; apply script builds points per declaration
- Schema/examples require `asset_pack_version: "0.1.0"`
- Compositor look deferred
