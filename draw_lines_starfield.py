# Starfield 3D

import harfang as hg

hg.InputInit()
hg.WindowSystemInit()

width, height = 1280, 720
window = hg.RenderInit('Harfang - Starfield', width, height, hg.RF_VSync | hg.RF_MSAA4X)

# vertex layout
vtx_layout = hg.VertexLayout()
vtx_layout.Begin()
vtx_layout.Add(hg.A_Position, 3, hg.AT_Float)
vtx_layout.Add(hg.A_Color0, 3, hg.AT_Float)
vtx_layout.End()

# simple shader program
shader = hg.LoadProgramFromFile('resources_compiled/shaders/pos_rgb')

# initialize stars
starfield_size = 10

max_stars = 1000
vtx = hg.Vertices(vtx_layout, max_stars * 2)

stars = []
for i in range(max_stars):
	stars.append(hg.RandomVec3(-starfield_size, starfield_size))

# main loop
while not hg.ReadKeyboard().Key(hg.K_Escape) and hg.IsWindowOpen(window):
	hg.SetViewClear(0, hg.CF_Color | hg.CF_Depth, hg.Color.Black, 1, 0)
	hg.SetViewRect(0, 0, 0, width, height)

	dt = hg.TickClock()
	dt_f = hg.time_to_sec_f(dt)

	# update stars
	vtx.Clear()
	for i, star in enumerate(stars):
		star.z -= 2 * dt_f
		if star.z < starfield_size:
			star.z += starfield_size

		vtx.Begin(2 * i).SetPos(star * hg.Vec3(1 / star.z, 1 / star.z, 0)).SetColor0(hg.Color.Black).End()
		vtx.Begin(2 * i + 1).SetPos(star * hg.Vec3(1.04 / star.z, 1.04 / star.z, 0)).SetColor0(hg.Color.White).End()

	# draw stars as lines
	hg.DrawLines(0, vtx, shader)

	hg.Frame()
	hg.UpdateWindow(window)

hg.RenderShutdown()
hg.DestroyWindow(window)
