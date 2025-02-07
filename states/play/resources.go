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

// type icons

//go:embed type_ammo.png
var sourceAmmoPng []byte
var resourceAmmoPng = &fyne.StaticResource{
	StaticName:    "ammo",
	StaticContent: sourceAmmoPng,
}

//go:embed type_ranged.png
var sourceRangedPng []byte
var resourceRangedPng = &fyne.StaticResource{
	StaticName:    "ranged",
	StaticContent: sourceRangedPng,
}

//go:embed type_weapon.png
var sourceWeaponPng []byte
var resourceWeaponPng = &fyne.StaticResource{
	StaticName:    "weapon",
	StaticContent: sourceWeaponPng,
}

//go:embed type_bodyarmor.png
var sourceBodyarmorPng []byte
var resourceBodyarmorPng = &fyne.StaticResource{
	StaticName:    "bodyarmor",
	StaticContent: sourceBodyarmorPng,
}

//go:embed type_shield.png
var sourceShieldPng []byte
var resourceShieldPng = &fyne.StaticResource{
	StaticName:    "shield",
	StaticContent: sourceShieldPng,
}

//go:embed type_cloak.png
var sourceCloakPng []byte
var resourceCloakPng = &fyne.StaticResource{
	StaticName:    "cloak",
	StaticContent: sourceCloakPng,
}

//go:embed type_container.png
var sourceContainerPng []byte
var resourceContainerPng = &fyne.StaticResource{
	StaticName:    "container",
	StaticContent: sourceContainerPng,
}
