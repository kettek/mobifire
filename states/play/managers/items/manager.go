package items

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

// Manager provides functionality for managing items and inventories. This handles network messages pertaining to items as well as the displaying of dialogs for given inventories.
type Manager struct {
	window  fyne.Window
	conn    *net.Connection
	handler *messages.MessageHandler

	inventories []*Inventory
}

// NewManager creates a new items/inventory manager.
func NewManager() *Manager {
	return &Manager{}
}

// SetManager sets the managers for the manager.
func (mgr *Manager) SetWindow(w fyne.Window) {
	mgr.window = w
}

// SetConnection sets the connection for the manager.
func (mgr *Manager) SetConnection(c *net.Connection) {
	mgr.conn = c
}

// SetHandler sets the message handler for the manager.
func (mgr *Manager) SetHandler(h *messages.MessageHandler) {
	mgr.handler = h
}

// Init sets up broad message handling for inventories, such as making them, deleting them, or transferring items between inventories. Further message response functionality is defined in the Inventory type itself.
func (mgr *Manager) Init() {
	mgr.handler.On(&messages.MessagePlayer{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		// We also handle player message, as this gives us the player's name and tag.
		msg := m.(*messages.MessagePlayer)
		if msg.Name == "" {
			return
		}
		inv, _ := mgr.ensureInventory(msg.Tag)
		inv.Item.Name = msg.Name + "'s Inventory" // TODO: Maybe set a field to denote player inventory and determine the title on popup.
		inv.Item.Weight = msg.Weight
		inv.Item.TotalWeight = msg.Weight
	})
	mgr.handler.On(&messages.MessageItem2{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageItem2)
		inv, existed := mgr.ensureInventory(msg.Location)
		// Send the inventory the item2 command if it's only just created, as it'll miss the message otherwise.
		if !existed {
			inv.handleItem2(msg)
		}
	})
	mgr.handler.On(&messages.MessageDeleteItem{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageDeleteItem)
		// Remove any inventories that match the item.
		for _, tag := range msg.Tags {
			mgr.removeInventory(tag)
		}
	})
	mgr.handler.On(&messages.MessageUpdateItem{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageUpdateItem)
		// This feels gross, but if the item got a location update, we need to find it and move it to the appropriate inventory. This would be less iffy with a global objects design, but I've been preferring this inventory-centric approach...
		for _, mf := range msg.Fields {
			switch f := mf.(type) {
			case messages.MessageUpdateItemLocation:
				targetInv, _ := mgr.ensureInventory(int32(f))
				for _, inv := range mgr.inventories {
					if item := inv.getItemByTag(msg.Tag); item != nil {
						targetInv.addItem(item)
						targetInv.sortItems()
						inv.removeItemByTag(msg.Tag)
						break
					}
				}
			}
		}
	})
}

func (mgr *Manager) ensureInventory(tag int32) (*Inventory, bool) {
	for _, inv := range mgr.inventories {
		if inv.Item.Tag == tag {
			return inv, true
		}
	}
	inv := newInventory(tag)
	inv.setup(mgr.handler, mgr.conn)
	mgr.inventories = append(mgr.inventories, inv)
	return inv, false
}

func (mgr *Manager) removeInventory(tag int32) {
	for i, inv := range mgr.inventories {
		if inv.Item.Tag == tag {
			mgr.inventories = append(mgr.inventories[:i], mgr.inventories[i+1:]...)
			inv.cleanup(mgr.handler)
			return
		}
	}
}

func (mgr *Manager) ShowInventory(tag int32, onSelect func(item *Item) bool) {
	for _, inv := range mgr.inventories {
		if inv.Item.Tag == tag {
			inv.onSelect = onSelect
			inv.showPopup(mgr.window, mgr.conn, false)
			return
		}
	}
}

func (mgr *Manager) ShowLimitedInventory(tag int32, onSelect func(item *Item) bool) {
	for _, inv := range mgr.inventories {
		if inv.Item.Tag == tag {
			inv.onSelect = onSelect
			inv.showPopup(mgr.window, mgr.conn, true)
			return
		}
	}
}

func (mgr *Manager) CloseInventory(tag int32) {
	for _, inv := range mgr.inventories {
		if inv.Item.Tag == tag {
			inv.closePopup()
			return
		}
	}
}

// GetItemByTag returns an item by its tag.
func (mgr *Manager) GetItemByTag(tag int32) *Item {
	for _, inv := range mgr.inventories {
		if item := inv.getItemByTag(tag); item != nil {
			return item
		}
	}
	return nil
}

// GetItemByName returns an item by its name.
func (mgr *Manager) GetItemByName(name string) *Item {
	for _, inv := range mgr.inventories {
		if item := inv.getItemByName(name); item != nil {
			return item
		}
	}
	return nil
}
