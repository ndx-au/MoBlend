# Phase 8 — independent shot targets context

Date: 2026-11-04. Status: approved implementation contract.

## Outcome

Agents build/render one shot of a multi-shot project without compiling or launching Blender for other shots. Shared modules are imported explicitly; there is no implicit global mutable scene state.

## Frozen decisions

- Project manifest may include `shots: [{ "id", "entry" }]` where `entry` is a module id with `kind: entry`.
- Top-level `entry` remains the default whole-project entry when `shots` is absent (Phase 7 behavior).
- **Selection:** multi-shot projects require `--shot <id>` for validate/compile/preview/render. A project with exactly one shot may default to that shot. `--all-shots` validates every shot closure and aggregates diagnostics.
- **Closure:** resolving a shot sets the effective entry to that shot's entry module and loads only its import closure.
- Unselected broken modules/files are not loaded for a selected shot.
- Outputs: with `--out DIR --shot ID`, nest under `DIR/<shot>/`. Always use distinct job directories. Write `dependencies.json` with shot id, relative input paths, content SHA-256, pack SHA, resolved declaration hash, and effective settings.
- Concurrent jobs must not share writable job trees.
- Deferred: incremental cache, automatic all-shot render, FFmpeg assembly, parameterized instances.

## Definition of done

SHOT-01–03 pass with unit/CLI and real Blender evidence recorded in VERIFICATION.md.
