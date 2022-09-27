package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - ImGui Mouse Capture", res_x, res_y, hg.RFVSync)

	hg.AddAssetsFolder("resources_compiled")

	hg.ImGuiInit(10, hg.LoadProgramFromAssets("core/shader/imgui"), hg.LoadProgramFromAssets("core/shader/imgui_image"))
	text_value := "Clicking into this field will not clear the screen in red."

	mouse := hg.NewMouse()
	keyboard := hg.NewKeyboard()

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		mouse.Update()
		keyboard.Update()

		dt := hg.TickClock()
		hg.ImGuiBeginFrame(res_x, res_y, dt, mouse.GetState(), keyboard.GetState())

		clear_color := hg.ColorGetBlack()
		if hg.ImGuiWantCaptureMouse() {
			clear_color = hg.ColorGetBlack() // black background if ImGui has mouse capture
		} else {
			if mouse.Down(int32(hg.MB0)) {
				clear_color = hg.ColorGetRed()
			}
		}

		hg.SetView2DWithZnearZfarFlagsColorDepthStencil(0, 0, 0, res_x, res_y, -1, 0, hg.CFColor|hg.CFDepth, clear_color, 1, 0)

		hg.ImGuiSetNextWindowPosCenterWithCondition(hg.ImGuiCondOnce)
		hg.ImGuiSetNextWindowSizeWithCondition(hg.NewVec2WithXY(700, 96), hg.ImGuiCondOnce)

		if hg.ImGuiBegin("Detecting ImGui mouse capture") {
			hg.ImGuiTextWrapped("Click outside of the GUI to clear the screen in red.")
			hg.ImGuiSeparator()
			_, text_value_temp := hg.ImGuiInputText("Text Input", text_value, 4096)
			text_value = *text_value_temp
		}
		hg.ImGuiEnd()

		hg.ImGuiEndFrameWithViewId(0)

		hg.Frame()
		hg.UpdateWindow(win)
	}
	hg.DestroyWindow(win)
}
