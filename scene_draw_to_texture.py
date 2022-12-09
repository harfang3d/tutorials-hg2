# Draw scene to texture

import harfang as hg
import math

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1024, 1024
win = hg.RenderInit('Draw Scene to Texture', res_x, res_y, hg.RF_VSync | hg.RF_MSAA8X)

hg.AddAssetsFolder("resources_compiled")

# create pipeline
pipeline = hg.CreateForwardPipeline()
res = hg.PipelineResources()

# load the scene to draw to a texture
scene = hg.Scene()
hg.LoadSceneFromAssets("materials/materials.scn", scene, res, hg.GetForwardPipelineInfo())

# create a 512x512 frame buffer to draw the scene to
frame_buffer = hg.CreateFrameBuffer(512, 512, hg.TF_RGBA32F, hg.TF_D24, 4, 'framebuffer')  # 4x MSAA
color = hg.GetColorTexture(frame_buffer)

# create the cube model
vtx_layout = hg.VertexLayoutPosFloatTexCoord0UInt8()

cube_mdl = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
cube_ref = res.AddModel('cube', cube_mdl)

# prepare the cube shader program
cube_prg = hg.LoadProgramFromAssets('shaders/texture')

# main loop
angle = 0

while not hg.ReadKeyboard().Key(hg.K_Escape) and hg.IsWindowOpen(win):
	dt = hg.TickClock()
	angle = angle + hg.time_to_sec_f(dt)

	# update scene and render to the frame buffer
	scene.GetCurrentCamera().GetTransform().SetPos(hg.Vec3(0, 0, -(math.sin(angle) * 3 + 4)))  # animate the scene current camera on Z

	scene.Update(dt)

	view_id = 0
	view_id, pass_ids = hg.SubmitSceneToPipeline(view_id, scene, hg.IntRect(0, 0, 512, 512), True, pipeline, res, frame_buffer.handle)

	# draw a rotating cube in immediate mode using the texture the scene was rendered to
	hg.SetViewPerspective(view_id, 0, 0, res_x, res_y, hg.TranslationMat4(hg.Vec3(0, 0, -1.8)))

	val_uniforms = [hg.MakeUniformSetValue('color', hg.Vec4(1, 1, 1, 1))]  # note: these could be moved out of the main loop but are kept here for readability
	tex_uniforms = [hg.MakeUniformSetTexture('s_tex', color, 0)]

	hg.DrawModel(view_id, cube_mdl, cube_prg, val_uniforms, tex_uniforms, hg.TransformationMat4(hg.Vec3(0, 0, 0), hg.Vec3(angle * 0.1, angle * 0.05, angle * 0.2)))

	#
	hg.Frame()
	hg.UpdateWindow(win)

hg.RenderShutdown()
hg.WindowSystemShutdown()
