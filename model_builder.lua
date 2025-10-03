-- Model builder usage

hg = require("harfang")

function create_grid_model_with_model_builder(vtx_layout, origin_pos, quad_size, range_x, range_z, center_on_origin, time)
    --[[
    vtx_layout : VertexLayout from harfang. Need the position and the normal.
    origin_pos : position where the model will be created
    quad_size : size of a quad (the grid is composed of multiple quads with the same size)
    range_x : number of quads on the x axis
    range_z : number of quads on the z axis
    center_on_origin : True if the model need to be centered on the origin_pos
    time : current clock in seconds to add movement in the mesh (waves movement) 
    ]]

    -- Get the grid vertex position 
    local vertex_positions = compute_grid_vertex_position(origin_pos, quad_size, range_x, range_z, center_on_origin, time)

    -- Get the grid triangles
    local grid_triangles = compute_triangles_for_grid_vertex(range_x, range_z)

    -- Get the grid vertex normals
    local vertex_normals = compute_vertex_normals(grid_triangles, vertex_positions)

    -- Create a model builder
    local mdl_builder = hg.ModelBuilder()

    -- Add all vertex to the model builder
    local vertex_indices = {}
    for i=1, #vertex_positions do
        local v = hg.Vertex()
        v.pos = vertex_positions[i]
        v.normal = vertex_normals[i]
        local v_id = mdl_builder:AddVertex(v)
        table.insert(vertex_indices, v_id)
    end

    -- Add the triangles to the model builder
    for i=1, #grid_triangles do
        local i0, i1, i2 = grid_triangles[i][1], grid_triangles[i][2], grid_triangles[i][3]
        mdl_builder:AddTriangle(vertex_indices[i0], vertex_indices[i1], vertex_indices[i2])
    end

    mdl_builder:EndList(0)

    -- Create athe model from the model builder
    local mdl = mdl_builder:MakeModel(vtx_layout)

    return mdl
end


function compute_grid_vertex_position(origin_pos, quad_size, range_x, range_z, center_on_origin, time)
    -- Create offset if the model need to be center on the origin_pos
    local offset = hg.Vec3(0, 0, 0)
    if center_on_origin then
        offset = hg.Vec3(range_x * quad_size / 2, 0, range_z * quad_size / 2)
    end
    
    -- Create the positions list
    local positions = {}

    -- For all the vertex of the grid (add 1 to the range to get the last vertex on each axes, because the range corresponding to the number of quad to draw)
    for iz=0, range_z do
        for ix=0, range_x do
            -- Create it's position
            local pos_x = origin_pos.x + ix * quad_size - offset.x
            local pos_z = origin_pos.z + iz * quad_size - offset.z
            local pos_y = math.sin(pos_x) * math.sin(pos_z) * math.sin(time) -- use time to add some movement
            table.insert(positions, hg.Vec3(pos_x, pos_y, pos_z))
        end
    end    

    return positions
end


function compute_triangles_for_grid_vertex(range_x, range_z)
    -- Create the triangles list
    local triangles = {}

    -- Use to get the vertex id in the positions list according to the x and z position in the grid
    local function get_vertex_id(ix, iz)
        return iz * (range_x + 1) + ix + 1
    end

    -- For all the quad of the grid
    for iz=0, range_z - 1 do
        for ix=0, range_x - 1 do
            -- Get the 4 vertex that compose the quad
            local a = get_vertex_id(ix, iz) -- bottom left
            local b = get_vertex_id(ix, iz + 1) -- top left
            local c = get_vertex_id(ix + 1, iz + 1) -- top right
            local d = get_vertex_id(ix + 1, iz) -- bottom right

            table.insert(triangles, {d, c, b})
            table.insert(triangles, {b, a, d})
        end
    end

    return triangles
end


