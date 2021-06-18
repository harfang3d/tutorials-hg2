-- ImGui basics

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - ImGui Basics', res_x, res_y, hg.RF_VSync)

-- initialize ImGui
hg.AddAssetsFolder('resources_compiled')

imgui_prg = hg.LoadProgramFromAssets('core/shader/imgui')
imgui_img_prg = hg.LoadProgramFromAssets('core/shader/imgui_image')

hg.ImGuiInit(10, imgui_prg, imgui_img_prg)

-- main loop
while not hg.ReadKeyboard():Key(hg.K_Escape) do
	hg.ImGuiBeginFrame(res_x, res_y, hg.TickClock(), hg.ReadMouse(), hg.ReadKeyboard())

	if hg.ImGuiBegin('Window') then
		hg.ImGuiText('Hello World!')
	end
	hg.ImGuiEnd()

	hg.SetView2D(0, res_x, res_y, -1, 1, hg.CF_Color | hg.CF_Depth, hg.Color.Black, 1, 0)
	hg.ImGuiEndFrame(0)

	hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)
