package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.WindowSystemInit()

	resX, resY := int32(256), int32(256)
	win := hg.RenderInitWithWindowTitleWidthHeightResetFlags("Harfang - Load from Assets", resX, resY, hg.RFVSync)

	// mount folder as an assets source and load texture from the assets system
	hg.AddAssetsFolder("resources_compiled")

	tex, tex_info := hg.LoadTextureFromAssets("pictures/owl.jpg", 0)

	if hg.IsValidWithT(tex) {
		fmt.Printf("Texture dimensions: %dx%d", tex_info.GetWidth(), tex_info.GetHeight())
	} else {
		fmt.Printf("Failed to load texture")
	}

	hg.RenderShutdown()
	hg.DestroyWindow(win)
}
