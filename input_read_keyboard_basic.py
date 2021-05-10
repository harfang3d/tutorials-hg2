# Reading basic keyboard state

import harfang as hg

hg.InputInit()

while True:
	state = hg.ReadKeyboard('raw')  # note: the 'raw' device can be queried without an open window, use 'default' otherwise

	for key in range(hg.K_Last):
		if state.Key(key):
			print('Key down: %d' % key)

	if state.Key(hg.K_Escape):
		break
