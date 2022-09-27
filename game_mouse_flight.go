package main

import (
	"math"

	hg "github.com/harfang3d/harfang-go/v3"
)

// gameplay settings
var (
	settingCameraChaseOffset   = hg.NewVec3WithXYZ(0, 0.2, 0)
	settingCameraChaseDistance = float32(1)

	settingPlaneSpeed            = float32(0.05)
	settingPlaneMouseSensitivity = float32(0.5)
	vtxLayout                    = hg.VertexLayoutPosFloatColorFloat()
	draw2DProgram                *hg.ProgramHandle
	draw2DRenderState            *hg.RenderState

	cameraNode *hg.Node
	planeNode  *hg.Node
)

func draw_circle(viewId uint16, center *hg.Vec3, radius float32, color *hg.Color) {
	segmentCount := 32.0
	step := 2.0 * math.Pi / segmentCount
	p0 := hg.NewVec3WithXYZ(center.GetX()+radius, center.GetY(), 0)
	p1 := hg.NewVec3WithXYZ(0, 0, 0)

	vtx := hg.NewVertices(vtxLayout, int32(segmentCount*2+2))

	for i := int32(0); i < int32(segmentCount)+1; i += 1 {
		p1.SetX(radius*float32(math.Cos(float64(i)*step)) + center.GetX())
		p1.SetY(radius*float32(math.Sin(float64(i)*step)) + center.GetY())
		vtx.Begin(2 * i).SetPos(p0).SetColor0(color).End()
		vtx.Begin(2*i + 1).SetPos(p1).SetColor0(color).End()
		p0.SetX(p1.GetX())
		p0.SetY(p1.GetY())
	}

	hg.DrawLinesWithRenderState(viewId, vtx, draw2DProgram, draw2DRenderState)

}

func updatePlane(mouseXNormd float32, mouseYNormd float32) {
	planeTransform := planeNode.GetTransform()

	planePos := planeTransform.GetPos()
	planePos = planePos.Add(hg.NormalizeWithVec3V(hg.GetZWithM(planeTransform.GetWorld())).MulWithK(settingPlaneSpeed))
	planePos.SetY(hg.ClampWithV(planePos.GetY(), 0.1, 50)) // floor/ceiling

	planeRot := planeTransform.GetRot()

	nextPlaneRot := hg.NewVec3WithVec3V(planeRot) // make a copy of the plane rotation
	nextPlaneRot.SetX(hg.ClampWithV(nextPlaneRot.GetX()+mouseYNormd*-0.03, -0.75, 0.75))
	nextPlaneRot.SetY(nextPlaneRot.GetY() + mouseXNormd*0.03)
	nextPlaneRot.SetZ(hg.ClampWithV(mouseXNormd*-0.75, -1.2, 1.2))

	planeRot = planeRot.Add((nextPlaneRot.Sub(planeRot)).MulWithK(settingPlaneMouseSensitivity))

	planeTransform.SetPos(planePos)
	planeTransform.SetRot(planeRot)

}

func updateChaseCamera(targetPos *hg.Vec3) {
	cameraTransform := cameraNode.GetTransform()
	cameraToTarget := hg.NormalizeWithVec3V(targetPos.Sub(cameraTransform.GetPos()))

	cameraTransform.SetPos(targetPos.Sub(cameraToTarget.MulWithK(settingCameraChaseDistance))) // camera is "distance" away from its target
	cameraTransform.SetRot(hg.ToEulerWithM(hg.Mat3LookAt(cameraToTarget)))

}

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	resX, resY := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Mouse Flight", resX, resY, hg.RFVSync|hg.RFMSAA8X)

	res := hg.NewPipelineResources()
	pipeline := hg.CreateForwardPipeline()

	keyboard := hg.NewKeyboard()
	mouse := hg.NewMouse()

	// access to compiled resources
	hg.AddAssetsFolder("resources_compiled")

	// 2D drawing helpers

	draw2DProgram = hg.LoadProgramFromAssets("shaders/pos_rgb")
	draw2DRenderState = hg.ComputeRenderStateWithDepthTestCulling(hg.BMAlpha, hg.DTLess, hg.FCDisabled)

	// setup game world
	scene := hg.NewScene()
	hg.LoadSceneFromAssets("playground/playground.scn", scene, res, hg.GetForwardPipelineInfo())

	planeNode, _ = hg.CreateInstanceFromAssets(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 4, 0)), "paper_plane/paper_plane.scn", res, hg.GetForwardPipelineInfo())
	cameraNode = hg.CreateCamera(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(0, 4, -5)), 0.01, 1000)

	scene.SetCurrentCamera(cameraNode)

	// game loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		dt := hg.TickClock() // tick clock, retrieve elapsed clock since last call

		// update mouse/keyboard devices
		keyboard.Update()
		mouse.Update()

		// compute ratio corrected normalized mouse position
		mouse_x, mouse_y := mouse.X(), mouse.Y()

		aspect_ratio := hg.ComputeAspectRatioX(float32(resX), float32(resY))
		mouseXNormd, mouseYNormd := (float32(mouse_x)/float32(resX)-0.5)*aspect_ratio.GetX(), (float32(mouse_y)/float32(resY)-0.5)*aspect_ratio.GetY()

		// update gameplay elements (plane & camera)
		updatePlane(mouseXNormd, mouseYNormd)
		updateChaseCamera(hg.NewVec3WithVec4V(planeNode.GetTransform().GetWorld().MulWithVec4V(hg.NewVec4WithVec3V(settingCameraChaseOffset))))

		// update scene and submit it to render pipeline
		scene.Update(dt)

		viewId := uint16(0)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontal(&viewId, scene, hg.NewIntRectWithSxSyExEy(0, 0, resX, resY), true, pipeline, res)

		// draw 2D GUI
		hg.SetView2DWithZnearZfarFlagsColorDepthStencilYUp(viewId, 0, 0, resX, resY, -1, 1, hg.CFDepth, hg.ColorGetBlack(), 1, 0, true)
		draw_circle(viewId, hg.NewVec3WithXYZ(float32(mouse_x), float32(mouse_y), 0), 20, hg.ColorGetWhite()) // display mouse cursor

		// end of frame
		hg.Frame()
		hg.UpdateWindow(win)

	}

	hg.RenderShutdown()
	hg.DestroyWindow(win)

	hg.WindowSystemShutdown()
	hg.InputShutdown()
}
