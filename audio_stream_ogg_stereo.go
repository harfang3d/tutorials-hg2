package main

import (
	"math"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.AudioInit()

	src_ref := hg.StreamOGGFileStereo("resources_compiled/sounds/metro_announce.ogg", hg.NewStereoSourceStateWithVolumeRepeat(1, hg.SRLoop)) // OGG 44.1kHz 16bit mono

	angle := 0.0

	for !hg.ReadKeyboardWithName("raw").Key(hg.KEscape) {
		angle += float64(hg.TimeToSecF(hg.TickClock())) * 0.5
		hg.SetSourcePanning(hg.SourceRef(src_ref), float32(math.Sin(angle))) // panning left = -1, panning right = 1
	}
	hg.AudioShutdown()
	hg.InputShutdown()
}
