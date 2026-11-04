"""Trusted MoBlend apply script — reads scene_data.json; never eval user strings."""

from __future__ import annotations

import json
import math
import struct
import sys
from pathlib import Path

import bpy

DATA_NAME = "scene_data.json"


def srgb_hex_to_linear(hex_color: str) -> tuple[float, float, float, float]:
    h = hex_color.lstrip("#")
    r, g, b = (int(h[i : i + 2], 16) / 255.0 for i in (0, 2, 4))

    def chan(c: float) -> float:
        return c / 12.92 if c <= 0.04045 else ((c + 0.055) / 1.055) ** 2.4

    return (chan(r), chan(g), chan(b), 1.0)


def clear_to_empty() -> None:
    bpy.ops.wm.read_factory_settings(use_empty=True)


def ensure_camera(scene, cam: dict) -> bpy.types.Object:
    data = bpy.data.cameras.new(cam["id"])
    obj = bpy.data.objects.new(cam["id"], data)
    bpy.context.collection.objects.link(obj)
    obj.location = cam["location"]
    obj.rotation_euler = [math.radians(v) for v in cam["rotation_euler_deg"]]
    fl = cam.get("focal_length_mm")
    if fl:
        data.lens = float(fl)
    scene.camera = obj
    return obj


def ensure_light(light: dict) -> bpy.types.Object:
    data = bpy.data.lights.new(light["id"], type=light["type"])
    data.energy = float(light["energy"])
    if light.get("color"):
        data.color = srgb_hex_to_linear(light["color"])[:3]
    if light.get("size") is not None and hasattr(data, "size"):
        data.size = float(light["size"])
    obj = bpy.data.objects.new(light["id"], data)
    bpy.context.collection.objects.link(obj)
    obj.location = light["location"]
    if light.get("rotation_euler_deg"):
        obj.rotation_euler = [math.radians(v) for v in light["rotation_euler_deg"]]
    return obj


def assign_material(obj: bpy.types.Object, mat_spec: dict, name: str) -> None:
    mat = bpy.data.materials.new(name)
    mat.use_nodes = True
    nodes = mat.node_tree.nodes
    links = mat.node_tree.links
    nodes.clear()
    out = nodes.new("ShaderNodeOutputMaterial")
    bsdf = nodes.new("ShaderNodeBsdfPrincipled")
    bsdf.inputs["Base Color"].default_value = srgb_hex_to_linear(mat_spec["base_color"])
    rough = mat_spec.get("roughness")
    if rough is not None:
        bsdf.inputs["Roughness"].default_value = float(rough)
    links.new(bsdf.outputs["BSDF"], out.inputs["Surface"])
    if obj.data.materials:
        obj.data.materials[0] = mat
    else:
        obj.data.materials.append(mat)


def make_form(form: str, name: str) -> bpy.types.Object:
    if form == "cube":
        mesh = bpy.data.meshes.new(name + "Mesh")
        obj = bpy.data.objects.new(name, mesh)
        bpy.context.collection.objects.link(obj)
        verts = [
            (-0.5, -0.5, -0.5),
            (0.5, -0.5, -0.5),
            (0.5, 0.5, -0.5),
            (-0.5, 0.5, -0.5),
            (-0.5, -0.5, 0.5),
            (0.5, -0.5, 0.5),
            (0.5, 0.5, 0.5),
            (-0.5, 0.5, 0.5),
        ]
        faces = [
            (0, 1, 2, 3),
            (4, 5, 6, 7),
            (0, 1, 5, 4),
            (2, 3, 7, 6),
            (1, 2, 6, 5),
            (0, 3, 7, 4),
        ]
        mesh.from_pydata(verts, [], faces)
        mesh.update()
        return obj
    if form == "plane":
        mesh = bpy.data.meshes.new(name + "Mesh")
        obj = bpy.data.objects.new(name, mesh)
        bpy.context.collection.objects.link(obj)
        mesh.from_pydata(
            [(-1, -1, 0), (1, -1, 0), (1, 1, 0), (-1, 1, 0)],
            [],
            [(0, 1, 2, 3)],
        )
        mesh.update()
        return obj
    if form == "cylinder":
        # Use mesh primitive operator with explicit context override avoided:
        # create via bpy.ops with scene update — set active after.
        bpy.ops.mesh.primitive_cylinder_add(radius=0.5, depth=1.0, location=(0, 0, 0))
        obj = bpy.context.active_object
        obj.name = name
        return obj
    raise RuntimeError(f"unsupported form {form}")


