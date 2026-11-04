# Native release packaging extension

Scope: extend existing Linux amd64 Bash/Go script with versioned archives, checksums, bundled onboarding material and explicit snapshot/tagged-build modes. Document recommended alpha/beta/RC/version policy. Do not create tags or publish releases.

| Order | Owner/files | Tasks | Requirements/checks |
|---|---|---|---|
| 1 | primary: scripts/release-linux-amd64.sh | Isolated staging, validated versions, clean annotated-tag release gate, archive + checksums + build info + docs/examples/notices | DIST-01/02, no Node, no overwrite |
| 2 | primary: scripts/test-release-linux-amd64.sh | Real local fixture builds, reproducibility, extracted package checks, tag/dirty/version/output rejection | Bash syntax + packaging acceptance |
| 3 | primary: docs/RELEASING.md, quick start | Snapshot/tag/version/release workflow and evidence-driven maturity recommendation | Authoritative SemVer/GitHub docs; no release publication |
| 4 | primary: planning/progress | Exact verification commands + real extracted-binary Blender preview | fmt/vet/tests, phase 6 addendum, append Dot Progress |

Preserve unrelated README and G000008 edits already present at session start. Avoid changing scene/schema/asset pack versions. No new package manager or publishing dependency.
