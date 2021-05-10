# Reading advanced mouse state

import harfang as hg

hg.InputInit()

mouse = hg.Mouse('raw')

while not hg.ReadKeyboard('raw').Key(hg.K_Escape):  # note: the 'raw' device can be queried without an open window, use 'default' otherwise
	mouse.Update()

	dt_x = mouse.DtX()
	dt_y = mouse.DtY()

	if dt_x != 0 or dt_y != 0:
		print('Mouse delta X=%r delta Y=%r' % (dt_x, dt_y))

	for i in range(3):
		if mouse.Pressed(hg.MB_0 + i):
			print('Mouse button %r pressed' % i)
		if mouse.Released(hg.MB_0 + i):
			print('Mouse button %r released' % i)
