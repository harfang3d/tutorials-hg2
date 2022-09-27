package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	resX, resY := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Draw Text", resX, resY, hg.RFVSync)

	// load font and shader program
	font := hg.LoadFontFromFileWithSize("resources_compiled/font/default.ttf", 96)
	font_prg := hg.LoadProgramFromFile("resources_compiled/core/shader/font")

	// text uniforms and render state
	text_uniform_values := hg.GoSliceOfUniformSetValue{hg.MakeUniformSetValueWithVec4V("u_color", hg.NewVec4WithXYZ(1, 1, 0))}
	text_render_state := hg.ComputeRenderStateWithDepthTestCulling(hg.BMAlpha, hg.DTAlways, hg.FCDisabled)

	// main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		hg.SetView2DWithZnearZfarFlagsColorDepthStencil(0, 0, 0, resX, resY, -1, 1, hg.CFColor|hg.CFDepth, hg.ColorI(32, 32, 32), 0, 1)

		hg.DrawTextWithPosHalignValignSliceOfValuesSliceOfTexturesState(0, font, "Hello world!", font_prg, "u_tex", 0, hg.Mat4GetIdentity(), hg.NewVec3WithXYZ(float32(resX/2), float32(resY/2), 0), hg.DTHACenter, hg.DTVACenter, text_uniform_values, hg.GoSliceOfUniformSetTexture{}, text_render_state)

		hg.Frame()
		hg.UpdateWindow(win)
	}

	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
