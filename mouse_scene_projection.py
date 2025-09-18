import harfang as hg

# draw_line will draw lines in the viewport
def draw_line(pos_a, pos_b, line_color, vid, vtx_line_layout, line_shader):
	vtx = hg.Vertices(vtx_line_layout, 2)
	vtx.Begin(0).SetPos(pos_a).SetColor0(line_color).End()
	vtx.Begin(1).SetPos(pos_b).SetColor0(line_color).End()
	hg.DrawLines(vid, vtx, line_shader)


# add_debug_cross allow to add lines to draw a cross and append them to lines structure. Those lines will be drawn at the end of the main loop by draw_line()
def add_debug_cross(lines, pos, world, size):
    lines.append([pos + hg.GetX(world) * size, pos - hg.GetX(world) * size, hg.Color.Red]) 
    lines.append([pos + hg.GetY(world) * size, pos - hg.GetY(world) * size, hg.Color.Green])
    lines.append([pos + hg.GetZ(world) * size, pos - hg.GetZ(world) * size, hg.Color.Blue])
    return lines


hg.InputInit()
hg.WindowSystemInit()

# Monitor resolution
res_x, res_y = 1920, 1080

# get the actual monitor size from the window system
mon_list = hg.GetMonitors()
mon_list_size = len(mon_list)
if mon_list_size >= 1:
    for _idx in range(mon_list_size):
        _mon_rect = hg.GetMonitorRect(mon_list[_idx])
        res_x = _mon_rect.ex - _mon_rect.sx
        res_y = _mon_rect.ey - _mon_rect.sy
        print("Found monitor size: {} x {}".format(res_x, res_y))
        break

