package main

import (
	"math"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.AudioInit()

	hg.AddAssetsFolder("resources_compiled")

	snd_ref := hg.LoadWAVSoundAsset("sounds/metro_announce.wav")
	src_ref := hg.PlaySpatialized(hg.SoundRef(snd_ref), hg.NewSpatializedSourceStateWithMtxVolumeRepeat(hg.Mat4GetIdentity(), 1, hg.SRLoop))

	angle := 0.0

	for !hg.ReadKeyboard().Key(hg.KEscape) {
		dt := hg.TickClock()
		dt_sec_f := hg.TimeToSecF(dt) // delta frame in seconds

		// compute the source old and new position in world
		src_old_pos := hg.NewVec3WithXYZ(float32(math.Sin(angle)), 0, float32(math.Cos(angle))).MulWithK(5.0)
		angle += float64(hg.TimeToSecF(hg.TickClock())) * 45.0
		src_new_pos := hg.NewVec3WithXYZ(float32(math.Sin(angle)), 0, float32(math.Cos(angle))).MulWithK(5.0)

		// source velocity in m.s-1
		src_vel := (src_new_pos.Sub(src_old_pos)).DivWithK(dt_sec_f)

		// update source properties
		hg.SetSourceTransform(hg.SourceRef(src_ref), hg.TranslationMat4(src_new_pos), src_vel)
	}
	hg.AudioShutdown()
	hg.InputShutdown()
}
