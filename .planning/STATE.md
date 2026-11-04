# Project state

Updated: 2026-11-04
Status: Roadmap Phases 0–6 complete for Linux amd64 RC `0.1.0-rc1` with documented blocks (LLM eval, visual contact sheet, non-Linux platforms).
Active phase: none (milestone complete pending owner publish authorization / blocked reviews).
Next action: owner visual review + authorize release publication if desired; optional LLM eval harness.

## Completed phases

| Phase | VERIFICATION | Progress |
|---|---|---|
| 0 | `phases/00-blender-boundary/VERIFICATION.md` | `M000.*` |
| 1 | `phases/01-declaration/VERIFICATION.md` | `M001.*` |
| 2 | `phases/02-title-composition/VERIFICATION.md` | `M002.*` |
| 3 | `phases/03-motion/VERIFICATION.md` | `M003.*` |
| 4 | `phases/04-render-cli/VERIFICATION.md` | `M004.*` |
| 5 | `phases/05-agent-usability/VERIFICATION.md` | `M005.*` |
| 6 | `phases/06-package/VERIFICATION.md` | `M006.*` |

## RC artifact

- `dist/linux-amd64/blendctl` (gitignored build output) sha256 `9a62f469ce3e3a24a84de5b595a564b5888692ac9200ec722ed31be6b5288a09`
- Build: `VERSION=0.1.0-rc1 scripts/release-linux-amd64.sh`

## Blocked / partial

- LLM authoring evaluation (Phase 5)
- Human contact-sheet visual review (QA-01)
- RUN-01 concurrent stress automation
- macOS/Windows support claims
