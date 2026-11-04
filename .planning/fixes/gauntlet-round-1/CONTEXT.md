# Repair decisions

- Compile creates unique `moblend-job-*` directories, nested under the resolved shot ID when present. CLI prints actual artifact paths; never overwrites an earlier compilation.
- Selected shots retain closure-only content validation. Non-shot projects and explicit whole-project checks validate every listed dependency graph; unused modules do not contribute entities to the selected scene.
- Dependency digests describe the exact bytes read during loading, including the manifest. They never re-read inputs after rendering. This is a per-file snapshot, not a filesystem-wide transaction.
- Linux process groups are the tested cancellation boundary. Other operating systems remain unclaimed.
- Phase 5 model authoring gate remains unmet; fixture rendering cannot satisfy it.
