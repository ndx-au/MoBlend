# Limitations and compatibility

- **Platforms:** Linux amd64 is the verified integration target. macOS/Windows are experimental until smoke-tested.
- **Blender:** Family `5.3` required by `doctor`. Beta builds used in development are not a release LTS claim.
- **Determinism:** Same declaration → byte-identical compile artifacts. Pixel-identical renders across GPUs/OS/versions are **not** promised.
- **Language:** No arbitrary Python, user `.blend` import, physics, rigs, or custom node graphs.
- **Outputs:** PNG frames + editable `.blend` + manifest. No FFmpeg/video in v1.
- **Assets:** Exact-ID catalogue only (`base_cube`, `column_chart`). No RAG.
- **Fonts:** Bundled DejaVu Sans; see `DejaVuSans.LICENSE`.
