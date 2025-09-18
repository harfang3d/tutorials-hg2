hg = require("harfang")

-- draw_line will draw lines in the viewport
function draw_line(pos_a, pos_b, line_color, vid, vtx_line_layout, line_shader)
	local vtx = hg.Vertices(vtx_line_layout, 2)
	vtx:Begin(0):SetPos(pos_a):SetColor0(line_color):End()
	vtx:Begin(1):SetPos(pos_b):SetColor0(line_color):End()
	hg.DrawLines(vid, vtx, line_shader)
end

-- add_debug_cross allows to create lines to draw a cross and append them to the "lines" structure. These lines will be drawn at the end of the main loop by draw_line()
function add_debug_cross(lines, pos, world, size)
    table.insert(lines, {pos_a = pos + hg.GetX(world) * size, pos_b = pos - hg.GetX(world) * size, color = hg.Color.Red})
    table.insert(lines, {pos_a = pos + hg.GetY(world) * size, pos_b = pos - hg.GetY(world) * size, color = hg.Color.Green})
    table.insert(lines, {pos_a = pos + hg.GetZ(world) * size, pos_b = pos - hg.GetZ(world) * size, color = hg.Color.Blue})
    return lines
end


hg.InputInit()
hg.WindowSystemInit()

-- Monitor resolution
res_x, res_y = 1920, 1080

-- get the actual monitor size from the window system
mon_list = hg.GetMonitors()
mon_list_size = mon_list:size()
if mon_list_size >= 1 then
    local _idx
    for _idx = 0, mon_list_size - 1 do
        local _mon_rect = hg.GetMonitorRect(mon_list:at(_idx))
        res_x = _mon_rect.ex - _mon_rect.sx
        res_y = _mon_rect.ey - _mon_rect.sy
        print("Found monitor size: " .. res_x .. " x " .. res_y)
        break
    end
end

