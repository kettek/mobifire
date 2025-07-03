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
	lastDir       int8 // Last direction the player issued a movement in.
}

// NewManager creates a new action manager.
func NewManager() *Manager {
	return &Manager{}
}

// Init initializes the action manager.
func (m *Manager) Init() {
	m.loadActions()
}

// loadActions loads the actions from the app preferences.
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

// saveActions saves the actions to the app preferences.
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

// SetApp sets the app for the action manager.
func (m *Manager) SetApp(app fyne.App) {
	m.app = app
}

// SetAction sets the action for the given index.
func (m *Manager) SetAction(i int, entry Entry) {
	// Grow as needed.
	if i >= len(m.entries) {
		m.entries = append(m.entries, make([]Entry, i-len(m.entries)+1)...)
	}

	entry.widget = m.entries[i].widget
	if entry.Image == nil {
		entry.Image = m.entries[i].Image
	}

	if entry.widget != nil {
		entry.widget.SetIcon(entry.Image)
	}

	m.entries[i] = entry
	m.saveActions()
}

// ClearAction clears the action for the given index.
func (m *Manager) ClearAction(i int) {
	m.SetAction(i, Entry{
		Image: data.GetResource("icon_action_blank.png"),
	})
}

// Action returns the action for the given index.
func (m *Manager) Action(i int) *Entry {
	if i < 0 || i >= len(m.entries) {
		return nil
	}
	return &m.entries[i]
}

// TriggerAction triggers the action for the given index.
func (m *Manager) TriggerAction(i int) {
	if i < 0 || i >= len(m.entries) {
		return
	}
	entry := m.entries[i]
	entry.Trigger(m)
}

// SetManagers sets the managers for the action manager.
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

// SetConnection sets the connection for the action manager.
func (m *Manager) SetConnection(conn *net.Connection) {
	m.conn = conn
}

// SetWindow sets the window for the action manager.
func (m *Manager) SetWindow(window fyne.Window) {
	m.window = window
}

// AcquireButton acquires a button for the given index, creating it if it does not exist.
func (m *Manager) AcquireButton(index int) *cfwidgets.AssignableButton {
	entry := m.Action(index)
	if entry == nil {
		m.SetAction(index, Entry{
			Image: data.GetResource("icon_action_blank.png"),
		})
		entry = m.Action(index)
	}
	if entry.widget != nil {
		return entry.widget
	}
	var button *cfwidgets.AssignableButton
	button = cfwidgets.NewAssignableButton(entry.Image, func() {
		if entry.Kind == nil {
			// Trigger the longpress/set action if there is no kind set.
			button.TriggerSecondary()
		} else {
			// Otherwise charge on.
			m.TriggerAction(index)
		}
	}, func() {
		var currentItem *fyne.MenuItem
		if entry.Kind != nil {
			currentItem = fyne.NewMenuItem(m.Action(index).TypeString(), nil)

			actionItems := m.getActionMenuItems(func(e Entry) {
				last := entry
				for last.Next != nil {
					last = last.Next
				}
				last.Next = &e
				m.SetAction(index, *entry) // FIXME: I don't like this deref.
			})

			if entry.Next != nil {
				actionItems = append(actionItems, fyne.NewMenuItemSeparator())
				for e := entry.Next; e != nil; e = e.Next {
					actionItems = append(actionItems, fyne.NewMenuItem(e.TypeString(), nil))
				}
			}

			actionItems = append([]*fyne.MenuItem{fyne.NewMenuItem("clear", func() {
				m.ClearAction(index)
				button.SetIcon(data.GetResource("icon_action_blank.png"))
			}), fyne.NewMenuItemSeparator()}, actionItems...)

			currentItem.ChildMenu = fyne.NewMenu("Sub Actions", actionItems...)
		}

		actionItems := m.getActionMenuItems(func(entry Entry) {
			m.SetAction(index, entry)
		})

		// Prepend current item entry if the entry is empty.
		if entry.Kind != nil {
			actionItems = append([]*fyne.MenuItem{currentItem, fyne.NewMenuItemSeparator()}, actionItems...)
		}

		actions := fyne.NewMenu("Actions", actionItems...)
		popup := widget.NewPopUpMenu(actions, m.window.Canvas())
		bpos := button.Position()
		bsize := button.Size()
		popup.ShowAtPosition(fyne.NewPos(bpos.X, bpos.Y+bsize.Height))
	})
	entry.widget = button
	return entry.widget
}

