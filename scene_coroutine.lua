hg = require("harfang")

function create_material(diffuse, specular, self, prg_ref)
    mat = hg.CreateMaterial(prg_ref, 'uDiffuseColor', diffuse, 'uSpecularColor', specular)
    hg.SetMaterialValue(mat, 'uSelfColor', self)
    return mat
end

--make the cube rotate in coroutine
co_rotating_cube = coroutine.create(function(cube)
    local rotation = cube:GetRot()
    while true do
        local dt_sec = hg.time_to_sec_f(dt)
        rotation.y = rotation.y + (math.pi * dt_sec)
        cube:SetRot(rotation)
        coroutine.yield()
    end
end)


--make cube translate in coroutine during the cube_rotation is running
co_translating_cube = coroutine.create(function(cube, steps, speed)
    local position = cube:GetPos()
    while true do
        --Translation on +X
        for i = 1, steps do
            local dt_sec = hg.time_to_sec_f(dt)
            position.x = position.x + (dt_sec * speed)
            cube:SetPos(position)
            coroutine.yield()
        end
        --Translation on Z+
        for i = 1, steps do
            dt_sec = hg.time_to_sec_f(dt)
            position.z = position.z + (dt_sec * speed)
            cube:SetPos(position)
            coroutine.yield()
        end
        --Translation on -X
        for i = 1, steps do
            dt_sec = hg.time_to_sec_f(dt)
            position.x = position.x - (dt_sec * speed)
            cube:SetPos(position)
            coroutine.yield()
        end
        --Translation on -Z
        for i = 1, steps do
            dt_sec = hg.time_to_sec_f(dt)
            position.z = position.z - (dt_sec * speed)
            cube:SetPos(position)
            coroutine.yield()
        end
    end
end)

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1024, 720
win = hg.RenderInit('Harfang - Scene Coroutine', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

--Link precompiled assets folder to the project
hg.AddAssetsFolder('resources_compiled')

--Create Pipeline
pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

--Create Materials
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()
prg_ref = hg.LoadPipelineProgramRefFromAssets('core/shader/default.hps', res, hg.GetForwardPipelineInfo())
mat_objects = create_material(hg.Vec4(0.5, 0.5, 0.5), hg.Vec4(1, 1, 1), hg.Vec4(0, 0, 0), prg_ref)

--Setup Scene
scene = hg.Scene()
scene.canvas.color = hg.ColorI(22, 56, 76)
scene.environment.fog_color = scene.canvas.color
scene.environment.fog_near = 20
scene.environment.fog_far = 80

cam_mtx = hg.TransformationMat4(hg.Vec3(0, 7, -10), hg.Deg3(30, 0, 0))
cam = hg.CreateCamera(scene, cam_mtx, 0.01, 3000)
scene:SetCurrentCamera(cam)

hg.CreateLinearLight(scene, hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Deg3(30, 59, 0)),
hg.Color(1, 0.8, 0.7), hg.Color(1, 0.8, 0.7), 10, hg.LST_Map, 0.002,
hg.Vec4(50, 100, 200, 400))
hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(0, 10, 10)), 100, hg.ColorI(94, 155, 228),
hg.ColorI(94, 255, 228))

--Create Models
cube_model_1 = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
cube_ref_1 = res:AddModel('cube', cube_model_1)
cube_instance_1 = hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 0, 0)), cube_ref_1, {mat_objects})
cube_transform_1 = cube_instance_1:GetTransform()

cube_model_2 = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
cube_ref_2 = res:AddModel('cube', cube_model_2)
cube_instance_2 = hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 3, 0)), cube_ref_2, {mat_objects})
cube_transform_2 = cube_instance_2:GetTransform()

--Main Loop
while not hg.ReadKeyboard():Key(hg.K_Escape) and hg.IsWindowOpen(win) do
    dt = hg.TickClock()
    
    --Update Scene
    scene:Update(dt)
    
    coroutine.resume(co_rotating_cube, cube_transform_1)
    coroutine.resume(co_translating_cube, cube_transform_2, 100, 2)

    hg.SubmitSceneToPipeline(0, scene, hg.IntRect(0, 0, res_x, res_y), true, pipeline, res)

    hg.Frame()
    hg.UpdateWindow(win)
end

hg.DestroyForwardPipeline(pipeline)
hg.RenderShutdown()
hg.DestroyWindow(win)

