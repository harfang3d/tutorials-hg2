package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(940), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("PBR Scene", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	hg.AddAssetsFolder("resources_compiled")

	// load scene
	scene := hg.NewScene()
	hg.LoadSceneFromAssets("materials/materials.scn", scene, res, hg.GetForwardPipelineInfo())

	// main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		dt := hg.TickClock()

		scene.Update(dt)
		viewId := uint16(0)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontal(&viewId, scene, hg.NewIntRectWithSxSyExEy(0, 0, res_x, res_y), true, pipeline, res)

		hg.Frame()
		hg.UpdateWindow(win)
	}
	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
