package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Node Physics Setup", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	// create models
	vtx_mdl := hg.VertexLayoutPosFloatNormUInt8()

	cube_mdl := hg.CreateCubeModel(vtx_mdl, 1, 1, 1)
	cube_ref := res.AddModel("cube", cube_mdl)

	ground_mdl := hg.CreateCubeModel(vtx_mdl, 50, 0.01, 50)
	ground_ref := res.AddModel("ground", ground_mdl)

	// create materials
	prg_ref := hg.LoadPipelineProgramRefFromFile("resources_compiled/core/shader/default.hps", res, hg.GetForwardPipelineInfo())
	mat := hg.CreateMaterialWithValueName0Value0ValueName1Value1(prg_ref, "uDiffuseColor", hg.NewVec4WithXYZW(0.5, 0.5, 0.5, 1), "uSpecularColor", hg.NewVec4WithXYZW(0.0, 0.0, 0.0, 0.1))

	// setup scene
	scene := hg.NewScene()

	cam := hg.CreateCamera(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 1, -5), hg.Deg3(5, 0, 0)), 0.01, 1000)
	scene.SetCurrentCamera(cam)

	hg.CreatePointLight(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(6, 4, -6)), 0)
	hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(100, 0.02, 100), hg.TranslationMat4(hg.NewVec3WithXYZ(0, -0.005, 0)), ground_ref, hg.GoSliceOfMaterial{mat}, 0)

	clocks := hg.NewSceneClocks()

	// setup physic cube
	cube_node := hg.CreateObjectWithSliceOfMaterials(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 2.5, 0), hg.NewVec3WithXYZ(0, 0, 0)), cube_ref, hg.GoSliceOfMaterial{mat})

	rb := scene.CreateRigidBody()
	rb.SetType(hg.RBTDynamic)

	collision := scene.CreateCollision()
	collision.SetType(hg.CTCube)
	collision.SetSize(hg.NewVec3WithXYZ(1, 1, 1))
	collision.SetMass(1)

	cube_node.SetRigidBody(rb)
	cube_node.SetCollision(0, collision)

	// scene physics
	physics := hg.NewSceneBullet3Physics()
	physics.SceneCreatePhysicsFromAssets(scene)

	// main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		dt := hg.TickClock()

		hg.SceneUpdateSystemsWithPhysicsStepMaxPhysicsStep(scene, clocks, dt, physics, hg.TimeFromSecF(1.0/60.0), 1)
		viewId := uint16(0)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontal(&viewId, scene, hg.NewIntRectWithSxSyExEy(0, 0, res_x, res_y), true, pipeline, res)

		hg.Frame()
		hg.UpdateWindow(win)
	}
	hg.RenderShutdown()
	hg.DestroyWindow(win)

	hg.WindowSystemShutdown()
	hg.InputShutdown()
}
