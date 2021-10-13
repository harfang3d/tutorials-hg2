# Draw text

import harfang as hg

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Draw Text', res_x, res_y, hg.RF_VSync)

# load font and shader program
font = hg.LoadFontFromFile('resources_compiled/font/default.ttf', 96)
font_prg = hg.LoadProgramFromFile('resources_compiled/core/shader/font')

# text uniforms and render state
text_uniform_values = [hg.MakeUniformSetValue('u_color', hg.Vec4(1, 1, 0))]
text_render_state = hg.ComputeRenderState(hg.BM_Alpha, hg.DT_Always, hg.FC_Disabled)

# main loop
while not hg.ReadKeyboard('default').Key(hg.K_Escape):
	hg.SetView2D(0, 0, 0, res_x, res_y, -1, 1, hg.CF_Color | hg.CF_Depth, hg.ColorI(32, 32, 32), 0, 1)

	hg.DrawText(0, font, 'Hello world!', font_prg, 'u_tex', 0, hg.Mat4.Identity, hg.Vec3(res_x / 2, res_y / 2, 0), hg.DTHA_Center, hg.DTVA_Center, text_uniform_values, [], text_render_state)

	hg.Frame()
	hg.UpdateWindow(win)

hg.RenderShutdown()
hg.DestroyWindow(win)
