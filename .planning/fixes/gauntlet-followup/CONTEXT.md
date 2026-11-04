# Repair decisions

- Import aliases resolve directly to module IDs; unused modules never participate in reference lookup. Bare entry camera ID stays unqualified. Cross-module references require alias:local.
- Source mappings remain separate artifacts and translate the entire original pointer prefix, including entry-owned scene/camera fields.
- Animated object or shader node-tree owner and exact property path select curves within the active action slot; no action-wide interpolation writes.
- Version probes and render execution share process-tree cancellation, isolated environment, bounded diagnostic capture. PNGs must decode completely with declared dimensions before success.
- Root inputs are bounded regular-file snapshots shared between detection and load. This does not add transactional multi-file snapshots or a hostile-filesystem sandbox.
- Existing user-authored documentation and prior review artifacts are included in the requested commit/push. No release will be published by this task.
- Frozen Gauntlet artifacts contain literal unified diffs: their context prefixes are significant whitespace. Scope whitespace-check exceptions to `.gauntlet/**` in `.gitattributes`, preserving hashes and keeping production source checks active.
