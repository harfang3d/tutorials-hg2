package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	resX, resY := int32(1280), int32(720)

	window := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Draw Text over Models", resX, resY, hg.RFVSync|hg.RFMSAA4X)

	// vertex layout and models
	vtxLayout := hg.VertexLayoutPosFloatNormUInt8()

	cubeMdl := hg.CreateCubeModel(vtxLayout, 1, 1, 1)
	groundMdl := hg.CreatePlaneModel(vtxLayout, 5, 5, 1, 1)

	shader := hg.LoadProgramFromFile("resources_compiled/shaders/mdl")

	// load font and shader program
	font := hg.LoadFontFromFileWithSize("resources_compiled/font/default.ttf", 96)
	font_prg := hg.LoadProgramFromFile("resources_compiled/core/shader/font")

	// text uniforms and render state
	text_uniform_values := hg.GoSliceOfUniformSetValue{hg.MakeUniformSetValueWithVec4V("u_color", hg.NewVec4WithXYZ(1, 1, 0))}
	text_render_state := hg.ComputeRenderStateWithDepthTestCulling(hg.BMAlpha, hg.DTAlways, hg.FCDisabled)

	// main loop
	angle := float32(0.0)

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(window) {
		dt := hg.TickClock()
		angle = angle + hg.TimeToSecF(dt)

		// 3D view
		viewpoint := hg.TranslationMat4(hg.NewVec3WithXYZ(0, 1, -3))
		hg.SetViewPerspectiveWithZnearZfar(0, 0, 0, resX, resY, viewpoint, 0.01, 5000)

		hg.DrawModelWithSliceOfValuesSliceOfTextures(0, cubeMdl, shader, hg.GoSliceOfUniformSetValue{}, hg.GoSliceOfUniformSetTexture{}, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 1, 0), hg.NewVec3WithXYZ(angle, angle, angle)))
		hg.DrawModelWithSliceOfValuesSliceOfTextures(0, groundMdl, shader, hg.GoSliceOfUniformSetValue{}, hg.GoSliceOfUniformSetTexture{}, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0, 0)))

		// 2D view, note that only the depth buffer is cleared
		hg.SetView2DWithZnearZfarFlagsColorDepthStencil(1, 0, 0, resX, resY, -1, 1, hg.CFDepth, hg.ColorI(32, 32, 32), 1, 0)

		hg.DrawTextWithPosHalignValignSliceOfValuesSliceOfTexturesState(1, font, "Hello world!", font_prg, "u_tex", 0, hg.Mat4GetIdentity(), hg.NewVec3WithXYZ(float32(resX/2), float32(resY/2), 0), hg.DTHACenter, hg.DTVACenter, text_uniform_values, hg.GoSliceOfUniformSetTexture{}, text_render_state)

		hg.Frame()
		hg.UpdateWindow(window)
	}

	hg.RenderShutdown()
	hg.DestroyWindow(window)
}
