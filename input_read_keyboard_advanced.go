package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()

	keyboard := hg.NewKeyboardWithName("raw") // note: the "raw" device can be queried without an open window, use "default" otherwise

	for !keyboard.Pressed(hg.KEscape) {
		keyboard.Update()

		for key := hg.Key(0); key < hg.KLast; key += 1 {
			if keyboard.Released(key) { // will react on key release using the current and previous keyboard state
				fmt.Println("Key released: ", key)
			}
		}
	}
}