-- Main screen and window init
mode_list = {hg.WV_Windowed, hg.WV_Fullscreen, hg.WV_Undecorated, hg.WV_FullscreenMonitor1, hg.WV_FullscreenMonitor2, hg.WV_FullscreenMonitor3}
win = hg.NewWindow('Mouse scene projection', res_x, res_y, 32, mode_list[3])
hg.RenderInit(win)
hg.RenderReset(res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

-- Create and configure the pipeline for rendering
pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

-- Assets folder
hg.AddAssetsFolder("resources_compiled")

-- Load scene
scene = hg.Scene()
hg.LoadSceneFromAssets("mouse_scene_projection/mouse_scene_projection.scn", scene, res, hg.GetForwardPipelineInfo())

-- Set camera
camera = scene:GetNode("Camera")
scene:SetCurrentCamera(camera)

-- Get sphere and rectangle nodes in the scene
rectangle_node = scene:GetNode("rectangle")

-- Create the shader to draw some 3D lines
vtx_line_layout = hg.VertexLayoutPosFloatColorUInt8()
shader_for_line = hg.LoadProgramFromAssets("shaders/pos_rgb")

-- Input init
keyboard = hg.Keyboard()
mouse = hg.Mouse()

-- Screen strategic position to emulate cursor pos
screen_pos_middle = hg.Vec3(res_x / 2, res_y / 2, 1.0)
screen_pos_up_left = hg.Vec3(0, res_y, 1.0)
screen_pos_down_right = hg.Vec3(res_x, 0, 1.0)

-- Main render loop
while not keyboard:Pressed(hg.K_Escape) and hg.IsWindowOpen(win) do
    dt = hg.TickClock()

    -- List containing the lines to draw in the viewport
    lines = {}

    -- Input updates and get mouse cursor pos
    keyboard:Update()
    mouse:Update()

    mouse_x, mouse_y = mouse:X(), mouse:Y()

    cursor_screen_pos = hg.Vec3(mouse_x, mouse_y, 1)
    resolution = hg.Vec2(res_x, res_y)

    -- Get the view state from the current camera
    view_state = hg.ComputePerspectiveViewState(camera:GetTransform():GetWorld(), camera:GetCamera():GetFov(), camera:GetCamera():GetZNear(), camera:GetCamera():GetZFar(), hg.ComputeAspectRatioX(res_x, res_y))
    
    -- Inverse proj and view mtx 
    inv_proj, flag = hg.Inverse(view_state.proj)
    flag, inv_view = hg.Inverse(view_state.view)

    -- get the translation of the invert view mtx
    ray_o = hg.GetT(inv_view)

    -- unproject the cursor screen pos to get the corresponding position in the 3D space
    flag, view_pos = hg.UnprojectFromScreenSpace(inv_proj, cursor_screen_pos, resolution)
    -- get the normalize position to display a debug cross at the screen
    view_pos_normalize = hg.Normalize(view_pos)
    view_pos_normalize = view_pos_normalize + ray_o

    -- unproject the middle screen pos to get the corresponding position in the 3D space
    flag, view_pos_middle = hg.UnprojectFromScreenSpace(inv_proj, screen_pos_middle, resolution)
    -- get the normalize position to display a debug cross at the screen
    view_pos_middle_normalize = hg.Normalize(view_pos_middle)
    view_pos_middle_normalize = view_pos_middle_normalize + ray_o

    -- unproject the top left screen pos to get the corresponding position in the 3D space
    flag, view_pos_up_left = hg.UnprojectFromScreenSpace(inv_proj, screen_pos_up_left, resolution)
    -- get the normalize position to display a debug cross at the screen
    view_pos_up_left_normalize = hg.Normalize(view_pos_up_left)
    view_pos_up_left_normalize = view_pos_up_left_normalize + ray_o

    -- unproject the bottom right screen pos to get the corresponding position in the 3D space
    flag, view_pos_down_right = hg.UnprojectFromScreenSpace(inv_proj, screen_pos_down_right, resolution)
    -- get the normalize position to display a debug cross at the screen
    view_pos_down_right_normalize = hg.Normalize(view_pos_down_right)
    view_pos_down_right_normalize = view_pos_down_right_normalize + ray_o

    -- add a debug cross at the screen cursor position
    lines = add_debug_cross(lines, view_pos_normalize, hg.TransformationMat4(view_pos_normalize, hg.Vec3(0, 0, 0)), 0.01)
    -- add debug crosses at the center screen, top left screen and bottom right screen position to debug
    lines = add_debug_cross(lines, view_pos_middle_normalize, hg.TransformationMat4(view_pos_middle_normalize, hg.Vec3(0, 0, 0)), 0.1)
    lines = add_debug_cross(lines, view_pos_up_left_normalize, hg.TransformationMat4(view_pos_up_left_normalize, hg.Vec3(0, 0, 0)), 0.1)
    lines = add_debug_cross(lines, view_pos_down_right_normalize, hg.TransformationMat4(view_pos_down_right_normalize, hg.Vec3(0, 0, 0)), 0.1)

    -- add a line to display the raycast from the cursor screen pos to the 3d position
    direction_line = {pos_a = hg.Vec3(0, 1.5, -5), pos_b = view_pos, color = hg.Color.Blue} 
    table.insert(lines, direction_line)

    -- Get the direction mtx to apply the rotation to a 3d rectangle 
    mat_look_at = hg.Mat4LookAt(rectangle_node:GetTransform():GetPos(), view_pos)
    rectangle_node:GetTransform():SetWorld(mat_look_at)

    -- Update scene
    scene:Update(dt)

    -- Render pass
    -- SubmitSceneToPipeline and get the view id to draw the debug lines
    view_id, pass_ids = hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res)

    -- Draw debug lines
    opaque_view_id = hg.GetSceneForwardPipelinePassViewId(pass_ids, hg.SFPP_Opaque)
    for i=1, #lines do
        draw_line(lines[i].pos_a, lines[i].pos_b, lines[i].color, opaque_view_id, vtx_line_layout, shader_for_line)
    end

    -- Update frame and window
    hg.Frame()
    hg.UpdateWindow(win)
end

-- Cleanup and shutdown operations
hg.RenderShutdown()
hg.DestroyWindow(win)
