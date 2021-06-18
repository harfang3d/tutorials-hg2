-- Reading advanced gamepad state

hg = require("harfang")

hg.InputInit()

hg.WindowSystemInit()
win = hg.NewWindow('Harfang - Read Gamepad', 320, 200)

gamepad = hg.Gamepad()

while not hg.ReadKeyboard():Key(hg.K_Escape) do
	gamepad:Update()

	if gamepad:Connected() then
		print('Gamepad slot 0 was just connected')
	end
	if gamepad:Disconnected() then
		print('Gamepad slot 0 was just disconnected')
	end
	if gamepad:Pressed(hg.GB_ButtonA) then
		print('Gamepad button A pressed')
	end
	if gamepad:Pressed(hg.GB_ButtonB) then
		print('Gamepad button B pressed')
	end
	if gamepad:Pressed(hg.GB_ButtonX) then
		print('Gamepad button X pressed')
	end
	if gamepad:Pressed(hg.GB_ButtonY) then
		print('Gamepad button Y pressed')
	end

	axis_left_x = gamepad:Axes(hg.GA_LeftX)
	if math.abs(axis_left_x) > 0.1 then
		print(string.format('Gamepad axis left X: %f' , axis_left_x))
	end
	axis_left_y = gamepad:Axes(hg.GA_LeftY)
	if math.abs(axis_left_y) > 0.1 then
		print(string.format('Gamepad axis left Y: %f' , axis_left_y))
	end

	hg.UpdateWindow(win)
end
hg.DestroyWindow(win)
