# Project definition

## Goal

Make visually compelling Blender motion graphics reliably controllable by coding agents without requiring them to manipulate Blender's full stateful API. Optimize for predictable outcomes, repeatable builds, and editable production artifacts.

A user describes a brief to their existing agent. The agent reads a compact local guide and catalogue, writes JSON, validates it, renders a preview, inspects the output, and revises the declaration. A Go tool owns validation, deterministic translation, asset extraction, Blender execution, and diagnostics.

## Origin and decision precedence

Source: [shared Gemini conversation](https://share.gemini.google/Pg4AEGlNs6K9), read on 2026-11-04. The discussion progressed from Blender video editing to a declarative bridge over Blender's Python API and reusable asset interfaces.

Explicit owner decisions override that discussion: **Go is the transpiler; no RAG in v1.** All other detailed boundaries here are proposed implementation defaults, chosen to make an achievable v1. The conversation's comparative performance, release-version, and LLM reliability claims are not verified requirements.

## V1 scope

- Local CLI usable directly by humans and coding agents.
- Strict versioned JSON scene specification with stable logical IDs.
- A small bundled `.blend` asset pack, schema, catalogue, examples, and Python templates.
- Text, basic solid forms, one parameterized Geometry Nodes asset, materials, camera, lighting, transforms, and bounded keyframe animation.
- One scene per declaration; frame range, resolution, frame rate, render settings, and a small fixed compositor preset.
- Preview stills, PNG animation frames, editable `.blend`, structured diagnostics, and provenance manifest.
- Three reference compositions: title reveal, product pedestal using a bundled geometric object, and data-driven column chart using the Geometry Nodes asset.
- Offline operation after Blender and the release executable are installed.

## Explicit non-goals

No RAG, embeddings, vector store, semantic asset search, remote model calls, built-in chatbot, browser UI, arbitrary Python execution, arbitrary node graph generation, arbitrary user `.blend` loading, plugin marketplace, or automatic asset downloads. No full NLE, VSE editing, audio synchronization, physics simulation, rig authoring, Grease Pencil authoring, or live scene mutation. No video codec or FFmpeg dependency for the v1 release gate. No promise to model all of Blender.

## Distribution and quality

Ship a Go executable with embedded resources; Blender remains a separately installed runtime. Extract resources to a private job directory because Blender needs filesystem paths. Use a pinned supported Blender version initially and expand only with tests. Linux amd64 is the first integration target; macOS arm64 and Windows amd64 are release targets only when actually tested.

The core quality claim is deterministic scene construction under a documented environment. Stable source does not guarantee byte-identical `.blend` files or cross-machine pixel identity. Visual ambition should come from carefully authored assets and composition, not a huge initial language.

## Success

A fresh user with supported Blender can use a release executable and shipped examples to render without installing language tooling or contacting a service. An agent supplied only the local guide, catalogue, and example can produce a valid changed composition. All three reference compositions build and render; repeated runs preserve the intended scene graph without duplication. The saved scene opens and remains editable in Blender.
