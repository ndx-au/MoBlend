#!/usr/bin/env bash
# Build a native, self-contained Linux amd64 package. No Node or upload tools.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

usage() {
  cat <<'HELP'
Usage: [VERSION=0.2.0-alpha.1] scripts/release-linux-amd64.sh [options]
  --release        Require VERSION, clean source and annotated vVERSION at HEAD
  --snapshot       Build a labelled development snapshot (default)
  --out DIR        New package destination; must not already exist
  --verify-render  Run doctor and bundled tiny preview before packaging
  --help           Show this help

Outputs: unpacked blendctl_VERSION_linux_amd64/, matching .tar.gz,
         SHA256SUMS (archive and binary). Nothing is tagged, pushed or uploaded.
HELP
}
fail() { printf 'release: %s\n' "$*" >&2; exit 1; }
MODE=snapshot
DEST=""
VERIFY_RENDER=0
while (($#)); do
  case "$1" in
    --release) MODE=release; shift ;;
    --snapshot) MODE=snapshot; shift ;;
    --verify-render) VERIFY_RENDER=1; shift ;;
    --out) (($# >= 2)) || fail '--out requires a directory'; DEST="$2"; [[ -n "$DEST" ]] || fail 'empty --out'; shift 2 ;;
    --help|-h) usage; exit 0 ;;
    *) fail "unknown argument: $1" ;;
  esac
done
[[ "$(uname -s)" == Linux && "$(uname -m)" == x86_64 ]] || fail 'build/run this script on Linux amd64'
for tool in go git tar gzip sha256sum find; do command -v "$tool" >/dev/null || fail "missing prerequisite: $tool"; done
[[ "$(tar --version)" == *'GNU tar'* ]] || fail 'GNU tar is required'

COMMIT="$(git rev-parse HEAD)"
SHORT_COMMIT="${COMMIT:0:12}"
DIRTY=false
[[ -z "$(git status --porcelain --untracked-files=all)" ]] || DIRTY=true
BASE_VERSION="${VERSION:-dev}"
if [[ "$BASE_VERSION" != dev ]]; then
  # SemVer syntax, including no leading zeroes in numeric prerelease identifiers.
  version_re='^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$'
  [[ "$BASE_VERSION" =~ $version_re ]] || fail 'VERSION must be SemVer without the v prefix'
  core="${BASE_VERSION%%+*}"
  if [[ "$core" == *-* ]]; then
    IFS=. read -r -a identifiers <<< "${core#*-}"
    for identifier in "${identifiers[@]}"; do
      [[ ! "$identifier" =~ ^0[0-9]+$ ]] || fail 'numeric prerelease identifiers cannot have leading zeroes'
    done
  fi
fi

if [[ "$MODE" == release ]]; then
  [[ "$BASE_VERSION" != dev ]] || fail '--release requires VERSION'
  [[ "$DIRTY" == false ]] || fail '--release requires a clean checkout (including untracked files)'
  TAG="v${BASE_VERSION}"
  [[ "$(git cat-file -t "refs/tags/$TAG" 2>/dev/null || true)" == tag ]] || fail "missing annotated tag $TAG"
  [[ "$(git rev-parse "refs/tags/${TAG}^{commit}")" == "$COMMIT" ]] || fail "$TAG does not point to HEAD"
  BUILD_VERSION="$BASE_VERSION"
else
  [[ "$BASE_VERSION" != dev ]] || BASE_VERSION=0.0.0-dev
  if [[ "$BASE_VERSION" == *+* ]]; then separator=.; else separator=+; fi
  BUILD_VERSION="${BASE_VERSION}${separator}git.${SHORT_COMMIT}"
  [[ "$DIRTY" == false ]] || BUILD_VERSION+='.dirty'
fi
EPOCH="${SOURCE_DATE_EPOCH:-$(git show -s --format=%ct "$COMMIT")}"
[[ "$EPOCH" =~ ^[0-9]+$ ]] || fail 'SOURCE_DATE_EPOCH must be nonnegative integer seconds'
STEM="blendctl_${BUILD_VERSION}_linux_amd64"
DEST="${DEST:-${ROOT}/dist/releases/${STEM}}"
[[ ! -e "$DEST" && ! -L "$DEST" ]] || fail "destination already exists: $DEST"
PARENT="$(dirname "$DEST")"
mkdir -p "$PARENT"
PARENT="$(cd "$PARENT" && pwd)"
DEST="${PARENT}/$(basename "$DEST")"
STAGE="$(mktemp -d "${PARENT}/.moblend-package.XXXXXX")"
trap 'rm -rf -- "$STAGE"' EXIT
PACKAGE="${STAGE}/${STEM}"
mkdir -p "$PACKAGE/docs" "$PACKAGE/testdata/compile"
if [[ "$MODE" == release ]]; then
  # Build tagged artifacts from an immutable Git snapshot, even if an editor
  # changes the original checkout after the cleanliness check.
  mkdir "$STAGE/source"
  git archive "$COMMIT" | tar -xf - -C "$STAGE/source"
  cd "$STAGE/source"