func (m *Manager) getActionMenuItems(setAction func(Entry)) []*fyne.MenuItem {
	itemsMenu := fyne.NewMenuItem("items", nil)
	itemsMenu.ChildMenu = fyne.NewMenu("Sub Actions",
		fyne.NewMenuItem("apply (or unapply)", func() {
			m.itemsManager.ShowLimitedInventory(m.itemsManager.GetPlayerTag(), func(item *items.Item) bool {
				action := Entry{
					Image: data.GetResource("icon_apply.png"),
					Kind: EntryApplyKind{
						ObjectName:      item.Name,
						OnlyIfUnapplied: false,
					},
				}
				if img, ok := data.GetFace(int(item.Face)); ok {
					action.Image = img
				}
				setAction(action)
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
					action.Image = img
				}
				setAction(action)
				m.itemsManager.CloseInventory(m.itemsManager.GetPlayerTag())
				return true
			})
		}),
		fyne.NewMenuItem("apply and fire", func() {
			m.itemsManager.ShowLimitedInventory(m.itemsManager.GetPlayerTag(), func(item *items.Item) bool {
				action := Entry{
					Image: data.GetResource("icon_apply.png"),
					Kind: EntryApplyKind{
						ObjectName: item.Name,
						Fire:       true,
					},
				}
				if img, ok := data.GetFace(int(item.Face)); ok {
					action.Image = img
				}
				setAction(action)
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
					action.Image = img
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
							setAction(action)
							m.spellsManager.CloseSpellsList()
						}
					}, m.window)
				} else {
					setAction(action)
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
					action.Image = img
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
							setAction(action)
							m.spellsManager.CloseSpellsList()
						}
					}, m.window)
				} else {
					setAction(action)
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
					action.Image = img
				}
				setAction(action)
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
					action.Image = img
				}
				setAction(action)
			})
		}),
	)
	faceMenu := fyne.NewMenuItem("face", nil)
	faceMenuItems := []*fyne.MenuItem{}
	var dirs = []string{"down", "north", "northeast", "east", "southeast", "south", "southwest", "west", "northwest"}
	for _, dir := range dirs {
		faceMenuItems = append(faceMenuItems, fyne.NewMenuItem(dir, func() {
			action := Entry{
				Kind: EntryFaceKind{
					Dir: dir,
				},
			}
			setAction(action)
		}))
	}
	faceMenuItems = append([]*fyne.MenuItem{fyne.NewMenuItem("last step", func() {
		action := Entry{
			Kind: EntryFaceKind{},
		}
		setAction(action)
	})}, faceMenuItems...)
	faceMenu.ChildMenu = fyne.NewMenu("Sub Actions",
		faceMenuItems...,
	)
	commandsMenu := fyne.NewMenuItem("command", nil)
	commandsMenu.ChildMenu = fyne.NewMenu("Sub Actions",
		fyne.NewMenuItem("step forward", func() {
			action := Entry{
				Kind: EntryStepForwardKind{},
			}
			setAction(action)
		}),
		faceMenu,
		fyne.NewMenuItem("custom", func() {
			entryWidget := widget.NewEntry()
			dialog.ShowForm("Command", "Submit", "Cancel", []*widget.FormItem{
				{Widget: entryWidget},
			}, func(b bool) {
				if b {
					action := Entry{
						Kind: EntryCommandKind{
							Command: entryWidget.Text,
							Repeat:  1,
						},
					}
					setAction(action)
				}
			}, m.window)
		}),
	)
	return []*fyne.MenuItem{itemsMenu, spellsMenu, skillsMenu, commandsMenu}
}

// SetDirectionFromString sets the last direction from a string.
func (m *Manager) SetDirectionFromString(str string) {
	switch str {
	case "northwest":
		m.lastDir = 8
	case "north":
		m.lastDir = 1
	case "northeast":
		m.lastDir = 2
	case "east":
		m.lastDir = 3
	case "southeast":
		m.lastDir = 4
	case "south":
		m.lastDir = 5
	case "southwest":
		m.lastDir = 6
	case "west":
		m.lastDir = 7
	}
}

// GetStringFromDirection returns the last direction as a string.
func (m *Manager) GetStringFromDirection() string {
	switch m.lastDir {
	case 1:
		return "north"
	case 2:
		return "northeast"
	case 3:
		return "east"
	case 4:
		return "southeast"
	case 5:
		return "south"
	case 6:
		return "southwest"
	case 7:
		return "west"
	case 8:
		return "northwest"
	}
	return ""
}
