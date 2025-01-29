package play

import (
	"fyne.io/fyne/v2"

	// We're using embed you dumb warning.
	_ "embed"
)

//go:embed blank.base.111.png
var sourceBlankPng []byte
var resourceBlankPng = &fyne.StaticResource{
	StaticName:    "blank",
	StaticContent: sourceBlankPng,
}

//go:embed mark.base.111.png
var sourceMarkPng []byte
var resourceMarkPng = &fyne.StaticResource{
	StaticName:    "mark",
	StaticContent: sourceMarkPng,
}

//go:embed icon_inventory.png
var sourceInventoryPng []byte
var resourceInventoryPng = &fyne.StaticResource{
	StaticName:    "inventory",
	StaticContent: sourceInventoryPng,
}
