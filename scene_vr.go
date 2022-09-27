package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

// Create materials
func create_material(ubc *hg.Vec4, orm *hg.Vec4, prg_ref *hg.PipelineProgramRef) *hg.Material {
	mat := hg.NewMaterial()
	hg.SetMaterialProgram(mat, prg_ref)
	hg.SetMaterialValueWithVec4V(mat, "uBaseOpacityColor", ubc)
	hg.SetMaterialValueWithVec4V(mat, "uOcclusionRoughnessMetalnessColor", orm)
	return mat
}

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - OpenVR Scene", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	hg.AddAssetsFolder("resources_compiled")

	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	render_data := hg.NewSceneForwardPipelineRenderData() // this object is used by the low-level scene rendering API to share view-independent data with both eyes

	// OpenVR initialization
	if !hg.OpenVRInit() {
		panic("Can't initialize OpenVR")
	}

	vr_left_fb := hg.OpenVRCreateEyeFrameBufferWithAa(hg.OVRAAMSAA4x)
	vr_right_fb := hg.OpenVRCreateEyeFrameBufferWithAa(hg.OVRAAMSAA4x)

	// Create models
	vtx_layout := hg.VertexLayoutPosFloatNormUInt8()

	cube_mdl := hg.CreateCubeModel(vtx_layout, 1, 1, 1)
	cube_ref := res.AddModel("cube", cube_mdl)
	ground_mdl := hg.CreateCubeModel(vtx_layout, 50, 0.01, 50)
	ground_ref := res.AddModel("ground", ground_mdl)

	// Load shader
	prg_ref := hg.LoadPipelineProgramRefFromAssets("core/shader/pbr.hps", res, hg.GetForwardPipelineInfo())

	// Create scene
	scene := hg.NewScene()
	scene.GetCanvas().SetColor(hg.NewColorWithRGBA(255/255, 255/255, 217/255, 1))
	scene.GetEnvironment().SetAmbient(hg.NewColorWithRGBA(15/255, 12/255, 9/255, 1))

	hg.CreateSpotLightWithDiffuseSpecularPriorityShadowTypeShadowBias(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(-8, 4, -5), hg.NewVec3WithXYZ(hg.Deg(19), hg.Deg(59), 0)), 0, hg.Deg(5), hg.Deg(30), hg.ColorGetWhite(), hg.ColorGetWhite(), 10, hg.LSTMap, 0.0001)
	hg.CreatePointLightWithDiffuseSpecularPriority(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(2.4, 1, 0.5)), 10, hg.NewColorWithRGBA(94.0/255.0, 255.0/255.0, 228.0/255.0, 1.0), hg.NewColorWithRGBA(94.0/255.0, 1.0, 228.0/255.0, 1.0), 0)

	mat_cube := create_material(hg.NewVec4WithXYZW(255/255, 230/255, 20/255, 1), hg.NewVec4WithXYZW(1, 0.658, 0., 1), prg_ref)
	hg.CreateObjectWithSliceOfMaterials(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 0.5, 0), hg.NewVec3WithXYZ(0, hg.Deg(70), 0)), cube_ref, hg.GoSliceOfMaterial{mat_cube})

	mat_ground := create_material(hg.NewVec4WithXYZW(255/255, 120/255, 147/255, 1), hg.NewVec4WithXYZW(1, 1, 0.1, 1), prg_ref)
	hg.CreateObjectWithSliceOfMaterials(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0, 0)), ground_ref, hg.GoSliceOfMaterial{mat_ground})

	// Setup 2D rendering to display eyes textures
	quad_layout := hg.NewVertexLayout()
	quad_layout.Begin().Add(hg.APosition, 3, hg.ATFloat).Add(hg.ATexCoord0, 3, hg.ATFloat).End()

	quad_model := hg.CreatePlaneModel(quad_layout, 1, 1, 1, 1)
	quad_render_state := hg.ComputeRenderStateWithDepthTestCulling(hg.BMAlpha, hg.DTDisabled, hg.FCDisabled)

	eye_t_size := float32(res_x) / 2.5
	eye_t_x := (float32(res_x)-2*eye_t_size)/6 + eye_t_size/2
	quad_matrix := hg.TransformationMat4WithScale(hg.NewVec3WithXYZ(0, 0, 0), hg.NewVec3WithXYZ(hg.Deg(90), hg.Deg(0), hg.Deg(0)), hg.NewVec3WithXYZ(eye_t_size, 1, eye_t_size))

	tex0_program := hg.LoadProgramFromAssets("shaders/sprite")

	quad_uniform_set_value_list := hg.NewUniformSetValueList()
	quad_uniform_set_value_list.Clear()
	quad_uniform_set_value_list.PushBack(hg.MakeUniformSetValueWithVec4V("color", hg.NewVec4WithXYZW(1, 1, 1, 1)))

	quad_uniform_set_texture_list := hg.NewUniformSetTextureList()

	// Main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {

		dt := hg.TickClock()

		scene.Update(dt)

		actor_body_mtx := hg.TransformationMat4(hg.NewVec3WithXYZ(-1.3, .45, -2), hg.NewVec3WithXYZ(0, 0, 0))

		vr_state := hg.OpenVRGetState(actor_body_mtx, 0.01, 1000)
		left, right := hg.OpenVRStateToViewState(vr_state)

		vid := uint16(0) // keep track of the next free view id
		passId := hg.NewSceneForwardPipelinePassViewId()

		// Prepare view-independent render data once
		hg.PrepareSceneForwardPipelineCommonRenderData(&vid, scene, render_data, pipeline, res, passId)
		vr_eye_rect := hg.NewIntRectWithSxSyExEy(0, 0, int32(vr_state.GetWidth()), int32(vr_state.GetHeight()))

		// Prepare the left eye render data then draw to its framebuffer
		hg.PrepareSceneForwardPipelineViewDependentRenderData(&vid, left, scene, render_data, pipeline, res, passId)
		hg.SubmitSceneToForwardPipelineWithFrameBuffer(&vid, scene, vr_eye_rect, left, pipeline, render_data, res, vr_left_fb.GetHandle())

		// Prepare the right eye render data then draw to its framebuffer
		hg.PrepareSceneForwardPipelineViewDependentRenderData(&vid, right, scene, render_data, pipeline, res, passId)
		hg.SubmitSceneToForwardPipelineWithFrameBuffer(&vid, scene, vr_eye_rect, right, pipeline, render_data, res, vr_right_fb.GetHandle())

		// Display the VR eyes texture to the backbuffer
		hg.SetViewRect(vid, 0, 0, uint16(res_x), uint16(res_y))
		vs := hg.ComputeOrthographicViewState(hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0, 0)), float32(res_y), 0.1, 100, hg.ComputeAspectRatioX(float32(res_x), float32(res_y)))
		hg.SetViewTransform(vid, vs.GetView(), vs.GetProj())

		quad_uniform_set_texture_list.Clear()
		quad_uniform_set_texture_list.PushBack(hg.MakeUniformSetTexture("s_tex", hg.OpenVRGetColorTexture(vr_left_fb), 0))
		hg.SetT(quad_matrix, hg.NewVec3WithXYZ(eye_t_x, 0, 1))
		hg.DrawModelWithRenderState(vid, quad_model, tex0_program, quad_uniform_set_value_list, quad_uniform_set_texture_list, quad_matrix, quad_render_state)

		quad_uniform_set_texture_list.Clear()
		quad_uniform_set_texture_list.PushBack(hg.MakeUniformSetTexture("s_tex", hg.OpenVRGetColorTexture(vr_right_fb), 0))
		hg.SetT(quad_matrix, hg.NewVec3WithXYZ(-eye_t_x, 0, 1))
		hg.DrawModelWithRenderState(vid, quad_model, tex0_program, quad_uniform_set_value_list, quad_uniform_set_texture_list, quad_matrix, quad_render_state)

		hg.Frame()
		hg.OpenVRSubmitFrame(vr_left_fb, vr_right_fb)

		hg.UpdateWindow(win)
	}

	hg.DestroyForwardPipeline(pipeline)
	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
