# Spatialized sound

import harfang as hg
import math

hg.InputInit()
hg.AudioInit()

hg.AddAssetsFolder("resources_compiled")

snd_ref = hg.LoadWAVSoundAsset("sounds/metro_announce.wav")
src_ref = hg.PlaySpatialized(snd_ref, hg.SpatializedSourceState(hg.Mat4.Identity, 1, hg.SR_Loop))

angle = 0

while not hg.ReadKeyboard('raw').Key(hg.K_Escape):
	dt = hg.TickClock()
	dt_sec_f = hg.time_to_sec_f(dt)  # delta frame in seconds

	# compute the source old and new position in world
	src_old_pos = hg.Vec3(math.sin(angle), 0, math.cos(angle)) * 5
	angle += hg.time_to_sec_f(hg.TickClock()) * 45
	src_new_pos = hg.Vec3(math.sin(angle), 0, math.cos(angle)) * 5

	# source velocity in m.s-1
	src_vel = (src_new_pos - src_old_pos) / dt_sec_f

	# update source properties
	hg.SetSourceTransform(src_ref, hg.TranslationMat4(src_new_pos), src_vel)

hg.AudioShutdown()
hg.InputShutdown()
