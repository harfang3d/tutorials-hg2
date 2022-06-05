-- Scene using the PBR shader

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 940, 720
win = hg.RenderInit('PBR Scene', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

--
pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

hg.AddAssetsFolder("resources_compiled")

-- load scene
scene = hg.Scene()
hg.LoadSceneFromAssets("materials/materials.scn", scene, res, hg.GetForwardPipelineInfo())

-- main loop
while not hg.ReadKeyboard():Key(hg.K_Escape) and hg.IsWindowOpen(win) do
	dt = hg.TickClock()

	scene:Update(dt)
	hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res)

	hg.Frame()
	hg.UpdateWindow(win)
end
hg.RenderShutdown()
hg.DestroyWindow(win)
