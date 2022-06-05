# Basic loop

import harfang as hg

hg.InputInit()
hg.WindowSystemInit()

width, height = 1280, 720
window = hg.RenderInit('Harfang - Basic Loop', width, height, hg.RF_VSync)

while not hg.ReadKeyboard().Key(hg.K_Escape) and hg.IsWindowOpen(window):
	hg.SetViewClear(0, hg.CF_Color | hg.CF_Depth, hg.Color.Green, 1, 0)
	hg.SetViewRect(0, 0, 0, width, height)

	hg.Touch(0)  # force the view to be processed as it would be ignored since nothing is drawn to it (a clear does not count)

	hg.Frame()
	hg.UpdateWindow(window)

hg.RenderShutdown()
hg.DestroyWindow(window)
