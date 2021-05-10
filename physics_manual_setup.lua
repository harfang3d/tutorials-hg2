-- Manually setup node physics

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Node Physics Setup', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

-- create models
vtx_mdl = hg.VertexLayoutPosFloatNormUInt8()

cube_mdl = hg.CreateCubeModel(vtx_mdl, 1, 1, 1)
cube_ref = res:AddModel('cube', cube_mdl)

ground_mdl = hg.CreateCubeModel(vtx_mdl, 50, 0.01, 50)
ground_ref = res:AddModel('ground', ground_mdl)

-- create materials
prg_ref = hg.LoadPipelineProgramRefFromFile('resources_compiled/core/shader/default.hps', res, hg.GetForwardPipelineInfo())
mat = hg.CreateMaterial(prg_ref, 'uDiffuseColor', hg.Vec4(0.5, 0.5, 0.5, 1), 'uSpecularColor', hg.Vec4(0.0, 0.0, 0.0, 0.1))

-- setup scene
scene = hg.Scene()

cam = hg.CreateCamera(scene, hg.TransformationMat4(hg.Vec3(0, 1, -5), hg.Deg3(5, 0, 0)), 0.01, 1000)
scene:SetCurrentCamera(cam)

hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(6, 4, -6)), 0)
hg.CreatePhysicCube(scene, hg.Vec3(100, 0.02, 100), hg.TranslationMat4(hg.Vec3(0, -0.005, 0)), ground_ref, {mat}, 0)

clocks = hg.SceneClocks()

-- setup physic cube
cube_node = hg.CreateObject(scene, hg.TransformationMat4(hg.Vec3(0, 2.5, 0), hg.Vec3(0, 0, 0)), cube_ref, {mat})

rb = scene:CreateRigidBody()
rb:SetType(hg.RBT_Dynamic)

collision = scene:CreateCollision()
collision:SetType(hg.CT_Cube)
collision:SetSize(hg.Vec3(1, 1, 1))
collision:SetMass(1)

cube_node:SetRigidBody(rb)
cube_node:SetCollision(0, collision)

-- scene physics
physics = hg.SceneNewtonPhysics()
physics:SceneCreatePhysicsFromAssets(scene)

-- main loop
while not hg.ReadKeyboard():Key(hg.K_Escape) do
	dt = hg.TickClock()

	hg.SceneUpdateSystems(scene, clocks, dt, physics, hg.time_from_sec_f(1 / 60), 1)
	hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res)

	hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)

hg.WindowSystemShutdown()
hg.InputShutdown()
