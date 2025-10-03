-- Model builder with iso surface

--[[
Known issue : 2 axis are inverted in hg.IsoSurfaceSphere() function. "height"/"depth" parameters are inverted, and the axis "y"/"z" for the sphere position are also inverted ("y" correspond to the depth and "z" to the height).
Ex: 
iso_surface = hg.NewIsoSurface(width, height, depth) -- parameters not inverted
hg.IsoSurfaceSphere(iso_surface, width, depth, height, x, z, y, ...) -- parameters inverted
]]

hg = require("harfang")
-- from math import pi, sin, cos

-- Launch or not the movements on the iso surface spheres (cf imgui window and create_iso_surface_with_spheres_list())
anim_spheres = true

-- Create a simple sphere with default params
function create_iso_surface_sphere_default()
    local iso_sphere = {
        iso_sphere_pos = hg.Vec3(0, 0, 0),
        iso_sphere_radius = 1.0,
        iso_sphere_value = 1.0,
        iso_sphere_exponent = 1.0
    }
    return iso_sphere
end


-- Create a sphere with pos and radius params
function create_iso_surface_sphere(pos, radius)
    local iso_sphere = {    
        iso_sphere_pos = pos,
        iso_sphere_radius = radius,
        iso_sphere_value = 1.0,
        iso_sphere_exponent = 1.0
    }
    return iso_sphere
end


-- Create a list of sphere in circle / snake
function create_iso_surface_spheres_list_circle(origin)
    local sphere_list = {}
    local sphere_radius = 8.0
    local radius = 15
    -- create 10 spheres
    for i = 0, 9 do
        local angle = i / 10 * math.pi * 2
        local x = radius * math.cos(angle)
        local y = radius * math.sin(angle) 
        local z = i * radius / 6
        table.insert(sphere_list, create_iso_surface_sphere(hg.Vec3(x, y, z) + origin, sphere_radius))
    end
    return sphere_list
end


-- Create a model from an iso surface containing some spheres
function create_iso_surface_with_spheres_list(iso_surface_bounds, iso_level, iso_scale, iso_spheres_list, clock_sec)
    -- Create iso surface
    local iso_surface_width = iso_surface_bounds[1]
    local iso_surface_height = iso_surface_bounds[2]
    local iso_surface_depth = iso_surface_bounds[3]
    local iso_surface = hg.NewIsoSurface(iso_surface_width, iso_surface_height, iso_surface_depth)

    -- Add spheres in the iso surface
    -- Note : The functions related to iso_surface_sphere have inverted height/depth parameters (in the following code, z axis corresponding to the height)
    for i=1, #iso_spheres_list do
        local iso_sphere = iso_spheres_list[i]
        local iso_sphere_pos = iso_sphere.iso_sphere_pos
        local iso_sphere_radius = iso_sphere.iso_sphere_radius
        local iso_sphere_value = iso_sphere.iso_sphere_value
        local iso_sphere_exponent = iso_sphere.iso_sphere_exponent

        -- If the animation is activate, compute the position on the height through the time (clock_sec corresponding to the actual clock in second)
        local pos_z
        if anim_spheres then
            local clock = math.cos(clock_sec + i) / 4
            pos_z = ((iso_surface_height / 2 - iso_sphere_radius) * clock) + iso_surface_height / 4
        else
            -- If the animation is not activate, just keep the height position of the sphere
            pos_z = iso_sphere_pos.z
        end

        hg.IsoSurfaceSphere(
            iso_surface, 
            iso_surface_width, iso_surface_depth, iso_surface_height, 
            iso_sphere_pos.x, iso_sphere_pos.y, pos_z, 
            iso_sphere_radius, iso_sphere_value, iso_sphere_exponent
        )
    end

    -- Create a model builder
    local mdl_builder = hg.ModelBuilder()

    -- Create vertex layout
    local vtx_layout = hg.VertexLayoutPosFloatNormUInt8() 

    -- Convert iso surface to model
    local material_idx = 0
    _ = hg.IsoSurfaceToModel(
        mdl_builder, 
        iso_surface, 
        iso_surface_width, iso_surface_height, iso_surface_depth, 
        material_idx, 
        iso_level, 
        iso_scale.x, iso_scale.y, iso_scale.z
    )

    -- Get the model from the isosurface sphere
    local iso_surface_mdl = mdl_builder:MakeModel(vtx_layout)

    return iso_surface_mdl
end


-- Init render and resources
hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Model builder iso surface', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

-- Init pipeline
pipeline = hg.CreateForwardPipeline(1024, true)
res = hg.PipelineResources()

-- Add assets folder
hg.AddAssetsFolder('resources_compiled')

-- Init ImGui
imgui_prg = hg.LoadProgramFromAssets('core/shader/imgui')
imgui_img_prg = hg.LoadProgramFromAssets('core/shader/imgui_image')

hg.ImGuiInit(10, imgui_prg, imgui_img_prg)

-- Setup scene
scene = hg.Scene()
hg.LoadSceneFromAssets("probe_scene/scene_iso_surface.scn", scene, res, hg.GetForwardPipelineInfo())

