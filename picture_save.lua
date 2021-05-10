-- Load a JPG picture and save it as PNG.

hg = require("harfang")

pic = hg.Picture()

ok = hg.LoadJPG(pic, 'resources/pictures/owl.jpg')
if not ok then
	os.exit('Failed to load picture!')
end

ok = hg.SavePNG(pic, 'owl.png')
if not ok then
	os.exit('Failed to save picture!')
end

print('Conversion complete')
