package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Modify material pipeline shader ", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

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

	mat_cube := hg.CreateMaterialWithValueName0Value0ValueName1Value1(shader, "uDiffuseColor", hg.NewVec4WithXYZ(1, 1, 1), "uSpecularColor", hg.NewVec4WithXYZ(1, 1, 1))
	mat_ground := hg.CreateMaterialWithValueName0Value0ValueName1Value1(shader, "uDiffuseColor", hg.NewVec4WithXYZ(1, 1, 1), "uSpecularColor", hg.NewVec4WithXYZ(0.1, 0.1, 0.1))

	// setup scene
	scene := hg.NewScene()

	cam := hg.CreateCamera(scene, hg.Mat4LookAt(hg.NewVec3WithXYZ(-1.30, 0.27, -2.47), hg.NewVec3WithXYZ(0, 0.5, 0)), 0.01, 1000)
	scene.SetCurrentCamera(cam)

	hg.CreateLinearLightWithDiffuseSpecularPriority(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 2, 0), hg.Deg3(27.5, -97.6, 16.6)), hg.ColorI(64, 64, 64), hg.ColorI(64, 64, 64), 10)
	hg.CreateSpotLightWithDiffuseSpecularPriorityShadowTypeShadowBias(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(5, 4, -5), hg.Deg3(19, -45, 0)), 0, hg.Deg(5), hg.Deg(30), hg.ColorI(255, 255, 255), hg.ColorI(255, 255, 255), 10, hg.LSTMap, 0.0001)

	cube_node := hg.CreateObjectWithSliceOfMaterials(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0.5, 0)), cube_ref, hg.GoSliceOfMaterial{mat_cube})
	hg.CreateObjectWithSliceOfMaterials(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0, 0)), ground_ref, hg.GoSliceOfMaterial{mat_ground})

	// material update states
	mat_has_texture := false
	mat_update_delay := int64(0)

	texture_ref := hg.LoadTextureFromAssetsWithFlagsResources("textures/squares.png", 0, res)

	// main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		dt := hg.TickClock()

		mat_update_delay = mat_update_delay - dt

		if mat_update_delay <= 0 {
			// set or remove cube node material texture
			mat := cube_node.GetObject().GetMaterial(0)

			if mat_has_texture {
				hg.SetMaterialTexture(mat, "uDiffuseMap", &hg.InvalidTextureRef, 0)
			} else {
				hg.SetMaterialTexture(mat, "uDiffuseMap", texture_ref, 0)
			}

			// update the pipeline shader variant according to the material uniform values
			hg.UpdateMaterialPipelineProgramVariant(mat, res)

			// reset delay and flip flag
			mat_update_delay = mat_update_delay + hg.TimeFromSec(1)
			mat_has_texture = !mat_has_texture
		}

		scene.Update(dt)

		viewId := uint16(0)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontal(&viewId, scene, hg.NewIntRectWithSxSyExEy(0, 0, res_x, res_y), true, pipeline, res)

		hg.Frame()
		hg.UpdateWindow(win)
	}
	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
