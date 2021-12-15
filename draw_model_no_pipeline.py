# Draw models without a pipeline

import harfang as hg

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Draw Models no Pipeline', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

# vertex layout and models
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

cube_mdl = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
ground_mdl = hg.CreatePlaneModel(vtx_layout, 5, 5, 1, 1)

shader = hg.LoadProgramFromFile('resources_compiled/shaders/mdl')

# main loop
angle = 0

while not hg.ReadKeyboard().Key(hg.K_Escape):
	dt = hg.TickClock()
	angle = angle + hg.time_to_sec_f(dt)

	viewpoint = hg.TranslationMat4(hg.Vec3(0, 1, -3))
	hg.SetViewPerspective(0, 0, 0, res_x, res_y, viewpoint)

	hg.DrawModel(0, cube_mdl, shader, [], [], hg.TransformationMat4(hg.Vec3(0, 1, 0), hg.Vec3(angle, angle, angle)))
	hg.DrawModel(0, ground_mdl, shader, [], [], hg.TranslationMat4(hg.Vec3(0, 0, 0)))

	hg.Frame()
	hg.UpdateWindow(win)

hg.RenderShutdown()
