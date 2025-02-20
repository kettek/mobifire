package action

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play/cfwidgets"
	"github.com/kettek/mobifire/states/play/managers"
	"github.com/kettek/mobifire/states/play/managers/items"
	"github.com/kettek/mobifire/states/play/managers/skills"
	"github.com/kettek/mobifire/states/play/managers/spells"
)

// Manager provides management of running and setting actions. Actions are tied to buttons that do some sort of dynamic action that is determined by the player, such as equipping an item, drinking a potion, readying a skill, etc.
type Manager struct {
	app           fyne.App
	window        fyne.Window
	conn          *net.Connection
	entries       []Entry
	skillsManager *skills.Manager
	itemsManager  *items.Manager
	spellsManager *spells.Manager
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Init() {
	m.loadActions()
}

func (m *Manager) loadActions() error {
	m.entries = nil
	entryStrings := m.app.Preferences().StringList("actions")
	for i, entryString := range entryStrings {
		var entry Entry
		err := entry.Unmarshal(i, []byte(entryString))
		if err != nil {
			fmt.Println("Error unmarshalling entry:", err)
			continue
		}
		m.entries = append(m.entries, entry)
	}
	return nil
}

func (m *Manager) saveActions() error {
	entryStrings := make([]string, len(m.entries))
	for i, entry := range m.entries {
		b, err := entry.MarshalJSON()
		if err != nil {
			fmt.Println("Error marshalling entry:", err)
			continue
		}
		entryStrings[i] = string(b)
	}
	m.app.Preferences().SetStringList("actions", entryStrings)
	return nil
}

func (m *Manager) SetApp(app fyne.App) {
	m.app = app
}

func (m *Manager) SetAction(i int, entry Entry) {
	// Grow as needed.
	if i >= len(m.entries) {
		m.entries = append(m.entries, make([]Entry, i-len(m.entries)+1)...)
	}
	m.entries[i] = entry
	m.saveActions()
}

func (m *Manager) Action(i int) *Entry {
	if i < 0 || i >= len(m.entries) {
		return nil
	}
	return &m.entries[i]
}

func (m *Manager) TriggerAction(i int) {
	if i < 0 || i >= len(m.entries) {
		return
	}
	entry := m.entries[i]
	entry.Trigger(m)
}

func (m *Manager) SetManagers(managers *managers.Managers) {
	for _, manager := range *managers {
		if manager, ok := manager.(*skills.Manager); ok {
			m.skillsManager = manager
		}
		if manager, ok := manager.(*items.Manager); ok {
			m.itemsManager = manager
		}
		if manager, ok := manager.(*spells.Manager); ok {
			m.spellsManager = manager
		}
	}
}

func (m *Manager) SetConnection(conn *net.Connection) {
	m.conn = conn
}

func (m *Manager) SetWindow(window fyne.Window) {
	m.window = window
}

func (m *Manager) AcquireButton(index int) *cfwidgets.AssignableButton {
	entry := m.Action(index)
	if entry == nil {
		m.SetAction(index, Entry{
			Image: data.GetResource("icon_blank.png"),
		})
		entry = m.Action(index)
	}
	if entry.widget != nil {
		return entry.widget
	}
	var button *cfwidgets.AssignableButton
	button = cfwidgets.NewAssignableButton(entry.Image, func() {
		m.TriggerAction(index)
	}, func() {
		currentItem := fyne.NewMenuItem(m.Action(index).TypeString(), nil)
		currentItem.Disabled = true

		itemsMenu := fyne.NewMenuItem("items", nil)
		itemsMenu.ChildMenu = fyne.NewMenu("Sub Actions",
			fyne.NewMenuItem("apply (always)", func() {
				m.itemsManager.ShowLimitedInventory(m.itemsManager.GetPlayerTag(), func(item *items.Item) bool {
					action := Entry{
						Image: data.GetResource("icon_apply.png"),
						Kind: EntryApplyKind{
							ObjectName:      item.Name,
							OnlyIfUnapplied: false,
						},
					}
					if img, ok := data.GetFace(int(item.Face)); ok {
						action.Image = &img
						button.SetIcon(&img)
					}
					m.SetAction(index, action)
					m.itemsManager.CloseInventory(m.itemsManager.GetPlayerTag())
					return true
				})
			}),
			fyne.NewMenuItem("apply (if unapplied)", func() {
				m.itemsManager.ShowLimitedInventory(m.itemsManager.GetPlayerTag(), func(item *items.Item) bool {
					action := Entry{
						Image: data.GetResource("icon_apply.png"),
						Kind: EntryApplyKind{
							ObjectName:      item.Name,
							OnlyIfUnapplied: true,
						},
					}

					if img, ok := data.GetFace(int(item.Face)); ok {
						action.Image = &img
						button.SetIcon(&img)
					}
					m.SetAction(index, action)
					m.itemsManager.CloseInventory(m.itemsManager.GetPlayerTag())
					return true
				})
			}),
			fyne.NewMenuItem("auto-apply and fire", func() {
				m.itemsManager.ShowLimitedInventory(m.itemsManager.GetPlayerTag(), func(item *items.Item) bool {
					action := Entry{
						Image: data.GetResource("icon_apply.png"),
						Kind: EntryApplyKind{
							ObjectName: item.Name,
							Fire:       true,
						},
					}
					if img, ok := data.GetFace(int(item.Face)); ok {
						action.Image = &img
						button.SetIcon(&img)
					}
					m.SetAction(index, action)
					m.itemsManager.CloseInventory(m.itemsManager.GetPlayerTag())
					return true
				})
			}),
		)
		spellsMenu := fyne.NewMenuItem("spells", nil)
		spellsMenu.ChildMenu = fyne.NewMenu("Sub Actions",
			fyne.NewMenuItem("invoke", func() {
				m.spellsManager.ShowSpellsList(func(spell spells.Spell) bool {
					action := Entry{
						Image: data.GetResource("icon_apply.png"),
						Kind: EntrySpellKind{
							Spell: int32(spell.Tag),
							Name:  spell.Name,
						},
					}
					if img, ok := data.GetFace(int(spell.Face)); ok {
						action.Image = &img
						button.SetIcon(&img)
						button.Refresh()
					}
					if spell.Usage > 0 {
						entryWidget := widget.NewEntry()
						dialog.ShowForm("Spell Parameter", "Submit", "Cancel", []*widget.FormItem{
							{Text: "Parameter", Widget: entryWidget},
						}, func(b bool) {
							if b {
								kind := action.Kind.(EntrySpellKind)
								kind.Extra = entryWidget.Text
								action.Kind = kind
								m.SetAction(index, action)
								m.spellsManager.CloseSpellsList()
							}
						}, m.window)
					} else {
						m.SetAction(index, action)
						m.spellsManager.CloseSpellsList()
					}
					return true
				})
			}),
			fyne.NewMenuItem("ready", func() {
				m.spellsManager.ShowSpellsList(func(spell spells.Spell) bool {
					action := Entry{
						Image: data.GetResource("icon_apply.png"),
						Kind: EntrySpellKind{
							Spell: int32(spell.Tag),
							Name:  spell.Name,
							Ready: true,
						},
					}
					if img, ok := data.GetFace(int(spell.Face)); ok {
						action.Image = &img
						button.SetIcon(&img)
					}
					// TODO: Check if readied spells can have parameters... I presume they can?
					if spell.Usage > 0 {
						entryWidget := widget.NewEntry()
						dialog.ShowForm("Spell Parameter", "Submit", "Cancel", []*widget.FormItem{
							{Text: "Parameter", Widget: entryWidget},
						}, func(b bool) {
							if b {
								kind := action.Kind.(EntrySpellKind)
								kind.Extra = entryWidget.Text
								action.Kind = kind
								m.SetAction(index, action)
								m.spellsManager.CloseSpellsList()
							}
						}, m.window)
					} else {
						m.SetAction(index, action)
						m.spellsManager.CloseSpellsList()
					}
					return true
				})
			}),
		)
		skillsMenu := fyne.NewMenuItem("skills", nil)
		skillsMenu.ChildMenu = fyne.NewMenu("Sub Actions",
			fyne.NewMenuItem("use", func() {
				m.skillsManager.ShowSimpleSkillsList(func(id int) {
					action := Entry{
						Image: data.GetResource("icon_apply.png"),
						Kind: EntrySkillKind{
							Skill: int32(id),
							Name:  m.skillsManager.Skill(uint16(id)).Name,
							Ready: false,
						},
					}
					skill := m.skillsManager.Skill(uint16(id))
					if img, ok := data.GetFace(int(skill.Face)); ok {
						action.Image = &img
						button.SetIcon(&img)
					}
					m.SetAction(index, action)
				})
			}),
			fyne.NewMenuItem("ready", func() {
				m.skillsManager.ShowSimpleSkillsList(func(id int) {
					action := Entry{
						Image: data.GetResource("icon_apply.png"),
						Kind: EntrySkillKind{
							Skill: int32(id),
							Name:  m.skillsManager.Skill(uint16(id)).Name,
							Ready: true,
						},
					}
					skill := m.skillsManager.Skill(uint16(id))
					if img, ok := data.GetFace(int(skill.Face)); ok {
						action.Image = &img
						button.SetIcon(&img)
					}
					m.SetAction(index, action)
				})
			}),
		)
		commandsMenu := fyne.NewMenuItem("command", func() {
			entryWidget := widget.NewEntry()
			dialog.ShowForm("Command", "Submit", "Cancel", []*widget.FormItem{
				{Widget: entryWidget},
			}, func(b bool) {
				if b {
					action := Entry{
						Image: entry.Image, // Re-use last image for now... is this a bad idea?
						Kind: EntryCommandKind{
							Command: entryWidget.Text,
							Repeat:  1,
						},
					}
					m.SetAction(index, action)
				}
			}, m.window)
		})

		actions := fyne.NewMenu("Actions", currentItem, fyne.NewMenuItemSeparator(), itemsMenu, spellsMenu, skillsMenu, commandsMenu)
		popup := widget.NewPopUpMenu(actions, m.window.Canvas())
		bpos := button.Position()
		bsize := button.Size()
		popup.ShowAtPosition(fyne.NewPos(bpos.X, bpos.Y+bsize.Height))
	})
	entry.widget = button
	return entry.widget
}
