# Troubleshooting

| Symptom | Check |
|---|---|
| `ENV_BLENDER_MISSING` | Install Blender 5.3 or set `MOBLEND_BLENDER` |
| `ENV_BLENDER_UNSUPPORTED` | Wrong major.minor; need 5.3 family |
| `validate` fails with pointer | Fix JSON at the reported pointer; see guide diagnostics table |
| `manifest.failed.json` present | Job did not complete; read `code`/`message`; no `status: complete` |
| Empty/missing PNG | Ensure Cycles CPU render finished; check Blender stderr |
| Chart looks wrong | `params.values` in 0..1; `bar_count` ≤ 12; pack `0.1.0` |
| Text missing glyphs | Bundled DejaVu Sans; exotic scripts may need a later font pack |

Rebuild pack assets:

```bash
scripts/build-base-pack.sh
scripts/build-column-chart.sh
```
