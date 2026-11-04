# MoBlend
A transpiler for agent-generated declarative JSON to be processed into Blender's python API.

A planning kit for an agent-friendly motion graphics tool. Working CLI name: `blendctl` (provisional; check naming before publication).

**Status: Phases 0–6 complete for Linux amd64 RC `0.1.0-rc1`.** Dot Progress in `.progress/`; plans in `.planning/`.

```bash
go test ./...
VERSION=0.1.0-rc1 scripts/release-linux-amd64.sh
./dist/linux-amd64/blendctl doctor
./dist/linux-amd64/blendctl preview testdata/compile/title_small.json --frame 24 --out /tmp/moblend-preview
```

Requires Blender 5.3 on `PATH` or `MOBLEND_BLENDER`. See `docs/QUICKSTART.md`. LLM/visual review evidence is explicitly blocked — do not claim them without review.

## Product in one sentence

Agents describe motion graphics in a small, versioned JSON format; a Go compiler converts that description into deterministic Blender Python, applies it to bundled assets, and produces an editable scene and rendered frames.

## Read order

1. `AGENTS.md` — execution rules.
2. `.planning/PROJECT.md` — scope and decisions.
3. `.planning/ARCHITECTURE.md` — contracts and semantics.
4. `.planning/REQUIREMENTS.md` — acceptance criteria.
5. `.planning/ROADMAP.md` — ordered implementation slices.
6. `.planning/STATE.md` — current handoff.

V1 needs an installed, supported Blender executable. Users do not need a separate Python installation, Go toolchain, Node runtime, vector database, model account, or hosted service. Development requires Go and Blender. The tool does not provide its own natural-language model; an existing coding agent authors the JSON.

Primary outputs are an editable `.blend`, PNG frames, and a machine-readable manifest. Video encoding, MCP, and broader Blender capabilities are later extensions.
