# Phase 7 — modular project composition context

Date: 2026-11-04. Status: implemented; see VERIFICATION.md.

## Reader and outcome

An implementation agent should be able to add multi-file project input to the existing RC without replacing its compiler or breaking single-file declarations. Agents need file-sized ownership boundaries and enough dependency information to edit one part of a production without loading every declaration into context.

## Fixed decisions

- Go owns resolution and compilation. No RAG, remote imports, script evaluation, raw Python fields or generic JSON patch/merge engine.
- Existing scene schema `1.0.0` remains compatible. Introduce separate project/module schema versions; project input is detected by an unambiguous version/type field, never a filename guess.
- Project root is the manifest directory. Imports resolve relative to their declaring file, with real-path containment checks. Snapshot bounded file bytes before validation/compilation.
- One resolved scene in Phase 7. Modules contribute objects, lights and animation. Entry owns scene, camera and asset-pack version. Inline materials retain their existing semantics.
- Explicit namespaces/imports; no directory scans, implicit override order or duplicate-ID overwrite. A module can depend on another through an explicit alias. Reusing a module as a parameterized instance is deferred.
- A module has a local namespace; composition assigns stable fully qualified IDs and rewrites typed target references. Reserve a collision-free internal qualification scheme and test compatibility with the existing ID validator before freezing schema syntax.
- Source maps retain the original module file and pointer after qualification. Add diagnostic source context without deleting/changing meanings of current fields. Sort diagnostics deterministically by source and pointer.
- Import order is not semantic. Do not sort keyframes or any other intentionally ordered collection blindly. Keep existing single-file canonicalization behavior intact.
- Existing per-scene object/light/track/render limits apply after composition. Define additional count/depth/byte limits in implementation context before accepting module input; file splitting is not permission for unbounded scenes.

## Integration observations

The existing CLI reads a file and calls `validate.ParseAndValidate` / `compile.CompileJSON`; the compiler consumes `spec.Declaration`. Add file-aware resolution above these entry points while preserving their existing APIs for single-file callers. Current diagnostics expose code, stage, pointer, message and hint, with no source file. Current materials are inline, so a standalone material-reference language would be extra scope.

## Definition of done

MOD-01–07 pass with recorded unit, CLI and real Blender evidence. The shipped guide describes how separate agents own modules, import contracts and shared-entry changes. Phase 8 multi-shot execution remains a separate gate. Existing blocked model evaluation stays separate from this phase's automated authoring-documentation checks.

## Frozen implementation contract (2026-11-04)

### Detection

- Project input: root object with `"kind":"project"` and `"project_version":"1.0.0"`.
- Module/entry files: `"kind":"module"|"entry"` and `"module_version":"1.0.0"`.
- Single-file scenes remain `"schema_version":"1.0.0"` without `kind`/`project_version`.

### Project manifest fields

| Field | Meaning |
|---|---|
| `project_version` | `1.0.0` |
| `kind` | `project` |
| `entry` | Module id of the entry module |
| `modules` | Explicit list `{ "id", "path" }` — paths relative to the manifest directory |

No directory scanning. Module `id` is the graph node name (not the namespace).

### Entry vs module

- **entry**: owns `schema_version`, `asset_pack_version`, `scene`, `camera`; may also declare `objects`/`lights`/`animation` and `imports`.
- **module**: owns `objects`/`lights`/`animation` and `imports` only (no scene/camera/versions).
- Both require `namespace` matching `^[A-Za-z][A-Za-z0-9]*$` (no underscore — reserved for qualification).

### Imports and qualification

- `imports`: `{ "alias", "module" }` where `module` is a project module id; `alias` matches the namespace regex.
- Local IDs keep scene id pattern `^[A-Za-z][A-Za-z0-9_]*$`.
- Fully qualified ID: `{namespace}__{local_id}` (double underscore).
- Same-module refs use bare local ids; cross-module refs use `{alias}:{local_id}` in module source, rewritten to the imported module's qualified ID.
- Camera id is not namespaced (entry-owned).

### Loader limits (reject, do not clamp)

| Limit | Value |
|---|---|
| Max modules | 32 |
| Max import depth | 8 |
| Max bytes per file | 256 KiB |
| Max total project bytes | 1 MiB |

Existing scene object/light/track/render limits apply to the **resolved** declaration.

### Provenance

- Diagnostics gain optional `source` (project-relative path); existing fields unchanged.
- Compile emits `source_map.json` beside script/data (not part of compile identity for script/data bytes).

## Deferred beyond Phase 7

Parameter substitution, module instancing, shared material libraries, incremental caching and shot scheduling (Phase 8).
