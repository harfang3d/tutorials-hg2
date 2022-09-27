package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	pic := hg.NewPicture()

	if !hg.LoadJPG(pic, "resources/pictures/owl.jpg") {
		panic("Failed to load picture!")
	}

	if !hg.SavePNG(pic, "owl.png") {
		panic("Failed to save picture!")
	}

	print("Conversion complete")
}
