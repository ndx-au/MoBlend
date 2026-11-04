# Release decisions

- Version: CLI 0.2.0-alpha.1, annotated tag v0.2.0-alpha.1. Scene/project/module schema 1.0.0 and asset pack 0.1.0 unchanged.
- Linux amd64 only. Supported tested runtime: Blender 5.3.0 Beta 7e7c173d45c8, CPU render baseline. Blender separately installed; no Go/Node/system Python needed by users.
- Publish a GitHub prerelease without marking Latest. Keep repository private; release authorization does not authorize a visibility change.
- Include pending packaging work and owner-authored token-cost documentation in the release source commit. Preserve all previous committed progress entries.
- Retain currently undeclared project-license status; include existing font license/notice. No new OSS license or legal grant is selected by the agent.
- Exact tagged package will be built after the clean commit, tested from extraction and rechecked after draft upload. Publication audit may follow the immutable tag as a separate main-branch documentation commit.
