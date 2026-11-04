# Composite leftover plan: MoBlend Gauntlet after-pass — fresh review after round 1 repairs

Round 1/3. Verdict: **revise**.

## Coverage

- correctness: 40 (0 critical, 4 high, 7 findings)
- tests: 19 (0 critical, 0 high, 5 findings)
- failure: 12 (0 critical, 0 high, 3 findings)
- concurrency: 8 (0 critical, 0 high, 2 findings)
- security: 8 (0 critical, 0 high, 2 findings)

## Must fix

### duplicate-import-alias-overwrite (correctness, high)

- Location:
```
internal/project/resolve.go:resolveManifest aliasToNS construction
```
- Bug:
```
Import aliases are assigned into a map without rejecting duplicates. A CLI probe with imports [{alias:title,module:title},{alias:title,module:other}] and a title:text scale track validates successfully. Compiling with the imports reversed changes that track from other__text to title__text. This violates the Phase 7 rejection of ambiguity and import-order independence.
```
- Predicted:
```
An agent can accidentally reuse an alias and silently animate the wrong module; reordering imports changes the resulting scene rather than producing a diagnostic.
```
- Task:
```
Reject duplicate aliases within each module before constructing aliasToNS, with a source-aware diagnostic and an import-reordering regression.
```

### entry-camera-track-unresolvable (correctness, high)

- Location:
```
internal/project/resolve.go:resolveRef / moduleMaps construction
```
- Bug:
```
Camera IDs remain entry-owned and unqualified, but the entry local reference map includes only objects and lights. Adding an entry animation target_id cam to the shipped modular title fixture fails with SEM_TRACK_TARGET (unknown target_id cam); the equivalent single-file declaration validates. Camera transform animation is supported by the existing scene language, so modular composition loses an existing capability.
```
- Predicted:
```
A supported camera move cannot be expressed in an entry or selected shot using the documented bare camera ID, blocking modular equivalents of valid single-file scenes.
```
- Task:
```
Resolve the entry camera ID explicitly to its unchanged ID, and cover modular/flat camera-animation equivalence.
```

### later-track-overwrites-interpolation (correctness, high)

- Location:
```
internal/resources/templates/apply_scene.py:apply_animation interpolation loop
```
- Bug:
```
After inserting each track, interpolation is assigned to every FCurve in the object action, including previously inserted properties. A bounded real probe calling the shipped function on Blender 5.3.0 Beta 7e7c173d45c8 inserted location step keys (frame 1 x=0; frame 9 x=8), followed by scale linear keys. At frame 3 location.x was 2.0 rather than the required held value 0; all location curves had been changed to LINEAR.
```
- Predicted:
```
Valid mixed-interpolation tracks on one object render the wrong motion, and reordering independent animation tracks changes their meaning.
```
- Task:
```
Apply interpolation only to the FCurves belonging to the current track property and action slot, then assert interior-frame values with mixed step/linear tracks in both orders.
```

### material-track-interpolation-never-applied (correctness, high)

- Location:
```
internal/resources/templates/apply_scene.py:apply_animation material keyframing and obj.animation_data guard
```
- Bug:
```
Material tracks insert keyframes into mat.node_tree, but interpolation is only assigned through obj.animation_data.action. On pinned Blender 5.3.0 Beta 7e7c173d45c8, a material-only object has no object animation_data and its node-tree curves remain BEZIER. A real shipped-function probe with frame 1/9 roughness step 0->1 and base-color linear black->white evaluated frame 3 as roughness 0.15625 (expected 0) and linear-red 0.15625 (expected 0.25).
```
- Predicted:
```
Validated material step and linear tracks silently render Blender default easing instead of their declared interpolation; adding a material track after an object track can also overwrite unrelated object interpolation.
```
- Task:
```
Resolve the actual animated datablock/action for material tracks and set interpolation on only that material property curves, with interior-frame Blender checks for step and linear material tracks.
```

## Residual risks

- **medium** (security)
- Id:
```
dispatch-unbounded-input-read
```
- Bug:
```
The public validate/compile/preview/render loader still calls os.ReadFile and DetectKind on the entire root input before the bounded project loader runs. Thus the round-one bounded-read repair does not bound the public entry point. Neither path requires a regular input file. A CLI probe with a 32 MiB whitespace-padded manifest returned LIM_FILE_BYTES only after reaching 37,760 KiB peak RSS; a FIFO root input remained blocked until the probe killed it after one second.
```
- Task:
```
Use one regular-file, size-bounded root-input read before DetectKind, then pass the same snapshot into project resolution; reject nonregular module inputs as well.
```

