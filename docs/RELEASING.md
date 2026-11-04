# Versions, packages and releases

The recommended current maturity is **alpha**. The scene/compiler/render pipeline,
modules and shots work on the tested Linux host and have substantial automated
coverage. Independent authoring/onboarding evidence, a stable Blender baseline
and stronger visual qualification remain incomplete. This is an assessment of
the evidence, not a claim that every supported path is defective.

The earlier `0.1.0-rc1` label records a local packaging milestone. It was not a
published release or proof that product qualification was finished. Preserve
that history; recommend **`v0.2.0-alpha.1`** for the first public alpha rather than
renaming artifacts or moving an existing tag. The owner has authorized that alpha
prerelease for the existing private GitHub repository. Publication does not change
repository visibility. Tags are created explicitly during the release workflow.

## Version policy

Use [Semantic Versioning](https://semver.org/), with a `v` prefix only on Git tags:

| Item | Example | Meaning |
|---|---|---|
| Snapshot | `0.0.0-dev+git.abcdef123456.dirty` | Local build; never a stable release |
| Alpha | `0.2.0-alpha.1`, `.2`, ... | Supported features under active validation; contracts may change |
| Beta | `0.2.0-beta.1` | Intended scope frozen, supported workflow qualified, external testing underway |
| Candidate | `0.2.0-rc.1` | Release criteria met; only release-blocking fixes expected |
| Stable | `0.2.0` | Qualified release of that scope; major zero still signals an evolving public contract |
| Public API commitment | `1.0.0` | Deliberate long-term compatibility commitment, not a phase number |

Within prerelease qualification, increment the numeric suffix. Before 1.0,
use a new minor version for incompatible contract changes or substantial scope,
and patch versions for compatible fixes after a numbered stable release. Document
migrations explicitly. After 1.0 use the standard major/minor/patch compatibility
rules. Schema, project/module and asset-pack versions are separate contracts;
do not bump them merely because the CLI version changed.

Treat tags and release assets as immutable. A correction gets a new version.
A tag names source; a GitHub Release adds notes and downloadable artifacts to
that tag. GitHub's automatic source tarball is not the compiled binary package.
See [GitHub's release model](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases).

## Native packaging

Required build tools: Linux amd64, installed Go meeting go.mod, Bash, Git, GNU tar,
gzip, find and sha256sum/coreutils. No Node/npm, GoReleaser or publishing CLI is
required. The script disables module downloads; install the Go toolchain first.
Blender is required only for `--verify-render` or later runtime use.

Run from the repository root (or invoke the script by its absolute path):

```bash
# Default: snapshot, explicitly marked with source revision and dirty state.
scripts/release-linux-amd64.sh

# Snapshot based on a proposed release version; still labelled +git.REVISION.
VERSION=0.2.0-alpha.1 scripts/release-linux-amd64.sh --verify-render

# Choose a NEW output directory (existing destinations are never replaced).
VERSION=0.2.0-alpha.1 scripts/release-linux-amd64.sh --out /tmp/moblend-package
```

The default output is `dist/releases/blendctl_VERSION_linux_amd64/`, containing:

```text
blendctl_VERSION_linux_amd64.tar.gz
SHA256SUMS
blendctl_VERSION_linux_amd64/
  blendctl
  README.txt
  GUIDE.md
  BUILD-INFO.txt
  DejaVuSans.LICENSE
  DejaVuSans.NOTICE
  docs/{QUICKSTART,LIMITATIONS,TROUBLESHOOTING,RELEASING}.md
  examples/...
  testdata/compile/{title_small,title_render_tiny}.json
```

Any root LICENSE/COPYING/NOTICE files are also copied when present. No project
license is invented when none is declared; BUILD-INFO states what was found.
Choose a project license/terms before public distribution. Font notices already
travel with the package. Blender is installed separately, not bundled.

`SHA256SUMS` covers the archive and unpacked executable. `BUILD-INFO.txt` records
mode, version, commit, dirty status, Go toolchain, source timestamp, schema/pack,
Blender compatibility and whether this packaging run verified a render.
Tagged packaging does not imply every release gate was run.

File ordering, ownership, permissions, archive timestamps and gzip metadata are
normalized. SOURCE_DATE_EPOCH defaults to the source commit timestamp. Repeated
builds require the same input snapshot, Go toolchain, mode and settings for the
same archive bytes; do not infer cross-toolchain or rendered-pixel identity.
Packaging stages privately and refuses an existing destination. For another
build, choose a fresh destination; keep published assets unchanged.

## Qualification and publication workflow

1. Freeze the scope and update notes, compatibility and requirement evidence.
   Run formatting, vet, full Go/Blender checks and packaging acceptance:

   ```bash
   scripts/check.sh
   scripts/test-release-linux-amd64.sh
   ```

   Confirm Blender integration tests actually ran; skips are not qualification.
2. Commit the exact source/docs to release and verify the checkout is clean.
   Create an **annotated tag** such as `v0.2.0-alpha.1` on that tested commit;
   sign it if signing is configured. This is an explicit maintainer action.
3. Package that exact commit:

   ```bash
   VERSION=0.2.0-alpha.1 scripts/release-linux-amd64.sh --release --verify-render
   ```

   `--release` refuses a dirty checkout, missing/lightweight tag or a tag pointing
   at another commit. It builds from an immutable Git archive of that commit, so
   later checkout edits cannot enter the tagged package. Version must be SemVer
   without its tag's `v` prefix.
4. Verify checksums, extract the archive into a fresh directory, validate bundled
   examples and render using the extracted binary with separately installed
   Blender. Review output/editability; retain exact build/runtime evidence.
5. Push the approved tag and prepare a **draft GitHub Release** for it. Attach
   the tar.gz and SHA256SUMS produced by that build. Include tested platforms,
   exact Blender build, changes, migrations, limitations and verification.
6. Review and publish deliberately. Alpha/beta/RC are GitHub **prereleases**;
   do not mark them as the stable/latest release. Do not rebuild or replace
   assets after publication; ship a new version instead.

The packaging script never creates/pushes tags, uploads assets or publishes.
Its `--release` option is a local provenance gate, not publication authorization.
The current work does not add a CI publishing workflow.

## Promotion criteria

- **Alpha now:** useful for controlled evaluation on Linux amd64, with explicit
  limitations and evolving interfaces.
- **Beta:** intended scope frozen; reproducible qualification on an exact stable
  Blender build; fresh-machine onboarding and external agent exercises; reviewed
  reference visuals; no known high-severity correctness or process/output bugs.
- **RC:** beta feedback addressed and every advertised release gate met on each
  advertised platform. A Linux-only release need not wait for other platforms,
  but it must not advertise them.
- **Stable:** the exact candidate archive passes final install/render/editability
  checks and has approved notes/compatibility. Actual-model claims need their
  own measured evidence; they do not follow from rendering prepared fixtures.