-- Create camera
camera_rot_x = 0
camera_rot_y = 0
camera_distance = 45
camera_new_mtx = hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Vec3(camera_rot_x, camera_rot_y, 0)) * hg.TransformationMat4(hg.Vec3(0, 0, -camera_distance), hg.Vec3(0, 0, 0))
camera_node = hg.CreateCamera(scene, camera_new_mtx, 0.5, 800)
scene:SetCurrentCamera(camera_node)

-- Get the ground in the scene
ground_node = scene:GetNode("ground")

-- Create a spot light
-- light_mtx = hg.TransformationMat4(hg.Vec3(5, 35, 12.5), hg.Vec3(hg.DegreeToRadian(85), hg.DegreeToRadian(90), hg.DegreeToRadian(0)))
spot_light_mtx = hg.TransformationMat4(hg.Vec3(12.5, 35, 12.5), hg.Vec3(hg.DegreeToRadian(90), hg.DegreeToRadian(90), hg.DegreeToRadian(0)))
inner_angle = hg.DegreeToRadian(0.1)
outer_angle = hg.DegreeToRadian(35)
spot_light_color = hg.Color(1, 1, 1, 1)
spot_diffuse_intensity = 10
spot_specular_intensity = 10
spot_light = hg.CreateSpotLight(scene, spot_light_mtx, 0, 
    inner_angle, outer_angle, 
    spot_light_color, spot_diffuse_intensity, 
    spot_light_color, spot_specular_intensity, 
    1, hg.LST_Map, 0.0001
)

-- Init input
keyboard = hg.Keyboard()

-- Create material for the iso surface spheres
prg_ref = hg.LoadPipelineProgramRefFromAssets('core/shader/pbr.hps', res, hg.GetForwardPipelineInfo())
iso_surface_material = hg.CreateMaterial(prg_ref, 'uBaseOpacityColor', hg.Vec4(1.0, 0.75, 0.15), 'uOcclusionRoughnessMetalnessColor', hg.Vec4(1, 0.2, 0.5))

-- Generate iso surface mdl
iso_surface_bounds = {50, 50, 50} -- Size of the iso surface (width, height, depth)

-- Create a dict list containing the spheres parameters to provide to the iso surface (simple version, create spheres at the origin of the world. They can be configured in the imgui window)
-- iso_sphere_list = {}
-- for x=1, 2 do
--     local iso_sphere = create_iso_surface_sphere_default()
--     table.insert(iso_sphere_list, iso_sphere)
-- end

-- Create a dict list containing the spheres parameters to provide to the iso surface (advanced version, create spheres in circle)
iso_sphere_list = create_iso_surface_spheres_list_circle(hg.Vec3(iso_surface_bounds[1] / 2, iso_surface_bounds[1] / 2, iso_surface_bounds[1] / 6))

iso_level = 0.8
iso_scale = hg.Vec3(1, 1, 1)

-- Generate the model corresponding to the spheres defined in iso_sphere_list
iso_surface_mdl = create_iso_surface_with_spheres_list(iso_surface_bounds, iso_level, iso_scale, iso_sphere_list, 1)

-- Add the iso surface model to the PipelineResources and get the corresponding model ref 
iso_surface_mdl_ref = res:AddModel('isosurface', iso_surface_mdl)

-- Create a node from the model ref
iso_surface_node_scale = hg.Vec3(0.5, 0.5, 0.5)
iso_surface_node = hg.CreateObject(scene, hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Vec3(0, 0, 0), iso_surface_node_scale), iso_surface_mdl_ref, {iso_surface_material})

-- Init AAA pipeline
pipeline_aaa_config = hg.ForwardPipelineAAAConfig()
pipeline_aaa_config.exposure = 1.3
pipeline_aaa_config.gamma = 0.8
pipeline_aaa_config.bloom_threshold = 0.05
pipeline_aaa_config.bloom_intensity = 0.7
pipeline_aaa_config.sample_count = 1
pipeline_aaa = hg.CreateForwardPipelineAAAFromAssets("core", pipeline_aaa_config, hg.BR_Equal, hg.BR_Equal)

-- main loop
frame = 0