function compute_vertex_normals(triangles, positions)
    -- Create normal list with 1 vec3 per vertex position
    local normals = {}
    for i=1, #positions do
        normals[i] = hg.Vec3(0, 0, 0)
    end

    -- For all the triangle
    for it=1, #triangles do
        local v0, v1, v2 = triangles[it][1], triangles[it][2], triangles[it][3]

        -- Get the corresponding vertex position
        local p0, p1, p2 = positions[v0], positions[v1], positions[v2]
        local face_normal = hg.Cross(p0 - p1, p2 - p1)

        -- Guard against degenerate triangles (zero area)
        if not (face_normal.x == 0 and face_normal.y == 0 and face_normal.z == 0) then
            -- Accumulate (area-weighted) face normal to each vertex
            normals[v0] = normals[v0] + face_normal
            normals[v1] = normals[v1] + face_normal
            normals[v2] = normals[v2] + face_normal
        end

    end

    -- Normalize all vertex normals
    for i=1, #normals do
        if hg.Len(normals[i]) > 0 then
            normals[i] = hg.Normalize(normals[i]) 
        else
            normals[i] = hg.Vec3(0, 1, 0)
        end
    end

    return normals
end


-- Init render and resources
hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Model builder', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

hg.AddAssetsFolder('resources_compiled')

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

-- Create materials
prg_ref = hg.LoadPipelineProgramRefFromAssets('core/shader/pbr.hps', res, hg.GetForwardPipelineInfo())
plane_material = hg.CreateMaterial(prg_ref, 'uBaseOpacityColor', hg.Vec4(0.5, 0.5, 0.5), 'uOcclusionRoughnessMetalnessColor', hg.Vec4(1, 1, 0.25))

-- Setup scene
scene = hg.Scene()
hg.LoadSceneFromAssets("probe_scene/pbr.scn", scene, res, hg.GetForwardPipelineInfo())

cam = hg.CreateCamera(scene, hg.TransformationMat4(hg.Vec3(0, 6, -12), hg.Vec3(hg.DegreeToRadian(30), 0, 0)), 0.01, 1000)
scene:SetCurrentCamera(cam)

-- light = hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(0, 5, 0)), 0)
light_mtx = hg.TransformationMat4(hg.Vec3(8, 5, 0), hg.Vec3(hg.DegreeToRadian(50), hg.DegreeToRadian(-90), hg.DegreeToRadian(-90)))
inner_angle = hg.DegreeToRadian(30)
outer_angle = hg.DegreeToRadian(45)
light_color = hg.Color(1, 1, 1, 1)
light = hg.CreateSpotLight(scene, light_mtx, 0, inner_angle, outer_angle, light_color, 1, light_color, 1, 1, hg.LST_Map, 0.0)

sphere_node = scene:GetNode("sphere")
-- sphere_node.Disable()

-- Create plane model
vtx_layout = hg.VertexLayoutPosFloatNormUInt8() 
-- vtx_layout = hg.VertexLayoutPosFloatColorFloat()

grid_start_pos = hg.Vec3(0, 0, 0)
quad_size = 0.25
-- Create the grid model from a model builder
grid_mdl = create_grid_model_with_model_builder(vtx_layout, grid_start_pos, quad_size, 40, 40, true, 1)
-- Add the grid model to the PipelineResources and get the corresponding model ref 
grid_mdl_ref = res:AddModel('grid', grid_mdl)
-- Create a node from the model ref
grid_node = hg.CreateObject(scene, hg.TransformationMat4(grid_start_pos, hg.Vec3(0, 0, 0)), grid_mdl_ref, {plane_material})

-- Init input
keyboard = hg.Keyboard()

-- main loop
while not keyboard:Down(hg.K_Escape) and hg.IsWindowOpen(win) do
    dt = hg.TickClock()
    dts = hg.time_to_sec_f(dt)
    current_time = hg.time_to_sec_f(hg.GetClock())

    -- Rotate the grid node
    -- grid_rot = grid_node.GetTransform().GetRot()
    -- rot = grid_rot + hg.Vec3(0, 0.5 * dts, 0)
    -- grid_node.GetTransform().SetRot(rot)

    new_grid_mdl = create_grid_model_with_model_builder(vtx_layout, grid_start_pos, quad_size, 40, 40, true, current_time)
    res:UpdateModel(grid_mdl_ref, new_grid_mdl)

    scene:Update(dt)
    hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res)

    hg.Frame()
    hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)