package main

import (
	"fmt"

	hg "github.com/harfang3d/harfang-go/v3"
)

func EntryTypeToString(type_ int32) string {
	return map[hg.DirEntryType]string{hg.DEFile: "file", hg.DEDir: "directory", hg.DELink: "link"}[hg.DirEntryType(type_)]
}

func main() {
	entries := hg.ListDirRecursive("resources", hg.DEAll)

	for i := int32(0); i < entries.Size(); i += 1 {
		entry := entries.At(i)
		fmt.Printf("- %s is a %s\n", entry.GetName(), EntryTypeToString(entry.GetType()))
	}
}
