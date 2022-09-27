package main

import (
	"fmt"
	"math"

	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	hg.InputInit()
	hg.WindowSystemInit()

	win := hg.NewWindowWithTitleWidthHeight("Harfang - Read Gamepad", 320, 200)

	gamepad := hg.NewGamepad()

	for !hg.ReadKeyboard().Key(hg.KEscape) && hg.IsWindowOpen(win) {
		gamepad.Update()

		if gamepad.Connected() {
			fmt.Println("Gamepad slot 0 was just connected")
		}
		if gamepad.Disconnected() {
			fmt.Println("Gamepad slot 0 was just disconnected")
		}

		if gamepad.Pressed(hg.GBButtonA) {
			fmt.Println("Gamepad button A pressed")
		}
		if gamepad.Pressed(hg.GBButtonB) {
			fmt.Println("Gamepad button B pressed")
		}
		if gamepad.Pressed(hg.GBButtonX) {
			fmt.Println("Gamepad button X pressed")
		}
		if gamepad.Pressed(hg.GBButtonY) {
			fmt.Println("Gamepad button Y pressed")
		}

		axis_left_x := gamepad.Axes(hg.GALeftX)
		if math.Abs(float64(axis_left_x)) > 0.1 {
			fmt.Println("Gamepad axis left X: %v", axis_left_x)
		}

		axis_left_y := gamepad.Axes(hg.GALeftY)
		if math.Abs(float64(axis_left_y)) > 0.1 {
			fmt.Println("Gamepad axis left Y: %v", axis_left_y)
		}

		hg.UpdateWindow(win)
	}
	hg.DestroyWindow(win)
}
