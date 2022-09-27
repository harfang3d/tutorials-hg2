package main

import (
	"runtime"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	width, height := int32(1280), int32(720)

	window := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Basic Loop", width, height, hg.RFVSync)

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(window) {
		hg.SetViewClearWithColDepthStencil(0, hg.CFColor|hg.CFDepth, hg.ColorGetGreen(), 1, 0)
		hg.SetViewRect(0, 0, 0, uint16(width), uint16(height))

		hg.Touch(0) // force the view to be processed as it would be ignored since nothing is drawn to it (a clear does not count)

		hg.Frame()
		hg.UpdateWindow(window)
		runtime.GC()
	}

	hg.RenderShutdown()
	hg.DestroyWindow(window)

}
