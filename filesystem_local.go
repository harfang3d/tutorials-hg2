package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	file := hg.Open("resources/pictures/owl.jpg")

	if hg.IsValid(file) {
		fmt.Printf("File is %d bytes long", hg.GetSize(file))
	} else {
		fmt.Printf("Failed to open file")
	}
	hg.Close(file)
}
