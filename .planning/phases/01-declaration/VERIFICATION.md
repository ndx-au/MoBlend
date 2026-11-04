# Phase 1 — verification evidence

**Date:** 2026-11-04  
**Status:** gate met (local, no Blender required for validate)

## Commands

```bash
gofmt -l -w .
go vet ./...
go test ./... -count=1 -timeout 10m
go build -o /tmp/blendctl ./cmd/blendctl
/tmp/blendctl schema >/dev/null
/tmp/blendctl assets list
/tmp/blendctl assets describe base_cube
/tmp/blendctl validate examples/title.json
/tmp/blendctl validate examples/pedestal.json
/tmp/blendctl validate examples/chart.json
/tmp/blendctl validate testdata/spec/negative/unknown_asset.json --json
```

## Gate checks

| Gate | Result |
|---|---|
| Positive fixtures / examples validate | PASS — title, pedestal, chart |
| Negative fixtures stable code + pointer | PASS — see `testdata/spec/negative/meta.json` + `TestNegativeFixtures` |
| Schema/model parity | PASS — `TestSchemaModelParityRequired` |
| validate does not launch Blender | PASS — pure Go path |
| Deferred fields recorded | PASS — `CONTEXT.md` |

## Negative expectation table

| Fixture | Code | Pointer |
|---|---|---|
| duplicate_key.json | SPEC_DUPLICATE_KEY | /schema_version |
| unknown_field.json | SPEC_UNKNOWN_FIELD | (empty; message names field) |
| bad_version.json | SPEC_SCHEMA_VERSION | /schema_version |
| unknown_asset.json | CAT_UNKNOWN_ASSET | /objects/0/asset_id |
| duplicate_id.json | SEM_DUPLICATE_ID | /objects/0/id |

## Notes

- `column_chart` is catalogue-valid for declarations; pack binding remains Phase 3.
- JSON Schema file is embedded and mirrored by Go checks; full draft-2020-12 engine not yet a hard dependency.
