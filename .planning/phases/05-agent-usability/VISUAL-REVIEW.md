# Visual review — reference contact sheet

**Date:** 2026-11-04  
**Reviewer:** agent (image inspection of regenerated fixtures)  
**Artifacts:** `testdata/visual/contact/references_sheet.png`, `testdata/visual/frames/{title,pedestal,chart}.png`  
**Regenerate:** `scripts/make-contact-sheet.sh`

## Checklist

| Composition | Framing | Typography / shape | Lighting | Notes |
|---|---|---|---|---|
| Title | Pass | Pass — full word **MoBlend** readable | Pass — key+fill, mild grain at 16 samples | Extruded centred text; wide 384×160 plate |
| Pedestal | Pass | Pass — blue cube on dark cylinder | Pass | Product reads clearly |
| Chart | Pass | Pass — two independent bar groups, height variation visible | Pass | Isolation intentional (chartA / chartB) |

## Issues found and fixed this pass

1. 64×64 title was illegible → raised resolution, pulled framing, centred/extruded text.
2. Intermediate framing cropped to “oBlen” → wider plate + farther camera + scale 1.0.
3. ApplyResult JSON used Go field names (`PNG`) → snake_case tags for CLI consumers.

## Remaining visual debt (non-blocking)

- Full-res `examples/*.json` (1280×720) not on the contact sheet (slow); smoke fixtures are the gate images.
- Cycles noise at low samples; intentional for CI speed.
- Chart bars are uncoloured mesh (material on GN object may not tint realized instances) — acceptable for structure; colour binding polish later.