while not keyboard:Down(hg.K_Escape) and hg.IsWindowOpen(win) do
    dt = hg.TickClock()
    current_time = hg.time_to_sec_f(hg.GetClock())

    -- Get the iso surface size
    iso_surface_width = iso_surface_bounds[1]
    iso_surface_height = iso_surface_bounds[2]
    iso_surface_depth = iso_surface_bounds[3]

    -- Scale factor to auto resize the ground node and set the position of the camera
    scale_factor = iso_scale * iso_surface_node_scale

    -- Update camera pos input
    if hg.ReadKeyboard():Key(hg.K_W) then
        camera_rot_x = camera_rot_x + (1 * math.pi / 180)
    elseif hg.ReadKeyboard():Key(hg.K_S) then
        camera_rot_x = camera_rot_x - (1 * math.pi / 180)
    elseif hg.ReadKeyboard():Key(hg.K_A) then
        camera_rot_y = camera_rot_y + (1 * math.pi / 180)
    elseif hg.ReadKeyboard():Key(hg.K_D) then
        camera_rot_y = camera_rot_y - (1 * math.pi / 180)
    end
    -- Update camera node pos
    camera_new_mtx = hg.TransformationMat4(hg.Vec3(iso_surface_width / 2 * scale_factor.x, iso_surface_height / 3 * scale_factor.y, iso_surface_depth / 2 * scale_factor.z), hg.Vec3(camera_rot_x, camera_rot_y, 0)) * hg.TransformationMat4(hg.Vec3(0, 0, -camera_distance), hg.Vec3(0, 0, 0))
    camera_node:GetTransform():SetWorld(camera_new_mtx)

    -- Update ground node scale
    -- ground_node_debug_mtx = hg.TransformationMat4(hg.Vec3(iso_surface_width / 2 * scale_factor.x, 0, iso_surface_depth / 2 * scale_factor.z), hg.Vec3(0, 0, 0), hg.Vec3(iso_surface_width * scale_factor.x, 1, iso_surface_depth * scale_factor.z)) -- ground at the exact size of the iso surface width/depth
    ground_node_preview_mtx = hg.TransformationMat4(
        hg.Vec3(iso_surface_width / 2 * scale_factor.x, 0, iso_surface_depth / 2 * scale_factor.z), 
        hg.Vec3(0, 0, 0), 
        hg.Vec3(iso_surface_width * scale_factor.x * 2, 1, iso_surface_depth * scale_factor.z * 2)) -- ground larger than the iso surface width/depth for beautiful preview
    ground_node:GetTransform():SetWorld(ground_node_preview_mtx)
    ground_node_scale = ground_node:GetTransform():GetScale()

    -- Update the iso surface if some settings have been changed in the imgui window or by the animation 
    new_iso_surface_mdl = create_iso_surface_with_spheres_list(iso_surface_bounds, iso_level, iso_scale, iso_sphere_list, current_time)
    res:UpdateModel(iso_surface_mdl_ref, new_iso_surface_mdl)
    -- Update the iso surface node scale
    iso_surface_node:GetTransform():SetScale(iso_surface_node_scale)
    iso_surface_node_scale = iso_surface_node:GetTransform():GetScale()

    -- Update and draw the scene
    scene:Update(dt)
    -- vid, passid = hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res)
    vid, passid = hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res, pipeline_aaa, pipeline_aaa_config, frame)

    -- Draw the imgui window
    hg.ImGuiBeginFrame(res_x, res_y, hg.TickClock(), hg.ReadMouse(), hg.ReadKeyboard())

    if hg.ImGuiBegin('IsoSurface setting', true, hg.ImGuiWindowFlags_AlwaysAutoResize) then
        -- Settings about the scene
        changed, camera_distance = hg.ImGuiInputInt("Camera distance", camera_distance)
        changed, iso_surface_node_scale = hg.ImGuiInputVec3("Model node scale", iso_surface_node_scale)
        hg.ImGuiText(string.format("ground size : x = %.2f, y = %.2f, z = %.2f", ground_node_scale.x, ground_node_scale.y, ground_node_scale.z))

        hg.ImGuiNewLine()

        -- Checkbox to activate or not the animation on the iso surface spheres
        changed, anim_spheres = hg.ImGuiCheckbox("Anim spheres", anim_spheres)

        hg.ImGuiNewLine()

        -- Settings about the iso surface
        if hg.ImGuiCollapsingHeader("Iso surface") then
            hg.ImGuiText("iso_surface_bounds")
            changed, iso_surface_bounds[1] = hg.ImGuiInputInt("width", iso_surface_bounds[1])
            changed, iso_surface_bounds[2] = hg.ImGuiInputInt("height", iso_surface_bounds[2])
            changed, iso_surface_bounds[3] = hg.ImGuiInputInt("depth", iso_surface_bounds[3])

            changed, iso_level = hg.ImGuiInputFloat("iso_level", iso_level)

            changed, iso_scale = hg.ImGuiInputVec3("iso_scale", iso_scale)
        end

        hg.ImGuiNewLine()

        -- Settings about all the spheres in the iso surface
        for i=1, #iso_sphere_list do
            local iso_sphere = iso_sphere_list[i]
            local str_i = tostring(i)
            if hg.ImGuiCollapsingHeader("Iso sphere" .. str_i) then
                changed, iso_sphere.iso_sphere_pos = hg.ImGuiInputVec3("iso_sphere_pos_" .. str_i, iso_sphere.iso_sphere_pos)
                changed, iso_sphere.iso_sphere_radius = hg.ImGuiInputFloat("iso_sphere_radius_" .. str_i, iso_sphere.iso_sphere_radius)
                changed, iso_sphere.iso_sphere_value = hg.ImGuiInputFloat("iso_sphere_value_" .. str_i, iso_sphere.iso_sphere_value)
                changed, iso_sphere.iso_sphere_exponent = hg.ImGuiInputFloat("iso_sphere_exponent_" .. str_i, iso_sphere.iso_sphere_exponent)
            end
        end
    end

    hg.ImGuiEnd()
    hg.ImGuiEndFrame(vid)

    frame = hg.Frame()
    hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)