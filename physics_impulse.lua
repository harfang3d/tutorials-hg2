-- Physics Impulse

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Physics Impulse', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

-- create models
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

cube_mdl = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
cube_ref = res:AddModel('cube', cube_mdl)

ground_mdl = hg.CreateCubeModel(vtx_layout, 50, 0.01, 50)
ground_ref = res:AddModel('ground', ground_mdl)

-- create material
prg_ref = hg.LoadPipelineProgramRefFromFile('resources_compiled/core/shader/default.hps', res, hg.GetForwardPipelineInfo())
mat = hg.CreateMaterial(prg_ref, 'uDiffuseColor', hg.Vec4(1, 1, 1), 'uSpecularColor', hg.Vec4(1, 1, 1))

-- setup scene
scene = hg.Scene()

cam = hg.CreateCamera(scene, hg.TransformationMat4(hg.Vec3(0, 1.5, -5), hg.Vec3(hg.Deg(10), 0, 0)), 0.01, 1000)
scene:SetCurrentCamera(cam)

lgt = hg.CreateLinearLight(scene, hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Vec3(hg.Deg(30), hg.Deg(59), 0)), hg.Color(1, 1, 1), hg.Color(1, 1, 1), 10, hg.LST_Map, 0.002, hg.Vec4(2, 4, 10, 16))

cube_node = hg.CreatePhysicCube(scene, hg.Vec3(1, 1, 1), hg.TranslationMat4(hg.Vec3(0, 2.5, 0)), cube_ref, {mat}, 2)
ground_node = hg.CreatePhysicCube(scene, hg.Vec3(100, 0.02, 100), hg.TranslationMat4(hg.Vec3(0, -0.005, 0)), ground_ref, {mat}, 0)

clocks = hg.SceneClocks()

-- scene physics
physics = hg.SceneBullet3Physics()
physics:SceneCreatePhysicsFromAssets(scene)
physics_step = hg.time_from_sec_f(1 / 60)

-- main loop
keyboard = hg.Keyboard()

while not keyboard:Down(hg.K_Escape) do
	keyboard:Update()

	dt = hg.TickClock()

	if keyboard:Pressed(hg.K_I) then
		physics:NodeWake(cube_node)
		world_pos = hg.GetT(cube_node:GetTransform():GetWorld())
		physics:NodeAddImpulse(cube_node, hg.Vec3(0, 8, 0), world_pos)
	end
	if keyboard:Pressed(hg.K_F) then
		physics:NodeWake(cube_node)
		world_pos = hg.GetT(cube_node:GetTransform():GetWorld())
		physics:NodeAddForce(cube_node, hg.Vec3(0, 8, 0), world_pos)
	end

	hg.SceneUpdateSystems(scene, clocks, dt, physics, physics_step, 3)
	hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res)

	hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)

hg.WindowSystemShutdown()
hg.InputShutdown()
