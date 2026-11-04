#!/usr/bin/env bash
# Build the embedded Phase 0 base.blend using the pinned Blender binary.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="${ROOT}/internal/resources/pack/base.blend"
BLENDER="${MOBLEND_BLENDER:-${BLENDER:-blender}}"

mkdir -p "$(dirname "$OUT")"
"$BLENDER" \
  --background \
  --factory-startup \
  --disable-autoexec \
  --python-exit-code 1 \
  --python "${ROOT}/assets/source/build_base_pack.py" \
  -- \
  --out "$OUT"

sha256sum "$OUT" | tee "${OUT}.sha256"
ls -la "$OUT"
