package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	for {
		state := hg.ReadKeyboardWithName("raw") // note: the "raw" device can be queried without an open window, use "default" otherwise

		for key := hg.Key(0); key < hg.KLast; key += 1 {
			if state.Key(key) {
				fmt.Println("Key down: ", key)
			}
		}

		if state.Key(hg.KEscape) {
			break
		}
	}
}
