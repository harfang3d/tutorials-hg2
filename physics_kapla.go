package main

import (
	"math"
	"runtime"
	"runtime/debug"

	hg "github.com/harfang3d/harfang-go/v3"
)

func fillRing(scn *hg.Scene, material *hg.Material, kaplaRef *hg.ModelRef, width float32, height float32, length float32, r float32, ringY float32, size float32, rAdjust float32, yOff float32, x float32, y float32, z float32) {
	step := math.Asin(float64((size*1.01)/2.0/(r-rAdjust))) * 2.0
	cubeCount := math.Floor((2.0 * math.Pi) / step)
	error := 2.0*math.Pi - step*cubeCount
	step += error / cubeCount // distribute error

	for a := 0.0; a < 2.0*math.Pi-error; a += step {
		world := hg.TransformationMat4(hg.NewVec3WithXYZ(float32(math.Cos(a))*r+x, ringY, float32(math.Sin(a))*r+z), hg.NewVec3WithXYZ(0, float32(-a)+yOff, 0))
		hg.CreatePhysicCubeWithSliceOfMaterialsMass(scn, hg.NewVec3WithXYZ(width, height, length), world, kaplaRef, hg.GoSliceOfMaterial{material}, 0.1)
	}
}

func addKaplaTower(scn *hg.Scene, ressources *hg.PipelineResources, width float32, height float32, length float32, radius float32, material *hg.Material, levelCount int, x float32, y float32, z float32, vtxLayout *hg.VertexLayout) {
	// Create a Kapla tower, return a list of created nodes
	levelY := y + height/2.0

	kaplaMdl := hg.CreateCubeModel(vtxLayout, width, height, length)
	kaplaRef := ressources.AddModel("kapla", kaplaMdl)

	for i := 0; i < levelCount/2; i++ {
		fillRing(scn, material, kaplaRef, width, height, length, radius-length/2.0, levelY, width, length/2.0, math.Pi/2.0, x, y, z)
		levelY += height
		fillRing(scn, material, kaplaRef, width, height, length, radius-length+width/2.0, levelY, length, width/2.0, 0, x, y, z)
		fillRing(scn, material, kaplaRef, width, height, length, radius-width/2.0, levelY, length, width/2.0, 0, x, y, z)
		levelY += height
	}
}

