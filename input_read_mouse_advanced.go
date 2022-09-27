package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	mouse := hg.NewMouseWithName("raw")

	for !hg.ReadKeyboard().Key(hg.KEscape) { // note: the "raw" device can be queried without an open window, use "default" otherwise
		mouse.Update()

		dt_x := mouse.DtX()
		dt_y := mouse.DtY()

		if dt_x != 0 || dt_y != 0 {
			fmt.Println("Mouse delta X=%d delta Y=%d", dt_x, dt_y)
		}

		for i := int32(0); i < 3; i += 1 {
			if mouse.Pressed(int32(hg.MB0) + i) {
				fmt.Println("Mouse button %d pressed", i)
			}
			if mouse.Released(int32(hg.MB0) + i) {
				fmt.Println("Mouse button %d released", i)
			}
		}
	}
}
