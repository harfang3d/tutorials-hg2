# Instantiating scenes

import harfang as hg

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Scene instances', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

hg.AddAssetsFolder("resources_compiled")

# rendering pipeline
pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

# load host scene
scene = hg.Scene()
hg.LoadSceneFromAssets('playground/playground.scn', scene, res, hg.GetForwardPipelineInfo())


# declare the biped actor class
class BipedActor:
	def __init__(self, pos):
		self.__node = hg.CreateInstanceFromAssets(scene, hg.Mat4.Identity, "biped/biped.scn", res, hg.GetForwardPipelineInfo())
		self.__node.GetTransform().SetPosRot(pos, hg.Deg3(0, hg.FRand(360), 0))

		self.__delay = 0
		self.__state = None
		self.__playing_anim_ref = None

	def __start_anim(self, name):
		anim = self.__node.GetInstanceSceneAnim(name)  # get instance specific animation
		if self.__playing_anim_ref is not None:
			scene.StopAnim(self.__playing_anim_ref)
		self.__playing_anim_ref = scene.PlayAnim(anim, hg.ALM_Loop)

	def update(self, dt):
		# check for state change
		self.__delay = self.__delay - dt

		if self.__delay <= 0:
			states = ['idle', 'walk', 'run']
			self.__state = states[hg.Rand(len(states))]
			self.__delay = self.__delay + hg.time_from_sec_f(hg.FRRand(2, 6))  # 2 to 6 seconds before next state change
			self.__start_anim(self.__state)

		# apply motion
		dt_sec_f = hg.time_to_sec_f(dt)

		transform = self.__node.GetTransform()
		pos, rot = transform.GetPosRot()

		if self.__state == 'walk':
			pos = pos - hg.GetZ(transform.GetWorld()) * hg.Mtr(1.15) * dt_sec_f  # 1.15 m/sec
			rot.y = rot.y + hg.Deg(50) * dt_sec_f
		elif self.__state == 'run':
			pos = pos - hg.GetZ(transform.GetWorld()) * hg.Mtr(4.5) * dt_sec_f  # 4.5 m/sec
			rot.y = rot.y - hg.Deg(70) * dt_sec_f

		# confine actor to playground
		pos = hg.Clamp(pos, hg.Vec3(-10, 0, -10), hg.Vec3(10, 0, 10))

		transform.SetPosRot(pos, rot)

	def destroy(self):
		scene.DestroyNode(self.__node)


# spawn initial actors
actors = []
for i in range(20):
	actors.append(BipedActor(hg.RandomVec3(hg.Vec3(-10, 0, -10), hg.Vec3(10, 0, 10))))

print('%d nodes in scene' % (scene.GetAllNodeCount()))

# main loop
keyboard = hg.Keyboard()

while not keyboard.Pressed(hg.K_Escape):
	_, res_x, res_y = hg.RenderResetToWindow(win, res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X | hg.RF_MaxAnisotropy)

	keyboard.Update()

	if keyboard.Pressed(hg.K_S):
		actors.append(BipedActor(hg.RandomVec3(hg.Vec3(-10, 0, -10), hg.Vec3(10, 0, 10))))
	if keyboard.Pressed(hg.K_D):
		if len(actors) > 0:
			actors.pop().destroy()
			scene.GarbageCollect()

	dt = hg.TickClock()

	for actor in actors:
		actor.update(dt)

	scene.Update(dt)

	view_state = hg.ComputePerspectiveViewState(hg.Mat4LookAt(hg.Vec3(0, 10, -14), hg.Vec3(0, 1, -4)), hg.Deg(45), 0.01, 1000, hg.ComputeAspectRatioX(res_x, res_y))
	vid, passId = hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), view_state, pipeline, res)

	hg.Frame()
	hg.UpdateWindow(win)

hg.RenderShutdown()
hg.DestroyWindow(win)
