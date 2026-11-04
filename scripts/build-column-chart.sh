#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PACK="${ROOT}/internal/resources/pack/base.blend"
BLENDER="${MOBLEND_BLENDER:-${BLENDER:-blender}}"
"$BLENDER" \
  --background \
  --factory-startup \
  --disable-autoexec \
  --python-exit-code 1 \
  --python "${ROOT}/assets/source/build_column_chart.py" \
  -- \
  --pack "$PACK" \
  --out "$PACK"
sha256sum "$PACK" | tee "${PACK}.sha256"
