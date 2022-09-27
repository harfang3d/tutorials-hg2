package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

// declare the biped actor class
type BipedActor struct {
	__node             *hg.Node
	__scene            *hg.Scene
	__delay            int64
	__state            string
	__playing_anim_ref *hg.ScenePlayAnimRef
}

func NewBipedActor(scene *hg.Scene, res *hg.PipelineResources, pos *hg.Vec3) BipedActor {
	b := BipedActor{}

	b.__node, _ = hg.CreateInstanceFromAssets(scene, hg.Mat4GetIdentity(), "biped/biped.scn", res, hg.GetForwardPipelineInfo())
	b.__node.GetTransform().SetPosRot(pos, hg.Deg3(0, hg.FRandWithRange(360), 0))

	b.__scene = scene
	b.__delay = 0
	b.__state = ""
	b.__playing_anim_ref = nil

	return b
}

func (b *BipedActor) __start_anim(name string) {
	anim := b.__node.GetInstanceSceneAnim(name) // get instance specific animation
	if b.__playing_anim_ref != nil {
		b.__scene.StopAnim(b.__playing_anim_ref)
	}
	b.__playing_anim_ref = b.__scene.PlayAnimWithLoopMode(anim, hg.ALMLoop)
}
func (b *BipedActor) update(dt int64) {
	// check for state change
	b.__delay = b.__delay - dt

	if b.__delay <= 0 {
		states := []string{"idle", "walk", "run"}
		b.__state = states[hg.RandWithRange(uint32(len(states)))]
		b.__delay = b.__delay + hg.TimeFromSecF(hg.FRRandWithRangeStartRangeEnd(2, 6)) // 2 to 6 seconds before next state change
		b.__start_anim(b.__state)
	}
	// apply motion
	dt_sec_f := hg.TimeToSecF(dt)

	transform := b.__node.GetTransform()
	pos, rot := transform.GetPosRot()

	if b.__state == "walk" {
		pos = pos.Sub(hg.GetZWithM(transform.GetWorld()).MulWithK(hg.Mtr(1.15) * dt_sec_f)) // 1.15 m/sec
		rot.SetY(rot.GetY() + hg.Deg(50)*dt_sec_f)
	} else if b.__state == "run" {
		pos = pos.Sub(hg.GetZWithM(transform.GetWorld()).MulWithK(hg.Mtr(4.5) * dt_sec_f)) // 4.5 m/sec
		rot.SetY(rot.GetY() - hg.Deg(70)*dt_sec_f)
	}

	// confine actor to playground
	pos = hg.ClampWithMinMax(pos, hg.NewVec3WithXYZ(-10, 0, -10), hg.NewVec3WithXYZ(10, 0, 10))

	transform.SetPosRot(pos, rot)
}

func (b *BipedActor) destroy() {
	b.__scene.DestroyNode(b.__node)
}

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Scene instances", res_x, res_y, hg.RFVSync|hg.RFMSAA4X)

	hg.AddAssetsFolder("resources_compiled")

	// rendering pipeline
	pipeline := hg.CreateForwardPipeline()
	res := hg.NewPipelineResources()

	// load host scene
	scene := hg.NewScene()
	hg.LoadSceneFromAssets("playground/playground.scn", scene, res, hg.GetForwardPipelineInfo())

	// spawn initial actors
	actors := []BipedActor{}
	for i := 0; i < 20; i += 1 {
		actors = append(actors, NewBipedActor(scene, res, hg.RandomVec3WithMinMax(hg.NewVec3WithXYZ(-10, 0, -10), hg.NewVec3WithXYZ(10, 0, 10))))
	}

	fmt.Printf("%d nodes in scene\n", (scene.GetAllNodeCount()))

	// main loop
	keyboard := hg.NewKeyboard()

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		hg.RenderResetToWindowWithResetFlags(win, &res_x, &res_y, uint32(hg.RFVSync|hg.RFMSAA4X|hg.RFMaxAnisotropy))

		keyboard.Update()

		if keyboard.Pressed(hg.KS) {
			actors = append(actors, NewBipedActor(scene, res, hg.RandomVec3WithMinMax(hg.NewVec3WithXYZ(-10, 0, -10), hg.NewVec3WithXYZ(10, 0, 10))))
		}
		if keyboard.Pressed(hg.KD) {
			if len(actors) > 0 {
				actors[0].destroy()
				actors = actors[1:]
				scene.GarbageCollect()
			}
		}
		dt := hg.TickClock()

		for id := range actors {
			actors[id].update(dt)
		}

		scene.Update(dt)

		view_state := hg.ComputePerspectiveViewState(hg.Mat4LookAt(hg.NewVec3WithXYZ(0, 10, -14), hg.NewVec3WithXYZ(0, 1, -4)), hg.Deg(45), 0.01, 1000, hg.ComputeAspectRatioX(float32(res_x), float32(res_y)))
		viewId := uint16(0)
		hg.SubmitSceneToPipeline(&viewId, scene, hg.NewIntRectWithSxSyExEy(0, 0, res_x, res_y), view_state, pipeline, res)

		hg.Frame()
		hg.UpdateWindow(win)
	}

	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