fi
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS= GOWORK=off GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local \
  go build -trimpath -buildvcs=false -ldflags "-s -w -X main.version=${BUILD_VERSION}" \
  -o "$PACKAGE/blendctl" ./cmd/blendctl
cp agent/GUIDE.md "$PACKAGE/GUIDE.md"
cp -R examples "$PACKAGE/examples"
cp testdata/compile/title_small.json "$PACKAGE/testdata/compile/title_small.json"
cp testdata/compile/title_render_tiny.json "$PACKAGE/testdata/compile/title_render_tiny.json"
cp docs/QUICKSTART.md docs/LIMITATIONS.md docs/TROUBLESHOOTING.md docs/RELEASING.md "$PACKAGE/docs/"
cp internal/resources/fonts/DejaVuSans.LICENSE internal/resources/fonts/DejaVuSans.NOTICE "$PACKAGE/"
PROJECT_LICENSE=not_declared
shopt -s nullglob
for notice in LICENSE* COPYING* NOTICE*; do
  [[ -f "$notice" ]] || continue
  cp "$notice" "$PACKAGE/"
  [[ "$notice" != LICENSE* && "$notice" != COPYING* ]] || PROJECT_LICENSE="$notice"
done
shopt -u nullglob
cat > "$PACKAGE/README.txt" <<'README'
MoBlend — Linux amd64

Install: copy blendctl onto your PATH, or use ./blendctl from this directory.
Prerequisite: Blender 5.3 family, separately installed; the currently tested
build is 5.3.0 Beta 7e7c173d45c8. Other Blender builds are not established by
this package. Set MOBLEND_BLENDER to its executable if it is not on PATH.
No Go, Node/npm, separate Python, model account or network is needed to run.

  ./blendctl version
  ./blendctl doctor
  ./blendctl validate examples/title.json
  ./blendctl preview testdata/compile/title_render_tiny.json --frame 1 --out /tmp/moblend-preview

Inspect the returned job paths; only a complete manifest establishes success.
GUIDE.md and docs/ explain declarations, limits, troubleshooting and releases.
Schemas, catalogue, fonts and first-party Blender assets are embedded.
DejaVuSans.LICENSE and DejaVuSans.NOTICE cover the embedded third-party font.
BUILD-INFO.txt identifies this package and its verification scope.
README
VERSION_LINE="$("$PACKAGE/blendctl" version)"
printf '%s\n' "$VERSION_LINE"
[[ "$VERSION_LINE" =~ schema=([^[:space:]]+) ]] || fail 'binary did not report schema version'
SCHEMA_VERSION="${BASH_REMATCH[1]}"
[[ "$VERSION_LINE" =~ pack=([^[:space:]]+) ]] || fail 'binary did not report asset pack version'
PACK_VERSION="${BASH_REMATCH[1]}"
"$PACKAGE/blendctl" validate "$PACKAGE/examples/title.json" --json
RENDER_CHECK=not_run
if [[ "$VERIFY_RENDER" == 1 ]]; then
  "$PACKAGE/blendctl" doctor
  mkdir "$STAGE/render-check"
  "$PACKAGE/blendctl" preview "$PACKAGE/testdata/compile/title_render_tiny.json" \
    --frame 1 --out "$STAGE/render-check" > "$STAGE/render-result.json"
  RENDER_CHECK=passed
fi
cat > "$PACKAGE/BUILD-INFO.txt" <<INFO
version=${BUILD_VERSION}
mode=${MODE}
commit=${COMMIT}
dirty=${DIRTY}
source_date_epoch=${EPOCH}
go_toolchain=$(go version)
platform=linux/amd64
schema_version=${SCHEMA_VERSION}
asset_pack_version=${PACK_VERSION}
blender_family=5.3
blender_tested_build=5.3.0 Beta 7e7c173d45c8
render_check=${RENDER_CHECK}
project_license=${PROJECT_LICENSE}
INFO
# Normalize only files owned by this temporary package, preserving source files.
find "$PACKAGE" -type d -exec chmod 755 {} +
find "$PACKAGE" -type f -exec chmod 644 {} +
chmod 755 "$PACKAGE/blendctl"
tar --sort=name --format=ustar --mtime="@${EPOCH}" --owner=0 --group=0 --numeric-owner \
  -C "$STAGE" -cf - "$STEM" | gzip -n > "$STAGE/${STEM}.tar.gz"
(
  cd "$STAGE"
  sha256sum "${STEM}.tar.gz" "${STEM}/blendctl" > SHA256SUMS
  sha256sum -c SHA256SUMS
)
# Only the archive, checksums and unpacked package are delivered, not render jobs.
cd "$ROOT"
rm -rf -- "$STAGE/render-check" "$STAGE/render-result.json" "$STAGE/source"
mv -T -n -- "$STAGE" "$DEST"
[[ ! -d "$STAGE" ]] || fail 'destination appeared during packaging; nothing replaced'
printf '\nPackage: %s\nArchive: %s/%s.tar.gz\nChecksums: %s/SHA256SUMS\n' "$DEST" "$DEST" "$STEM" "$DEST"
