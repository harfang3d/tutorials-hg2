package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Scene Draw to Multiple Viewports", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	hg.AddAssetsFolder("resources_compiled")

	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	// create models
	vtx_layout := hg.VertexLayoutPosFloatNormUInt8()

	cube_mdl := hg.CreateCubeModel(vtx_layout, 1, 1, 1)
	cube_ref := res.AddModel("cube", cube_mdl)
	ground_mdl := hg.CreateCubeModel(vtx_layout, 100, 0.01, 100)
	ground_ref := res.AddModel("ground", ground_mdl)

	// create materials
	shader := hg.LoadPipelineProgramRefFromAssets("core/shader/default.hps", res, hg.GetForwardPipelineInfo())

	mat_yellow_cube := hg.CreateMaterialWithValueName0Value0ValueName1Value1(shader, "uDiffuseColor", hg.Vec4I(255, 220, 64), "uSpecularColor", hg.Vec4I(255, 220, 64))
	mat_red_cube := hg.CreateMaterialWithValueName0Value0ValueName1Value1(shader, "uDiffuseColor", hg.Vec4I(255, 0, 0), "uSpecularColor", hg.Vec4I(255, 0, 0))
	mat_ground := hg.CreateMaterialWithValueName0Value0ValueName1Value1(shader, "uDiffuseColor", hg.Vec4I(128, 128, 128), "uSpecularColor", hg.Vec4I(128, 128, 128))

	// setup scene (note that we do not create any camera)
	scene := hg.NewScene()

	hg.CreateSpotLightWithDiffuseSpecularPriorityShadowTypeShadowBias(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(-8, 4, -5), hg.Deg3(19, 59, 0)), 0, hg.Deg(5), hg.Deg(30), hg.ColorGetWhite(), hg.ColorGetWhite(), 10, hg.LSTMap, 0.00005)
	hg.CreatePointLightWithDiffuseSpecularPriority(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(3, 1, 2.5)), 5, hg.ColorI(128, 192, 255), hg.ColorGetBlack(), 0)

	yellow_cube := hg.CreateObjectWithSliceOfMaterials(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(1, 0.5, 0), hg.NewVec3WithXYZ(0, hg.Deg(0), 0)), cube_ref, hg.GoSliceOfMaterial{mat_yellow_cube})
	hg.CreateObjectWithSliceOfMaterials(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(-1, 0.5, 0)), cube_ref, hg.GoSliceOfMaterial{mat_red_cube})
	hg.CreateObjectWithSliceOfMaterials(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0, 0)), ground_ref, hg.GoSliceOfMaterial{mat_ground})

	// define viewports
	viewports := []struct {
		rect             *hg.IntRect
		cam_pos, cam_rot *hg.Vec3
	}{
		{rect: hg.NewIntRectWithSxSyExEy(0, 0, res_x/2, res_y/2), cam_pos: hg.NewVec3WithXYZ(-4.015, 2.368, -3.484), cam_rot: hg.NewVec3WithXYZ(0.35, 0.87, 0.0)},
		{rect: hg.NewIntRectWithSxSyExEy(res_x/2, 0, res_x, res_y/2), cam_pos: hg.NewVec3WithXYZ(-4.143, 2.976, 4.127), cam_rot: hg.NewVec3WithXYZ(0.423, 2.365, 0.0)},
		{rect: hg.NewIntRectWithSxSyExEy(0, res_y/2, res_x/2, res_y), cam_pos: hg.NewVec3WithXYZ(4.020, 2.374, 3.469), cam_rot: hg.NewVec3WithXYZ(0.353, 4.016, 0.0)},
		{rect: hg.NewIntRectWithSxSyExEy(res_x/2, res_y/2, res_x, res_y), cam_pos: hg.NewVec3WithXYZ(3.469, 2.374, -4.020), cam_rot: hg.NewVec3WithXYZ(0.353, -0.695, 0.0)},
	}

	// main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		dt := hg.TickClock()

		// animate yellow cube & update scene once for all viewports
		rot := yellow_cube.GetTransform().GetRot()
		rot.SetY(rot.GetY() + hg.TimeToSecF(dt))
		yellow_cube.GetTransform().SetRot(rot)

		scene.Update(dt)

		// prepare view-independent render data (eg. spot shadow maps)
		render_data := hg.NewSceneForwardPipelineRenderData()

		views := hg.NewSceneForwardPipelinePassViewId()
		vid := uint16(0)
		hg.PrepareSceneForwardPipelineCommonRenderData(&vid, scene, render_data, pipeline, res, views)

		for _, viewport := range viewports {
			// compute viewport specific view state
			view_state := hg.ComputePerspectiveViewState(hg.TransformationMat4(viewport.cam_pos, viewport.cam_rot), hg.Deg(45), 0.01, 1000, hg.ComputeAspectRatioX(float32(res_x), float32(res_y)))
			// prepare view-dependent render data & submit draw
			hg.PrepareSceneForwardPipelineViewDependentRenderData(&vid, view_state, scene, render_data, pipeline, res, views)
			hg.SubmitSceneToForwardPipeline(&vid, scene, viewport.rect, view_state, pipeline, render_data, res)
		}
		hg.Frame()
		hg.UpdateWindow(win)
	}
	hg.DestroyForwardPipeline(pipeline)
	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
