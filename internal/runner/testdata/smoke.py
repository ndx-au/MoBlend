"""Trusted Phase 0 smoke script: mutate pack, save scene, render PNG.

Invoked by the Go runner; not part of the declaration language.
"""

from __future__ import annotations

import argparse
import json
import struct
import sys
from pathlib import Path

import bpy

OBJECT_NAME = "MoBlend_BaseCube"
EXPECTED_PACK = "0.1.0"


def parse_args(argv: list[str]) -> argparse.Namespace:
    if "--" in argv:
        argv = argv[argv.index("--") + 1 :]
    else:
        argv = []
    p = argparse.ArgumentParser()
    p.add_argument("--out-dir", required=True)
    p.add_argument("--scale", type=float, default=1.5)
    p.add_argument("--fail", action="store_true", help="raise deliberately")
    p.add_argument("--width", type=int, default=64)
    p.add_argument("--height", type=int, default=64)
    return p.parse_args(argv)


def png_size(path: Path) -> tuple[int, int]:
    data = path.read_bytes()
    if data[:8] != b"\x89PNG\r\n\x1a\n":
        raise RuntimeError(f"not a PNG: {path}")
    w, h = struct.unpack(">II", data[16:24])
    return int(w), int(h)


def main() -> None:
    args = parse_args(sys.argv)
    if args.fail:
        raise RuntimeError("deliberate smoke failure")

    out = Path(args.out_dir).resolve()
    out.mkdir(parents=True, exist_ok=True)

    scene = bpy.context.scene
    pack_ver = scene.get("moblend_pack_version")
    if pack_ver != EXPECTED_PACK:
        raise RuntimeError(f"pack version {pack_ver!r} != {EXPECTED_PACK!r}")

    obj = bpy.data.objects.get(OBJECT_NAME)
    if obj is None:
        raise RuntimeError(f"missing object {OBJECT_NAME}")
    obj.scale = (args.scale, args.scale, args.scale)

    scene.render.engine = "CYCLES"
    scene.cycles.device = "CPU"
    scene.cycles.samples = 8
    scene.cycles.use_denoising = False
    scene.render.resolution_x = args.width
    scene.render.resolution_y = args.height
    scene.render.resolution_percentage = 100
    scene.render.image_settings.file_format = "PNG"
    scene.render.filepath = str(out / "smoke.png")

    blend_path = out / "smoke.blend"
    bpy.ops.wm.save_as_mainfile(filepath=str(blend_path))
    bpy.ops.render.render(write_still=True)

    png = out / "smoke.png"
    if not png.is_file() or png.stat().st_size < 64:
        raise RuntimeError("smoke.png missing or empty")
    w, h = png_size(png)
    if w != args.width or h != args.height:
        raise RuntimeError(f"png size {w}x{h} != {args.width}x{args.height}")

    # Nonempty content: at least one non-zero byte in pixel data region.
    raw = png.read_bytes()
    if not any(b != 0 for b in raw[100:]):
        raise RuntimeError("png appears empty/blank")

    info = {
        "pack_version": pack_ver,
        "scale": list(obj.scale),
        "png": str(png),
        "blend": str(blend_path),
        "width": w,
        "height": h,
        "png_bytes": png.stat().st_size,
        "blender": bpy.app.version_string,
        "build_hash": bpy.app.build_hash.decode()
        if isinstance(bpy.app.build_hash, (bytes, bytearray))
        else str(bpy.app.build_hash),
        "engine": scene.render.engine,
        "cycles_device": scene.cycles.device,
    }
    (out / "result.json").write_text(json.dumps(info, indent=2) + "\n")
    print("SMOKE_OK " + json.dumps(info))


if __name__ == "__main__":
    main()
