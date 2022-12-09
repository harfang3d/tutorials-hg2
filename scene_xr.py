# Display a scene in XR

import harfang as hg
import sys

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit("Harfang - OpenXR Scene", res_x, res_y, hg.RF_None)

hg.AddAssetsFolder("resources_compiled")

pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

render_data = hg.SceneForwardPipelineRenderData()  # this object is used by the low-level scene rendering API to share view-independent data with both eyes

# OpenXR initialization
if not hg.OpenXRInit():
    sys.exit()

eye_framebuffers = hg.OpenXRCreateEyeFrameBuffer(hg.OXRAA_MSAA4x)

# Create models
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

cube_mdl = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
cube_ref = res.AddModel('cube', cube_mdl)
ground_mdl = hg.CreateCubeModel(vtx_layout, 50, 0.01, 50)
ground_ref = res.AddModel('ground', ground_mdl)

# Load shader
prg_ref = hg.LoadPipelineProgramRefFromAssets('core/shader/pbr.hps', res, hg.GetForwardPipelineInfo())

# Create materials
def create_material(ubc, orm):
    mat = hg.Material()
    hg.SetMaterialProgram(mat, prg_ref)
    hg.SetMaterialValue(mat, "uBaseOpacityColor", ubc)
    hg.SetMaterialValue(mat, "uOcclusionRoughnessMetalnessColor", orm)
    return mat


# Create scene
scene = hg.Scene()
scene.canvas.color = hg.Color(255 / 255, 255 / 255, 217 / 255, 1)
scene.environment.ambient = hg.Color(15 / 255, 12 / 255, 9 / 255, 1)

lgt = hg.CreateSpotLight(scene, hg.TransformationMat4(hg.Vec3(-8, 4, -5), hg.Vec3(hg.Deg(19), hg.Deg(59), 0)), 0, hg.Deg(5), hg.Deg(30), hg.Color.White, hg.Color.White, 10, hg.LST_Map, 0.0001)
back_lgt = hg.CreatePointLight(scene, hg.TranslationMat4(hg.Vec3(2.4, 1, 0.5)), 10, hg.Color(94 / 255, 255 / 255, 228 / 255, 1), hg.Color(94 / 255, 1, 228 / 255, 1), 0)

mat_cube = create_material(hg.Vec4(255 / 255, 230 / 255, 20 / 255, 1), hg.Vec4(1, 0.658, 0., 1))
hg.CreateObject(scene, hg.TransformationMat4(hg.Vec3(0, 0.5, 0), hg.Vec3(0, hg.Deg(70), 0)), cube_ref, [mat_cube])

mat_ground = create_material(hg.Vec4(255 / 255, 120 / 255, 147 / 255, 1), hg.Vec4(1, 1, 0.1, 1))
hg.CreateObject(scene, hg.TranslationMat4(hg.Vec3(0, 0, 0)), ground_ref, [mat_ground])

# Setup 2D rendering to display eyes textures
quad_layout = hg.VertexLayout()
quad_layout.Begin().Add(hg.A_Position, 3, hg.AT_Float).Add(hg.A_TexCoord0, 3, hg.AT_Float).End()

quad_model = hg.CreatePlaneModel(quad_layout, 1, 1, 1, 1)
quad_render_state = hg.ComputeRenderState(hg.BM_Alpha, hg.DT_Disabled, hg.FC_Disabled)

eye_t_size = res_x / 2.5
eye_t_x = (res_x - 2 * eye_t_size) / 6 + eye_t_size / 2
quad_matrix = hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Vec3(hg.Deg(90), hg.Deg(0), hg.Deg(0)), hg.Vec3(eye_t_size, 1, eye_t_size))

tex0_program = hg.LoadProgramFromAssets("shaders/sprite")

quad_uniform_set_value_list = hg.UniformSetValueList()
quad_uniform_set_value_list.clear()
quad_uniform_set_value_list.push_back(hg.MakeUniformSetValue("color", hg.Vec4(1, 1, 1, 1)))

quad_uniform_set_texture_list = hg.UniformSetTextureList()

# Main loop
while not hg.ReadKeyboard().Key(hg.K_Escape) and hg.IsWindowOpen(win):
    dt = hg.TickClock()

    scene.Update(dt)

    actor_body_mtx = hg.TransformationMat4(hg.Vec3(-1.3, .45, -2), hg.Vec3(0, 0, 0))

    vid = 0  # keep track of the next free view id
    passId = hg.SceneForwardPipelinePassViewId()

    # Prepare view-independent render data once
    vid, passId = hg.PrepareSceneForwardPipelineCommonRenderData(vid, scene, render_data, pipeline, res, passId)

    def update_controllers(head):
        pass

    def draw_scene(rect, view_state, view_id, fb):
        passId = hg.SceneForwardPipelinePassViewId()
        view_id, passId = hg.PrepareSceneForwardPipelineViewDependentRenderData(view_id, view_state, scene, render_data, pipeline, res, passId)
        view_id, passId = hg.SubmitSceneToForwardPipeline(view_id, scene, rect, view_state, pipeline, render_data, res, fb)

        return view_id

    openxrFrameInfo, vid = hg.OpenXRSubmitSceneToForwardPipeline(hg.TranslationMat4(hg.Vec3(0, 0, 0)), update_controllers, draw_scene, vid, 0.1, 100)
    
    # Display the VR eyes texture to the backbuffer
    hg.SetViewRect(vid, 0, 0, res_x, res_y)
    vs = hg.ComputeOrthographicViewState(hg.TranslationMat4(hg.Vec3(0, 0, 0)), res_y, 0.1, 100, hg.ComputeAspectRatioX(res_x, res_y))
    hg.SetViewTransform(vid, vs.view, vs.proj)

    if openxrFrameInfo.id_fbs.size() > 0 :
        quad_uniform_set_texture_list.clear()
        quad_uniform_set_texture_list.push_back(hg.MakeUniformSetTexture("s_tex", hg.OpenXRGetColorTextureFromId(eye_framebuffers, openxrFrameInfo, 0), 0))
        hg.SetT(quad_matrix, hg.Vec3(eye_t_x, 0, 1))
        hg.DrawModel(vid, quad_model, tex0_program, quad_uniform_set_value_list, quad_uniform_set_texture_list, quad_matrix, quad_render_state)

        quad_uniform_set_texture_list.clear()    
        quad_uniform_set_texture_list.push_back(hg.MakeUniformSetTexture("s_tex", hg.OpenXRGetColorTextureFromId(eye_framebuffers, openxrFrameInfo, 1), 0))
        hg.SetT(quad_matrix, hg.Vec3(-eye_t_x, 0, 1))
        hg.DrawModel(vid, quad_model, tex0_program, quad_uniform_set_value_list, quad_uniform_set_texture_list, quad_matrix, quad_render_state)

    hg.Frame()

    hg.OpenXRFinishSubmitFrameBuffer(openxrFrameInfo)

    hg.UpdateWindow(win)

hg.DestroyForwardPipeline(pipeline)
hg.RenderShutdown()
hg.DestroyWindow(win)
