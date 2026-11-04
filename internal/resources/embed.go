package resources

import "embed"

//go:embed pack/base.blend
var packFS embed.FS

//go:embed schema/v1.json catalogue/v1.json templates/apply_scene.py fonts/DejaVuSans.ttf fonts/DejaVuSans.LICENSE agent/GUIDE.md
var contentFS embed.FS
