-- Draw models without a pipeline

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Draw Text over Models', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

-- vertex layout and models
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

cube_mdl = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
ground_mdl = hg.CreatePlaneModel(vtx_layout, 5, 5, 1, 1)

shader = hg.LoadProgramFromFile('resources_compiled/shaders/mdl')

-- load font and shader program
font = hg.LoadFontFromFile('resources_compiled/font/default.ttf', 96)
font_prg = hg.LoadProgramFromFile('resources_compiled/core/shader/font')

-- text uniforms and render state
text_uniform_values = {hg.MakeUniformSetValue('u_color', hg.Vec4(1, 1, 0))}
text_render_state = hg.ComputeRenderState(hg.BM_Alpha, hg.DT_Always, hg.FC_Disabled)

-- main loop
angle = 0

while not hg.ReadKeyboard():Key(hg.K_Escape) do
	dt = hg.TickClock()
	angle = angle + hg.time_to_sec_f(dt)

	-- 3D view
	viewpoint = hg.TranslationMat4(hg.Vec3(0, 1, -3))
	hg.SetViewPerspective(0, 0, 0, res_x, res_y, viewpoint, 0.01, 5000)

	hg.DrawModel(0, cube_mdl, shader, {}, {}, hg.TransformationMat4(hg.Vec3(0, 1, 0), hg.Vec3(angle, angle, angle)))
	hg.DrawModel(0, ground_mdl, shader, {}, {}, hg.TranslationMat4(hg.Vec3(0, 0, 0)))

	-- 2D view, note that only the depth buffer is cleared
	hg.SetView2D(1, 0, 0, res_x, res_y, -1, 1, hg.CF_Depth, hg.ColorI(32, 32, 32), 1, 0)

	hg.DrawText(1, font, 'Hello world!', font_prg, 'u_tex', 0, hg.Mat4.Identity, hg.Vec3(res_x / 2, res_y / 2, 0), hg.DTHA_Center, hg.DTVA_Center, text_uniform_values, {}, text_render_state)

	hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
