# Reading advanced keyboard state

import harfang as hg

hg.InputInit()

keyboard = hg.Keyboard('raw')  # note: the 'raw' device can be queried without an open window, use 'default' otherwise

while not keyboard.Pressed(hg.K_Escape):
	keyboard.Update()

	for key in range(hg.K_Last):
		if keyboard.Released(key):  # will react on key release using the current and previous keyboard state
			print('Key released: %d' % key)