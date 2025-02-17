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

// SetManagers sets the managers for the manager.
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
	inv.setup(mgr.handler)
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

func (mgr *Manager) ShowPopup(tag int32) {
	for _, inv := range mgr.inventories {
		if inv.Item.Tag == tag {
			inv.ShowPopup(mgr.window, mgr.conn)
			return
		}
	}
}
