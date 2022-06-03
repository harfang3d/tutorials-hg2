# Physics overrides node matrix.
# Manual changes of position/rotation/scale of a physics node won't affect its matrix.

import harfang as hg

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Physics Matrix Interaction', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

hg.ImGuiInit(10, hg.LoadProgramFromFile('resources_compiled/core/shader/imgui'), hg.LoadProgramFromFile('resources_compiled/core/shader/imgui_image'))

# create models
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

cube_mdl = hg.CreateCubeModel(vtx_layout, 0.5, 0.5, 0.5)
cube_ref = res.AddModel('cube', cube_mdl)
ground_mdl = hg.CreateCubeModel(vtx_layout, 50, 0.01, 50)
ground_ref = res.AddModel('ground', ground_mdl)

prg_ref = hg.LoadPipelineProgramRefFromFile('resources_compiled/core/shader/default.hps', res, hg.GetForwardPipelineInfo())

# create material
mat = hg.CreateMaterial(prg_ref, 'uDiffuseColor', hg.Vec4(0.5, 0.5, 0.5), 'uSpecularColor', hg.Vec4(1, 1, 1))

# setup scene
scene = hg.Scene()

cam = hg.CreateCamera(scene, hg.TransformationMat4(hg.Vec3(0, 1.5, -5), hg.Deg3(5, 0, 0)), 0.01, 1000)
scene.SetCurrentCamera(cam)

lgt = hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(3, 4, -6)), 0)

cube_node = hg.CreatePhysicCube(scene, hg.Vec3(1, 1, 1), hg.TranslationMat4(hg.Vec3(1.25, 2.5, 0)), cube_ref, [mat], 2)
ground_node = hg.CreatePhysicCube(scene, hg.Vec3(100, 0.02, 100), hg.TranslationMat4(hg.Vec3(0, -0.005, 0)), ground_ref, [mat], 0)

clocks = hg.SceneClocks()

# scene physics
physics = hg.SceneBullet3Physics()
physics.SceneCreatePhysicsFromAssets(scene)

# main loop
mouse, keyboard = hg.Mouse(), hg.Keyboard()

while not keyboard.Pressed(hg.K_Escape) and hg.IsWindowOpen(win):
	keyboard.Update()
	mouse.Update()

	# ImGui view
	dt = hg.TickClock()
	hg.ImGuiBeginFrame(res_x, res_y, dt, mouse.GetState(), keyboard.GetState())

	if hg.ImGuiBegin('Transform and Physics', True, hg.ImGuiWindowFlags_AlwaysAutoResize):
		hg.ImGuiTextWrapped('This tutorial demonstrates the interaction between the physics system and the Transform component of a node. The node position, rotation and scale are overriden by an active node rigid body.')

		hg.ImGuiSeparator()

		r, v = hg.ImGuiInputVec3('Transform.pos', cube_node.GetTransform().GetPos(), 2)
		if r:
			cube_node.GetTransform().SetPos(v)

		if hg.ImGuiButton('Press to reset position using Transform.SetPos'):
			cube_node.GetTransform().SetPos(hg.Vec3(1.25, 2.5, 0))

		hg.ImGuiSeparator()

		hg.ImGuiInputVec3('Transform.GetWorld().T', hg.GetT(cube_node.GetTransform().GetWorld()), 2, hg.ImGuiInputTextFlags_ReadOnly)

		if physics.NodeHasBody(cube_node):
			if hg.ImGuiButton('Press to destroy the cube node physics'):
				physics.NodeDestroyPhysics(cube_node)
				physics.GarbageCollect(scene)
		else:
			if hg.ImGuiButton('Create to create the cube node physics'):
				physics.NodeCreatePhysicsFromAssets(cube_node)

	hg.ImGuiEnd()

	hg.ImGuiEndFrame(255)

	# scene view
	hg.SceneUpdateSystems(scene, clocks, dt, physics, hg.time_from_sec_f(1 / 60), 1)
	hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), True, pipeline, res)

	hg.Frame()
	hg.UpdateWindow(win)

hg.RenderShutdown()
hg.DestroyWindow(win)

hg.WindowSystemShutdown()
hg.InputShutdown()
