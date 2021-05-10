-- Dynamically assign lights to the fixed pipeline slots by adjusting their priority

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Light priority relative to a specific world position', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

hg.AddAssetsFolder("resources_compiled")

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

-- create models
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

light_mdl = hg.CreateSphereModel(vtx_layout, 0.05, 8, 16)
light_ref = res:AddModel('light', light_mdl)
orb_mdl = hg.CreateSphereModel(vtx_layout, 1, 16, 32)
orb_ref = res:AddModel('orb', orb_mdl)
ground_mdl = hg.CreateCubeModel(vtx_layout, 100, 0.01, 100)
ground_ref = res:AddModel('ground', ground_mdl)

-- create materials
shader = hg.LoadPipelineProgramRefFromAssets('core/shader/default.hps', res, hg.GetForwardPipelineInfo())

mat_light = hg.CreateMaterial(shader, 'uDiffuseColor', hg.Vec4(0, 0, 0), 'uSpecularColor', hg.Vec4(0, 0, 0))
hg.SetMaterialValue(mat_light, "uSelfColor", hg.Vec4(1, 0.9, 0.75))
mat_orb = hg.CreateMaterial(shader, 'uDiffuseColor', hg.Vec4(1, 1, 1), 'uSpecularColor', hg.Vec4(1, 1, 1))
hg.SetMaterialValue(mat_orb, "uSelfColor", hg.Vec4(0, 0, 0))
mat_ground = hg.CreateMaterial(shader, 'uDiffuseColor', hg.Vec4(1, 1, 1), 'uSpecularColor', hg.Vec4(1, 1, 1))
hg.SetMaterialValue(mat_ground, "uSelfColor", hg.Vec4(0, 0, 0))

-- setup scene
scene = hg.Scene()

cam = hg.CreateCamera(scene, hg.Mat4LookAt(hg.Vec3(5, 4, -7), hg.Vec3(0, 1.5, 0)), 0.01, 1000)
scene:SetCurrentCamera(cam)

orb_node = hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 1, 0)), orb_ref, {mat_orb})
hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 0, 0)), ground_ref, {mat_ground})

-- create an array of dynamic lights
light_obj = scene:CreateObject(light_ref, {mat_light})  -- sphere model to visualize lights

light_nodes = {}
for i=0, 15 do
	node = hg.CreatePointLight(scene, hg.Mat4.Identity, 1.5, hg.Color(1, 0.85, 0.25, 1), hg.Color(1, 0.9, 0.5, 1))
	node:SetObject(light_obj)
	table.insert(light_nodes, node)
end

-- main loop
angle = 0

while not hg.ReadKeyboard():Key(hg.K_Escape) do
	dt = hg.TickClock()

	-- animate lights
	angle = angle + hg.time_to_sec_f(dt)

	for i, node in ipairs(light_nodes) do
		a = angle + i * hg.Deg(15)
		node:GetTransform():SetPos(hg.Vec3(math.cos(a * -0.6) * math.sin(a) * 5, math.cos(a * 1.25) * 2 + 2.15, math.sin(a * 0.5) * math.cos(-a * 0.8) * 5))
	end
	-- update light priorities according to their distance to the orb
	for i, node in ipairs(light_nodes) do
		priority = hg.Dist(orb_node:GetTransform():GetPos(), node:GetTransform():GetPos())
		--priority = node:GetTransform():GetPos().y  -- uncomment to prioritize lights near the ground
		node:GetLight():SetPriority(-priority)
	end

	scene:Update(dt)
	vid, passId = hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res)

	hg.Frame()
	hg.UpdateWindow(win)
end

hg.DestroyForwardPipeline(pipeline)
hg.RenderShutdown()
hg.DestroyWindow(win)
