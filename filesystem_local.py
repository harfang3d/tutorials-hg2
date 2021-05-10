import harfang as hg

file = hg.Open('resources/pictures/owl.jpg')

if hg.IsValid(file):
	print('File is %d bytes long' % hg.GetSize(file))
else:
	print('Failed to open file')

hg.Close(file)
