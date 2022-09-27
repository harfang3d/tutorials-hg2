package main

import (
	"math"
	"runtime"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	resX, resY := int32(1280), int32(720)

	window := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Draw Lines", resX, resY, hg.RFVSync|hg.RFMSAA4X)

	lineCount := int32(1000)

	shader := hg.LoadProgramFromFile("resources_compiled/shaders/white")

	// vertices
	vtxLayout := hg.NewVertexLayout()
	vtxLayout.Begin()
	vtxLayout.Add(hg.APosition, 3, hg.ATFloat)
	vtxLayout.End()

	vtx := hg.NewVertices(vtxLayout, lineCount*2)

	// main loop
	angle := 0.0

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(window) {
		hg.SetViewClearWithColDepthStencil(0, hg.CFColor|hg.CFDepth, hg.ColorI(64, 64, 64), 1.0, 0)
		hg.SetViewRect(0, 0, 0, uint16(resX), uint16(resY))

		vtx.Clear()
		for i := 0.0; i < float64(lineCount); i++ {
			vtx.Begin(int32(2 * i)).SetPos(hg.NewVec3WithXYZ(float32(math.Sin(angle+i*0.005)), float32(math.Cos(angle+i*0.01)), 0.0)).End()
			vtx.Begin(int32(2*i + 1)).SetPos(hg.NewVec3WithXYZ(float32(math.Sin(angle+i*-0.005)), float32(math.Cos(angle+i*0.005)), 0.0)).End()
		}

		hg.DrawLines(0, vtx, shader) // submit all lines in a single call

		angle = angle + float64(hg.TimeToSecF(hg.TickClock()))

		hg.Frame()
		hg.UpdateWindow(window)
		runtime.GC()
	}

	hg.RenderShutdown()
	hg.DestroyWindow(window)

}
