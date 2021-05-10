-- Play a mono sound with stereo panning

hg = require("harfang")

hg.InputInit()
hg.AudioInit()

snd_ref = hg.LoadWAVSoundFile('resources_compiled/sounds/metro_announce.wav')  -- WAV 44.1kHz 16bit mono
src_ref = hg.PlayStereo(snd_ref, hg.StereoSourceState(1, hg.SR_Loop))

angle = 0

while not hg.ReadKeyboard('raw'):Key(hg.K_Escape) do
	angle = angle + hg.time_to_sec_f(hg.TickClock()) * 0.5
	hg.SetSourcePanning(src_ref, math.sin(angle))  -- panning left = -1, panning right = 1
end

hg.AudioShutdown()
hg.InputShutdown()
