-- How to resize the render window.

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 512, 512
win = hg.RenderInit('Harfang - Render Resize to Window', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

-- create model
vtx_layout = hg.VertexLayoutPosFloatNormUInt8()

cube_mdl = hg.CreateCubeModel(vtx_layout, 1, 1, 1)
cube_prg = hg.LoadProgramFromFile('resources_compiled/shaders/mdl')

-- main loop
while not hg.ReadKeyboard():Key(hg.K_Escape) do
	render_was_reset, res_x, res_y = hg.RenderResetToWindow(win, res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X | hg.RF_MaxAnisotropy)
	if render_was_reset then
		print(string.format('Render reset to %dx%d', res_x, res_y))
	end

	viewpoint = hg.TransformationMat4(hg.Vec3(1, 1, -2), hg.Deg3(24, -27, 0))
	hg.SetViewPerspective(0, 0, 0, res_x, res_y, viewpoint, 0.01, 100, 1.8, hg.CF_Color | hg.CF_Depth, hg.ColorI(64, 64, 64), 1, 0)

	hg.DrawModel(0, cube_mdl, cube_prg, {}, {}, hg.TranslationMat4(hg.Vec3(0, 0, 0)))

	hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
