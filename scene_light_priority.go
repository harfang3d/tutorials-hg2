package main

import (
	"math"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Light priority relative to a specific world position", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	hg.AddAssetsFolder("resources_compiled")

	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	// create models
	vtx_layout := hg.VertexLayoutPosFloatNormUInt8()

	light_mdl := hg.CreateSphereModel(vtx_layout, 0.05, 8, 16)
	light_ref := res.AddModel("light", light_mdl)
	orb_mdl := hg.CreateSphereModel(vtx_layout, 1, 16, 32)
	orb_ref := res.AddModel("orb", orb_mdl)
	ground_mdl := hg.CreateCubeModel(vtx_layout, 100, 0.01, 100)
	ground_ref := res.AddModel("ground", ground_mdl)

	// create materials
	shader := hg.LoadPipelineProgramRefFromAssets("core/shader/default.hps", res, hg.GetForwardPipelineInfo())

	mat_light := hg.CreateMaterialWithValueName0Value0ValueName1Value1(shader, "uDiffuseColor", hg.NewVec4WithXYZ(0, 0, 0), "uSpecularColor", hg.NewVec4WithXYZ(0, 0, 0))
	hg.SetMaterialValueWithVec4V(mat_light, "uSelfColor", hg.NewVec4WithXYZ(1, 0.9, 0.75))
	mat_orb := hg.CreateMaterialWithValueName0Value0ValueName1Value1(shader, "uDiffuseColor", hg.NewVec4WithXYZ(1, 1, 1), "uSpecularColor", hg.NewVec4WithXYZ(1, 1, 1))
	hg.SetMaterialValueWithVec4V(mat_orb, "uSelfColor", hg.NewVec4WithXYZ(0, 0, 0))
	mat_ground := hg.CreateMaterialWithValueName0Value0ValueName1Value1(shader, "uDiffuseColor", hg.NewVec4WithXYZ(1, 1, 1), "uSpecularColor", hg.NewVec4WithXYZ(1, 1, 1))
	hg.SetMaterialValueWithVec4V(mat_ground, "uSelfColor", hg.NewVec4WithXYZ(0, 0, 0))

	// setup scene
	scene := hg.NewScene()

	cam := hg.CreateCamera(scene, hg.Mat4LookAt(hg.NewVec3WithXYZ(5, 4, -7), hg.NewVec3WithXYZ(0, 1.5, 0)), 0.01, 1000)
	scene.SetCurrentCamera(cam)

	orb_node := hg.CreateObjectWithSliceOfMaterials(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 1, 0)), orb_ref, hg.GoSliceOfMaterial{mat_orb})
	hg.CreateObjectWithSliceOfMaterials(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0, 0)), ground_ref, hg.GoSliceOfMaterial{mat_ground})

	// create an array of dynamic lights
	light_obj := scene.CreateObjectWithModelSliceOfMaterials(light_ref, hg.GoSliceOfMaterial{mat_light}) // sphere model to visualize lights

	light_nodes := []*hg.Node{}
	for i := 0; i < 16; i++ {
		node := hg.CreatePointLightWithDiffuseSpecular(scene, hg.Mat4GetIdentity(), 1.5, hg.NewColorWithRGBA(1, 0.85, 0.25, 1), hg.NewColorWithRGBA(1, 0.9, 0.5, 1))
		node.SetObject(light_obj)
		light_nodes = append(light_nodes, node)
	}
	// main loop
	angle := float32(0.0)

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		dt := hg.TickClock()

		// animate lights
		angle = angle + hg.TimeToSecF(dt)

		for i, node := range light_nodes {
			a := float64(angle + float32(i)*hg.Deg(15))
			node.GetTransform().SetPos(hg.NewVec3WithXYZ(float32(math.Cos(a*-0.6)*math.Sin(a)*5.0), float32(math.Cos(a*1.25)*2+2.15), float32(math.Sin(a*0.5)*math.Cos(-a*0.8)*5.0)))
		}
		// update light priorities according to their distance to the orb
		for _, node := range light_nodes {
			priority := hg.Dist2WithVec3AVec3B(orb_node.GetTransform().GetPos(), node.GetTransform().GetPos())
			//priority = node.GetTransform().GetPos().y  // uncomment to prioritize lights near the ground
			node.GetLight().SetPriority(-priority)
		}

		scene.Update(dt)
		viewId := uint16(0)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontal(&viewId, scene, hg.NewIntRectWithSxSyExEy(0, 0, res_x, res_y), true, pipeline, res)

		hg.Frame()
		hg.UpdateWindow(win)
	}

	hg.DestroyForwardPipeline(pipeline)
	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
