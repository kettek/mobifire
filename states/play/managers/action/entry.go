package action

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/termfire/messages"
)

type Kind int

const (
	KindNone Kind = iota
	KindApplyFromInventory
	KindReadySkill
	KindUseSkill
	KindReadySpell
	KindCastSpell
)

type Entry struct {
	Kind         Kind
	Image        fyne.Resource
	InventoryTag int32  // The inventory to target. This will generally always be the player.
	ObjectTag    int32  // The ID of the object. This is cached and used when possible. If the ObjectTag does not exist, the Name will be used to do a lookup.
	Name         string // Item name, skill name, or spell name
	OnTrigger    func(messages.Message)
}

func (e Entry) Trigger() {
	// TODO: ???
}
