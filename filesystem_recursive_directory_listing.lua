-- Recursive directory listing

hg = require("harfang")

function entry_type_to_string(type)
	tp={}
	tp[hg.DE_File], tp[hg.DE_Dir], tp[hg.DE_Link] = 'file', 'directory', 'link'
	return tp[type]
end

entries = hg.ListDirRecursive('resources', hg.DE_All)

for i=0, entries:size() do
	entry = entries:at(i)
	print(string.format('- %s is a %s', entry.name, entry_type_to_string(entry.type)))
end