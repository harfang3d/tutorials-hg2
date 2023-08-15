-- Toyota 2JZ-GTE Engine model by Serhii Denysenko (CGTrader: serhiidenysenko8256)
-- URL : https://www.cgtrader.com/3d-models/vehicle/part/toyota-2jz-gte-engine-2932b715-2f42-4ecd-93ce-df9507c67ce8

hg = require("harfang")
math.randomseed(os.time())

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('AAA Scene', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

hg.AddAssetsFolder("resources_compiled")

-- load scene
scene = hg.Scene()
hg.LoadSceneFromAssets("car_engine/engine.scn", scene, res, hg.GetForwardPipelineInfo())

-- AAA pipeline
pipeline_aaa_config = hg.ForwardPipelineAAAConfig()
pipeline_aaa = hg.CreateForwardPipelineAAAFromAssets("core", pipeline_aaa_config, hg.BR_Equal, hg.BR_Equal)
pipeline_aaa_config.sample_count = 1
target_dof_focus_point = 3.5
target_dof_focus_length = 2.0
pipeline_aaa_config.dof_focus_point = target_dof_focus_point -- Distance to the focus point (in meters)
pipeline_aaa_config.dof_focus_length = target_dof_focus_length -- Depth of field (in meters); smaller values result in a narrower focused area.

-- main loop
frame = 0

while not hg.ReadKeyboard():Key(hg.K_Escape) and hg.IsWindowOpen(win) do
	dt = hg.TickClock()

	trs = scene:GetNode('engine_master'):GetTransform()
	trs:SetRot(trs:GetRot() + hg.Vec3(0, hg.Deg(15) * hg.time_to_sec_f(dt), 0))

	-- change DOF randomly
	if frame % 250 == 0 then
		target_dof_focus_point = math.random() + 2.5 -- random value between 2.5 and 3.5
		target_dof_focus_length = math.random() * 1.5 + 0.5 -- random value between 0.5 and 2.0
	end
	
	pipeline_aaa_config.dof_focus_point = hg.Lerp(pipeline_aaa_config.dof_focus_point, target_dof_focus_point, 0.1)
	pipeline_aaa_config.dof_focus_length = hg.Lerp(pipeline_aaa_config.dof_focus_length, target_dof_focus_length, 0.1)

	scene:Update(dt)
	hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res, pipeline_aaa, pipeline_aaa_config, frame)

	frame = hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)
