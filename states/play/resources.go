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

//go:embed icon_applied.png
var sourceAppliedPng []byte
var resourceAppliedPng = &fyne.StaticResource{
	StaticName:    "applied",
	StaticContent: sourceAppliedPng,
}

//go:embed icon_locked.png
var sourceLockedPng []byte
var resourceLockedPng = &fyne.StaticResource{
	StaticName:    "locked",
	StaticContent: sourceLockedPng,
}

//go:embed icon_unlocked.png
var sourceUnlockedPng []byte
var resourceUnlockedPng = &fyne.StaticResource{
	StaticName:    "unlocked",
	StaticContent: sourceUnlockedPng,
}

//go:embed icon_marked.png
var sourceMarkedPng []byte
var resourceMarkedPng = &fyne.StaticResource{
	StaticName:    "marked",
	StaticContent: sourceMarkedPng,
}

//go:embed icon_magic.png
var sourceMagicPng []byte
var resourceMagicPng = &fyne.StaticResource{
	StaticName:    "magic",
	StaticContent: sourceMagicPng,
}

//go:embed icon_cursed.png
var sourceCursedPng []byte
var resourceCursedPng = &fyne.StaticResource{
	StaticName:    "cursed",
	StaticContent: sourceCursedPng,
}

//go:embed icon_damned.png
var sourceDamnedPng []byte
var resourceDamnedPng = &fyne.StaticResource{
	StaticName:    "damned",
	StaticContent: sourceDamnedPng,
}

//go:embed icon_blessed.png
var sourceBlessedPng []byte
var resourceBlessedPng = &fyne.StaticResource{
	StaticName:    "blessed",
	StaticContent: sourceBlessedPng,
}

//go:embed icon_unpaid.png
var sourceUnpaidPng []byte
var resourceUnpaidPng = &fyne.StaticResource{
	StaticName:    "unpaid",
	StaticContent: sourceUnpaidPng,
}

//go:embed icon_unidentified.png
var sourceUnidentifiedPng []byte
var resourceUnidentifiedPng = &fyne.StaticResource{
	StaticName:    "unidentified",
	StaticContent: sourceUnidentifiedPng,
}
