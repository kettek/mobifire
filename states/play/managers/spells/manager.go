package spells

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/states/play/cfwidgets"
	"github.com/kettek/mobifire/states/play/layouts"
	"github.com/kettek/termfire/messages"
)

type Spell = messages.Spell

type Manager struct {
	window  fyne.Window
	handler *messages.MessageHandler
	spells  []Spell
}

func NewManager() *Manager {
	return &Manager{}
}

func (mgr *Manager) SetWindow(window fyne.Window) {
	mgr.window = window
}

func (mgr *Manager) SetHandler(handler *messages.MessageHandler) {
	mgr.handler = handler
}

func (mgr *Manager) Init() {
	mgr.handler.On(&messages.MessageAddSpell{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageAddSpell)
		for _, spell := range msg.Spells {
			mgr.spells = append(mgr.spells, spell)
		}
	})
	mgr.handler.On(&messages.MessageUpdateSpell{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageUpdateSpell)
		if spell := mgr.getSpell(msg.Tag); spell != nil {
			for _, update := range msg.Fields {
				switch u := update.(type) {
				case messages.MessageUpdateSpellMana:
					spell.Mana = int16(u)
				case messages.MessageUpdateSpellGrace:
					spell.Grace = int16(u)
				case messages.MessageUpdateSpellDamage:
					spell.Damage = int16(u)
				}
			}
		}
	})
	mgr.handler.On(&messages.MessageDeleteSpell{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageDeleteSpell)
		mgr.deleteSpell(msg.Tag)
	})

}

func (mgr *Manager) getSpell(tag int32) *Spell {
	for _, spell := range mgr.spells {
		if spell.Tag == tag {
			return &spell
		}
	}
	return nil
}

func (mgr *Manager) deleteSpell(tag int32) {
	for i, spell := range mgr.spells {
		if spell.Tag == tag {
			mgr.spells = append(mgr.spells[:i], mgr.spells[i+1:]...)
			return
		}
	}
}

func (mgr *Manager) ShowSpellsList() {
	var popup *cfwidgets.PopUp

	info := widget.NewRichTextWithText("...")
	info.Wrapping = fyne.TextWrapWord
	infoScroll := container.NewVScroll(info)

	list := widget.NewList(
		func() int {
			return len(mgr.spells)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("...")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			spell := mgr.spells[i]
			o.(*widget.Label).SetText(spell.Name)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		spell := mgr.spells[id]
		info.Segments = data.TextToRichTextSegments(spell.Description)
		info.Refresh()
		infoScroll.ScrollToTop()
	}

	cnt := container.New(&layouts.Inventory{}, list, infoScroll)

	content := container.NewBorder(nil, nil, nil, nil, cnt)

	dialog := layouts.NewDialog(mgr.window)
	dialog.Full = true

	popup = cfwidgets.NewPopUp(container.New(dialog, content), mgr.window.Canvas())

	popup.ShowCentered(mgr.window.Canvas())
}
