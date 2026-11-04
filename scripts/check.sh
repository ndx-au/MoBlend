#!/usr/bin/env bash
# Local CI-equivalent checks for Phase 0+.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "== gofmt =="
test -z "$(gofmt -l .)"

echo "== go vet =="
go vet ./...

echo "== go test (unit) =="
go test ./...

if [[ "${MOBLEND_INTEGRATION:-0}" == "1" ]]; then
  echo "== go test (integration) =="
  go test ./internal/runner -count=1 -timeout 10m
fi

echo "OK"
