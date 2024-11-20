from asyncio import create_task

import harfang as hg
import asyncio
import math

#make the cube rotate in async function
async def async_cube_rotation(rotating_cube):
	dt = 1.0/100.0
	rot = rotating_cube.GetRot()
	while True:
		rot.y += math.pi * dt / 2.0
		rotating_cube.SetRot(rot)
		await asyncio.sleep(dt)

#make cube translate in async function during the cube_rotation is running
async def async_cube_translation(translating_cube, steps=120, speed=1.0):
	dt = 1.0/100.0
	pos = translating_cube.GetPos()
	while True:
		for i in range(steps):
			pos.x += dt * speed
			translating_cube.SetPos(pos)
			await asyncio.sleep(dt)

		for i in range(steps):
			pos.z += dt * speed
			translating_cube.SetPos(pos)
			await asyncio.sleep(dt)

		for i in range(steps):
			pos.x -= dt * speed
			translating_cube.SetPos(pos)
			await asyncio.sleep(dt)

		for i in range(steps):
			pos.z -= dt * speed
			translating_cube.SetPos(pos)
			await asyncio.sleep(dt)

def create_material(diffuse, specular, self, prg_ref):
    mat = hg.CreateMaterial(prg_ref, 'uDiffuseColor', diffuse, 'uSpecularColor', specular)
    hg.SetMaterialValue(mat, 'uSelfColor', self)
    return mat

async def async_main():
    hg.InputInit()
    hg.WindowSystemInit()

    res_x, res_y = 1280, 720
    win = hg.RenderInit('Harfang - Scene Coroutine', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

    # Link precompiled assets folder to the project
    hg.AddAssetsFolder('resources_compiled')

    # Create Pipeline
    pipeline = hg.CreateForwardPipeline()
    res = hg.PipelineResources()

    # create materials
    vtx_layout = hg.VertexLayoutPosFloatNormUInt8()
    prg_ref = hg.LoadPipelineProgramRefFromAssets('core/shader/default.hps', res, hg.GetForwardPipelineInfo())
    mat_objects = create_material(hg.Vec4(0.5, 0.5, 0.5), hg.Vec4(1, 1, 1), hg.Vec4(0, 0, 0), prg_ref)

    # setup scene
    scene = hg.Scene()
    scene.canvas.color = hg.ColorI(22, 56, 76)
    scene.environment.fog_color = scene.canvas.color
    scene.environment.fog_near = 20
    scene.environment.fog_far = 80

    cam_mtx = hg.TransformationMat4(hg.Vec3(0, 10, -15), hg.Deg3(30, 0, 0))
    cam = hg.CreateCamera(scene, cam_mtx, 0.01, 3000)
    scene.SetCurrentCamera(cam)

    hg.CreateLinearLight(scene, hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Deg3(30, 59, 0)),
        hg.Color(1, 0.8, 0.7), hg.Color(1, 0.8, 0.7), 10, hg.LST_Map, 0.002, hg.Vec4(50, 100, 200, 400))
    hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(0, 10, 10)), 100, hg.ColorI(94, 155, 228),
        hg.ColorI(94, 255, 228))

    # create models
    rotating_cube_model = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
    rotating_cube_ref = res.AddModel('cube', rotating_cube_model)
    rotating_cube_instance = hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 0, 0)), rotating_cube_ref, [mat_objects])
    rotating_cube_transform = rotating_cube_instance.GetTransform()

    translating_cube_model = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
    translating_cube_ref = res.AddModel('cube', translating_cube_model)
    translating_cube_instance = hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 3, 0)), translating_cube_ref, [mat_objects])
    rotating_cube_tranform = translating_cube_instance.GetTransform()

    clocks = hg.SceneClocks()

    #set up physicss
    physics = hg.SceneBullet3Physics()
    physics.SceneCreatePhysicsFromAssets(scene)

    #create async tasks
    asyncio.create_task(async_cube_rotation(rotating_cube_transform))
    asyncio.create_task(async_cube_translation(rotating_cube_tranform, steps=100, speed=1.0))
    asyncio.get_event_loop()

    # main loop
    while not hg.ReadKeyboard().Key(hg.K_Escape) and hg.IsWindowOpen(win):
        dt = hg.TickClock()

        #update the scene
        scene.Update(dt)

        hg.SceneUpdateSystems(scene, clocks, dt, physics, hg.time_from_sec_f(1 / 60), 4)
        view_id, pass_id = hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), True, pipeline, res)

        #give async tasks to do their job
        await asyncio.sleep(0)

        hg.Frame()
        hg.UpdateWindow(win)

    hg.DestroyForwardPipeline(pipeline)
    hg.RenderShutdown()
    hg.DestroyWindow(win)


asyncio.run(async_main())


