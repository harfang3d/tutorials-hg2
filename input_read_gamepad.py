# Reading advanced gamepad state

import harfang as hg

hg.InputInit()

hg.WindowSystemInit()
win = hg.NewWindow('Harfang - Read Gamepad', 320, 200)

gamepad = hg.Gamepad()

while not hg.ReadKeyboard().Key(hg.K_Escape) and hg.IsWindowOpen(win):
	gamepad.Update()

	if gamepad.Connected():
		print('Gamepad slot 0 was just connected')
	if gamepad.Disconnected():
		print('Gamepad slot 0 was just disconnected')

	if gamepad.Pressed(hg.GB_ButtonA):
		print('Gamepad button A pressed')
	if gamepad.Pressed(hg.GB_ButtonB):
		print('Gamepad button B pressed')
	if gamepad.Pressed(hg.GB_ButtonX):
		print('Gamepad button X pressed')
	if gamepad.Pressed(hg.GB_ButtonY):
		print('Gamepad button Y pressed')

	axis_left_x = gamepad.Axes(hg.GA_LeftX)
	if abs(axis_left_x) > 0.1:
		print('Gamepad axis left X: %r' % axis_left_x)

	axis_left_y = gamepad.Axes(hg.GA_LeftY)
	if abs(axis_left_y) > 0.1:
		print('Gamepad axis left Y: %r' % axis_left_y)

	hg.UpdateWindow(win)

hg.DestroyWindow(win)
