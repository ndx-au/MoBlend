#!/usr/bin/env bash
# Build a Linux amd64 release binary with checksum.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT="${ROOT}/dist/linux-amd64"
VERSION="${VERSION:-dev}"
mkdir -p "$OUT"
export CGO_ENABLED=0
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" \
  -o "${OUT}/blendctl" "${ROOT}/cmd/blendctl"
cp "${ROOT}/internal/resources/fonts/DejaVuSans.LICENSE" "${OUT}/DejaVuSans.LICENSE"
cp "${ROOT}/agent/GUIDE.md" "${OUT}/GUIDE.md"
(
  cd "$OUT"
  sha256sum blendctl | tee blendctl.sha256
)
ls -la "${OUT}/blendctl"
"${OUT}/blendctl" version
"${OUT}/blendctl" doctor || true
