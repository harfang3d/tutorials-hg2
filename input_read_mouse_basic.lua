-- Reading basic mouse state

hg = require("harfang")

hg.InputInit()

while not hg.ReadKeyboard('raw'):Key(hg.K_Escape) do  -- note: the 'raw' device can be queried without an open window, use 'default' otherwise
	state = hg.ReadMouse('raw')

	x = state:X()
	y = state:Y()

	button_0 = state:Button(hg.MB_0)
	button_1 = state:Button(hg.MB_1)
	button_2 = state:Button(hg.MB_2)

	wheel = state:Wheel()

	print(string.format('Mouse state: X=%d Y=%d B0=%s B1=%s B2=%s Wheel=%d', x, y, tostring(button_0), tostring(button_1), tostring(button_2), wheel))
end
