-- Mouse flight

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Mouse Flight', res_x, res_y, hg.RF_VSync | hg.RF_MSAA8X)

res = hg.PipelineResources()
pipeline = hg.CreateForwardPipeline()

keyboard = hg.Keyboard()
mouse = hg.Mouse()

-- access to compiled resources
hg.AddAssetsFolder('resources_compiled')

-- 2D drawing helpers
vtx_layout = hg.VertexLayoutPosFloatColorFloat()

draw2D_program = hg.LoadProgramFromAssets('shaders/pos_rgb')
draw2D_render_state = hg.ComputeRenderState(hg.BM_Alpha, hg.DT_Less, hg.FC_Disabled)


function draw_circle(view_id, center, radius, color)
	segment_count = 32
	step = 2 * math.pi / segment_count
	p0 = hg.Vec3(center.x + radius, center.y, 0)
	p1 = hg.Vec3(0, 0, 0)

	vtx = hg.Vertices(vtx_layout, segment_count * 2 + 2)

	for i = 0, segment_count do
		p1.x = radius * math.cos(i * step) + center.x
		p1.y = radius * math.sin(i * step) + center.y
		vtx:Begin(2 * i):SetPos(p0):SetColor0(color):End()
		vtx:Begin(2 * i + 1):SetPos(p1):SetColor0(color):End()
		p0.x, p0.y = p1.x, p1.y
	end
	hg.DrawLines(view_id, vtx, draw2D_program, draw2D_render_state)
end

-- gameplay settings
setting_camera_chase_offset = hg.Vec3(0, 0.2, 0)
setting_camera_chase_distance = 1

setting_plane_speed = 0.05
setting_plane_mouse_sensitivity = 0.5

-- setup game world
scene = hg.Scene()
hg.LoadSceneFromAssets('playground/playground.scn', scene, res, hg.GetForwardPipelineInfo())

plane_node = hg.CreateInstanceFromAssets(scene, hg.TranslationMat4(hg.Vec3(0, 4, 0)), 'paper_plane/paper_plane.scn', res, hg.GetForwardPipelineInfo())
camera_node = hg.CreateCamera(scene, hg.TranslationMat4(hg.Vec3(0, 4, -5)), 0.01, 1000)

scene:SetCurrentCamera(camera_node)


function update_plane(mouse_x_normd, mouse_y_normd)
	plane_transform = plane_node:GetTransform()

	plane_pos = plane_transform:GetPos()
	plane_pos = plane_pos + hg.Normalize(hg.GetZ(plane_transform:GetWorld())) * setting_plane_speed
	plane_pos.y = hg.Clamp(plane_pos.y, 0.1, 50)  -- floor/ceiling

	plane_rot = plane_transform:GetRot()

	next_plane_rot = hg.Vec3(plane_rot)  -- make a copy of the plane rotation
	next_plane_rot.x = hg.Clamp(next_plane_rot.x + mouse_y_normd * -0.03, -0.75, 0.75)
	next_plane_rot.y = next_plane_rot.y + mouse_x_normd * 0.03
	next_plane_rot.z = hg.Clamp(mouse_x_normd * -0.75, -1.2, 1.2)

	plane_rot = plane_rot + (next_plane_rot - plane_rot) * setting_plane_mouse_sensitivity

	plane_transform:SetPos(plane_pos)
	plane_transform:SetRot(plane_rot)
end

function update_chase_camera(target_pos)
	camera_transform = camera_node:GetTransform()
	camera_to_target = hg.Normalize(target_pos - camera_transform:GetPos())

	camera_transform:SetPos(target_pos - camera_to_target * setting_camera_chase_distance)  -- camera is 'distance' away from its target
	camera_transform:SetRot(hg.ToEuler(hg.Mat3LookAt(camera_to_target)))
end

-- game loop
while not keyboard:Down(hg.K_Escape) and hg.IsWindowOpen(win) do
	dt = hg.TickClock()  -- tick clock, retrieve elapsed clock since last call

	-- update mouse/keyboard devices
	keyboard:Update()
	mouse:Update()

	-- compute ratio corrected normalized mouse position
	mouse_x, mouse_y = mouse:X(), mouse:Y()

	aspect_ratio = hg.ComputeAspectRatioX(res_x, res_y)
	mouse_x_normd, mouse_y_normd = (mouse_x / res_x - 0.5) * aspect_ratio.x, (mouse_y / res_y - 0.5) * aspect_ratio.y

	-- update gameplay elements (plane & camera)
	update_plane(mouse_x_normd, mouse_y_normd)
	update_chase_camera(plane_node:GetTransform():GetWorld() * setting_camera_chase_offset)

	-- update scene and submit it to render pipeline
	scene:Update(dt)

	view_id = 0
	view_id, passes_id = hg.SubmitSceneToPipeline(view_id, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res)

	-- draw 2D GUI
	hg.SetView2D(view_id, 0, 0, res_x, res_y, -1, 1, hg.CF_Depth, hg.Color.Black, 1, 0, true)
	draw_circle(view_id, hg.Vec3(mouse_x, mouse_y, 0), 20, hg.Color.White)  -- display mouse cursor

	-- end of frame
	hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)

hg.WindowSystemShutdown()
hg.InputShutdown()
