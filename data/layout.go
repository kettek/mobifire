package data

import (
	"embed"
	"fmt"
	"path"

	"github.com/kettek/rebui"
	"gopkg.in/yaml.v3"
)

//go:embed layout/*.yaml
var layoutFS embed.FS
var layouts []Layout

// Layout is a layout for mobifire.
type Layout struct {
	Name  string
	Nodes []rebui.Node
}

func GetLayout(n string) (Layout, error) {
	var layout Layout
	var found bool
	for _, l := range layouts {
		if l.Name == n {
			layout = l
			found = true
			break
		}
	}
	if !found {
		return layout, fmt.Errorf("no such layout %s", n)
	}

	// Duplicate nodes so we don't mess with the original slice.
	nodes := make([]rebui.Node, len(layout.Nodes))
	copy(nodes, layout.Nodes)

	return Layout{
		Name:  layout.Name,
		Nodes: nodes,
	}, nil
}

func init() {
	entries, err := layoutFS.ReadDir("layout")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		b, err := layoutFS.ReadFile(path.Join("layout", entry.Name()))
		if err != nil {
			panic(err)
		}

		var layout Layout

		if err := yaml.Unmarshal(b, &layout.Nodes); err != nil {
			panic(err)
		}

		layout.Name = (entry.Name())[:len(entry.Name())-len(path.Ext(entry.Name()))]

		fmt.Println(" i read layout", layout.Name)

		layouts = append(layouts, layout)
	}
}
