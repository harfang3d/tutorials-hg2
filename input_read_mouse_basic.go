package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	for !hg.ReadKeyboardWithName("raw").Key(hg.KEscape) { // note: the "raw" device can be queried without an open window, use "default" otherwise
		state := hg.ReadMouseWithName("raw")

		x := state.X()
		y := state.Y()

		button_0 := state.Button(hg.MB0)
		button_1 := state.Button(hg.MB1)
		button_2 := state.Button(hg.MB2)

		wheel := state.Wheel()

		fmt.Println("Mouse state: X=%d Y=%d B0=%v B1=%v B2=%v Wheel=%d", x, y, button_0, button_1, button_2, wheel)
	}
}