def make_text(spec: dict, font_path: Path) -> bpy.types.Object:
    curve = bpy.data.curves.new(spec["id"], type="FONT")
    curve.body = spec["text"]
    if font_path.is_file():
        font = bpy.data.fonts.load(str(font_path))
        curve.font = font
    curve.size = float(spec.get("size", 1.0))
    obj = bpy.data.objects.new(spec["id"], curve)
    bpy.context.collection.objects.link(obj)
    return obj


def _append_object(pack_blend: Path, object_name: str) -> bpy.types.Object:
    before = {o.name for o in bpy.data.objects}
    with bpy.data.libraries.load(str(pack_blend), link=False) as (data_from, data_to):
        if object_name not in data_from.objects:
            raise RuntimeError(f"pack missing object {object_name}")
        data_to.objects = [object_name]
    appended = [o for o in bpy.data.objects if o.name not in before]
    if not appended:
        # Name collision: find by name
        obj = bpy.data.objects.get(object_name)
        if obj is None:
            raise RuntimeError(f"failed to append {object_name}")
        return obj
    obj = appended[0]
    if obj.name not in bpy.context.collection.objects:
        bpy.context.collection.objects.link(obj)
    return obj


def _append_node_group(pack_blend: Path, group_name: str) -> bpy.types.NodeTree:
    before = {g.name for g in bpy.data.node_groups}
    with bpy.data.libraries.load(str(pack_blend), link=False) as (data_from, data_to):
        if group_name not in data_from.node_groups:
            raise RuntimeError(f"pack missing node group {group_name}")
        data_to.node_groups = [group_name]
    added = [g for g in bpy.data.node_groups if g.name not in before]
    if added:
        return added[0]
    ng = bpy.data.node_groups.get(group_name)
    if ng is None:
        raise RuntimeError(f"failed to append node group {group_name}")
    return ng


def load_asset(asset_id: str, name: str, pack_blend: Path, params: dict | None = None) -> bpy.types.Object:
    params = params or {}
    if asset_id == "base_cube":
        obj = _append_object(pack_blend, "MoBlend_BaseCube")
        # Duplicate mesh/material data so instances do not share mutations.
        if obj.data:
            obj.data = obj.data.copy()
        obj.name = name
        return obj

    if asset_id == "column_chart":
        values = params.get("values") or [0.2, 0.5, 0.8, 0.4, 0.65]
        if not isinstance(values, list) or not values:
            raise RuntimeError("column_chart requires params.values array")
        if len(values) > 12:
            raise RuntimeError("column_chart supports at most 12 values")
        values = [max(0.0, min(1.0, float(v))) for v in values]
        bar_count = int(params.get("bar_count") or len(values))
        bar_count = max(1, min(12, bar_count))
        values = values[:bar_count]
        while len(values) < bar_count:
            values.append(0.0)
        bar_width = float(params.get("bar_width") or 0.4)

        spacing = max(bar_width * 1.5, 0.5)
        start = -0.5 * spacing * (len(values) - 1)
        coords = [(start + i * spacing, 0.0, 0.0) for i in range(len(values))]
        mesh = bpy.data.meshes.new(name + "Mesh")
        mesh.from_pydata(coords, [], [])
        mesh.update()
        attr = mesh.attributes.new(name="bar_value", type="FLOAT", domain="POINT")
        for i, v in enumerate(values):
            attr.data[i].value = v

        obj = bpy.data.objects.new(name, mesh)
        bpy.context.collection.objects.link(obj)
        ng = _append_node_group(pack_blend, "MoBlend_ColumnChart")
        ng = ng.copy()
        ng.name = f"MoBlend_ColumnChart_{name}"
        mod = obj.modifiers.new("MoBlendChart", "NODES")
        mod.node_group = ng
        # Blender 5.3+: modifier inputs live on mod.properties.inputs.<identifier>.value
        inputs = mod.properties.inputs
        if hasattr(inputs, "Socket_2"):
            inputs.Socket_2.value = bar_width
        else:
            # Per-instance fallback: set group interface default before evaluation.
            for item in ng.interface.items_tree:
                if getattr(item, "in_out", None) == "INPUT" and item.name == "Bar Width":
                    item.default_value = bar_width
                    break
        obj["moblend_asset_id"] = "column_chart"
        obj["moblend_bar_values"] = values
        return obj

    raise RuntimeError(f"asset {asset_id} not bindable")


