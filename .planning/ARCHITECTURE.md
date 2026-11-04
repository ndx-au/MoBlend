# Architecture and proposed contracts

## Pipeline

Brief → external coding agent → JSON → strict decoder/schema checks → semantic validation → normalized typed IR → deterministic Python emission → isolated Blender process + extracted asset pack → scene, frames, manifest.

The CLI works without an agent. Natural-language generation is outside its trusted core. Asset lookup is exact ID lookup in a shipped catalogue; there is no retrieval service.

## Suggested repository layout

```text
cmd/blendctl/             CLI entry point
internal/spec/           decoding, schema, normalization
internal/validate/       references, types, capabilities, limits
internal/compile/        typed operations and deterministic emitter
internal/runner/         extraction, subprocess, cancellation, manifests
internal/resources/      go:embed root: schema, catalogue, assets, templates
assets/source/           first-party asset generation sources and notes
examples/                three validated example declarations
agent/                   compact usage guide and example walkthrough
scripts/                 repeatable asset build and integration harness
testdata/               invalid, golden, scene and visual fixtures
.planning/               roadmap, phase plans, verification, decisions
```

Keep runtime resources below their Go package so `go:embed` can include them. Embedding support is documented by [Go's embed package](https://pkg.go.dev/embed). Track small runtime `.blend` assets in the repository with source generation scripts and checksums; avoid an LFS setup requirement for the initial pack.

## Versioning and asset contract

Use independent explicit `schema_version`, `asset_pack_version`, compiler version, and supported Blender version. Reject unsupported versions with a precise diagnostic; no silent migrations in v1. The catalogue describes each asset ID, type, input names, value types, units, defaults, bounds, animatable properties, and internal binding identifiers. Never resolve assets using fuzzy display-name matching.

An integration probe opens the pack and checks every declared object, material, node group, and exposed socket against the catalogue. Bind Geometry Nodes inputs through the API identifiers verified on the supported Blender version, not assumed positional indices. Duplicate mutable datablocks when necessary to avoid one instance changing another.

## Language semantics to freeze in Phase 1

- Root: versions, scene settings, objects, camera, lights, animation tracks, render preset.
- JSON Schema with unknown fields rejected; decoder also rejects duplicate keys, trailing data, and non-finite numeric values.
- Stable unique object IDs; asset IDs are catalogue references. References must resolve; reject cycles where parenting is supported.
- Coordinates: right-handed Blender space, Z up; distances in metres; angles in degrees at the declaration boundary, converted explicitly to radians.
- Opaque colours: sRGB hex at the boundary; convert to linear values for shader inputs. Output colour-management preset is explicit and version tested.
- Integer frames, inclusive start and end. V1 uses integer FPS. Duplicate keyframe times and conflicting tracks for the same property are errors.
- Only catalogue-approved typed property targets are animatable. Step and linear interpolation first; one documented ease-in-out preset after verified handle mapping. No user expressions or arbitrary paths.
- A fresh scene is reconstructed each run. Omitted properties use documented defaults; no inherited previous-job state.
- Explicit limits for input bytes, objects, text length, frames, resolution, and samples; exact limits are recorded in Phase 1. Exceeding them produces errors instead of clamping silently.
- External asset paths are excluded in v1. Outputs and executable paths are CLI arguments, not fields that get interpolated into Python.

## Compiler and runner

Use Go structs for the normalized representation and typed operations for generation. Templates contain trusted Python structure. Pass variable values in serialized JSON data read by the generated script, or a rigorously encoded data literal; never insert user text directly into Python source. Normalization resolves defaults and stable ordering. Emit a readable Python entry point plus a canonical data file; compilation requires no Blender process.

Use `os/exec` with cancellation and a timeout. Phase 0 must verify the exact command-line options and ordering against the chosen Blender build, including background operation, startup isolation, disabling auto-execution of embedded scripts, and nonzero exit on Python failure. Disabling auto-execution must not prevent the explicitly supplied trusted entry script from running.

Run each job in a unique private directory with an unmodified pack, explicit environment/configuration, bounded logs, and a completion manifest. Do not reuse mutable global Blender state. Refuse a nonempty destination unless an explicit overwrite mode is requested; restrict any replacement to tool-owned artifacts. Never delete an arbitrary destination tree. Failed runs retain useful diagnostics and are never marked complete. Cancellation terminates the child process tree using tested platform behavior.

This is execution hygiene, not a hostile-code sandbox. V1 only loads its own trusted pack and generated code. A future custom asset/plugin feature needs a separate trust design.

## Proposed CLI surface

The following surface was the baseline design for Phases 0–6; consult current CLI help and phase evidence for implemented behavior. Project/shot extensions below remain planned:

```text
blendctl doctor --blender /path/to/blender
blendctl schema
blendctl assets list --json
blendctl assets describe <asset-id> --json
blendctl guide
blendctl validate scene.json --json
blendctl compile scene.json --out build/
blendctl preview scene.json --frame 30 --out preview/
blendctl render scene.json --out render/
```

`compile` emits script, canonical data, and resource/version metadata. Preview and render also save the scene and manifest. Preview renders one selected in-range frame with explicit preview settings; render produces the full configured frame range. Any preview override appears in the manifest.

Diagnostics include stable code, stage, JSON pointer, message, and remediation hint. Machine output goes to stdout; progress/logging goes to stderr. Specify stable exit categories for validation, environment, execution, and cancellation. A manifest includes input hash, compiler version, resource hashes, Blender version, effective render settings, status, outputs, and timing. Timing and run IDs do not enter the compile identity.

## Determinism and verification boundaries

1. Same normalized input and compiler/resources → byte-identical generated script and data.
2. Same supported Blender environment → equivalent inspected scene structure, bindings, and animation samples.
3. Render regression checks use a pinned CPU-capable configuration, tolerant comparisons, and explicit environment metadata; do not compare raw `.blend` hashes or promise GPU equality.

Reference render aesthetic approval is separate from structural correctness. Preserve a contact sheet and require a named visual review result. Consult the [Blender command-line manual](https://docs.blender.org/manual/en/latest/advanced/command_line/arguments.html) and version-matched Python API during Phase 0; the linked manual was not retrievable during preparation of this kit.

## Approved extension: module resolver (Phase 7, not implemented)

The current implementation compiles a validated `spec.Declaration`; add a file-aware project resolver before that boundary, not another Blender backend. Keep raw single-file entry points compatible. Proposed flow:

Project manifest + explicit module graph → strict decoding → bounded dependency resolution + source map → assembled declaration → existing semantic validation/normalization → existing compiler/runner.

The manifest has its own version and explicit entry/module dependencies. The existing scene schema version is not silently redefined to mean a project. Modules contribute objects, lights and animation tracks; the entry owns camera, scene settings and asset-pack version. Material values stay inline as supported by the current scene model; shared named material libraries and parameter substitution are not Phase 7 requirements.

The phase CONTEXT defines logical namespaces and explicit import aliases. Qualify IDs deterministically and rewrite references through typed fields, not recursive string replacement. Reject ambiguity and cycles; no last-writer-wins merging, implicit overrides, templating, expressions, remote imports or wildcard scans. Import order may not determine scene meaning. Keep intentionally ordered lists (for example keyframes) ordered; canonicalize only unordered module/entity collections.

Resolve paths relative to the declaring file under the manifest directory. Normalize paths, resolve symlinks and reject escape from that root. Read bounded input snapshots once per job so validation and compilation use the same bytes. Bound modules, depth and total bytes; scene limits apply after assembly. JSON imports are declarations, not permission to load arbitrary Blender assets or execute files.

Carry source file and original JSON pointer for each declaration field/entity through reference rewriting and validator errors. Source filenames are project-relative. Existing diagnostic fields stay compatible; source context is additive. Source maps and provenance live in separate artifacts so absolute checkout paths never alter generated script/data. Existing duplicate-ID and reference validation still applies after assembly.

## Approved extension: independent shots (Phase 8, not implemented)

A project may declare multiple shot IDs, each with a scene entry and explicit module dependencies. A selected build resolves that shot's transitive dependency closure into one declaration and uses the current fresh-process rendering model. The loader validates the outer manifest for every command, but only selected-shot contents for a selected-shot build. An explicit project-wide validation traverses every shot.

Outputs and manifests are isolated per shot. Record content hashes for the resolved input snapshot, dependency graph and tool/resource versions; do not implement cache reuse until invalidation semantics have been designed and tested. A shared module is immutable input during a run. Changing it requires rebuilding dependent shots but does not justify silently rerendering every shot.

No multi-scene Blender runtime, NLE timeline assembly or distributed scheduler is implied. Parameterized component instantiation will require a later separate contract for types, defaults, bounds and reference binding.
