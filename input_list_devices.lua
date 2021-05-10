-- List input devices

hg = require("harfang")

hg.InputInit()

names = hg.GetMouseNames()
print('Mouse device names: ' .. table.concat(names, ','))

names = hg.GetKeyboardNames()
print('Keyboard device names: ' .. table.concat(names, ','))

names = hg.GetGamepadNames()
print('Gamepad device names: ' .. table.concat(names, ','))