- **medium** (failure, concurrency)
- Id:
```
doctor-cancel-pipe-hang
```
- Bug:
```
Run now kills its process group and bounds pipe waiting, but the preceding --version probe uses a separate exec.CommandContext with neither configureProcessTree nor WaitDelay. A descendant retaining stdout/stderr keeps CombinedOutput blocked after cancellation or the nominal 15-second probe timeout. Bounded reproduction against the current built CLI: configure MOBLEND_BLENDER to a wrapper whose --version starts sleep 60 and waits, start preview, send SIGTERM after the child is ready. The CLI remained blocked after one second and the child was still alive; manually killing that child allowed cancellation and the failed marker to finish. The existing Run/CLI cancellation tests passed but exercise a fast --version branch.
```
- Task:
```
Apply the same dedicated process-group cancellation and bounded pipe wait to Doctor subprocesses, and extend the CLI signal regression to cancel while --version is running with a descendant.
```

- **medium** (tests)
- Id:
```
import-permutation-test-noop
```
- Bug:
```
The test credited with import-order determinism rewrites the one-element imports array to the same title import. It tests relocation/key serialization, but never changes import order; MOD-03 and the Phase 7 gate explicitly require this permutation.
```
- Task:
```
Use a valid graph with at least two imports and cross-module animation references, reverse imports and manifest module order, and assert identical normalized data/script bytes.
```

- **medium** (failure)
- Id:
```
logs-still-unbounded
```
- Bug:
```
Phase 4 explicitly requires bounded logs, but both process paths still accumulate all stdout/stderr with CombinedOutput. Run copies the accumulated output into Result and its returned error, and RunApply serializes the full error into the failed marker. A bounded fake-Blender probe emitting 4 MiB then exiting 1 produced a 4,194,402-byte failed marker and 4,194,344-byte CLI stderr without truncation. The repair evidence acknowledges this gap, but it remains implemented behavior.
```
- Task:
```
Replace CombinedOutput with a capped draining collector shared by execution and version probes, retain a bounded diagnostic tail with an explicit truncation marker, and test output beyond the cap.
```

- **medium** (tests)
- Id:
```
mixed-animation-channels-unprobed
```
- Bug:
```
Blender interpolation probes cover isolated location tracks and a rotation endpoint; none combines different interpolation modes on one object or samples animated material sockets. These exercise different action owners and curve-selection paths in apply_animation, so single-property probes do not cover the supported animation contract.
```
- Task:
```
Add a reopened Blender probe with step location plus linear rotation on one object and step/linear material tracks; assert interior samples independently and repeat with track order reversed.
```

- **medium** (correctness)
- Id:
```
module-entry-fields-silently-ignored
```
- Bug:
```
Typed decoding accepts scene, camera, schema_version and asset_pack_version on kind module, while the shipped module schema explicitly forbids all four. Only the selected entry fields are used. Adding scene:{fps:-10} to title.module.json still produces validate --json ok:true; that setting is silently discarded.
```
- Task:
```
Enforce the kind-dependent module/entry field contract during load and reject prohibited fields before assembly.
```

- **medium** (concurrency)
- Id:
```
orphan-namespace-map-order
```
- Bug:
```
The new all-listed graph check leaves orphan modules in loaded, but namespace uniqueness and moduleMaps are built only for the entry closure. resolveRef searches the entire loaded map by namespace, so an orphan with the namespace of an imported closure module can win randomized map iteration and return an unknown-ID error instead of resolving the explicitly imported module.
```
- Task:
```
Resolve alias references directly to the imported module ID and its moduleMaps entry; keep orphan validation separate from the selected closure lookup, and add a repeated-load regression with a same-namespace orphan.
```

- **medium** (failure)
- Id:
```
preview-empty-marked-complete
```
- Bug:
```
The repaired render path checks expected frames are regular and nonempty, while preview still accepts any existing preview.png via os.Stat. A zero-byte file or directory passes. In a bounded current-CLI failure injection, a fake supported-version Blender exited zero after writing a nonempty scene.blend and an empty preview.png; preview exited zero and wrote manifest.json with status complete. This is an output-failure injection, not a claim that the real Blender baseline produced an empty image.
```
- Task:
```
Require a regular nonempty preview PNG and validate its PNG header/dimensions before writing success; cover empty, directory and truncated preview outputs with failed-marker regressions.
```