# Main screen and window init
mode_list = [hg.WV_Windowed, hg.WV_Fullscreen, hg.WV_Undecorated, hg.WV_FullscreenMonitor1, hg.WV_FullscreenMonitor2, hg.WV_FullscreenMonitor3]
win = hg.NewWindow('Poc mouse raycast', res_x, res_y, 32, mode_list[2])
hg.RenderInit(win)
hg.RenderReset(res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

# Create and configure the pipeline for rendering
pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

# Assets folder    
hg.AddAssetsFolder("resources_compiled")

# Load scene
scene = hg.Scene()
hg.LoadSceneFromAssets("mouse_scene_projection/mouse_scene_projection.scn", scene, res, hg.GetForwardPipelineInfo())

# Set camera
camera = scene.GetNode("Camera")
scene.SetCurrentCamera(camera)

# Get sphere and rectangle nodes in the scene
rectangle_node = scene.GetNode("rectangle")

# Create the shader to draw some 3D lines
vtx_line_layout = hg.VertexLayoutPosFloatColorUInt8()
shader_for_line = hg.LoadProgramFromAssets("shaders/pos_rgb")

# Input init
keyboard = hg.Keyboard()
mouse = hg.Mouse()

# Screen strategic position to emulate cursor pos
screen_pos_middle = hg.Vec3(res_x / 2, res_y / 2, 1.0)
screen_pos_up_left = hg.Vec3(0, res_y, 1.0)
screen_pos_down_right = hg.Vec3(res_x, 0, 1.0)

# Main render loop
while not keyboard.Pressed(hg.K_Escape) and hg.IsWindowOpen(win) :
    dt = hg.TickClock()

    # Lines list containing the lines to draw in the viewport
    lines = []

    # Input updates and get mouse cursor pos
    keyboard.Update()
    mouse.Update()

    mouse_x, mouse_y = mouse.X(), mouse.Y()

    # Raycast the mouse cursor in the 3d space
    cursor_screen_pos = hg.Vec3(mouse_x, mouse_y, 1)
    resolution = hg.Vec2(res_x, res_y)

    # Get the view state from the current camera
    view_state = hg.ComputePerspectiveViewState(camera.GetTransform().GetWorld(), camera.GetCamera().GetFov(), camera.GetCamera().GetZNear(), camera.GetCamera().GetZFar(), hg.ComputeAspectRatioX(res_x, res_y))
    
    # Inverse proj and view mtx 
    inv_proj, flag = hg.Inverse(view_state.proj)
    flag, inv_view = hg.Inverse(view_state.view)

    # get the translation of the invert view mtx
    ray_o = hg.GetT(inv_view)

    # unproject the cursor screen pos to get the corresponding position in the 3D space
    flag, view_pos = hg.UnprojectFromScreenSpace(inv_proj, cursor_screen_pos, resolution)
    # get the normalize position to display a debug cross at the screen
    view_pos_normalize = hg.Normalize(view_pos)
    view_pos_normalize = view_pos_normalize + ray_o

    # unproject the middle screen pos to get the corresponding position in the 3D space
    flag, view_pos_middle = hg.UnprojectFromScreenSpace(inv_proj, screen_pos_middle, resolution)
    # get the normalize position to display a debug cross at the screen
    view_pos_middle_normalize = hg.Normalize(view_pos_middle)
    view_pos_middle_normalize = view_pos_middle_normalize + ray_o

    # unproject the top left screen pos to get the corresponding position in the 3D space
    flag, view_pos_up_left = hg.UnprojectFromScreenSpace(inv_proj, screen_pos_up_left, resolution)
    # get the normalize position to display a debug cross at the screen
    view_pos_up_left_normalize = hg.Normalize(view_pos_up_left)
    view_pos_up_left_normalize = view_pos_up_left_normalize + ray_o

    # unproject the bottom right screen pos to get the corresponding position in the 3D space
    flag, view_pos_down_right = hg.UnprojectFromScreenSpace(inv_proj, screen_pos_down_right, resolution)
    # get the normalize position to display a debug cross at the screen
    view_pos_down_right_normalize = hg.Normalize(view_pos_down_right)
    view_pos_down_right_normalize = view_pos_down_right_normalize + ray_o

    # add a debug cross at the screen cursor position
    lines = add_debug_cross(lines, view_pos_normalize, hg.TransformationMat4(view_pos_normalize, hg.Vec3(0, 0, 0)), 0.01)
    # add debug crosses at the center screen, top left screen and bottom roght screen position to debug
    lines = add_debug_cross(lines, view_pos_middle_normalize, hg.TransformationMat4(view_pos_middle_normalize, hg.Vec3(0, 0, 0)), 0.1)
    lines = add_debug_cross(lines, view_pos_up_left_normalize, hg.TransformationMat4(view_pos_up_left_normalize, hg.Vec3(0, 0, 0)), 0.1)
    lines = add_debug_cross(lines, view_pos_down_right_normalize, hg.TransformationMat4(view_pos_down_right_normalize, hg.Vec3(0, 0, 0)), 0.1)

    # add a line to display the raycast from the cursor screen pos to the 3d position
    direction_line = [hg.Vec3(0, 1.5, -5), view_pos, hg.Color.Blue]
    lines.append(direction_line)

    # Get the direction mtx to apply the rotation to a 3d rectangle 
    mat_look_at = hg.Mat4LookAt(rectangle_node.GetTransform().GetPos(), view_pos)
    rectangle_node.GetTransform().SetWorld(mat_look_at)

    # Update scene
    scene.Update(dt)

    # Render pass
    # SubmitSceneToPipeline and get the view id to draw the debug lines
    view_id, pass_ids = hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), True, pipeline, res)

    # Draw debug lines
    opaque_view_id = hg.GetSceneForwardPipelinePassViewId(pass_ids, hg.SFPP_Opaque)
    for line in lines:
        draw_line(line[0], line[1], line[2], opaque_view_id, vtx_line_layout, shader_for_line)

    # Update frame and window
    hg.Frame()
    hg.UpdateWindow(win)

# Cleanup and shutdown operations
hg.RenderShutdown()
hg.DestroyWindow(win)

