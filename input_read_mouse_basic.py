# Reading basic mouse state

import harfang as hg

hg.InputInit()

while not hg.ReadKeyboard('raw').Key(hg.K_Escape):  # note: the 'raw' device can be queried without an open window, use 'default' otherwise
	state = hg.ReadMouse('raw')

	x = state.X()
	y = state.Y()

	button_0 = state.Button(hg.MB_0)
	button_1 = state.Button(hg.MB_1)
	button_2 = state.Button(hg.MB_2)

	wheel = state.Wheel()

	print('Mouse state: X=%r Y=%r B0=%r B1=%r B2=%r Wheel=%r' % (x, y, button_0, button_1, button_2, wheel))
