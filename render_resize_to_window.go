package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(512), int32(512)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Render Resize to Window", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	// create model
	vtx_layout := hg.VertexLayoutPosFloatNormUInt8()

	cube_mdl := hg.CreateCubeModel(vtx_layout, 1, 1, 1)
	cube_prg := hg.LoadProgramFromFile("resources_compiled/shaders/mdl")

	// main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		render_was_reset := hg.RenderResetToWindowWithResetFlags(win, &res_x, &res_y, uint32(hg.RFVSync|hg.RFMSAA4X|hg.RFMaxAnisotropy))
		if render_was_reset {
			fmt.Printf("Render reset to %dx%d", res_x, res_y)
		}

		viewpoint := hg.TransformationMat4(hg.NewVec3WithXYZ(1, 1, -2), hg.Deg3(24, -27, 0))
		hg.SetViewPerspectiveWithZnearZfarZoomFactorFlagsColorDepthStencil(0, 0, 0, res_x, res_y, viewpoint, 0.01, 100, 1.8, hg.CFColor|hg.CFDepth, hg.ColorI(64, 64, 64), 1, 0)

		hg.DrawModelWithSliceOfValuesSliceOfTextures(0, cube_mdl, cube_prg, hg.GoSliceOfUniformSetValue{}, hg.GoSliceOfUniformSetTexture{}, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 0, 0)))

		hg.Frame()
		hg.UpdateWindow(win)
	}

	hg.RenderShutdown()
}
