-- Reading advanced keyboard state

hg = require("harfang")

hg.InputInit()

keyboard = hg.Keyboard('raw')  -- note: the 'raw' device can be queried without an open window, use 'default' otherwise

while not keyboard:Pressed(hg.K_Escape) do
	keyboard:Update()

	for key=0, hg.K_Last-1 do
		if keyboard:Released(key) then  -- will react on key release using the current and previous keyboard state
			print(string.format('Key released: %d' , key))
		end
	end
end