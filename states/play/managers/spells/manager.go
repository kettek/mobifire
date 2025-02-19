package spells

import (
	"fmt"
	"image/color"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/states/play/cfwidgets"
	"github.com/kettek/mobifire/states/play/layouts"
	"github.com/kettek/mobifire/states/play/managers"
	"github.com/kettek/mobifire/states/play/managers/skills"
	"github.com/kettek/termfire/messages"
)

type Spell = messages.Spell

type Manager struct {
	window        fyne.Window
	handler       *messages.MessageHandler
	skillsManager *skills.Manager
	spells        []Spell
	skills        []uint8
	popup         *cfwidgets.PopUp
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

func (mgr *Manager) SetManagers(managers *managers.Managers) {
	for _, manager := range *managers {
		if sm, ok := manager.(*skills.Manager); ok {
			mgr.skillsManager = sm
		}
	}
}

func (mgr *Manager) Init() {
	mgr.handler.On(&messages.MessageAddSpell{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageAddSpell)
		for _, spell := range msg.Spells {
			mgr.spells = append(mgr.spells, spell)
		}
		mgr.sortSpells()
		mgr.getSkills()
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

func (mgr *Manager) sortSpells() {
	// For now, sort by skill.
	slices.SortStableFunc(mgr.spells, func(a, b Spell) int {
		return int(a.Skill) - int(b.Skill)
	})
}

func (mgr *Manager) getSkills() {
	mgr.skills = nil
	for _, spell := range mgr.spells {
		found := false
		for _, skill := range mgr.skills {
			if skill == spell.Skill {
				found = true
				break
			}
		}
		if !found {
			mgr.skills = append(mgr.skills, spell.Skill)
		}
	}
}

func (mgr *Manager) getSpellsBySkill(skill uint8) []Spell {
	var spells []Spell
	for _, spell := range mgr.spells {
		if spell.Skill == skill {
			spells = append(spells, spell)
		}
	}
	return spells
}

func (mgr *Manager) ShowSpellsList(onSelect func(spell Spell) bool) {
	info := widget.NewRichTextWithText("...")
	info.Wrapping = fyne.TextWrapWord
	infoScroll := container.NewVScroll(info)

	makeListForSpells := func(spells []Spell) *widget.List {
		list := widget.NewList(
			func() int {
				return len(spells)
			},
			func() fyne.CanvasObject {
				rect := canvas.NewRectangle(color.NRGBA{255, 255, 255, 100})
				return container.New(&layouts.SpellEntry{IconSize: data.CurrentFaceSet().Width, Rect: rect}, rect, &canvas.Image{}, widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""))
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				spell := spells[i]
				rect := o.(*fyne.Container).Objects[0].(*canvas.Rectangle)
				icon := o.(*fyne.Container).Objects[1].(*canvas.Image)
				name := o.(*fyne.Container).Objects[2].(*widget.Label)
				level := o.(*fyne.Container).Objects[3].(*widget.Label)
				mana := o.(*fyne.Container).Objects[4].(*widget.Label)
				castingTime := o.(*fyne.Container).Objects[5].(*widget.Label)

				if face, ok := data.GetFace(int(spell.Face)); ok {
					icon.Resource = &face
				} else {
					icon.Resource = data.GetResource("blank.png")
				}
				icon.Refresh()
				name.SetText(spell.Name)
				level.SetText(fmt.Sprintf("%d", spell.Level))
				if spell.Mana > 0 {
					mana.SetText(fmt.Sprintf("%d", spell.Mana))
				} else {
					mana.SetText(fmt.Sprintf("%d", spell.Grace))
				}
				castingTime.SetText(fmt.Sprintf("%d", spell.CastingTime))

				// This is lame, but for now, hardcode check spell schools.
				skill := mgr.skillsManager.Skill(uint16(spell.Skill))
				if skill.Name == "pyromancy" {
					rect.FillColor = color.NRGBA{200, 0, 0, 100}
				} else if skill.Name == "evocation" {
					rect.FillColor = color.NRGBA{0, 0, 200, 100}
				} else if skill.Name == "sorcery" {
					rect.FillColor = color.NRGBA{200, 0, 200, 100}
				} else if skill.Name == "summoning" {
					rect.FillColor = color.NRGBA{0, 200, 0, 100}
				} else if skill.Name == "praying" {
					rect.FillColor = color.NRGBA{200, 200, 0, 100}
				} else {
					rect.FillColor = color.NRGBA{200, 200, 200, 100}
				}
				rect.Refresh()
			},
		)
		list.OnSelected = func(id widget.ListItemID) {
			if onSelect != nil && onSelect(spells[id]) {
				return
			}
			spell := spells[id]
			skill := mgr.skillsManager.Skill(uint16(spell.Skill))
			text := fmt.Sprintf("[b]%s[/b]\n\n[b]Skill:[/b] %s\n[b]Level:[/b] %d\n[b]Mana:[/b] %d\n[b]Casting Time:[/b] %d\n\n%s", spell.Name, skill.Name, spell.Level, spell.Mana, spell.CastingTime, spell.Description)
			info.Segments = data.TextToRichTextSegments(text)
			info.Refresh()
			infoScroll.ScrollToTop()
		}
		return list
	}

	var skillTabs []*container.TabItem
	for _, skillID := range mgr.skills {
		skill := mgr.skillsManager.Skill(uint16(skillID))
		skillTabs = append(skillTabs, container.NewTabItem("", makeListForSpells(mgr.getSpellsBySkill(skillID))))
		if face, ok := data.GetFace(int(skill.Face)); ok {
			skillTabs[len(skillTabs)-1].Icon = &face
		}
	}

	tabs := container.NewAppTabs(skillTabs...)
	/*tabs.OnSelected = func(tab *container.TabItem) {
		list := tab.Content.(*widget.List)
		list.Select(0)
	}*/

	cnt := container.New(&layouts.Inventory{}, tabs, infoScroll)

	content := container.NewBorder(nil, nil, nil, nil, cnt)

	dialog := layouts.NewDialog(mgr.window)
	dialog.Full = true

	mgr.popup = cfwidgets.NewPopUp(container.New(dialog, content), mgr.window.Canvas())

	mgr.popup.ShowCentered(mgr.window.Canvas())
}

func (mgr *Manager) CloseSpellsList() {
	if mgr.popup != nil {
		mgr.popup.Hide()
	}
}
