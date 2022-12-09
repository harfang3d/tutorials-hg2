package main

import (
	"math"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1024), int32(1024)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Draw Scene to Texture", res_x, res_y, hg.RFVSync|hg.RFMSAA8X)

	hg.AddAssetsFolder("resources_compiled")

	// create pipeline
	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	// load the scene to draw to a texture
	scene := hg.NewScene()
	hg.LoadSceneFromAssets("materials/materials.scn", scene, res, hg.GetForwardPipelineInfo())

	// create a 512x512 frame buffer to draw the scene to
	frame_buffer := hg.CreateFrameBufferWithWidthHeightColorFormatDepthFormatAaName(512, 512, hg.TFRGBA32F, hg.TFD24, 4, "framebuffer") // 4x MSAA
	color := hg.GetColorTexture(frame_buffer)

	// create the cube model
	vtx_layout := hg.VertexLayoutPosFloatTexCoord0UInt8()

	cube_mdl := hg.CreateCubeModel(vtx_layout, 1, 1, 1)
	res.AddModel("cube", cube_mdl)

	// prepare the cube shader program
	cube_prg := hg.LoadProgramFromAssets("shaders/texture")

	// main loop
	angle := 0.0

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		dt := hg.TickClock()
		angle = angle + float64(hg.TimeToSecF(dt))

		// update scene and render to the frame buffer
		scene.GetCurrentCamera().GetTransform().SetPos(hg.NewVec3WithXYZ(0, 0, float32(-(math.Sin(angle)*3.0 + 4.0)))) // animate the scene current camera on Z

		scene.Update(dt)

		view_id := uint16(0)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontalFb(&view_id, scene, hg.NewIntRectWithSxSyExEy(0, 0, 512, 512), true, pipeline, res, frame_buffer.GetHandle())

		// draw a rotating cube in immediate mode using the texture the scene was rendered to
		hg.SetViewPerspective(view_id, 0, 0, res_x, res_y, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0, -1.8)))

		val_uniforms := hg.GoSliceOfUniformSetValue{hg.MakeUniformSetValueWithVec4V("color", hg.NewVec4WithXYZW(1, 1, 1, 1))} // note: these could be moved out of the main loop but are kept here for readability
		tex_uniforms := hg.GoSliceOfUniformSetTexture{hg.MakeUniformSetTexture("s_tex", color, 0)}

		hg.DrawModelWithSliceOfValuesSliceOfTextures(view_id, cube_mdl, cube_prg, val_uniforms, tex_uniforms, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 0, 0), hg.NewVec3WithXYZ(float32(angle*0.1), float32(angle*0.05), float32(angle*0.2))))

		//
		hg.Frame()
		hg.UpdateWindow(win)
	}
	hg.RenderShutdown()
	hg.WindowSystemShutdown()
}
