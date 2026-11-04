#!/usr/bin/env bash
# Exercise real package builds and tag/output boundaries in a disposable checkout.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORK="$(mktemp -d)"
trap 'rm -rf -- "$WORK"' EXIT
mkdir "$WORK/source"
tar --exclude=.git --exclude=dist -C "$ROOT" -cf - . | tar -xf - -C "$WORK/source"
cd "$WORK/source"
git init -q -b main
git config user.name 'MoBlend packaging test'
git config user.email 'packaging-test@example.invalid'
git config commit.gpgsign false
git config tag.gpgsign false
git config core.hooksPath /dev/null
git add -A
git commit -qm 'Packaging acceptance fixture'
SCRIPT="$WORK/source/scripts/release-linux-amd64.sh"
expect_failure() {
  if "$@" > "$WORK/expected-failure.log" 2>&1; then
    printf 'expected failure: %s\n' "$*" >&2; exit 1
  fi
}
# Invalid labels must fail before creating a package, not enter linker/shell syntax.
expect_failure env VERSION='../escape' "$SCRIPT" --out "$WORK/invalid"
expect_failure env VERSION='0.2.0-alpha.01' "$SCRIPT" --out "$WORK/leading-zero"
expect_failure env VERSION='0.2.0-alpha.1' "$SCRIPT" --release --out "$WORK/missing-tag"
git tag v0.2.0-alpha.1
expect_failure env VERSION='0.2.0-alpha.1' "$SCRIPT" --release --out "$WORK/lightweight"
git tag -d v0.2.0-alpha.1 >/dev/null
git tag -a v0.2.0-alpha.1 -m 'Packaging acceptance tag'
touch dirty-untracked
expect_failure env VERSION='0.2.0-alpha.1' "$SCRIPT" --release --out "$WORK/dirty"
rm dirty-untracked
VERSION=0.2.0-alpha.1 "$SCRIPT" --release --out "$WORK/release-a"
VERSION=0.2.0-alpha.1 "$SCRIPT" --release --out "$WORK/release-b"
STEM=blendctl_0.2.0-alpha.1_linux_amd64
cmp "$WORK/release-a/$STEM.tar.gz" "$WORK/release-b/$STEM.tar.gz"
expect_failure env VERSION='0.2.0-alpha.1' "$SCRIPT" --release --out "$WORK/release-a"
(cd "$WORK/release-a" && sha256sum -c SHA256SUMS)
mkdir "$WORK/extracted"
tar -xzf "$WORK/release-a/$STEM.tar.gz" -C "$WORK/extracted"
PACKAGE="$WORK/extracted/$STEM"
[[ -x "$PACKAGE/blendctl" && -f "$PACKAGE/GUIDE.md" && -f "$PACKAGE/DejaVuSans.NOTICE" ]]
[[ -f "$PACKAGE/docs/RELEASING.md" && -f "$PACKAGE/examples/projects/three_shots/project.json" ]]
[[ -f "$PACKAGE/testdata/compile/title_render_tiny.json" ]]
grep -qx 'mode=release' "$PACKAGE/BUILD-INFO.txt"
grep -qx 'dirty=false' "$PACKAGE/BUILD-INFO.txt"
"$PACKAGE/blendctl" version | grep -q 'blendctl 0.2.0-alpha.1 schema=1.0.0 pack=0.1.0'
"$PACKAGE/blendctl" validate "$PACKAGE/examples/title.json" --json
# Tag drift, dirty tracked source, and snapshot identity are independent boundaries.
git commit --allow-empty -qm 'New untagged commit'
expect_failure env VERSION='0.2.0-alpha.1' "$SCRIPT" --release --out "$WORK/tag-drift"
printf '\n' >> docs/LIMITATIONS.md
expect_failure env VERSION='0.2.0-alpha.1' "$SCRIPT" --release --out "$WORK/tracked-dirty"
VERSION=0.2.0-alpha.1 "$SCRIPT" --snapshot --out "$WORK/snapshot"
snapshot_info=("$WORK/snapshot"/blendctl_*/BUILD-INFO.txt)
grep -qx 'mode=snapshot' "${snapshot_info[0]}"
grep -q '^version=0.2.0-alpha.1+git\..*\.dirty$' "${snapshot_info[0]}"
printf 'Packaging acceptance: PASS (real build/extract, reproducibility, checksums, versions, clean/tag/output guards)\n'
