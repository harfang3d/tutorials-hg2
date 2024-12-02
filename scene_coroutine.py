import harfang as hg
import asyncio
import math

# make the cube rotate in async function
async def async_cube_rotation(cube):
    global dt
    rot = cube.GetRot()
    while True:
        #need to convert dt to seconds to have a decent rotation
        dt_sec = hg.time_to_sec_f(dt) * 2
        rot.y += math.pi * dt_sec
        cube.SetRot(rot)
        await asyncio.sleep(dt_sec)


# make cube translate in async function during the cube_rotation is running
async def async_cube_translation(cube, steps, speed):
    global dt
    pos = cube.GetPos()
    while True:
        # Translation on +X
        for i in range(steps):
            dt_sec = hg.time_to_sec_f(dt)
            pos.x += dt_sec * speed
            cube.SetPos(pos)
            await asyncio.sleep(dt_sec)
        # Translation on +Z
        for i in range(steps):
            dt_sec = hg.time_to_sec_f(dt)
            pos.z += dt_sec * speed
            cube.SetPos(pos)
            await asyncio.sleep(dt_sec)
        # Translation on -X
        for i in range(steps):
            dt_sec = hg.time_to_sec_f(dt)
            pos.x -= dt_sec * speed
            cube.SetPos(pos)
            await asyncio.sleep(dt_sec)
        # Translation on -Z
        for i in range(steps):
            dt_sec = hg.time_to_sec_f(dt)
            pos.z -= dt_sec * speed
            cube.SetPos(pos)
            await asyncio.sleep(dt_sec)


def create_material(diffuse, specular, self, prg_ref):
    mat = hg.CreateMaterial(prg_ref, 'uDiffuseColor', diffuse, 'uSpecularColor', specular)
    hg.SetMaterialValue(mat, 'uSelfColor', self)
    return mat


async def async_main():
    global dt

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

    cam_mtx = hg.TransformationMat4(hg.Vec3(0, 7, -10), hg.Deg3(30, 0, 0))
    cam = hg.CreateCamera(scene, cam_mtx, 0.01, 3000)
    scene.SetCurrentCamera(cam)

    hg.CreateLinearLight(scene, hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Deg3(30, 59, 0)),
                         hg.Color(1, 0.8, 0.7), hg.Color(1, 0.8, 0.7), 10, hg.LST_Map, 0.002,
                         hg.Vec4(50, 100, 200, 400))
    hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(0, 10, 10)), 100, hg.ColorI(94, 155, 228),
                        hg.ColorI(94, 255, 228))

    # create models
    cube_model_1 = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
    cube_ref_1 = res.AddModel('cube', cube_model_1)
    cube_instance_1 = hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 0, 0)), cube_ref_1,
                                             [mat_objects])
    cube_transform_1 = cube_instance_1.GetTransform()

    cube_model_2 = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
    cube_ref_2 = res.AddModel('cube', cube_model_2)
    cube_instance_2 = hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 3, 0)), cube_ref_2,
                                                [mat_objects])
    cube_transform_2 = cube_instance_2.GetTransform()

    # create async tasks
    asyncio.create_task(async_cube_rotation(cube_transform_1))
    asyncio.create_task(async_cube_translation(cube_transform_2, steps=100, speed=2.0))
    asyncio.get_event_loop()

    # main loop
    while not hg.ReadKeyboard().Key(hg.K_Escape) and hg.IsWindowOpen(win):
        dt = hg.TickClock()

        # update the scene
        scene.Update(dt)

        hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), True, pipeline, res)

        # give async tasks time to do their job
        await asyncio.sleep(0)

        hg.Frame()
        hg.UpdateWindow(win)

    hg.DestroyForwardPipeline(pipeline)
    hg.RenderShutdown()
    hg.DestroyWindow(win)


asyncio.run(async_main())

