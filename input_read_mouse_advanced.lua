-- Reading advanced mouse state

hg = require("harfang")

hg.InputInit()

mouse = hg.Mouse('raw')

while not hg.ReadKeyboard('raw'):Key(hg.K_Escape) do  -- note: the 'raw' device can be queried without an open window, use 'default' otherwise
	mouse:Update()

	dt_x = mouse:DtX()
	dt_y = mouse:DtY()

	if dt_x ~= 0 or dt_y ~= 0 then
		print(string.format('Mouse delta X=%d delta Y=%d', dt_x, dt_y))
	end

	for i=0, 2 do
		if mouse:Pressed(hg.MB_0 + i) then
			print(string.format('Mouse button %d pressed' , i))
		end
		if mouse:Released(hg.MB_0 + i) then
			print(string.format('Mouse button %d released' , i))
		end
	end
end