def apply_transform(obj: bpy.types.Object, spec: dict) -> None:
    if spec.get("location") is not None:
        obj.location = spec["location"]
    if spec.get("rotation_euler_deg") is not None:
        obj.rotation_euler = [math.radians(v) for v in spec["rotation_euler_deg"]]
    if spec.get("scale") is not None:
        obj.scale = spec["scale"]


def configure_scene(scene, sc: dict) -> None:
    scene.frame_start = int(sc["frame_start"])
    scene.frame_end = int(sc["frame_end"])
    scene.frame_current = int(sc["frame_start"])
    scene.render.fps = int(sc["fps"])
    scene.render.resolution_x = int(sc["resolution"][0])
    scene.render.resolution_y = int(sc["resolution"][1])
    scene.render.resolution_percentage = 100
    scene.render.engine = "CYCLES"
    scene.cycles.device = "CPU"
    scene.cycles.samples = int(sc["samples"])
    scene.cycles.use_denoising = False
    scene.render.image_settings.file_format = "PNG"
    scene.view_settings.view_transform = "Standard"
    if sc.get("world_color"):
        world = bpy.data.worlds.new("MoBlendWorld")
        scene.world = world
        world.use_nodes = True
        bg = world.node_tree.nodes.get("Background")
        if bg:
            bg.inputs[0].default_value = srgb_hex_to_linear(sc["world_color"])


def iter_action_fcurves(action):
    """Yield FCurves from legacy or layered (Blender 5+) actions."""
    if hasattr(action, "fcurves") and action.fcurves:
        for fc in action.fcurves:
            yield fc
        return
    for layer in getattr(action, "layers", []) or []:
        for strip in getattr(layer, "strips", []) or []:
            for bag in getattr(strip, "channelbags", []) or []:
                for fc in getattr(bag, "fcurves", []) or []:
                    yield fc


def apply_animation(tracks: list) -> None:
    for track in tracks:
        obj = bpy.data.objects.get(track["target_id"])
        if obj is None:
            raise RuntimeError(f"missing animation target {track['target_id']}")
        prop = track["property"]
        interp = track.get("interpolation") or "linear"
        for key in track["keys"]:
            frame = int(key["frame"])
            value = key["value"]
            if prop == "location":
                obj.location = value
                obj.keyframe_insert(data_path="location", frame=frame)
            elif prop == "rotation_euler_deg":
                obj.rotation_euler = [math.radians(v) for v in value]
                obj.keyframe_insert(data_path="rotation_euler", frame=frame)
            elif prop == "scale":
                obj.scale = value
                obj.keyframe_insert(data_path="scale", frame=frame)
            elif prop == "material.base_color":
                mat = obj.data.materials[0]
                bsdf = mat.node_tree.nodes.get("Principled BSDF")
                if bsdf is None:
                    for n in mat.node_tree.nodes:
                        if n.type == "BSDF_PRINCIPLED":
                            bsdf = n
                            break
                bsdf.inputs["Base Color"].default_value = srgb_hex_to_linear(value)
                bsdf.inputs["Base Color"].keyframe_insert("default_value", frame=frame)
            elif prop == "material.roughness":
                mat = obj.data.materials[0]
                bsdf = next(n for n in mat.node_tree.nodes if n.type == "BSDF_PRINCIPLED")
                bsdf.inputs["Roughness"].default_value = float(value)
                bsdf.inputs["Roughness"].keyframe_insert("default_value", frame=frame)
            else:
                raise RuntimeError(f"unsupported anim property {prop}")
        # Interpolation: Blender 5.x layered actions expose fcurves via channelbags.
        if obj.animation_data and obj.animation_data.action:
            for fc in iter_action_fcurves(obj.animation_data.action):
                for kp in fc.keyframe_points:
                    if interp == "step":
                        kp.interpolation = "CONSTANT"
                    elif interp == "ease_in_out":
                        kp.interpolation = "BEZIER"
                        kp.handle_left_type = "AUTO_CLAMPED"
                        kp.handle_right_type = "AUTO_CLAMPED"
                    else:
                        kp.interpolation = "LINEAR"


