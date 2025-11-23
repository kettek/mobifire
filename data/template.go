package data

import (
	"embed"
	"fmt"
	"path"

	"github.com/kettek/rebui"
	"gopkg.in/yaml.v3"
)

//go:embed template/*.yaml
var templateFS embed.FS
var templates []Template

type Template struct {
	Name  string
	Nodes rebui.Nodes
}

func GetTemplate(n string) (Template, error) {
	var template Template
	var found bool
	for _, t := range templates {
		if t.Name == n {
			template = t
			found = true
			break
		}
	}
	if !found {
		return template, fmt.Errorf("no such template %s", n)
	}

	// Duplicate nodes so we don't mess with the original slice.
	nodes := make(rebui.Nodes, len(template.Nodes))
	copy(nodes, template.Nodes)

	return Template{
		Name:  template.Name,
		Nodes: nodes,
	}, nil
}

func init() {
	entries, err := templateFS.ReadDir("template")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		b, err := templateFS.ReadFile(path.Join("template", entry.Name()))
		if err != nil {
			panic(err)
		}

		var template Template

		if err := yaml.Unmarshal(b, &template.Nodes); err != nil {
			panic(err)
		}

		template.Name = (entry.Name())[:len(entry.Name())-len(path.Ext(entry.Name()))]

		fmt.Println(" i read template", template.Name)

		templates = append(templates, template)
	}
	rebui.SetTemplateLoader(func(path string) (rebui.Nodes, error) {
		template, err := GetTemplate(path)
		if err != nil {
			return nil, err
		}
		return template.Nodes, nil
	})
}
