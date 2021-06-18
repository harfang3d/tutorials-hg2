-- Reading basic keyboard state

hg = require("harfang")

hg.InputInit()

while true do
	state = hg.ReadKeyboard('raw')  -- note: the 'raw' device can be queried without an open window, use 'default' otherwise

	for key=0, hg.K_Last-1 do
		if state:Key(key) then
			print(string.format('Key down: %d' , key))
		end
	end
	if state:Key(hg.K_Escape) then
		break
	end
end