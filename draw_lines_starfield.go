package main

import (
	"runtime"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	width, height := int32(1280), int32(720)
	window := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Starfield", width, height, hg.RFVSync|hg.RFMSAA4X)

	// vertex layout
	vtxLayout := hg.NewVertexLayout()
	vtxLayout.Begin()
	vtxLayout.Add(hg.APosition, 3, hg.ATFloat)
	vtxLayout.Add(hg.AColor0, 3, hg.ATFloat)
	vtxLayout.End()

	// simple shader program
	shader := hg.LoadProgramFromFile("resources_compiled/shaders/pos_rgb")

	// initialize stars
	starfieldSize := float32(10.0)

	maxStars := 1000
	vtx := hg.NewVertices(vtxLayout, int32(maxStars*2))

	var stars []*hg.Vec3
	for i := 0; i < maxStars; i++ {
		stars = append(stars, hg.RandomVec3(-starfieldSize, starfieldSize))
	}

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(window) {
		hg.SetViewClearWithColDepthStencil(0, hg.CFColor|hg.CFDepth, hg.ColorGetBlack(), 1.0, 0)
		hg.SetViewRect(0, 0, 0, uint16(width), uint16(height))

		dt := hg.TickClock()
		dt_f := hg.TimeToSecF(dt)

		// update stars
		vtx.Clear()
		for i, star := range stars {
			star.SetZ(star.GetZ() - float32(2.0*dt_f))
			if star.GetZ() < starfieldSize {
				star.SetZ(star.GetZ() + starfieldSize)
			}

			// draw stars
			vtx.Begin(int32(2 * i)).SetPos(star.Mul(hg.NewVec3WithXYZ(1.0/star.GetZ(), 1.0/star.GetZ(), 0.0))).SetColor0(hg.ColorGetBlack()).End()
			vtx.Begin(int32(2*i + 1)).SetPos(star.Mul(hg.NewVec3WithXYZ(1.04/star.GetZ(), 1.04/star.GetZ(), 0.0))).SetColor0(hg.ColorGetWhite()).End()
		}

		hg.DrawLines(0, vtx, shader)

		hg.Frame()
		hg.UpdateWindow(window)
		runtime.GC()
	}

	hg.RenderShutdown()
	hg.DestroyWindow(window)

}
