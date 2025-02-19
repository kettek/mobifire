package action

import (
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play/managers"
	"github.com/kettek/mobifire/states/play/managers/items"
	"github.com/kettek/mobifire/states/play/managers/skills"
	"github.com/kettek/mobifire/states/play/managers/spells"
)

// Manager provides management of running and setting actions. Actions are tied to buttons that do some sort of dynamic action that is determined by the player, such as equipping an item, drinking a potion, readying a skill, etc.
type Manager struct {
	app           fyne.App
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
