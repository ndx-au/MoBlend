# Phase 2 — context

## Emitter contract

- Trusted template Python reads `scene_data.json` beside the script (never interpolates user text into source).
- Compile identity excludes paths, timestamps, run IDs.
- Stable key ordering when marshaling canonical data (`encoding/json` with sorted object keys via deterministic marshal helper).

## Text / fonts

- Use Blender font loading from an embedded TTF extracted next to the job pack.
- Candidate: DejaVu Sans (bitstream) if license NOTICE is included; confirm before commit of binary font bytes.
- Until font is embedded, text objects may use Blender's built-in default font for Phase 2 smoke if documented as interim — prefer bundled font if redistributable file is available locally.

## Build-from-pristine

Each compile/render job extracts embedded `base.blend` anew; never mutate the embedded source.
