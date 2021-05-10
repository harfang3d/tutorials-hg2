# Load a picture from a file.

import harfang as hg

pic = hg.Picture()

if hg.LoadPicture(pic, "resources/pictures/owl.jpg"):
	print("Picture dimensions: %rx%r" % (pic.GetWidth(), pic.GetHeight()))
else:
	print("Failed to load picture!")
