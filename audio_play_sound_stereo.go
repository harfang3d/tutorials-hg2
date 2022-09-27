package main

import (
	"math"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.AudioInit()

	snd_ref := hg.LoadWAVSoundFile("resources_compiled/sounds/metro_announce.wav") // WAV 44.1kHz 16bit mono
	src_ref := hg.PlayStereo(hg.SoundRef(snd_ref), hg.NewStereoSourceStateWithVolumeRepeat(1, hg.SRLoop))

	angle := 0.0

	for !hg.ReadKeyboardWithName("raw").Key(hg.KEscape) {
		angle += float64(hg.TimeToSecF(hg.TickClock())) * 0.5
		hg.SetSourcePanning(hg.SourceRef(src_ref), float32(math.Sin(angle))) // panning left = -1, panning right = 1
	}
	hg.AudioShutdown()
	hg.InputShutdown()
}
