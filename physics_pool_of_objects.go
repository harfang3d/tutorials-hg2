package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

var prg_ref *hg.PipelineProgramRef

func create_material(diffuse *hg.Vec4, specular *hg.Vec4, self *hg.Vec4) *hg.Material {
	mat := hg.CreateMaterialWithValueName0Value0ValueName1Value1(prg_ref, "uDiffuseColor", diffuse, "uSpecularColor", specular)
	hg.SetMaterialValueWithVec4V(mat, "uSelfColor", self)
	return mat
}

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Physics Pool", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	hg.AddAssetsFolder("resources_compiled")

	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	// create models
	vtx_layout := hg.VertexLayoutPosFloatNormUInt8()

	sphere_mdl := hg.CreateSphereModel(vtx_layout, 0.5, 12, 24)
	sphere_ref := res.AddModel("sphere", sphere_mdl)
	cube_mdl := hg.CreateCubeModel(vtx_layout, 1, 1, 1)
	cube_ref := res.AddModel("cube", cube_mdl)

	// create materials
	prg_ref = hg.LoadPipelineProgramRefFromAssets("core/shader/default.hps", res, hg.GetForwardPipelineInfo())

	mat_ground := create_material(hg.NewVec4WithXYZ(0.5, 0.5, 0.5), hg.NewVec4WithXYZ(0.1, 0.1, 0.1), hg.NewVec4WithXYZ(0, 0, 0))
	mat_walls := create_material(hg.NewVec4WithXYZ(0.5, 0.5, 0.5), hg.NewVec4WithXYZ(0.1, 0.1, 0.1), hg.NewVec4WithXYZ(0, 0, 0))
	mat_objects := create_material(hg.NewVec4WithXYZ(0.5, 0.5, 0.5), hg.NewVec4WithXYZ(1, 1, 1), hg.NewVec4WithXYZ(0, 0, 0))

	// setup scene
	scene := hg.NewScene()
	scene.GetCanvas().SetColor(hg.ColorI(22, 56, 76))
	scene.GetEnvironment().SetFogColor(scene.GetCanvas().GetColor())
	scene.GetEnvironment().SetFogNear(20)
	scene.GetEnvironment().SetFogFar(80)

	cam_mtx := hg.TransformationMat4(hg.NewVec3WithXYZ(0, 20, -30), hg.Deg3(30, 0, 0))
	cam := hg.CreateCamera(scene, cam_mtx, 0.01, 5000)
	scene.SetCurrentCamera(cam)

	hg.CreateLinearLightWithDiffuseSpecularPriorityShadowTypeShadowBiasPssmSplit(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 0, 0), hg.Deg3(30, 59, 0)), hg.NewColorWithRGB(1, 0.8, 0.7), hg.NewColorWithRGB(1, 0.8, 0.7), 10, hg.LSTMap, 0.002, hg.NewVec4WithXYZW(50, 100, 200, 400))
	hg.CreatePointLightWithDiffuseSpecular(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 10, 10)), 100, hg.ColorI(94, 155, 228), hg.ColorI(94, 255, 228))

	mdl_ref := res.AddModel("ground", hg.CreateCubeModel(vtx_layout, 100, 1, 100))
	hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(30, 1, 30), hg.TranslationMat4(hg.NewVec3WithXYZ(0, -0.5, 0)), mdl_ref, hg.GoSliceOfMaterial{mat_ground}, 0)
	mdl_ref = res.AddModel("wall", hg.CreateCubeModel(vtx_layout, 1, 11, 32))
	hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(1, 11, 32), hg.TranslationMat4(hg.NewVec3WithXYZ(-15.5, -0.5, 0)), mdl_ref, hg.GoSliceOfMaterial{mat_walls}, 0)
	hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(1, 11, 32), hg.TranslationMat4(hg.NewVec3WithXYZ(15.5, -0.5, 0)), mdl_ref, hg.GoSliceOfMaterial{mat_walls}, 0)
	mdl_ref = res.AddModel("wall2", hg.CreateCubeModel(vtx_layout, 32, 11, 1))
	hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(32, 11, 1), hg.TranslationMat4(hg.NewVec3WithXYZ(0, -0.5, -15.5)), mdl_ref, hg.GoSliceOfMaterial{mat_walls}, 0)
	hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(32, 11, 1), hg.TranslationMat4(hg.NewVec3WithXYZ(0, -0.5, 15.5)), mdl_ref, hg.GoSliceOfMaterial{mat_walls}, 0)

	clocks := hg.NewSceneClocks()

	// setup physics
	physics := hg.NewSceneBullet3Physics()
	physics.SceneCreatePhysicsFromAssets(scene)

	physic_nodes := []*hg.Node{} // keep track of dynamically created physic nodes

	// text rendering
	font := hg.LoadFontFromAssetsWithSize("font/default.ttf", 32)
	font_program := hg.LoadProgramFromAssetsWithVertexShaderNameFragmentShaderName("core/shader/font.vsb", "core/shader/font.fsb")

	text_uniform_values := hg.GoSliceOfUniformSetValue{hg.MakeUniformSetValueWithVec4V("u_color", hg.NewVec4WithXYZ(1, 1, 0.5))}
	text_render_state := hg.ComputeRenderStateWithDepthTestCulling(hg.BMAlpha, hg.DTAlways, hg.FCDisabled)

	// main loop
	for hg.IsWindowOpen(win) {

		state := hg.ReadKeyboard()

		if state.Key(hg.KS) {
			for i := 1; i < 8; i += 1 {
				hg.SetMaterialValueWithVec4V(mat_objects, "uDiffuseColor", hg.RandomVec4(0, 1))

				var node *hg.Node
				if hg.FRand() > 0.5 {
					node = hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.Vec3GetOne(), hg.TranslationMat4(hg.RandomVec3WithMinMax(hg.NewVec3WithXYZ(-10, 18, -10), hg.NewVec3WithXYZ(10, 18, 10))), cube_ref, hg.GoSliceOfMaterial{mat_objects}, 1)
				} else {
					node = hg.CreatePhysicSphereWithSliceOfMaterialsMass(scene, 0.5, hg.TranslationMat4(hg.RandomVec3WithMinMax(hg.NewVec3WithXYZ(-10, 18, -10), hg.NewVec3WithXYZ(10, 18, 10))), sphere_ref, hg.GoSliceOfMaterial{mat_objects}, 1)
				}
				physics.NodeCreatePhysicsFromAssets(node) // update physics state

				physic_nodes = append(physic_nodes, node)
			}
		} else if state.Key(hg.KD) {
			if len(physic_nodes) > 0 {
				for i := 1; i < 8; i += 1 {
					node := physic_nodes[0]

					if node != nil {
						scene.DestroyNode(node)
						physic_nodes = physic_nodes[1:]
					}
				}
			}

			hg.SceneGarbageCollectSystemsWithPhysics(scene, physics)
		} else if state.Key(hg.KEscape) {
			break
		}

		// update scene and physic system, synchronize physics to scene, submit scene to draw
		dt := hg.TickClock()

		hg.SceneUpdateSystemsWithPhysicsStepMaxPhysicsStep(scene, clocks, dt, physics, hg.TimeFromSecF(1.0/60.0), 4)
		view_id := uint16(0)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontal(&view_id, scene, hg.NewIntRectWithSxSyExEy(0, 0, res_x, res_y), true, pipeline, res)

		// on-screen usage text
		hg.SetView2DWithZnearZfarFlagsColorDepthStencil(view_id, 0, 0, res_x, res_y, -1, 1, hg.CFDepth, hg.ColorGetBlack(), 1, 0)
		hg.DrawTextWithPosHalignValignSliceOfValuesSliceOfTexturesState(view_id, font, "S: Add object - D: Destruct object", font_program, "u_tex", 0, hg.Mat4GetIdentity(), hg.NewVec3WithXYZ(460, float32(res_y-60), 0), hg.DTHALeft, hg.DTVABottom, text_uniform_values, hg.GoSliceOfUniformSetTexture{}, text_render_state)
		hg.DrawTextWithPosHalignValignSliceOfValuesSliceOfTexturesState(view_id, font, fmt.Sprint("%d Object", len(physic_nodes)), font_program, "u_tex", 0, hg.Mat4GetIdentity(), hg.NewVec3WithXYZ(float32(res_x-200), float32(res_y-60), 0), hg.DTHALeft, hg.DTVABottom, text_uniform_values, hg.GoSliceOfUniformSetTexture{}, text_render_state)

		hg.Frame()
		hg.UpdateWindow(win)
	}

	hg.RenderShutdown()
	hg.DestroyWindow(win)

	hg.WindowSystemShutdown()
	hg.InputShutdown()
}
