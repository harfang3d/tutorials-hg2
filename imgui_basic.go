package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - ImGui Basics", res_x, res_y, hg.RFVSync)

	// initialize ImGui
	hg.AddAssetsFolder("resources_compiled")

	imgui_prg := hg.LoadProgramFromAssets("core/shader/imgui")
	imgui_img_prg := hg.LoadProgramFromAssets("core/shader/imgui_image")

	hg.ImGuiInit(10, imgui_prg, imgui_img_prg)

	// main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		hg.RenderResetToWindow(win, &res_x, &res_y)

		hg.ImGuiBeginFrame(res_x, res_y, hg.TickClock(), hg.ReadMouse(), hg.ReadKeyboard())

		if hg.ImGuiBegin("Window") {
			hg.ImGuiText("Hello World!")
		}
		hg.ImGuiEnd()

		hg.SetView2DWithZnearZfarFlagsColorDepthStencil(0, 0, 0, res_x, res_y, -1, 1, hg.CFColor|hg.CFDepth, hg.ColorGetBlack(), 1, 0)
		hg.ImGuiEndFrameWithViewId(0)

		hg.Frame()
		hg.UpdateWindow(win)
	}
	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
