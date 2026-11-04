# Packaging decisions

- Recommendation: alpha today; next public release v0.2.0-alpha.1, then beta/RC only after stated evidence. Historic local 0.1.0-rc1 is not renamed.
- Default mode produces explicitly labelled snapshot versions. --release requires VERSION, a clean checkout, and an annotated vVERSION tag at HEAD. Packaging never tags/pushes/uploads.
- One new destination contains the unpacked package, versioned tar.gz and SHA256SUMS. Refuse existing destinations; stage privately and move on success.
- Use Go, Bash, GNU tar, gzip and sha256sum; no Node. Archive metadata uses the source commit timestamp unless SOURCE_DATE_EPOCH is supplied. Same source snapshot/toolchain/settings should reproduce the archive; no cross-toolchain/pixel identity claim.
- Optional --verify-render runs the bundled tiny scene through Blender before sealing the archive. Build info states whether rendering was checked; tagged packaging alone is not release qualification.
- Include guide, examples, the tiny smoke scene, selected user docs, build/runtime compatibility info, font license/notice and any project LICENSE/COPYING/NOTICE files present. Do not invent a missing project license.
- Tagged builds read the selected commit through git archive, rather than reading mutable checkout files during compilation/copying. Short revision and source timestamp use the same captured commit.
