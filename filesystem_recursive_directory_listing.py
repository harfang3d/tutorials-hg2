# Recursive directory listing

import harfang as hg

def entry_type_to_string(type):
	return {hg.DE_File: 'file', hg.DE_Dir: 'directory', hg.DE_Link: 'link'}[type]

entries = hg.ListDirRecursive('resources', hg.DE_All)

for entry in entries:
	print('- %s is a %s' % (entry.name, entry_type_to_string(entry.type)))