def png_ok(path: Path, w: int, h: int) -> None:
    data = path.read_bytes()
    if data[:8] != b"\x89PNG\r\n\x1a\n":
        raise RuntimeError("not png")
    pw, ph = struct.unpack(">II", data[16:24])
    if pw != w or ph != h:
        raise RuntimeError(f"size {pw}x{ph}")
    if path.stat().st_size < 64:
        raise RuntimeError("empty png")


def main() -> None:
    argv = sys.argv
    argv = argv[argv.index("--") + 1 :] if "--" in argv else []
    # Expected: --data scene_data.json --out-dir DIR [--render-frame N]
    args = {}
    i = 0
    while i < len(argv):
        if argv[i].startswith("--") and i + 1 < len(argv):
            args[argv[i][2:].replace("-", "_")] = argv[i + 1]
            i += 2
        else:
            i += 1

    job_dir = Path(args.get("out_dir", ".")).resolve()
    data_path = Path(args.get("data", job_dir / DATA_NAME)).resolve()
    data = json.loads(data_path.read_text())
    font_path = Path(args.get("font", job_dir / "fonts" / "DejaVuSans.ttf"))
    pack_blend = Path(args.get("pack", job_dir / "base.blend"))

    # Start from opened pack if provided as blend file on CLI; still clear extras carefully.
    # Runner opens base.blend then runs this script — keep pack objects only as asset source.
    # Rebuild scene empty then load assets from pack path.
    clear_to_empty()
    scene = bpy.context.scene
    configure_scene(scene, data["scene"])
    ensure_camera(scene, data["camera"])
    for light in data.get("lights") or []:
        ensure_light(light)

    for obj_spec in data["objects"]:
        otype = obj_spec["type"]
        if otype == "text":
            obj = make_text(obj_spec, font_path)
        elif otype == "form":
            obj = make_form(obj_spec["form"], obj_spec["id"])
        elif otype == "asset":
            obj = load_asset(
                obj_spec["asset_id"],
                obj_spec["id"],
                pack_blend,
                obj_spec.get("params") or {},
            )
        else:
            raise RuntimeError(f"bad type {otype}")
        apply_transform(obj, obj_spec)
        if obj_spec.get("material"):
            assign_material(obj, obj_spec["material"], obj_spec["id"] + "Mat")

    apply_animation(data.get("animation") or [])

    blend_out = job_dir / "scene.blend"
    bpy.ops.wm.save_as_mainfile(filepath=str(blend_out))

    mode = args.get("mode", "preview")
    w = int(data["scene"]["resolution"][0])
    h = int(data["scene"]["resolution"][1])
    if mode == "render":
        frames_dir = job_dir / "frames"
        frames_dir.mkdir(parents=True, exist_ok=True)
        start = int(data["scene"]["frame_start"])
        end = int(data["scene"]["frame_end"])
        written = []
        for frame in range(start, end + 1):
            scene.frame_set(frame)
            png = frames_dir / f"frame_{frame:04d}.png"
            scene.render.filepath = str(png)
            bpy.ops.render.render(write_still=True)
            png_ok(png, w, h)
            written.append(str(png))
        print(
            "APPLY_OK "
            + json.dumps({"blend": str(blend_out), "mode": "render", "frames": written})
        )
    else:
        frame = int(args.get("render_frame", data["scene"]["frame_start"]))
        scene.frame_set(frame)
        png = job_dir / "preview.png"
        scene.render.filepath = str(png)
        bpy.ops.render.render(write_still=True)
        png_ok(png, w, h)
        print(
            "APPLY_OK "
            + json.dumps({"blend": str(blend_out), "png": str(png), "frame": frame, "mode": "preview"})
        )


if __name__ == "__main__":
    main()
