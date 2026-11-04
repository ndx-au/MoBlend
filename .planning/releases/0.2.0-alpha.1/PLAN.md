# Publish v0.2.0-alpha.1

Owner explicitly authorized the alpha.1 release on 2026-11-05. Preserve the private repository's visibility. No software-license change, no unsupported platform/stable-Blender/model claims.

| Order | Owner/files | Work | Gate |
|---|---|---|---|
| 1 | primary: release docs/README/planning | Finalize alpha notes and commit pending packaging/docs | No source scope expansion; clean checkout |
| 2 | primary: Go + packaging checks | Fresh fmt/vet/full Go+Blender tests, packaging acceptance | No skipped Blender gates |
| 3 | primary: git/package | Annotated v0.2.0-alpha.1 at tested commit; --release --verify-render | Exact tag/commit/version, dirty=false |
| 4 | primary: extracted archive | Checksums, example/project validation, still + sequence, reopen saved scene | Correct outputs/editability and runtime pin |
| 5 | primary: GitHub | Push source/tag; create prerelease draft, attach tar.gz/SHA256SUMS; download/verify | Correct asset bytes, prerelease=true, latest=false |
| 6 | primary: publication/audit | Publish draft; confirm remote metadata; append immutable progress and final verification | URL, source tag, assets and checksum match; repo remains private |

Do not overwrite an existing published tag or asset. Build/render artifacts stay ignored; commit source/docs/audit. No automatic release service or Node dependency.
