#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
if [[ -z "${BLENDCTL:-}" ]]; then
  go build -o /tmp/blendctl "$ROOT/cmd/blendctl"
  BLENDCTL=/tmp/blendctl
fi
OUT="$ROOT/testdata/visual"
mkdir -p "$OUT/frames" "$OUT/contact"
preview() {
  local name="$1" scene="$2" frame="$3"
  local job="/tmp/moblend-visual-$name"
  rm -rf "$job"
  "$BLENDCTL" preview "$scene" --frame "$frame" --out "$job" >"/tmp/preview-$name.json"
  python3 - "$name" <<'PY'
import json, shutil, sys
name=sys.argv[1]
doc=json.load(open(f"/tmp/preview-{name}.json"))
shutil.copy(doc["png"], f"testdata/visual/frames/{name}.png")
print(name, doc["png"])
PY
}
cd "$ROOT"
preview title testdata/compile/title_small.json 24
preview pedestal testdata/compile/pedestal_small.json 24
preview chart testdata/compile/chart_small.json 1
montage -label '%t' -geometry 256x256+8+8 -tile 3x1 -background '#111111' -fill white \
  testdata/visual/frames/title.png \
  testdata/visual/frames/pedestal.png \
  testdata/visual/frames/chart.png \
  testdata/visual/contact/references_sheet.png
ls -la testdata/visual/frames testdata/visual/contact
