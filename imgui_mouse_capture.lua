-- Detect ImGui mouse capture

hg= require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - ImGui Mouse Capture', res_x, res_y, hg.RF_VSync)

hg.AddAssetsFolder('resources_compiled')

--
hg.ImGuiInit(10, hg.LoadProgramFromAssets('core/shader/imgui'), hg.LoadProgramFromAssets('core/shader/imgui_image'))
text_value = 'Clicking into this field will not clear the screen in red.'

--
mouse = hg.Mouse()
keyboard = hg.Keyboard()

while not hg.ReadKeyboard():Key(hg.K_Escape) do
	mouse:Update()
	keyboard:Update()

	dt = hg.TickClock()
	hg.ImGuiBeginFrame(res_x, res_y, dt, mouse:GetState(), keyboard:GetState())

	if hg.ImGuiWantCaptureMouse() then
		clear_color = hg.Color.Black  -- black background if ImGui has mouse capture
	else
		if mouse:Down(hg.MB_0) then
			clear_color = hg.Color.Red  
		else 
			clear_color = hg.Color.Black
		end
	end

	hg.SetView2D(0, res_x, res_y, -1, 0, hg.CF_Color | hg.CF_Depth, clear_color, 1, 0)

	hg.ImGuiSetNextWindowPosCenter(hg.ImGuiCond_Once)
	hg.ImGuiSetNextWindowSize(hg.Vec2(700, 96), hg.ImGuiCond_Once)

	if hg.ImGuiBegin('Detecting ImGui mouse capture') then
		hg.ImGuiTextWrapped('Click outside of the GUI to clear the screen in red.')
		hg.ImGuiSeparator()
		_, text_value = hg.ImGuiInputText('Text Input', text_value, 4096)
	end
	hg.ImGuiEnd()

	hg.ImGuiEndFrame(0)

	hg.Frame()
	hg.UpdateWindow(win)
end
hg.DestroyWindow(win)