func main() {
	debug.SetGCPercent(-1)

	hg.InputInit()
	hg.WindowSystemInit()

	resX, resY := int32(1280), int32(720)

	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Kapla", resX, resY, hg.RFVSync|hg.RFMSAA4X)

	hg.AddAssetsFolder("resources_compiled")

	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	// create models
	vtxLayout := hg.VertexLayoutPosFloatNormUInt8()

	sphereMdl := hg.CreateSphereModel(vtxLayout, 0.5, 12, 24)
	sphereRef := res.AddModel("sphere", sphereMdl)

	// Create material
	prgRef := hg.LoadPipelineProgramRefFromAssets("core/shader/pbr.hps", res, hg.GetForwardPipelineInfo())

	matCube := hg.CreateMaterialWithValueName0Value0ValueName1Value1(prgRef, "uBaseOpacityColor", hg.Vec4I(255, 255, 56), "uOcclusionRoughnessMetalnessColor", hg.NewVec4WithXYZ(1, 0.658, 1))
	matGround := hg.CreateMaterialWithValueName0Value0ValueName1Value1(prgRef, "uBaseOpacityColor", hg.Vec4I(171, 255, 175), "uOcclusionRoughnessMetalnessColor", hg.NewVec4WithXYZ(1, 1, 1))
	matSpheres := hg.CreateMaterialWithValueName0Value0ValueName1Value1(prgRef, "uBaseOpacityColor", hg.Vec4I(255, 71, 75), "uOcclusionRoughnessMetalnessColor", hg.NewVec4WithXYZ(1, 0.5, 0.1))

	// Setup scene
	scene := hg.NewScene()
	scene.GetCanvas().SetColor(hg.ColorI(200, 210, 208))
	scene.GetEnvironment().SetAmbient(hg.ColorGetBlack())

	cam := hg.CreateCamera(scene, hg.Mat4GetIdentity(), 0.01, 1000)
	scene.SetCurrentCamera(cam)

	hg.CreateLinearLightWithDiffuseSpecularPriorityShadowTypeShadowBiasPssmSplit(scene, hg.TransformationMat4(hg.NewVec3WithXYZ(0, 0, 0), hg.Deg3(19, 59, 0)), hg.NewColorWithRGBA(1.5, 0.9, 1.2, 1), hg.NewColorWithRGBA(1.5, 0.9, 1.2, 1), 10, hg.LSTMap, 0.002, hg.NewVec4WithXYZW(8, 20, 40, 120))
	hg.CreatePointLightWithDiffuseSpecularPriority(scene, hg.TranslationMat4(hg.NewVec3WithXYZ(30, 20, 25)), 100, hg.NewColorWithRGBA(0.8, 0.5, 0.4, 1), hg.NewColorWithRGBA(0.8, 0.5, 0.4, 1), 0)

	mdlRef := res.AddModel("ground", hg.CreateCubeModel(vtxLayout, 200, 0.1, 200))
	hg.CreatePhysicCubeWithSliceOfMaterialsMass(scene, hg.NewVec3WithXYZ(200, 0.1, 200), hg.TranslationMat4(hg.NewVec3WithXYZ(0, -0.5, 0)), mdlRef, hg.GoSliceOfMaterial{matGround}, 0)

	addKaplaTower(scene, res, 0.5, 2, 2, 6, matCube, 12, -12, 0, 0, vtxLayout)
	addKaplaTower(scene, res, 0.5, 2, 2, 6, matCube, 12, 12, 0, 0, vtxLayout)

	clocks := hg.NewSceneClocks()

	// Input devices and fps controller states
	keyboard := hg.NewKeyboard()
	mouse := hg.NewMouse()

	camPos := hg.NewVec3WithXYZ(28.3, 31.8, 26.9)
	camRot := hg.NewVec3WithXYZ(0.6, -2.38, 0)

	// Setup physics
	physics := hg.NewSceneBullet3Physics()
	physics.SceneCreatePhysicsFromAssets(scene)
	physicsStep := hg.TimeFromSecF(1.0 / 60.0)

	// Main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		keyboard.Update()
		mouse.Update()

		dt := hg.TickClock()

		speed := float32(8.0)
		if keyboard.Down(hg.KLShift) {
			speed = 20.0
		}
		hg.FpsController(keyboard, mouse, camPos, camRot, speed, dt)

		cam.GetTransform().SetPos(camPos)
		cam.GetTransform().SetRot(camRot)

		if keyboard.Pressed(hg.KSpace) {
			node := hg.CreatePhysicSphereWithSliceOfMaterialsMass(scene, 0.5, hg.TranslationMat4(camPos), sphereRef, hg.GoSliceOfMaterial{matSpheres}, 0.5)
			physics.NodeCreatePhysicsFromAssets(node)
			physics.NodeAddImpulseWithWorldPos(node, hg.GetZWithM(cam.GetTransform().GetWorld()).MulWithK(25), camPos)
		}

		// Display scene
		viewID := uint16(0)
		hg.SceneUpdateSystemsWithPhysicsStepMaxPhysicsStep(scene, clocks, dt, physics, physicsStep, 1)
		hg.SubmitSceneToPipelineWithFovAxisIsHorizontal(&viewID, scene, hg.NewIntRectWithSxSyExEy(0, 0, resX, resY), true, pipeline, res)

		hg.Frame()
		hg.UpdateWindow(win)
		runtime.GC()
	}

	hg.DestroyForwardPipeline(pipeline)
	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
