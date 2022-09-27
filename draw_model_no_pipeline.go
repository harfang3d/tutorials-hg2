package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	resX, resY := int32(1280), int32(720)

	window := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Draw Models no Pipeline", resX, resY, hg.RFVSync|hg.RFMSAA4X)

	// vertex layout and models
	vtxLayout := hg.VertexLayoutPosFloatNormUInt8()

	cubeMdl := hg.CreateCubeModel(vtxLayout, 1, 1, 1)
	groundMdl := hg.CreatePlaneModel(vtxLayout, 5, 5, 1, 1)

	shader := hg.LoadProgramFromFile("resources_compiled/shaders/mdl")

	// main loop
	angle := float32(0.0)

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(window) {
		dt := hg.TickClock()
		angle = angle + hg.TimeToSecF(dt)

		viewpoint := hg.TranslationMat4(hg.NewVec3WithXYZ(0, 1, -3))
		hg.SetViewPerspective(0, 0, 0, resX, resY, viewpoint)

		hg.DrawModelWithSliceOfValuesSliceOfTextures(0, cubeMdl, shader, hg.GoSliceOfUniformSetValue{}, hg.GoSliceOfUniformSetTexture{}, hg.TransformationMat4(hg.NewVec3WithXYZ(0.0, 1, 0.0), hg.NewVec3WithXYZ(angle, angle, angle)))
		hg.DrawModelWithSliceOfValuesSliceOfTextures(0, groundMdl, shader, hg.GoSliceOfUniformSetValue{}, hg.GoSliceOfUniformSetTexture{}, hg.TranslationMat4(hg.NewVec3WithXYZ(0.0, 0.0, 0.0)))

		hg.Frame()
		hg.UpdateWindow(window)
	}

	hg.RenderShutdown()
	hg.DestroyWindow(window)
}
