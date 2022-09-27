package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Physics Matrix Interaction", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	hg.ImGuiInit(10, hg.LoadProgramFromFile("resources_compiled/core/shader/imgui"), hg.LoadProgramFromFile("resources_compiled/core/shader/imgui_image"))

	// physics debug
	vtx_line_layout := hg.VertexLayoutPosFloatColorUInt8()
	line_shader := hg.LoadProgramFromFile("resources_compiled/shaders/pos_rgb")

	// create models
	vtx_layout := hg.VertexLayoutPosFloatNormUInt8()

	cube_mdl := hg.CreateCubeModel(vtx_layout, 1, 1, 1)
	cube_ref := res.AddModel("cube", cube_mdl)
	ground_mdl := hg.CreateCubeModel(vtx_layout, 100, 0.02, 100)
	ground_ref := res.AddModel("ground", ground_mdl)

	prg_ref := hg.LoadPipelineProgramRefFromFile("resources_compiled/core/shader/default.hps", res, hg.GetForwardPipelineInfo())

	// create material
	mat := hg.CreateMaterialWithValueName0Value0ValueName1Value1(prg_ref, "uDiffuseColor", hg.NewVec4WithXYZ(0.5, 0.5, 0.5), "uSpecularColor", hg.NewVec4WithXYZ(1, 1, 1))

	// setup scene
	scene := hg.NewScene()

	cam_mat := hg.TransformationMat4(hg.NewVec3WithXYZ(0, 1.5, -5), hg.Deg3(5, 0, 0))
	cam := hg.CreateCamera(scene, cam_mat, 0.01, 1000)
	scene.SetCurrentCamera(cam)
	view_matrix := hg.InverseFast(cam_mat)
	c := cam.GetCamera()
	projection_matrix := hg.ComputePerspectiveProjectionMatrix(c.GetZNear(), c.GetZFar(), hg.FovToZoomFactor(c.GetFov()), hg.NewVec2WithXY(float32(res_x)/float32(res_y), 1))

	hg.CreatePointLight(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(3, 4, -6)), 0)

	cube_node := hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(1, 1, 1), hg.TranslationMat4(hg.NewVec3WithXYZ(1.25, 2.5, 0)), cube_ref, hg.GoSliceOfMaterial{mat}, 2)
	hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(100, 0.02, 100), hg.TranslationMat4(hg.NewVec3WithXYZ(0, -0.005, 0)), ground_ref, hg.GoSliceOfMaterial{mat}, 0)

	clocks := hg.NewSceneClocks()

	// scene physics
	physics := hg.NewSceneBullet3Physics()
	physics.SceneCreatePhysicsFromAssets(scene)

	// main loop
	mouse, keyboard := hg.NewMouse(), hg.NewKeyboard()

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		keyboard.Update()
		mouse.Update()
		dt := hg.TickClock()

		// scene view
		view_id := uint16(0)
		hg.SceneUpdateSystemsWithPhysicsStepMaxPhysicsStep(scene, clocks, dt, physics, hg.TimeFromSecF(1.0/60.0), 1)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontal(&view_id, scene, hg.NewIntRectWithSxSyExEy(0, 0, res_x, res_y), true, pipeline, res)

		// Debug physics display
		hg.SetViewClearWithColDepthStencil(view_id, 0, hg.ColorGetBlack(), 1.0, 0)
		hg.SetViewRect(view_id, 0, 0, uint16(res_x), uint16(res_y))
		hg.SetViewTransform(view_id, view_matrix, projection_matrix)
		rs := hg.ComputeRenderStateWithDepthTestCulling(hg.BMOpaque, hg.DTDisabled, hg.FCDisabled)
		physics.RenderCollision(view_id, vtx_line_layout, line_shader, rs, 0)

		// ImGui view
		hg.ImGuiBeginFrame(res_x, res_y, dt, mouse.GetState(), keyboard.GetState())

		imguiOpen := true
		if hg.ImGuiBeginWithOpenFlags("Transform and Physics", &imguiOpen, hg.ImGuiWindowFlagsAlwaysAutoResize) {
			hg.ImGuiTextWrapped("This tutorial demonstrates the interaction between the physics system and the Transform component of a node. The node position, rotation and scale are overriden by an active node rigid body.")

			hg.ImGuiSeparator()
			v := cube_node.GetTransform().GetPos()
			r := hg.ImGuiInputVec3WithDecimalPrecision("Transform.pos", v, 2)
			if r {
				cube_node.GetTransform().SetPos(v)
			}

			if hg.ImGuiButton("Press to reset position using Transform.SetPos") {
				cube_node.GetTransform().SetPos(hg.NewVec3WithXYZ(1.25, 2.5, 0))
			}

			hg.ImGuiSeparator()

			hg.ImGuiInputVec3WithDecimalPrecisionFlags("Transform.GetWorld().T", hg.GetT(cube_node.GetTransform().GetWorld()), 2, hg.ImGuiInputTextFlagsReadOnly)

			if physics.NodeHasBody(cube_node) {
				if hg.ImGuiButton("Press to destroy the cube node physics") {
					physics.NodeDestroyPhysics(cube_node)
					physics.GarbageCollect(scene)
				}
			} else {
				if hg.ImGuiButton("Create to create the cube node physics") {
					physics.NodeCreatePhysicsFromAssets(cube_node)
				}
			}
		}

		hg.ImGuiEnd()

		hg.ImGuiEndFrameWithViewId(255)

		hg.Frame()
		hg.UpdateWindow(win)
	}
	hg.RenderShutdown()
	hg.DestroyWindow(win)

	hg.WindowSystemShutdown()
	hg.InputShutdown()
}
