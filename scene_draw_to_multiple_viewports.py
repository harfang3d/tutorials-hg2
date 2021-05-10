# Draw to multiple viewports

import harfang as hg

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Scene Draw to Multiple Viewports', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

hg.AddAssetsFolder("resources_compiled")

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

# create models
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

cube_mdl = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
cube_ref = res.AddModel('cube', cube_mdl)
ground_mdl = hg.CreateCubeModel(vtx_layout, 100, 0.01, 100)
ground_ref = res.AddModel('ground', ground_mdl)

# create materials
shader = hg.LoadPipelineProgramRefFromAssets('core/shader/default.hps', res, hg.GetForwardPipelineInfo())

mat_yellow_cube = hg.CreateMaterial(shader, 'uDiffuseColor', hg.Vec4I(255, 220, 64), 'uSpecularColor', hg.Vec4I(255, 220, 64))
mat_red_cube = hg.CreateMaterial(shader, 'uDiffuseColor', hg.Vec4I(255, 0, 0), 'uSpecularColor', hg.Vec4I(255, 0, 0))
mat_ground = hg.CreateMaterial(shader, 'uDiffuseColor', hg.Vec4I(128, 128, 128), 'uSpecularColor', hg.Vec4I(128, 128, 128))

# setup scene (note that we do not create any camera)
scene = hg.Scene()

hg.CreateSpotLight(scene, hg.TransformationMat4(hg.Vec3(-8, 4, -5), hg.Deg3(19, 59, 0)), 0, hg.Deg(5), hg.Deg(30), hg.Color.White, hg.Color.White, 10, hg.LST_Map, 0.00005)
hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(3, 1, 2.5)), 5, hg.ColorI(128, 192, 255), hg.Color.Black, 0)

yellow_cube = hg.CreateObject(scene, hg.TransformationMat4(hg.Vec3(1, 0.5, 0), hg.Vec3(0, hg.Deg(0), 0)), cube_ref, [mat_yellow_cube])
hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(-1, 0.5, 0)), cube_ref, [mat_red_cube])
hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 0, 0)), ground_ref, [mat_ground])

# define viewports
viewports = [
	{'rect': hg.IntRect(0, 0, res_x // 2, res_y // 2), 'cam_pos': hg.Vec3(-4.015, 2.368, -3.484), 'cam_rot': hg.Vec3(0.35, 0.87, 0.0)},
	{'rect': hg.IntRect(res_x // 2, 0, res_x, res_y // 2), 'cam_pos': hg.Vec3(-4.143, 2.976, 4.127), 'cam_rot': hg.Vec3(0.423, 2.365, 0.0)},
	{'rect': hg.IntRect(0, res_y // 2, res_x // 2, res_y), 'cam_pos': hg.Vec3(4.020, 2.374, 3.469), 'cam_rot': hg.Vec3(0.353, 4.016, 0.0)},
	{'rect': hg.IntRect(res_x // 2, res_y // 2, res_x, res_y), 'cam_pos': hg.Vec3(3.469, 2.374, -4.020), 'cam_rot': hg.Vec3(0.353, -0.695, 0.0)}
]

# main loop
while not hg.ReadKeyboard().Key(hg.K_Escape):
	dt = hg.TickClock()

	# animate yellow cube & update scene once for all viewports
	rot = yellow_cube.GetTransform().GetRot()
	rot.y = rot.y + hg.time_to_sec_f(dt)
	yellow_cube.GetTransform().SetRot(rot)

	scene.Update(dt)

	# prepare view-independent render data (eg. spot shadow maps)
	render_data = hg.SceneForwardPipelineRenderData()

	views = hg.SceneForwardPipelinePassViewId()
	vid = 0
	vid, pass_ids = hg.PrepareSceneForwardPipelineCommonRenderData(vid, scene, render_data, pipeline, res, views)

	for viewport in viewports:
		# compute viewport specific view state
		view_state = hg.ComputePerspectiveViewState(hg.TransformationMat4(viewport['cam_pos'], viewport['cam_rot']), hg.Deg(45), 0.01, 1000, hg.ComputeAspectRatioX(res_x, res_y))
		# prepare view-dependent render data & submit draw
		vid, pass_ids = hg.PrepareSceneForwardPipelineViewDependentRenderData(vid, view_state, scene, render_data, pipeline, res, views)
		vid, pass_ids = hg.SubmitSceneToForwardPipeline(vid, scene, viewport['rect'], view_state, pipeline, render_data, res)

	hg.Frame()
	hg.UpdateWindow(win)

hg.DestroyForwardPipeline(pipeline)
hg.RenderShutdown()
hg.DestroyWindow(win)
