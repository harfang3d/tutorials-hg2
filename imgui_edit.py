# ImGui basics

import harfang as hg

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - ImGui Edit', res_x, res_y, hg.RF_VSync)

hg.AddAssetsFolder('resources_compiled')

# initialize ImGui
imgui_prg = hg.LoadProgramFromAssets('core/shader/imgui')
imgui_img_prg = hg.LoadProgramFromAssets('core/shader/imgui_image')

hg.ImGuiInit(10, imgui_prg, imgui_img_prg)

imgui_output_view = 255
imgui_view_clear_color = hg.Color(0, 0, 0)
imgui_clear_color_preset = 0

# main loop
while not hg.ReadKeyboard().Key(hg.K_Escape):
	dt = hg.TickClock()

	# ImGui frame
	hg.ImGuiBeginFrame(res_x, res_y, dt, hg.ReadMouse(), hg.ReadKeyboard())

	hg.ImGuiSetNextWindowPosCenter(hg.ImGuiCond_Once)

	if hg.ImGuiBegin('ImGui Controls', True, hg.ImGuiWindowFlags_AlwaysAutoResize)[0]:
		val_modified, imgui_clear_color_preset = hg.ImGuiCombo('Set Clear Color', imgui_clear_color_preset, ['Red', 'Green', 'Blue'])

		# apply preset if a combo entry was selected
		if val_modified:
			if imgui_clear_color_preset == 0:
				imgui_view_clear_color = hg.Color(1, 0, 0)
			elif imgui_clear_color_preset == 1:
				imgui_view_clear_color = hg.Color(0, 1, 0)
			else:
				imgui_view_clear_color = hg.Color(0, 0, 1)

		# reset clear color to black on button click
		if hg.ImGuiButton('Reset Clear Color'):
			imgui_view_clear_color = hg.Color.Black

		# custom clear color edit
		val_modified, imgui_view_clear_color = hg.ImGuiColorEdit('Edit Clear Color', imgui_view_clear_color)

		# edit the ImGui output view
		val_modified, imgui_output_view = hg.ImGuiInputInt('ImGui Output View', imgui_output_view)
		if val_modified:
			imgui_output_view = max(0, min(imgui_output_view, 255))  # keep output view in [0;255] range

	hg.ImGuiEnd()

	hg.SetView2D(imgui_output_view, 0, 0, res_x, res_y, -1, 0, hg.CF_Color | hg.CF_Depth, imgui_view_clear_color, 1, 0)
	hg.ImGuiEndFrame(imgui_output_view)

	# rendering frame
	hg.Frame()

	hg.UpdateWindow(win)

hg.RenderShutdown()
hg.DestroyWindow(win)
