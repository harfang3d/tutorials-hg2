import harfang as hg
import asyncio
import math

#make the cube rotate in async function
async def get_cube_rotation(rotating_cube, dt):
    rot_speed = (math.pi / 100.0) * dt
    rot = rotating_cube.GetRot()
    rot.y += rot_speed
    return rot

async def get_cube_translation(translating_cube, dt, steps=120, speed=0.01):
    speed = dt * speed
    pos = translating_cube.GetPos()

    for _ in range(steps):
        pos.x += speed
        translating_cube.SetPos(pos)
        await asyncio.sleep(0)

    for _ in range(steps):
        pos.z += speed
        translating_cube.SetPos(pos)
        await asyncio.sleep(0)

    for _ in range(steps):
        pos.x -= speed
        translating_cube.SetPos(pos)
        await asyncio.sleep(0)
        
    for _ in range(steps):
        pos.z -= speed
        translating_cube.SetPos(pos)
        await asyncio.sleep(0)

    return pos
def create_material(diffuse, specular, self, prg_ref):
	mat = hg.CreateMaterial(prg_ref, 'uDiffuseColor', diffuse, 'uSpecularColor', specular)
	hg.SetMaterialValue(mat, 'uSelfColor', self)
	return mat

async def render_scene(scene, win, pipeline, res, res_x, res_y):
    dt = hg.TickClock()
    scene.Update(dt)

    view_id, pass_id = hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), True, pipeline, res)

    hg.Frame()
    hg.UpdateWindow(win)

async def async_main():
    hg.InputInit()
    hg.WindowSystemInit()

    res_x, res_y = 1280, 720
    win = hg.RenderInit('Harfang - Scene Coroutine', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

    hg.AddAssetsFolder('resources_compiled')

    pipeline = hg.CreateForwardPipeline()
    res = hg.PipelineResources()

    vtx_layout = hg.VertexLayoutPosFloatNormUInt8()
    prg_ref = hg.LoadPipelineProgramRefFromAssets('core/shader/default.hps', res, hg.GetForwardPipelineInfo())
    mat_objects = create_material(hg.Vec4(0.5, 0.5, 0.5), hg.Vec4(1, 1, 1), hg.Vec4(0, 0, 0), prg_ref)

    scene = hg.Scene()
    scene.canvas.color = hg.ColorI(22, 56, 76)
    scene.environment.fog_color = scene.canvas.color
    scene.environment.fog_near = 20
    scene.environment.fog_far = 80

    cam_mtx = hg.TransformationMat4(hg.Vec3(0, 10, -15), hg.Deg3(30, 0, 0))
    cam = hg.CreateCamera(scene, cam_mtx, 0.01, 3000)
    scene.SetCurrentCamera(cam)

    hg.CreateLinearLight(scene, hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Deg3(30, 59, 0)),
                         hg.Color(1, 0.8, 0.7), hg.Color(1, 0.8, 0.7), 10, hg.LST_Map, 0.002,
                         hg.Vec4(50, 100, 200, 400))
    hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(0, 10, 10)), 100, hg.ColorI(94, 155, 228),
                        hg.ColorI(94, 255, 228))

    # create models
    rotating_cube_model = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
    rotating_cube_ref = res.AddModel('rotating_cube', rotating_cube_model)
    rotating_cube_instance = hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 0, 0)), rotating_cube_ref,
                                             [mat_objects])
    rotating_cube_transform = rotating_cube_instance.GetTransform()

    translating_cube_model = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
    translating_cube_ref = res.AddModel('translating_cube', translating_cube_model)
    translating_cube_instance = hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 3, 0)), translating_cube_ref,
                                                [mat_objects])
    translating_cube_transform = translating_cube_instance.GetTransform()

    physics = hg.SceneBullet3Physics()
    physics.SceneCreatePhysicsFromAssets(scene)

    while not hg.ReadKeyboard().Key(hg.K_Escape) and hg.IsWindowOpen(win):
        dt = hg.TickClock()

        rotation = await get_cube_rotation(rotating_cube_transform, dt)
        translating_cube_pos = await get_cube_translation(translating_cube_transform, dt)

        rotating_cube_transform.SetRot(rotation)
        translating_cube_transform.SetPos(translating_cube_pos)

        await render_scene(scene, win, pipeline, res, res_x, res_y)

        await asyncio.sleep(0.01)

    hg.DestroyForwardPipeline(pipeline)
    hg.RenderShutdown()
    hg.DestroyWindow(win)

asyncio.run(async_main())