- **medium** (correctness)
- Id:
```
qualified-reference-import-bypass
```
- Bug:
```
Any reference containing __ is returned unchanged, bypassing the explicit alias-import contract. A helper module with no imports can animate title__text when title is merely a sibling dependency of the entry; an isolated CLI probe validates this successfully. The frozen Phase 7 contract requires alias:local_id for cross-module references and explicit imports.
```
- Task:
```
Reject qualified IDs in module source and require cross-module targets to resolve through a declared import alias.
```

- **medium** (tests)
- Id:
```
s04-failed-shot-missing
```
- Bug:
```
The concurrent shot integration still runs only three successful previews and checks artifact existence/distinct job directories. Separate timeout/failure tests do not exercise the required concurrent successful/failed-shot matrix or assert manifest shot attribution. The repair verification correctly acknowledges this remaining gap.
```
- Task:
```
Run concurrent jobs with one deliberately failing shot and assert each successful manifest/dependency manifest has its expected shot ID, while only the failed job has a failed marker and no complete manifest.
```

- **medium** (tests)
- Id:
```
source-original-pointer-unasserted
```
- Bug:
```
The sole source-aware diagnostic test asserts only the source filename. Its invalid object is at index zero in both its module and the assembled declaration, and it never checks Code or Pointer; no project test asserts original pointers after sorting/reindexing.
```
- Task:
```
Add nested object and animation failures whose original indexes differ from assembled indexes, asserting exact code, relative source file and original source pointer in the resolver and CLI JSON result.
```

- **medium** (correctness)
- Id:
```
source-pointer-still-flattened
```
- Bug:
```
Semantic diagnostics add Source but retain flattened declaration array indices. In a CLI probe with one object in main.entry.json and one invalid-color object in title.module.json, the latter is reported as source title.module.json, pointer /objects/1/material/base_color although that file only has /objects/0. Entry scene/camera errors also have no source mapping. Phase 7 explicitly requires original file and pointer preservation.
```
- Task:
```
Translate each flattened pointer prefix to the stored SourceRef.Pointer and map entry-owned scene/camera diagnostics to the entry source.
```

- **medium** (security)
- Id:
```
unbounded-child-log-capture
```
- Bug:
```
Execution still uses exec.Cmd.CombinedOutput with no output cap in both the render launcher and version probe. The repair explicitly acknowledges in-memory logs, but the roadmap requires bounded logs; no bound or truncation path exists. Run also copies the captured bytes into Result.Stdout and formatted errors.
```
- Task:
```
Replace CombinedOutput with a capped log sink that drains both streams while retaining a bounded diagnostic tail and explicit truncation metadata; use it for Doctor too.
```

- **low** (tests)
- Id:
```
ease-inout-probe-gap
```
- Bug:
```
SCENE-02 still credits easing at intermediate frames, but there is no dedicated Blender interior-sample assertion for ease_in_out. Existing step/linear samples and rotation endpoint checks cannot distinguish an eased curve from linear interpolation. The repair verification leaves this probe open.
```
- Task:
```
Add a pinned-Blender easing fixture and probe asymmetric interior frames against the documented easing tolerance, in addition to its endpoints.
```


## Sequence

1. [correctness]
```
Reject duplicate aliases within each module before constructing aliasToNS, with a source-aware diagnostic and an import-reordering regression.
```
2. [correctness]
```
Resolve the entry camera ID explicitly to its unchanged ID, and cover modular/flat camera-animation equivalence.
```
3. [correctness]
```
Apply interpolation only to the FCurves belonging to the current track property and action slot, then assert interior-frame values with mixed step/linear tracks in both orders.
```
4. [correctness]
```
Resolve the actual animated datablock/action for material tracks and set interpolation on only that material property curves, with interior-frame Blender checks for step and linear material tracks.
```

## Guardrails

- Do not drop a finding without saying why it is invalid.
- After the next impl round, run `node scripts/judge.mjs compose` on an updated review (plan + new diff).
- This is round 1. If compose still says revise, another after-pass is allowed (max 3 rounds).
