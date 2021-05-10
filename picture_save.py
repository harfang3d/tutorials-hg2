# Load a JPG picture and save it as PNG.

import harfang as hg
import sys

pic = hg.Picture()

ok = hg.LoadJPG(pic, 'resources/pictures/owl.jpg')
if not ok:
	sys.exit('Failed to load picture!')

ok = hg.SavePNG(pic, 'owl.png')
if not ok:
	sys.exit('Failed to save picture!')

print('Conversion complete')
