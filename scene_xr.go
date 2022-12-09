package main

/*
#include <stdint.h>
extern void update_controllersCallback(uintptr);
extern uint16_t draw_sceneCallback(uintptr, uintptr, uintptr, uintptr);
*/
import "C"

import (
	"unsafe"

	hg "github.com/harfang3d/harfang-go/v3"
)

var (
	prg_ref     *hg.PipelineProgramRef
	scene       *hg.Scene
	render_data *hg.SceneForwardPipelineRenderData

	pipeline *hg.ForwardPipeline
	res      *hg.PipelineResources
)

// Create materials
func create_material(ubc *hg.Vec4, orm *hg.Vec4) *hg.Material {
	mat := hg.NewMaterial()
	hg.SetMaterialProgram(mat, prg_ref)
	hg.SetMaterialValueWithVec4V(mat, "uBaseOpacityColor", ubc)
	hg.SetMaterialValueWithVec4V(mat, "uOcclusionRoughnessMetalnessColor", orm)
	return mat
}

//export update_controllersCallback
func update_controllersCallback(head uintptr) {
}

//export draw_sceneCallback
func draw_sceneCallback(rect uintptr, viewState uintptr, view_id *uint16, fb uintptr) uint16 {
	rect_ := hg.NewIntRectFromCPointer(unsafe.Pointer(rect))
	viewState_ := hg.NewViewStateFromCPointer(unsafe.Pointer(viewState))
	fb_ := hg.NewFrameBufferHandleFromCPointer(unsafe.Pointer(fb))
	passId := hg.NewSceneForwardPipelinePassViewId()
	hg.PrepareSceneForwardPipelineViewDependentRenderData(view_id, viewState_, scene, render_data, pipeline, res, passId)
	passId = hg.SubmitSceneToForwardPipelineWithFrameBuffer(view_id, scene, rect_, viewState_, pipeline, render_data, res, fb_)

	hg.Frame()
	return *view_id
}

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	hg.SetLogLevel(hg.LLAll)
	hg.SetLogDetailed(true)

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightTypeResetFlags("Harfang - OpenXR Scene", res_x, res_y, hg.RTDirect3D11, hg.RFNone)

	hg.AddAssetsFolder("resources_compiled")

	pipeline = hg.CreateForwardPipeline()
	res = hg.NewPipelineResources()

	render_data = hg.NewSceneForwardPipelineRenderData() // this object is used by the low-level scene rendering API to share view-independent data with both eyes

	// OpenVR initialization
	if !hg.OpenXRInit() {
		panic("Can't initialize OpenXR")
	}

	eye_framebuffer := hg.OpenXRCreateEyeFrameBufferWithAa(hg.OXRAAMSAA4x)

	// Create models
	vtx_layout := hg.VertexLayoutPosFloatNormUInt8()

	cube_mdl := hg.CreateCubeModel(vtx_layout, 1, 1, 1)
	cube_ref := res.AddModel("cube", cube_mdl)
	ground_mdl := hg.CreateCubeModel(vtx_layout, 50, 0.01, 50)
	ground_ref := res.AddModel("ground", ground_mdl)

	// Load shader
	prg_ref = hg.LoadPipelineProgramRefFromAssets("core/shader/pbr.hps", res, hg.GetForwardPipelineInfo())

	// Create scene
	scene = hg.NewScene()
	scene.GetCanvas().SetColor(hg.NewColorWithRGBA(255.0/255.0, 255.0/255.0, 217.0/255.0, 1.0))
	scene.GetEnvironment().SetAmbient(hg.NewColorWithRGBA(15.0/255.0, 12.0/255.0, 9.0/255.0, 1.0))

	hg.CreateSpotLightWithDiffuseSpecularPriorityShadowTypeShadowBias(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(-8, 4, -5), hg.NewVec3WithXYZ(hg.Deg(19), hg.Deg(59), 0)), 0, hg.Deg(5), hg.Deg(30), hg.ColorGetWhite(), hg.ColorGetWhite(), 10, hg.LSTMap, 0.0001)
	hg.CreatePointLightWithDiffuseSpecularPriority(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(2.4, 1, 0.5)), 10, hg.NewColorWithRGBA(94.0/255.0, 255.0/255.0, 228.0/255.0, 1.0), hg.NewColorWithRGBA(94.0/255.0, 1.0, 228.0/255.0, 1.0), 0)

	mat_cube := create_material(hg.NewVec4WithXYZW(255.0/255.0, 230.0/255.0, 255.0/255.0, 1), hg.NewVec4WithXYZW(1, 0.658, 0., 1))
	hg.CreateObjectWithSliceOfMaterials(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 0.5, 0), hg.NewVec3WithXYZ(0, hg.Deg(70), 0)), cube_ref, hg.GoSliceOfMaterial{mat_cube})

	mat_ground := create_material(hg.NewVec4WithXYZW(255/255.0, 120/255.0, 147/255.0, 1), hg.NewVec4WithXYZW(1, 1, 0.1, 1))
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

	//<-time.After(time.Second * 3)

	// Main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {

		dt := hg.TickClock()

		scene.Update(dt)

		vid := uint16(0) // keep track of the next free view id
		passId := hg.NewSceneForwardPipelinePassViewId()

		// Prepare view-independent render data once
		hg.PrepareSceneForwardPipelineCommonRenderData(&vid, scene, render_data, pipeline, res, passId)

		openxrFrameInfo := hg.OpenXRSubmitSceneToForwardPipeline(hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0, 0)), C.update_controllersCallback, C.draw_sceneCallback, &vid, 0.1, 100)

		// Display the VR eyes texture to the backbuffer
		hg.SetViewRect(vid, 0, 0, uint16(res_x), uint16(res_y))
		vs := hg.ComputeOrthographicViewState(hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0, 0)), float32(res_y), 0.1, 100, hg.ComputeAspectRatioX(float32(res_x), float32(res_y)))
		hg.SetViewTransform(vid, vs.GetView(), vs.GetProj())

		if openxrFrameInfo.GetIdFbs().Len() > 0 {
			quad_uniform_set_texture_list.Clear()
			quad_uniform_set_texture_list.PushBack(hg.MakeUniformSetTexture("s_tex", hg.OpenXRGetColorTextureFromId(eye_framebuffer, openxrFrameInfo, 0), 0))
			hg.SetT(quad_matrix, hg.NewVec3WithXYZ(eye_t_x, 0, 1))
			hg.DrawModelWithRenderState(vid, quad_model, tex0_program, quad_uniform_set_value_list, quad_uniform_set_texture_list, quad_matrix, quad_render_state)

			quad_uniform_set_texture_list.Clear()
			quad_uniform_set_texture_list.PushBack(hg.MakeUniformSetTexture("s_tex", hg.OpenXRGetColorTextureFromId(eye_framebuffer, openxrFrameInfo, 1), 0))
			hg.SetT(quad_matrix, hg.NewVec3WithXYZ(-eye_t_x, 0, 1))
			hg.DrawModelWithRenderState(vid, quad_model, tex0_program, quad_uniform_set_value_list, quad_uniform_set_texture_list, quad_matrix, quad_render_state)
		}
		hg.Frame()

		hg.OpenXRFinishSubmitFrameBuffer(openxrFrameInfo)

		hg.UpdateWindow(win)
	}

	hg.DestroyForwardPipeline(pipeline)
	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
