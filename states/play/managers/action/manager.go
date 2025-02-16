package action

import (
	"github.com/kettek/mobifire/states/play/managers"
	"github.com/kettek/mobifire/states/play/managers/inventory"
	"github.com/kettek/mobifire/states/play/managers/skills"
)

// Manager provides management of running and setting actions. Actions are tied to buttons that do some sort of dynamic action that is determined by the player, such as equipping an item, drinking a potion, readying a skill, etc.
type Manager struct {
	entries          []Entry
	skillsManager    *skills.Manager
	inventoryManager *inventory.Manager
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Init() {
	// TODO: Load actions from files?
}

func (m *Manager) Action(i int) *Entry {
	if i < 0 || i >= len(m.entries) {
		return nil
	}
	return &m.entries[i]
}

func (m *Manager) SetManagers(managers *managers.Managers) {
	for _, manager := range *managers {
		if manager, ok := manager.(*skills.Manager); ok {
			m.skillsManager = manager
		}
		if manager, ok := manager.(*inventory.Manager); ok {
			m.inventoryManager = manager
		}
	}
}
