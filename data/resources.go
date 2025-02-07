package data

import "fyne.io/fyne/v2"

var resources = map[int]*fyne.StaticResource{}

func addResource(name string, content []byte) int {
	resources[len(resources)] = &fyne.StaticResource{
		StaticName:    name,
		StaticContent: content,
	}
	return len(resources) - 1
}

// GetResourceByID returns the given static resource by ID.
func GetResourceByID(id int) *fyne.StaticResource {
	return resources[id]
}

// GetResource returns the given static resource by name.
func GetResource(name string) *fyne.StaticResource {
	for _, res := range resources {
		if res.StaticName == name {
			return res
		}
	}
	return nil
}
