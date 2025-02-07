package data

import (
	"embed"
	"path"
)

//go:embed images/*.png
var fs embed.FS

func init() {
	// Iterate them images.
	entries, err := fs.ReadDir("images")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		b, err := fs.ReadFile(path.Join("images", entry.Name()))
		if err != nil {
			panic(err)
		}
		addResource(entry.Name(), b)
	}
}
