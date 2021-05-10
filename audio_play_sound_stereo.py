# Play a mono sound with stereo panning

import harfang as hg
import math

hg.InputInit()
hg.AudioInit()

snd_ref = hg.LoadWAVSoundFile('resources_compiled/sounds/metro_announce.wav')  # WAV 44.1kHz 16bit mono
src_ref = hg.PlayStereo(snd_ref, hg.StereoSourceState(1, hg.SR_Loop))

angle = 0

while not hg.ReadKeyboard('raw').Key(hg.K_Escape):
	angle += hg.time_to_sec_f(hg.TickClock()) * 0.5
	hg.SetSourcePanning(src_ref, math.sin(angle))  # panning left = -1, panning right = 1

hg.AudioShutdown()
hg.InputShutdown()
