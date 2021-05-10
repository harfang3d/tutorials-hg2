-- Access to a file by mounting a folder as an assets source

hg = require("harfang")

hg.WindowSystemInit()

res_x, res_y = 256, 256
win = hg.RenderInit('Harfang - Load from Assets', res_x, res_y, hg.RF_VSync)

-- mount folder as an assets source and load texture from the assets system
hg.AddAssetsFolder('resources_compiled')

tex, tex_info = hg.LoadTextureFromAssets('pictures/owl.jpg', 0)

if hg.IsValid(tex) then
	print(string.format('Texture dimensions: %dx%d',tex_info.width, tex_info.height))
else
	print('Failed to load texture')
end

hg.RenderShutdown()
