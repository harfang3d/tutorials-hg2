// Toyota 2JZ-GTE Engine model by Serhii Denysenko (CGTrader: serhiidenysenko8256)
// URL : https://www.cgtrader.com/3d-models/vehicle/part/toyota-2jz-gte-engine-2932b715-2f42-4ecd-93ce-df9507c67ce8
package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("AAA Scene", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	hg.AddAssetsFolder("resources_compiled")

	// load scene
	scene := hg.NewScene()
	hg.LoadSceneFromAssets("car_engine/engine.scn", scene, res, hg.GetForwardPipelineInfo())

	// AAA pipeline
	pipeline_aaa_config := hg.NewForwardPipelineAAAConfig()
	pipeline_aaa := hg.CreateForwardPipelineAAAFromAssetsWithSsgiRatioSsrRatio("core", pipeline_aaa_config, hg.BREqual, hg.BREqual)
	pipeline_aaa_config.SetSampleCount(1)

	// main loop
	frame := int32(0)
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {

		dt := hg.TickClock()

		trs := scene.GetNode("engine_master").GetTransform()
		trs.SetRot(trs.GetRot().Add(hg.NewVec3WithXYZ(0, hg.Deg(15)*hg.TimeToSecF(dt), 0)))

		scene.Update(dt)
		viewId := uint16(0)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontalAaaAaaConfigFrame(&viewId, scene, hg.NewIntRectWithSxSyExEy(0, 0, res_x, res_y), true, pipeline, res, pipeline_aaa, pipeline_aaa_config, frame)

		frame = int32(hg.Frame())
		hg.UpdateWindow(win)
	}

	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
