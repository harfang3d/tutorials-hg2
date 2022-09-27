package main

import (
	"fmt"
	"strings"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()

	names := hg.GetMouseNames()
	fmt.Print("Mouse device names: ")
	namesStr := []string{}
	for i := int32(0); i < names.Size(); i += 1 {
		namesStr = append(namesStr, names.At(i))
	}
	fmt.Println(strings.Join(namesStr, ", "))

	names = hg.GetKeyboardNames()
	fmt.Print("Keyboard device names: ")
	namesStr = []string{}
	for i := int32(0); i < names.Size(); i += 1 {
		namesStr = append(namesStr, names.At(i))
	}
	fmt.Println(strings.Join(namesStr, ", "))

	names = hg.GetGamepadNames()
	fmt.Print("Gamepad device names: ")
	namesStr = []string{}
	for i := int32(0); i < names.Size(); i += 1 {
		namesStr = append(namesStr, names.At(i))
	}
	fmt.Println(strings.Join(namesStr, ", "))
}
