hg = require("harfang")

file = hg.Open('resources/pictures/owl.jpg')

if hg.IsValid(file) then
	print(string.format('File is %d bytes long',hg.GetSize(file)))
else
	print('Failed to open file')
end

hg.Close(file)
