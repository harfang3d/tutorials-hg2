package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Physics Force/Impulse (Press space to alternate)", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	// create models
	vtx_layout := hg.VertexLayoutPosFloatNormUInt8()

	cube_mdl := hg.CreateCubeModel(vtx_layout, 1, 1, 1)
	cube_ref := res.AddModel("cube", cube_mdl)

	ground_mdl := hg.CreateCubeModel(vtx_layout, 50, 0.01, 50)
	ground_ref := res.AddModel("ground", ground_mdl)

	// create material
	prg_ref := hg.LoadPipelineProgramRefFromFile("resources_compiled/core/shader/default.hps", res, hg.GetForwardPipelineInfo())
	mat := hg.CreateMaterialWithValueName0Value0ValueName1Value1(prg_ref, "uDiffuseColor", hg.NewVec4WithXYZ(1, 1, 1), "uSpecularColor", hg.NewVec4WithXYZ(1, 1, 1))

	// setup scene
	scene := hg.NewScene()

	cam := hg.CreateCamera(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 1.5, -5), hg.NewVec3WithXYZ(hg.Deg(10), 0, 0)), 0.01, 1000)
	scene.SetCurrentCamera(cam)

	hg.CreateLinearLightWithDiffuseSpecularPriorityShadowTypeShadowBiasPssmSplit(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 0, 0), hg.NewVec3WithXYZ(hg.Deg(30), hg.Deg(59), 0)), hg.NewColorWithRGB(1, 1, 1), hg.NewColorWithRGB(1, 1, 1), 10, hg.LSTMap, 0.002, hg.NewVec4WithXYZW(2, 4, 10, 16))

	cube_node := hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(1, 1, 1), hg.TranslationMat4(hg.NewVec3WithXYZ(0, 1.5, 0)), cube_ref, hg.GoSliceOfMaterial{mat}, 2)
	hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(100, 0.02, 100), hg.TranslationMat4(hg.NewVec3WithXYZ(0, -0.005, 0)), ground_ref, hg.GoSliceOfMaterial{mat}, 0)

	clocks := hg.NewSceneClocks()

	// scene physics
	physics := hg.NewSceneBullet3Physics()
	physics.SceneCreatePhysicsFromAssets(scene)
	physics_step := hg.TimeFromSecF(1.0 / 60.0)

	// main loop
	keyboard := hg.NewKeyboard()

	use_force := true

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		keyboard.Update()

		dt := hg.TickClock()

		if keyboard.Pressed(hg.KSpace) {
			use_force = !use_force
		}

		world_pos := hg.GetT(cube_node.GetTransform().GetWorld())
		dist_to_ground := world_pos.GetY() - 0.5

		if dist_to_ground < 1.0 {
			k := -(dist_to_ground - 1.0)

			if use_force {
				F := hg.NewVec3WithXYZ(0, 1, 0).MulWithK(k * 80) // apply a force inversely proportional to the distance to the ground
				physics.NodeAddForceWithWorldPos(cube_node, F, world_pos)
			} else {
				stiffness := 10

				cur_velocity := physics.NodeGetLinearVelocity(cube_node)
				tgt_velocity := hg.NewVec3WithXYZ(0, 1, 0).MulWithK(k * float32(stiffness)) // compute a velocity that brings us to 1 meter above the ground

				I := tgt_velocity.Sub(cur_velocity) // an impulse is an instantaneous change in velocity
				physics.NodeAddImpulseWithWorldPos(cube_node, I, world_pos)
			}
		}
		physics.NodeWake(cube_node)

		hg.SceneUpdateSystemsWithPhysicsStepMaxPhysicsStep(scene, clocks, dt, physics, physics_step, 3)
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
