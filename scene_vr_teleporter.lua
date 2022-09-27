-- Display a scene in VR
require "math"
hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit("Harfang - OpenVR Scene", res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

hg.AddAssetsFolder("resources_compiled")

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

render_data = hg.SceneForwardPipelineRenderData()  -- this object is used by the low-level scene rendering API to share view-independent data with both eyes

-- OpenVR initialization
if not hg.OpenVRInit() then
	os.exit()
end

vr_left_fb = hg.OpenVRCreateEyeFrameBuffer(hg.OVRAA_MSAA4x)
vr_right_fb = hg.OpenVRCreateEyeFrameBuffer(hg.OVRAA_MSAA4x)

-- Create models
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

cube_mdl = hg.CreateCubeModel(vtx_layout, 0.1, 0.1, 0.1)
cube_ref = res:AddModel('cube', cube_mdl)
ground_mdl = hg.CreateCubeModel(vtx_layout, 50, 0.01, 50)
ground_ref = res:AddModel('ground', ground_mdl)

-- Teleporter spline
line_shader = hg.LoadProgramFromFile('resources_compiled/shaders/pos_rgb')
vtx_layout_spline = hg.VertexLayout()
vtx_layout_spline:Begin()
vtx_layout_spline:Add(hg.A_Position, 3, hg.AT_Float)
vtx_layout_spline:Add(hg.A_Color0, 3, hg.AT_Float)
vtx_layout_spline:End()

-- Load shader
prg_ref = hg.LoadPipelineProgramRefFromAssets('core/shader/pbr.hps', res, hg.GetForwardPipelineInfo())

-- Create materials
function create_material(ubc, orm)
	mat = hg.Material()
	hg.SetMaterialProgram(mat, prg_ref)
	hg.SetMaterialValue(mat, "uBaseOpacityColor", ubc)
	hg.SetMaterialValue(mat, "uOcclusionRoughnessMetalnessColor", orm)
	return mat
end

-- ----------------- Teleporter -----------------
function UpdateTeleporterPos(controller, actor_pos, head_pos, teleporter_pos, playground)
	if controller:IsConnected() then
		world = controller:World()
		T = hg.GetT(world)
		Z = hg.GetZ(world)

		min_Vy = -0.1

		if Z.y < min_Vy then
			if not playgroundthen then
				teleport_I = T + Z * (playground[1].y - T.y) / Z.y
				teleport_I = hg.Clamp(teleport_I, playground[1], playground[2])
			else
				teleport_I = T + Z * (0 - T.y) / Z.y
			end

			-- detect trigger
			if controller.Pressed(hg.VRCB_Axis1) then
				new_actor_pos = teleport_I + (actor_pos - head_pos)
				actor_diff = new_actor_pos - actor_pos
				teleporter_pos = teleport_I + actor_diff
				actor_pos = new_actor_pos
			else
				teleporter_pos = teleport_I
			end
		end
	end
	return actor_pos, teleporter_pos
end


function UpdateTeleporterNode(controller, actor_pos, head_pos, teleporter_node, playground)
	if controller:IsConnected() then
		world = controller:World()
		T = hg.GetT(world)
		Z = hg.GetZ(world)

		min_Vy = -0.1

		if Z.y < min_Vy then
			if  playground then
				teleport_I = T + Z * (playground[1].y - T.y) / Z.y
				teleport_I = hg.Clamp(teleport_I, playground[1], playground[2])
			else
				teleport_I = T + Z * (0 - T.y) / Z.y
			end

			-- detect trigger
			if controller:Pressed(hg.VRCB_Axis1) then
				new_actor_pos = teleport_I + (actor_pos - head_pos)
				actor_diff = new_actor_pos - actor_pos
				teleporter_pos = teleport_I + actor_diff
				actor_pos = new_actor_pos
			else
				teleporter_pos = teleport_I
			end

			teleporter_node:GetTransform():SetPos(teleporter_pos)
			teleporter_node:Enable()
		else
			teleporter_node:Disable()
		end
	end
	return actor_pos
end

function GetHeadPos(vr_state, actor_pos)
	headT = hg.GetT(vr_state.head)
	return hg.Vec3(headT.x, actor_pos.y, headT.z)
end


--  --------------------- VR CONTROLLER ---------------------
function InitVRControllers(vr_controller, vr_controller_idx)

	if not vr_controller or not vr_controller_idx then
		vr_controller = {hg.VRController(), hg.VRController()}
		vr_controller_idx = {1, 1}
	end

	name_template = "openvr_controller_"

	left_hand_connected = vr_controller[1]:IsConnected()
	right_hand_connected = vr_controller[2]:IsConnected()

	if not left_hand_connected then
		for i = 1, 16 do
			if not right_hand_connected or vr_controller_idx[2] ~= i then
				if hg.ReadVRController(name_template .. tostring(i)):IsConnected() then
					vr_controller[1] = hg.VRController(name_template .. tostring(i))
					vr_controller_idx[1] = i
					left_hand_connected = True
					print("Using controller %1 for left hand" .. tostring(i))
					break
				end
			end
		end
	end

	if not right_hand_connected then
		for i = 1, 16 do
			if not left_hand_connected or vr_controller_idx[1] ~= i then
				if hg.ReadVRController(name_template .. tostring(i)):IsConnected() then
					vr_controller[2] = hg.VRController(name_template .. tostring(i))
					vr_controller_idx[2] = i
					right_hand_connected = True
					print("Using controller %1 for right hand" .. tostring(i))
					break
				end
			end
		end
	end
	return vr_controller, vr_controller_idx
end


function UpdateVRControllers(vr_controller, vr_controller_idx)
	vr_controller, vr_controller_idx = InitVRControllers(vr_controller, vr_controller_idx)

	for i = 1, 2 do
		vr_controller[i]:Update()
	end

	return vr_controller, vr_controller_idx
end


function DrawTeleporterSpline(controller_mtx, ground_mtx, vid)
	dir_teleporter = hg.GetZ(controller_mtx)
	pos_start = hg.GetT(controller_mtx)
	cos_angle = hg.Dot(dir_teleporter, hg.Normalize(hg.Vec3(dir_teleporter.x, 0, dir_teleporter.z)))
	cos_angle = math.min(1.0, math.max(cos_angle, -1))
	angle = math.acos(cos_angle)
	strength_force = (((math.sin(angle) + 1) / 2)^2)*2

	DrawSpline(pos_start, pos_start + dir_teleporter * strength_force, ground_mtx + hg.Vec3(0, -strength_force, 0), ground_mtx, vid)
end


function DrawSpline(p1, p2, p3, p4, vid)
	step = 10
	prev_value = {p1.x, p1.y, p1.z}
	vtx = hg.Vertices(vtx_layout_spline, step * 2)
	for i = 0, step+1 do
		step_range = (1 / step) * i
		val_x = hg.CubicInterpolate(p2.x, p1.x, p4.x, p3.x, step_range)
		val_y = hg.CubicInterpolate(math.min(p2.y, 0), p1.y, p4.y, p3.y, step_range)
		val_z = hg.CubicInterpolate(p2.z, p1.z, p4.z, p3.z, step_range)
		val = {val_x, val_y, val_z}
		vtx:Begin(2 * i):SetPos(hg.Vec3(prev_value[1], prev_value[2], prev_value[3])):SetColor0(hg.Color.Blue):End()
		vtx:Begin(2 * i + 1):SetPos(hg.Vec3(val[1], val[2], val[3])):SetColor0(hg.Color.Blue):End()
		prev_value = val
	end
	hg.DrawLines(vid, vtx, line_shader)  -- submit all lines in a single call
end


function main()
	-- Create scene
	scene = hg.Scene()
	scene.canvas.color = hg.Color(255 / 255, 255 / 255, 217 / 255, 1)
	scene.environment.ambient = hg.Color(15 / 255, 12 / 255, 9 / 255, 1)

	lgt = hg.CreateSpotLight(scene, hg.TransformationMat4(hg.Vec3(-8, 4, -5), hg.Vec3(hg.Deg(19), hg.Deg(59), 0)), 0, hg.Deg(5), hg.Deg(30), hg.Color.White, hg.Color.White, 10, hg.LST_Map, 0.0001)
	back_lgt = hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(2.4, 1, 0.5)), 10, hg.Color(94 / 255, 255 / 255, 228 / 255, 1), hg.Color(94 / 255, 1, 228 / 255, 1), 0)

	mat_cube = create_material(hg.Vec4(255 / 255, 230 / 255, 20 / 255, 1), hg.Vec4(1, 0.658, 0., 1))
	hg.CreateObject(scene, hg.TransformationMat4(hg.Vec3(0, 0.5, 0), hg.Vec3(0, hg.Deg(70), 0)), cube_ref, {mat_cube})

	mat_ground = create_material(hg.Vec4(255 / 255, 120 / 255, 147 / 255, 1), hg.Vec4(1, 1, 0.1, 1))
	hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 0, 0)), ground_ref, {mat_ground})

	-- Setup 2D rendering to display eyes textures
	quad_layout = hg.VertexLayout()
	quad_layout:Begin():Add(hg.A_Position, 3, hg.AT_Float):Add(hg.A_TexCoord0, 3, hg.AT_Float):End()

	quad_model = hg.CreatePlaneModel(quad_layout, 1, 1, 1, 1)
	quad_render_state = hg.ComputeRenderState(hg.BM_Alpha, hg.DT_Disabled, hg.FC_Disabled)

	eye_t_size = res_x / 2.5
	eye_t_x = (res_x - 2 * eye_t_size) / 6 + eye_t_size / 2
	quad_matrix = hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Vec3(hg.Deg(90), hg.Deg(0), hg.Deg(0)), hg.Vec3(eye_t_size, 1, eye_t_size))

	tex0_program = hg.LoadProgramFromAssets("shaders/sprite")

	quad_uniform_set_value_list = hg.UniformSetValueList()
	quad_uniform_set_value_list:clear()
	quad_uniform_set_value_list:push_back(hg.MakeUniformSetValue("color", hg.Vec4(1, 1, 1, 1)))

	quad_uniform_set_texture_list = hg.UniformSetTextureList()

	-- Teleporter
	teleporter_node, boola = hg.CreateInstanceFromAssets(scene, hg.Mat4.Identity, 'teleporter/teleporter.scn', res, hg.GetForwardPipelineInfo())

	hand_left, boolb = hg.CreateInstanceFromAssets(scene, hg.Mat4.Identity, 'vr_controller/vr_controller.scn', res, hg.GetForwardPipelineInfo())
	hand_left:SetName("hand_left")
	hand_right, boolc = hg.CreateInstanceFromAssets(scene, hg.Mat4.Identity, 'vr_controller/vr_controller.scn', res, hg.GetForwardPipelineInfo())
	hand_right:SetName("hand_right")

	actor_pos = hg.Vec3(-1.3, 0, -2)
	head_pos = hg.Vec3()
	playground = {hg.Vec3(-3, 0, -3), hg.Vec3(3, 0, 3)}

	vr_controller, vr_controller_idx = InitVRControllers()

	-- Main loop
	while not hg.ReadKeyboard():Key(hg.K_Escape) do
		dt = hg.TickClock()

		-- Teleporter and VR controller update
		vr_controller, vr_controller_idx = UpdateVRControllers(vr_controller, vr_controller_idx)
		hand_left:GetTransform():SetWorld(vr_controller[1]:World())
		hand_right:GetTransform():SetWorld(vr_controller[2]:World())

		controller = vr_controller[1]
		actor_pos = UpdateTeleporterNode(controller, actor_pos, head_pos, teleporter_node)

		scene:Update(dt)

		vr_state = hg.OpenVRGetState(hg.TranslationMat4(actor_pos), 0.1, 200)

		head_pos = GetHeadPos(vr_state, actor_pos)

		left, right = hg.OpenVRStateToViewState(vr_state)

		vid = 0  -- keep track of the next free view id
		passId = hg.SceneForwardPipelinePassViewId()

		-- Prepare view-independent render data once
		vid, passId = hg.PrepareSceneForwardPipelineCommonRenderData(vid, scene, render_data, pipeline, res, passId)
		vr_eye_rect = hg.IntRect(0, 0, vr_state.width, vr_state.height)

		-- Prepare the left eye render data then draw to its framebuffer
		vid, passId = hg.PrepareSceneForwardPipelineViewDependentRenderData(vid, left, scene, render_data, pipeline, res, passId)
		vid, passId = hg.SubmitSceneToForwardPipeline(vid, scene, vr_eye_rect, left, pipeline, render_data, res, vr_left_fb:GetHandle())

		-- Prepare the right eye render data then draw to its framebuffer
		vid, passId = hg.PrepareSceneForwardPipelineViewDependentRenderData(vid, right, scene, render_data, pipeline, res, passId)
		vid, passId = hg.SubmitSceneToForwardPipeline(vid, scene, vr_eye_rect, right, pipeline, render_data, res, vr_right_fb:GetHandle())

		-- Display teleporter spline:
		hg.SetViewFrameBuffer(vid, vr_left_fb:GetHandle())
		hg.SetViewRect(vid, 0, 0, vr_state.width, vr_state.height)
		hg.SetViewClear(vid, 0, 0, 1.0, 0)
		hg.SetViewTransform(vid, left.view, left.proj)
		DrawTeleporterSpline(controller:World(), teleporter_node:GetTransform():GetPos(), vid)

		vid = vid + 1

		hg.SetViewFrameBuffer(vid, vr_right_fb:GetHandle())
		hg.SetViewRect(vid, 0, 0, vr_state.width, vr_state.height)
		hg.SetViewClear(vid, 0, 0, 1.0, 0)
		hg.SetViewTransform(vid, right.view, right.proj)
		DrawTeleporterSpline(controller:World(), teleporter_node:GetTransform():GetPos(), vid)

		vid = vid + 1

		-- Display the VR eyes texture to the backbuffer
		hg.SetViewRect(vid, 0, 0, res_x, res_y)
		vs = hg.ComputeOrthographicViewState(hg.TranslationMat4(hg.Vec3(0, 0, 0)), res_y, 0.1, 100, hg.ComputeAspectRatioX(res_x, res_y))
		hg.SetViewTransform(vid, vs.view, vs.proj)

		quad_uniform_set_texture_list:clear()
		quad_uniform_set_texture_list:push_back(hg.MakeUniformSetTexture("s_tex", hg.OpenVRGetColorTexture(vr_left_fb), 0))
		hg.SetT(quad_matrix, hg.Vec3(eye_t_x, 0, 1))
		hg.DrawModel(vid, quad_model, tex0_program, quad_uniform_set_value_list, quad_uniform_set_texture_list, quad_matrix, quad_render_state)

		quad_uniform_set_texture_list:clear()
		quad_uniform_set_texture_list:push_back(hg.MakeUniformSetTexture("s_tex", hg.OpenVRGetColorTexture(vr_right_fb), 0))
		hg.SetT(quad_matrix, hg.Vec3(-eye_t_x, 0, 1))
		hg.DrawModel(vid, quad_model, tex0_program, quad_uniform_set_value_list, quad_uniform_set_texture_list, quad_matrix, quad_render_state)

		hg.Frame()
		hg.OpenVRSubmitFrame(vr_left_fb, vr_right_fb)

		hg.UpdateWindow(win)
	end
	hg.DestroyForwardPipeline(pipeline)
	hg.RenderShutdown()
	hg.DestroyWindow(win)
end


main()