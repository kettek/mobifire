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

//go:embed icon_commands.png
var sourceCommandsPng []byte
var resourceCommandsPng = &fyne.StaticResource{
	StaticName:    "commands",
	StaticContent: sourceCommandsPng,
}

//go:embed icon_apply.png
var sourceApplyPng []byte
var resourceApplyPng = &fyne.StaticResource{
	StaticName:    "apply",
	StaticContent: sourceApplyPng,
}

//go:embed icon_get.png
var sourceGetPng []byte
var resourceGetPng = &fyne.StaticResource{
	StaticName:    "get",
	StaticContent: sourceGetPng,
}
