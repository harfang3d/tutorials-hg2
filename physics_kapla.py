# Physics kapla towers

import harfang as hg
from math import pi, cos, sin, asin

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Kapla - Press SPACEBAR', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

hg.AddAssetsFolder('resources_compiled')

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

# create models
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

sphere_mdl = hg.CreateSphereModel(vtx_layout, 0.5, 12, 24)
sphere_ref = res.AddModel('sphere', sphere_mdl)

# create materials
prg_ref = hg.LoadPipelineProgramRefFromAssets('core/shader/pbr.hps', res, hg.GetForwardPipelineInfo())

mat_cube = hg.CreateMaterial(prg_ref, 'uBaseOpacityColor', hg.Vec4I(255, 255, 56), 'uOcclusionRoughnessMetalnessColor', hg.Vec4(1, 0.658, 1))
mat_ground = hg.CreateMaterial(prg_ref, 'uBaseOpacityColor', hg.Vec4I(171, 255, 175), 'uOcclusionRoughnessMetalnessColor', hg.Vec4(1, 1, 1))
mat_spheres = hg.CreateMaterial(prg_ref, 'uBaseOpacityColor', hg.Vec4I(255, 71, 75), 'uOcclusionRoughnessMetalnessColor', hg.Vec4(1, 0.5, 0.1))

# setup scene
scene = hg.Scene()
scene.canvas.color = hg.ColorI(200, 210, 208)
scene.environment.ambient = hg.Color.Black

cam = hg.CreateCamera(scene, hg.Mat4.Identity, 0.01, 1000)
scene.SetCurrentCamera(cam)

lgt = hg.CreateLinearLight(scene, hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Deg3(19, 59, 0)), hg.Color(1.5, 0.9, 1.2, 1), hg.Color(1.5, 0.9, 1.2, 1), 10, hg.LST_Map, 0.002, hg.Vec4(8, 20, 40, 120))
back_lgt = hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(30, 20, 25)), 100, hg.Color(0.8, 0.5, 0.4, 1), hg.Color(0.8, 0.5, 0.4, 1), 0)

mdl_ref = res.AddModel('ground', hg.CreateCubeModel(vtx_layout, 200, 0.1, 200))
hg.CreatePhysicCube(scene, hg.Vec3(200, 0.1, 200), hg.TranslationMat4(hg.Vec3(0, -0.5, 0)), mdl_ref, [mat_ground], 0)


def add_kapla_tower(scn, resources, width, height, length, radius, material, level_count, x, y, z):
	"""Create a Kapla tower, return a list of created nodes"""
	level_y = y + height / 2

	kapla_mdl = hg.CreateCubeModel(vtx_layout, width, height, length)
	kapla_ref = resources.AddModel('kapla', kapla_mdl)

	nodes = []

	for i in range(level_count // 2):
		def fill_ring(r, ring_y, size, r_adjust, y_off):
			step = asin((size * 1.01) / 2 / (r - r_adjust)) * 2
			cube_count = (2 * pi) // step
			error = 2 * pi - step * cube_count
			step += error / cube_count  # distribute error

			a = 0
			while a < (2 * pi - error):
				world = hg.TransformationMat4(hg.Vec3(cos(a) * r + x, ring_y, sin(a) * r + z), hg.Vec3(0, -a + y_off, 0))
				node = hg.CreatePhysicCube(scn, hg.Vec3(width, height, length), world, kapla_ref, [material], 0.1)
				nodes.append(node)
				a += step

		fill_ring(radius - length / 2, level_y, width, length / 2, pi / 2)
		level_y += height
		fill_ring(radius - length + width / 2, level_y, length, width / 2, 0)
		fill_ring(radius - width / 2, level_y, length, width / 2, 0)
		level_y += height

	return nodes


add_kapla_tower(scene, res, 0.5, 2, 2, 6, mat_cube, 12, -12, 0, 0)
add_kapla_tower(scene, res, 0.5, 2, 2, 6, mat_cube, 12, 12, 0, 0)

clocks = hg.SceneClocks()

# input devices and fps controller states
keyboard = hg.Keyboard()
mouse = hg.Mouse()

cam_pos = hg.Vec3(28.3, 31.8, 26.9)
cam_rot = hg.Vec3(0.6, -2.38, 0)

# setup physics
physics = hg.SceneBullet3Physics()
physics.SceneCreatePhysicsFromAssets(scene)
physics_step = hg.time_from_sec_f(1 / 60)

# main loop
while not keyboard.Down(hg.K_Escape):
	keyboard.Update()
	mouse.Update()

	dt = hg.TickClock()

	hg.FpsController(keyboard, mouse, cam_pos, cam_rot, 20 if keyboard.Down(hg.K_LShift) else 8, dt)

	cam.GetTransform().SetPos(cam_pos)
	cam.GetTransform().SetRot(cam_rot)

	if keyboard.Pressed(hg.K_Space):
		node = hg.CreatePhysicSphere(scene, 0.5, hg.TranslationMat4(cam_pos), sphere_ref, [mat_spheres], 0.5)
		physics.NodeCreatePhysicsFromAssets(node)
		physics.NodeAddImpulse(node, hg.GetZ(cam.GetTransform().GetWorld()) * 25.0, cam_pos)

	hg.SceneUpdateSystems(scene, clocks, dt, physics, physics_step, 1)
	hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), True, pipeline, res)

	hg.Frame()
	hg.UpdateWindow(win)

hg.DestroyForwardPipeline(pipeline)

hg.RenderShutdown()
hg.DestroyWindow(win)
