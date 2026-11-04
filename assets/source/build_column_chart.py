#!/usr/bin/env python3
"""Author MoBlend_ColumnChart Geometry Nodes group into the base pack.

Contract:
  - Input geometry: point cloud / mesh verts (one point per bar)
  - Point float attribute ``bar_value`` in [0,1] scales instance Z
  - Group input ``Bar Width`` (float) controls cube X/Y size

  blender --background --factory-startup --disable-autoexec --python-exit-code 1 \\
    --python assets/source/build_column_chart.py -- \\
    --pack internal/resources/pack/base.blend \\
    --out internal/resources/pack/base.blend
"""

from __future__ import annotations

import argparse
import json
import sys

import bpy


GROUP_NAME = "MoBlend_ColumnChart"
OBJECT_NAME = "MoBlend_ColumnChart"
PACK_VERSION = "0.1.0"


def parse_args(argv: list[str]) -> argparse.Namespace:
    if "--" in argv:
        argv = argv[argv.index("--") + 1 :]
    else:
        argv = []
    p = argparse.ArgumentParser()
    p.add_argument("--pack", required=True)
    p.add_argument("--out", required=True)
    return p.parse_args(argv)


def ensure_group() -> bpy.types.GeometryNodeTree:
    existing = bpy.data.node_groups.get(GROUP_NAME)
    if existing:
        bpy.data.node_groups.remove(existing)
    ng = bpy.data.node_groups.new(GROUP_NAME, "GeometryNodeTree")
    iface = ng.interface
    iface.new_socket(name="Geometry", in_out="OUTPUT", socket_type="NodeSocketGeometry")
    iface.new_socket(name="Geometry", in_out="INPUT", socket_type="NodeSocketGeometry")
    width = iface.new_socket(name="Bar Width", in_out="INPUT", socket_type="NodeSocketFloat")
    width.default_value = 0.4
    width.min_value = 0.05
    width.max_value = 2.0

    nodes = ng.nodes
    links = ng.links
    n_in = nodes.new("NodeGroupInput")
    n_out = nodes.new("NodeGroupOutput")
    n_in.location = (-800, 0)
    n_out.location = (600, 0)

    cube = nodes.new("GeometryNodeMeshCube")
    cube.location = (-400, -200)
    # Size driven later via combine + bar width on X/Y and 1.0 on Z before scale.

    combine = nodes.new("ShaderNodeCombineXYZ")
    combine.location = (-600, -200)
    combine.inputs[2].default_value = 1.0
    links.new(n_in.outputs["Bar Width"], combine.inputs[0])
    links.new(n_in.outputs["Bar Width"], combine.inputs[1])
    links.new(combine.outputs[0], cube.inputs["Size"])

    inst = nodes.new("GeometryNodeInstanceOnPoints")
    inst.location = (-100, 0)
    links.new(n_in.outputs["Geometry"], inst.inputs["Points"])
    links.new(cube.outputs["Mesh"], inst.inputs["Instance"])

    # Scale Z by bar_value attribute; X/Y stay 1.
    attr = nodes.new("GeometryNodeInputNamedAttribute")
    attr.data_type = "FLOAT"
    attr.location = (-400, 150)
    attr.inputs["Name"].default_value = "bar_value"

    scale_xyz = nodes.new("ShaderNodeCombineXYZ")
    scale_xyz.location = (-100, 200)
    scale_xyz.inputs[0].default_value = 1.0
    scale_xyz.inputs[1].default_value = 1.0
    links.new(attr.outputs["Attribute"], scale_xyz.inputs[2])

    scale_inst = nodes.new("GeometryNodeScaleInstances")
    scale_inst.location = (200, 0)
    links.new(inst.outputs["Instances"], scale_inst.inputs["Instances"])
    links.new(scale_xyz.outputs[0], scale_inst.inputs["Scale"])

    # Lift bars so base sits on Z=0: translate by (0,0, value/2) after scale of height 1.
    # Approximate with Translate Instances using same attribute * 0.5.
    mul = nodes.new("ShaderNodeMath")
    mul.operation = "MULTIPLY"
    mul.location = (200, 200)
    mul.inputs[1].default_value = 0.5
    links.new(attr.outputs["Attribute"], mul.inputs[0])

    offset = nodes.new("ShaderNodeCombineXYZ")
    offset.location = (400, 200)
    links.new(mul.outputs[0], offset.inputs[2])

    translate = nodes.new("GeometryNodeTranslateInstances")
    translate.location = (400, 0)
    links.new(scale_inst.outputs["Instances"], translate.inputs["Instances"])
    links.new(offset.outputs[0], translate.inputs["Translation"])

    realize = nodes.new("GeometryNodeRealizeInstances")
    realize.location = (600, 0)
    links.new(translate.outputs["Instances"], realize.inputs["Geometry"])
    links.new(realize.outputs["Geometry"], n_out.inputs["Geometry"])
    return ng


def seed_object(ng: bpy.types.GeometryNodeTree) -> None:
    old = bpy.data.objects.get(OBJECT_NAME)
    if old:
        bpy.data.objects.remove(old, do_unlink=True)
    mesh = bpy.data.meshes.new(OBJECT_NAME + "Mesh")
    # Five demo points on X.
    coords = [(-2.0 + i, 0.0, 0.0) for i in range(5)]
    mesh.from_pydata(coords, [], [])
    mesh.update()
    attr = mesh.attributes.new(name="bar_value", type="FLOAT", domain="POINT")
    for i, v in enumerate([0.2, 0.5, 0.8, 0.4, 0.65]):
        attr.data[i].value = v
    obj = bpy.data.objects.new(OBJECT_NAME, mesh)
    bpy.context.collection.objects.link(obj)
    mod = obj.modifiers.new("MoBlendChart", "NODES")
    mod.node_group = ng
    obj["moblend_asset_id"] = "column_chart"


def main() -> None:
    args = parse_args(sys.argv)
    bpy.ops.wm.open_mainfile(filepath=args.pack)
    ng = ensure_group()
    seed_object(ng)
    bpy.context.scene["moblend_pack_version"] = PACK_VERSION
    bpy.ops.wm.save_as_mainfile(filepath=args.out)
    inputs = []
    for item in ng.interface.items_tree:
        if getattr(item, "in_out", None) == "INPUT":
            inputs.append(
                {
                    "name": item.name,
                    "identifier": getattr(item, "identifier", item.name),
                    "socket_type": item.socket_type,
                }
            )
    print(
        "COLUMN_CHART_OK "
        + json.dumps({"pack_version": PACK_VERSION, "group": GROUP_NAME, "inputs": inputs})
    )


if __name__ == "__main__":
    main()
