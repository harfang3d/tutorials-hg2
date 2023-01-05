package main

import (
	"math"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	res_x, res_y := int32(1280), int32(720)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - ImGui Edit", res_x, res_y, hg.RFVSync)

	hg.AddAssetsFolder("resources_compiled")

	// initialize ImGui
	imgui_prg := hg.LoadProgramFromAssets("core/shader/imgui")
	imgui_img_prg := hg.LoadProgramFromAssets("core/shader/imgui_image")

	hg.ImGuiInit(10, imgui_prg, imgui_img_prg)

	imgui_output_view := int32(255)
	imgui_view_clear_color := hg.NewColorWithRGB(0, 0, 0)
	imgui_clear_color_preset := int32(0)
	imguiOpen := true

	// main loop
	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		hg.RenderResetToWindow(win, &res_x, &res_y)
		dt := hg.TickClock()

		// ImGui frame
		hg.ImGuiBeginFrame(res_x, res_y, dt, hg.ReadMouse(), hg.ReadKeyboard())

		hg.ImGuiSetNextWindowPosCenterWithCondition(hg.ImGuiCondOnce)

		if hg.ImGuiBeginWithOpenFlags("ImGui Controls", &imguiOpen, hg.ImGuiWindowFlagsAlwaysAutoResize) {
			val_modified := hg.ImGuiComboWithSliceOfItems("Set Clear Color", &imgui_clear_color_preset, hg.GoSliceOfstring{"Red", "Green", "Blue"})

			// apply preset if a combo entry was selected
			if val_modified {
				if imgui_clear_color_preset == 0 {
					imgui_view_clear_color = hg.NewColorWithRGB(1, 0, 0)
				} else if imgui_clear_color_preset == 1 {
					imgui_view_clear_color = hg.NewColorWithRGB(0, 1, 0)
				} else {
					imgui_view_clear_color = hg.NewColorWithRGB(0, 0, 1)
				}
			}

			// reset clear color to black on button click
			if hg.ImGuiButton("Reset Clear Color") {
				imgui_view_clear_color = hg.ColorGetBlack()
			}

			// custom clear color edit
			val_modified = hg.ImGuiColorEdit("Edit Clear Color", imgui_view_clear_color)

			// edit the ImGui output view
			val_modified = hg.ImGuiInputInt("ImGui Output View", &imgui_output_view)
			if val_modified {
				imgui_output_view = int32(math.Max(0.0, math.Min(float64(imgui_output_view), 255.0))) // keep output view in [0;255] range}
			}
		}
		hg.ImGuiEnd()

		hg.SetView2DWithZnearZfarFlagsColorDepthStencil(uint16(imgui_output_view), 0, 0, res_x, res_y, -1, 0, hg.CFColor|hg.CFDepth, imgui_view_clear_color, 1, 0)
		hg.ImGuiEndFrameWithViewId(uint16(imgui_output_view))

		// rendering frame
		hg.Frame()

		hg.UpdateWindow(win)
	}
	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
