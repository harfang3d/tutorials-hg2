# List input devices

import harfang as hg

hg.InputInit()

names = hg.GetMouseNames()
print('Mouse device names: ' + ','.join(names))

names = hg.GetKeyboardNames()
print('Keyboard device names: ' + ','.join(names))

names = hg.GetGamepadNames()
print('Gamepad device names: ' + ','.join(names))
