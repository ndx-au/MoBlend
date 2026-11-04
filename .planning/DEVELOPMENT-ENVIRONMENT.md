# Local development environment

Verified on 2026-11-04, Ubuntu 26.04.1 LTS, Linux x86_64.

- Official Blender 5.3.0 Beta, build `7e7c173d45c8`, built 2026-11-04.
- Installed under `/home/terence/.local/opt/blender-5.3.0-beta+v53.7e7c173d45c8-linux.x86_64-release/`.
- `/home/terence/.local/bin/blender` and `blender-beta` resolve through `/home/terence/.local/opt/blender-beta/blender`.
- Existing apt package Blender 5.0.1 remains at `/usr/bin/blender`.
- Applications menu entry: Blender Beta.
- Download: https://cdn.builder.blender.org/download/daily/blender-5.3.0-beta+v53.7e7c173d45c8-linux.x86_64-release.tar.xz
- Published SHA-256 verified: `c1f78a674c59f1c618679858b71ae8b09f857ca2e49dbd01bb5d68baef184b61`.

Verification: version output matched the beta build. Background execution with factory startup, auto-execution disabled, and Python failure exit enabled saved a scene and rendered a 64×64 PNG with CPU Cycles. Artifacts and script are under `/home/terence/.cache/blender-beta-smoke/`.

Use `blender-beta` for development or configure its absolute executable path. This is a pinned snapshot of the latest official beta available at installation time; updates are manual so a new daily build does not silently change render behavior. It establishes a local candidate development baseline, not release compatibility. Phase 0 still needs the project's embedded-pack/compiler smoke test and its own evidence.
