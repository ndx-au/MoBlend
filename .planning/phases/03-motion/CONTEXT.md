# Phase 3 — context

## Pack versioning

Phase 0 pack is `0.0.1-phase0`. Adding `column_chart` node group requires a pack bump to `0.1.0` and matching `asset_pack_version` in schema/catalogue/examples (breaking for old declarations — acceptable pre-release).

## Chart contract (proposed)

| Input | Type | Bounds | Identifier strategy |
|---|---|---|---|
| values | float array (max 12) | 0..1 each | GN input socket name `Values` / identifier recorded after authoring |
| bar_count | int | 1..12 | socket `Bar Count` |

Bars are instances of a unit cube scaled in Z by value; spacing fixed. Material: single Principled colour override via object material or GN color input if simple.

## Animation

Phase 2 already inserts keyframes; Phase 3 adds frame-probe tests and pedestal/chart motion examples. Ease-in-out uses BEZIER + AUTO_CLAMPED handles on layered actions.
