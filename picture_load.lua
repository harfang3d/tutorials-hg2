-- Load a picture from a file.

hg = require("harfang")

pic = hg.Picture()

if hg.LoadPicture(pic, "resources/pictures/owl.jpg") then
	print(string.format("Picture dimensions: %dx%d" , pic:GetWidth(), pic:GetHeight()))
else
	print("Failed to load picture!")
end
