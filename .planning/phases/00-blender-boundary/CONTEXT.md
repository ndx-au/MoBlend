# Phase 0 — context and decisions

## Toolchain pins (evidence-driven)

| Tool | Pin | Evidence |
|---|---|---|
| Go | `go1.25.6` (`go` directive `1.25`) | `go version` on gandalf 2026-11-04 |
| Blender (dev baseline) | `5.3.0 Beta` build `7e7c173d45c8` | `DEVELOPMENT-ENVIRONMENT.md`; `blender --version` |
| Primary OS | Linux amd64 | First integration target |

Release Blender baseline remains open until a stable/LTS build is smoke-tested; beta is the local candidate for Phase 0 only.

## Invocation contract (from Blender help + manual)

Ordered argv pattern for trusted jobs:

1. `--background`
2. `--factory-startup`
3. `--disable-autoexec` (default; set explicitly for clarity)
4. `--python-exit-code 1`
5. `<extracted pack .blend>`
6. `--python <trusted smoke script>`
7. `--` then job-specific args (output paths, property overrides)

Spawn via `os/exec` argument vector only — never a shell string.

`--disable-autoexec` must not block the explicitly passed `--python` script (CLI scripts run regardless of autoexec).

## Asset pack v0

- Single `base.blend`: camera, area light, cube mesh, Principled material with known hex colour `#4A90D9`, Cycles CPU, 128×128 authoring default (smoke uses 64×64).
- Authored by checked-in `assets/source/build_base_pack.py`; runtime bytes live under `internal/resources/pack/` for `go:embed`.
- Pack version string: `0.0.1-phase0` until Phase 3 catalogue freeze.
- SHA-256: `882acdcadfaa1d5e180355a7fe2b1ab1d09865d95e2d9a8d531a88e97da083f5`
- Authoring note: Blender 5.3 deprecation warning on `Material.use_nodes` (removal expected in 6.0).

## Module / CLI naming

- Module path: `github.com/ndx-au/MoBlend` (matches origin).
- Provisional binary: `blendctl` (per ARCHITECTURE); Phase 0 exposes a smoke/doctor-seed path, not the full CLI surface.

## Deferred

Stable release Blender version; fonts packaging; compositor look; full `doctor` UX; JSON schema.