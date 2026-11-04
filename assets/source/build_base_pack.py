#!/usr/bin/env python3
"""Author the Phase 0 minimal first-party .blend pack.

Run under the pinned Blender binary, not system Python:

  blender --background --factory-startup --python-exit-code 1 \\
    --python assets/source/build_base_pack.py -- \\
    --out internal/resources/pack/base.blend
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

import bpy


PACK_VERSION = "0.0.1-phase0"
OBJECT_NAME = "MoBlend_BaseCube"
MATERIAL_NAME = "MoBlend_BaseMaterial"
CAMERA_NAME = "MoBlend_Camera"
LIGHT_NAME = "MoBlend_KeyLight"
# Opaque sRGB hex used as the known material colour contract for smoke probes.
BASE_COLOR_HEX = "#4A90D9"
BASE_COLOR_LINEAR = (0.066626, 0.274677, 0.686685, 1.0)  # approx sRGB #4A90D9


def parse_args(argv: list[str]) -> argparse.Namespace:
    if "--" in argv:
        argv = argv[argv.index("--") + 1 :]
    else:
        argv = []
    parser = argparse.ArgumentParser(description="Build MoBlend base.blend pack")
    parser.add_argument(
        "--out",
        required=True,
        help="Destination .blend path",
    )
    return parser.parse_args(argv)


def clear_scene() -> None:
    bpy.ops.wm.read_factory_settings(use_empty=True)


def build_scene() -> None:
    scene = bpy.context.scene
    scene.render.engine = "CYCLES"
    scene.cycles.device = "CPU"
    scene.cycles.samples = 16
    scene.cycles.use_denoising = False
    scene.render.resolution_x = 128
    scene.render.resolution_y = 128
    scene.render.resolution_percentage = 100
    scene.render.image_settings.file_format = "PNG"
    scene.render.image_settings.color_mode = "RGBA"
    scene.render.film_transparent = False
    # Explicit colour management for pack contract.
    scene.view_settings.view_transform = "Standard"
    scene.view_settings.look = "None"
    scene.display_settings.display_device = "sRGB"
    scene.frame_start = 1
    scene.frame_end = 1
    scene.frame_current = 1
    scene.render.fps = 24

    mesh = bpy.data.meshes.new(OBJECT_NAME + "Mesh")
    obj = bpy.data.objects.new(OBJECT_NAME, mesh)
    bpy.context.collection.objects.link(obj)
    bm_verts = [
        (-0.5, -0.5, -0.5),
        (0.5, -0.5, -0.5),
        (0.5, 0.5, -0.5),
        (-0.5, 0.5, -0.5),
        (-0.5, -0.5, 0.5),
        (0.5, -0.5, 0.5),
        (0.5, 0.5, 0.5),
        (-0.5, 0.5, 0.5),
    ]
    bm_faces = [
        (0, 1, 2, 3),
        (4, 5, 6, 7),
        (0, 1, 5, 4),
        (2, 3, 7, 6),
        (1, 2, 6, 5),
        (0, 3, 7, 4),
    ]
    mesh.from_pydata(bm_verts, [], bm_faces)
    mesh.update()
    obj.location = (0.0, 0.0, 0.0)

    mat = bpy.data.materials.new(MATERIAL_NAME)
    mat.use_nodes = True
    nodes = mat.node_tree.nodes
    links = mat.node_tree.links
    nodes.clear()
    out = nodes.new("ShaderNodeOutputMaterial")
    bsdf = nodes.new("ShaderNodeBsdfPrincipled")
    bsdf.inputs["Base Color"].default_value = BASE_COLOR_LINEAR
    bsdf.inputs["Roughness"].default_value = 0.4
    links.new(bsdf.outputs["BSDF"], out.inputs["Surface"])
    obj.data.materials.append(mat)

    cam_data = bpy.data.cameras.new(CAMERA_NAME)
    cam = bpy.data.objects.new(CAMERA_NAME, cam_data)
    bpy.context.collection.objects.link(cam)
    cam.location = (3.0, -3.0, 2.0)
    cam.rotation_euler = (1.1, 0.0, 0.785)
    scene.camera = cam

    light_data = bpy.data.lights.new(LIGHT_NAME, type="AREA")
    light_data.energy = 300.0
    light_data.size = 2.0
    light = bpy.data.objects.new(LIGHT_NAME, light_data)
    bpy.context.collection.objects.link(light)
    light.location = (2.0, -1.5, 3.0)

    # Custom properties for pack identity probes.
    scene["moblend_pack_version"] = PACK_VERSION
    scene["moblend_base_color_hex"] = BASE_COLOR_HEX
    obj["moblend_object_id"] = "base_cube"


def main() -> None:
    args = parse_args(sys.argv)
    out = Path(args.out).resolve()
    out.parent.mkdir(parents=True, exist_ok=True)

    clear_scene()
    build_scene()
    bpy.ops.wm.save_as_mainfile(filepath=str(out))

    meta = {
        "pack_version": PACK_VERSION,
        "out": str(out),
        "blender": bpy.app.version_string,
        "build_hash": bpy.app.build_hash.decode()
        if isinstance(bpy.app.build_hash, (bytes, bytearray))
        else str(bpy.app.build_hash),
        "objects": [OBJECT_NAME, CAMERA_NAME, LIGHT_NAME],
        "material": MATERIAL_NAME,
        "base_color_hex": BASE_COLOR_HEX,
    }
    print("BASE_PACK_OK " + json.dumps(meta))


if __name__ == "__main__":
    main()
