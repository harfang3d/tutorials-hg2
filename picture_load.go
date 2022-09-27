package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	pic := hg.NewPicture()

	if hg.LoadPicture(pic, "resources/pictures/owl.jpg") {
		fmt.Printf("Picture dimensions: %dx%d", pic.GetWidth(), pic.GetHeight())
	} else {
		fmt.Println("Failed to load picture!")
	}
}
