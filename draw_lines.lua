-- Draw Lines

hg = require("harfang")

hg.InputInit()
hg.WindowSystemInit()

res_x, res_y = 1280, 720
win = hg.RenderInit('Harfang - Draw Lines', res_x, res_y, hg.RF_VSync | hg.RF_MSAA4X)

line_count = 1000

shader = hg.LoadProgramFromFile('resources_compiled/shaders/white')

-- vertices
vtx_layout = hg.VertexLayout()
vtx_layout:Begin()
vtx_layout:Add(hg.A_Position, 3, hg.AT_Float)
vtx_layout:End()

vtx = hg.Vertices(vtx_layout, line_count * 2)

-- main loop
angle = 0

while not hg.ReadKeyboard():Key(hg.K_Escape) and hg.IsWindowOpen(win) do
	hg.SetViewClear(0, hg.CF_Color | hg.CF_Depth, hg.ColorI(64, 64, 64), 1, 0)
	hg.SetViewRect(0, 0, 0, res_x, res_y)

	vtx:Clear()
	for i = 0, line_count-1 do
		vtx:Begin(2 * i):SetPos(hg.Vec3(math.sin(angle + i * 0.005), math.cos(angle + i * 0.01), 0)):End()
		vtx:Begin(2 * i + 1):SetPos(hg.Vec3(math.sin(angle + i * -0.005), math.cos(angle + i * 0.005), 0)):End()
	end

	hg.DrawLines(0, vtx, shader)  -- submit all lines in a single call

	angle = angle + hg.time_to_sec_f(hg.TickClock())

	hg.Frame()
	hg.UpdateWindow(win)
end

hg.RenderShutdown()
hg.DestroyWindow(win)
