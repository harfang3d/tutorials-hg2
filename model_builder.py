# Model builder usage

import harfang as hg
from math import sin


def create_grid_model_with_model_builder(vtx_layout, origin_pos, quad_size, range_x, range_z, center_on_origin, time):
    """
    vtx_layout : VertexLayout from harfang. Need the position and the normal.
    origin_pos : position where the model will be created
    quad_size : size of a quad (the grid is composed of multiple quads with the same size)
    range_x : number of quads on the x axis
    range_z : number of quads on the z axis
    center_on_origin : True if the model need to be centered on the origin_pos
    time : current clock in seconds to add movement in the mesh (waves movement) 
    """

    # Get the grid vertex position 
    vertex_positions = compute_grid_vertex_position(origin_pos, quad_size, range_x, range_z, center_on_origin, time)

    # Get the grid triangles
    grid_triangles = compute_triangles_for_grid_vertex(range_x, range_z)

    # Get the grid vertex normals
    vertex_normals = compute_vertex_normals(grid_triangles, vertex_positions)

    # Create a model builder
    mdl_builder = hg.ModelBuilder()

    # Add all vertex to the model builder
    vertex_indices = []
    for pos, nrm in zip(vertex_positions, vertex_normals):
        v = hg.Vertex()
        v.pos = pos
        v.normal = nrm
        v_id = mdl_builder.AddVertex(v)
        vertex_indices.append(v_id)

    # Add the triangles to the model builder
    for (i0, i1, i2) in grid_triangles:
        mdl_builder.AddTriangle(vertex_indices[i0], vertex_indices[i1], vertex_indices[i2])

    mdl_builder.EndList(0)

    # Create athe model from the model builder
    mdl = mdl_builder.MakeModel(vtx_layout)

    return mdl


def compute_grid_vertex_position(origin_pos, quad_size, range_x, range_z, center_on_origin, time):
    # Create offset if the model need to be center on the origin_pos
    if center_on_origin:
        offset = hg.Vec3(range_x * quad_size / 2, 0, range_z * quad_size / 2)
    else:
        offset = hg.Vec3(0, 0, 0)
    
    # Create the positions list
    positions = []

    # For all the vertex of the grid (add 1 to the range to get the last vertex on each axes, because the range corresponding to the number of quad to draw)
    for iz in range(range_z + 1):
        for ix in range(range_x + 1):
            # Create it's position
            pos_x = origin_pos.x + ix * quad_size - offset.x
            pos_z = origin_pos.z + iz * quad_size - offset.z
            pos_y = sin(pos_x) * sin(pos_z) * sin(time) # use time to add some movement
            positions.append(hg.Vec3(pos_x, pos_y, pos_z))

    return positions


def compute_triangles_for_grid_vertex(range_x, range_z):
    # Create the triangles list
    triangles = []

    # Use to get the vertex id in the positions list according to the x and z position in the grid
    def get_vertex_id(ix, iz):
        return iz * (range_x + 1) + ix

    # For all the quad of the grid
    for iz in range(range_z):
        for ix in range(range_x):
            # Get the 4 vertex that compose the quad
            a = get_vertex_id(ix, iz) # bottom left
            b = get_vertex_id(ix, iz + 1) # top left
            c = get_vertex_id(ix + 1, iz + 1) # top right
            d = get_vertex_id(ix + 1, iz) # bottom right

            triangles.append((d, c, b))
            triangles.append((b, a, d))

    return triangles


def compute_vertex_normals(triangles, positions):
    # Create normal list with 1 vec3 per vertex position
    normals = [hg.Vec3(0, 0, 0) for _ in positions]

    # For all the triangle
    for v0, v1, v2 in triangles:
        # Get the corresponding vertex position
        p0, p1, p2 = positions[v0], positions[v1], positions[v2]
        face_normal = hg.Cross(p0 - p1, p2 - p1)

        # Guard against degenerate triangles (zero area)
        if face_normal.x == 0 and face_normal.y == 0 and face_normal.z == 0:
            continue

        # Accumulate (area-weighted) face normal to each vertex
        normals[v0] += face_normal
        normals[v1] += face_normal
        normals[v2] += face_normal
    
    # Normalize all vertex normals
    for i in range(len(normals)):
        if hg.Len(normals[i]) > 0:
            normals[i] = hg.Normalize(normals[i]) 
        else :
            normals[i] = hg.Vec3(0, 1, 0)

    return normals



# Init render and resources
hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Model builder', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

hg.AddAssetsFolder('resources_compiled')

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

# Create materials
prg_ref = hg.LoadPipelineProgramRefFromAssets('core/shader/pbr.hps', res, hg.GetForwardPipelineInfo())
plane_material = hg.CreateMaterial(prg_ref, 'uBaseOpacityColor', hg.Vec4(0.5, 0.5, 0.5), 'uOcclusionRoughnessMetalnessColor', hg.Vec4(1, 1, 0.25))

# Setup scene
scene = hg.Scene()
hg.LoadSceneFromAssets("probe_scene/pbr.scn", scene, res, hg.GetForwardPipelineInfo())

cam = hg.CreateCamera(scene, hg.TransformationMat4(hg.Vec3(0, 6, -12), hg.Vec3(hg.DegreeToRadian(30), 0, 0)), 0.01, 1000)
scene.SetCurrentCamera(cam)

# light = hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(0, 5, 0)), 0)
light_mtx = hg.TransformationMat4(hg.Vec3(8, 5, 0), hg.Vec3(hg.DegreeToRadian(50), hg.DegreeToRadian(-90), hg.DegreeToRadian(-90)))
inner_angle = hg.DegreeToRadian(30)
outer_angle = hg.DegreeToRadian(45)
light_color = hg.Color(1, 1, 1, 1)
light = hg.CreateSpotLight(scene, light_mtx, 0, inner_angle, outer_angle, light_color, 1, light_color, 1, 1, hg.LST_Map, 0.0)

sphere_node = scene.GetNode("sphere")
# sphere_node.Disable()

# Create plane model
vtx_layout = hg.VertexLayoutPosFloatNormUInt8() 
# vtx_layout = hg.VertexLayoutPosFloatColorFloat()

grid_start_pos = hg.Vec3(0, 0, 0)
quad_size = 0.25
# Create the grid model from a model builder
grid_mdl = create_grid_model_with_model_builder(vtx_layout, grid_start_pos, quad_size, 40, 40, True, 1)
# Add the grid model to the PipelineResources and get the corresponding model ref 
grid_mdl_ref = res.AddModel('grid', grid_mdl)
# Create a node from the model ref
grid_node = hg.CreateObject(scene, hg.TransformationMat4(grid_start_pos, hg.Vec3(0, 0, 0)), grid_mdl_ref, [plane_material])

# Init input
keyboard = hg.Keyboard()

# main loop
while not keyboard.Down(hg.K_Escape) and hg.IsWindowOpen(win):
    dt = hg.TickClock()
    dts = hg.time_to_sec_f(dt)
    current_time = hg.time_to_sec_f(hg.GetClock())

    # Rotate the grid node
    # grid_rot = grid_node.GetTransform().GetRot()
    # rot = grid_rot + hg.Vec3(0, 0.5 * dts, 0)
    # grid_node.GetTransform().SetRot(rot)

    new_grid_mdl = create_grid_model_with_model_builder(vtx_layout, grid_start_pos, quad_size, 40, 40, True, current_time)
    res.UpdateModel(grid_mdl_ref, new_grid_mdl)

    scene.Update(dt)
    hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), True, pipeline, res)

    hg.Frame()
    hg.UpdateWindow(win)

hg.RenderShutdown()
hg.DestroyWindow(win)