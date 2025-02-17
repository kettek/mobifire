package items

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

// Inventory handles item functionality. It receives handler messages directly and self-filters if the messages do not apply.
type Inventory struct {
	Item     Item // the inventory item itself... this is just used for weight and tag (I think weight is just really for player inventory, as that's the only "floating"/non-contained inventory that uses a weight value afaik).
	Items    []*Item
	handlers []*messages.Handler

	widget *InventoryWidget
}

func newInventory(tag int32) *Inventory {
	return &Inventory{
		Item: Item{
			ItemObject: messages.ItemObject{
				Tag: tag,
			},
		},
	}
}

func (inv *Inventory) clear() {
	inv.Items = nil
}

func (inv *Inventory) setup(handler *messages.MessageHandler) {
	inv.handlers = append(inv.handlers, handler.On(&messages.MessageItem2{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageItem2)
		if msg.Location != inv.Item.Tag {
			return
		}
		inv.handleItem2(msg)
	}))
	inv.handlers = append(inv.handlers, handler.On(&messages.MessageDeleteInventory{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageDeleteInventory)
		if msg.Tag != inv.Item.Tag {
			return
		}
		inv.handleDeleteInventory(msg)
	}))
	inv.handlers = append(inv.handlers, handler.On(&messages.MessageUpdateItem{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageUpdateItem)
		inv.handleUpdateItem(msg)
	}))
	inv.handlers = append(inv.handlers, handler.On(&messages.MessageDeleteItem{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageDeleteItem)
		inv.handleDeleteItem(msg)
	}))
}

func (inv *Inventory) handleItem2(msg *messages.MessageItem2) {
	for _, o := range msg.Objects {
		inv.addItem(&Item{ItemObject: o})
	}
	// Update UI
}

func (inv *Inventory) handleDeleteInventory(_ *messages.MessageDeleteInventory) {
	inv.clear()
	// Update UI
}

func (inv *Inventory) handleUpdateItem(msg *messages.MessageUpdateItem) {
	var changed bool
	if inv.Item.Tag == msg.Tag {
		// First handle if it's ourself -- I think this is just for weights?
		inv.Item.Update(msg)
		changed = true
	} else {
		// Otherwise check if we have the item and update it.
		for _, item := range inv.Items {
			if item.Tag == msg.Tag {
				item.Update(msg)
				changed = true
				break
			}
		}
	}

	if changed {
		// Update UI
	}
}

func (inv *Inventory) handleDeleteItem(msg *messages.MessageDeleteItem) {
	var changed bool
	for _, tag := range msg.Tags {
		if inv.removeItemByTag(tag) {
			changed = true
		}
	}
	if changed {
		// Update UI
	}
}

func (inv *Inventory) cleanup(handler *messages.MessageHandler) {
	for _, h := range inv.handlers {
		handler.Off(h)
	}
	inv.handlers = nil
}

func (inv *Inventory) addItem(item *Item) {
	inv.Items = append(inv.Items, item)
}

func (inv *Inventory) removeItemByTag(tag int32) bool {
	for i, item := range inv.Items {
		if item.Tag == tag {
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			return true
		}
	}
	return false
}

func (inv *Inventory) getItemByTag(tag int32) *Item {
	for _, item := range inv.Items {
		if item.Tag == tag {
			return item
		}
	}
	return nil
}

func (inv *Inventory) showPopup(window fyne.Window, conn *net.Connection) {
	if inv.widget == nil {
		inv.widget = newInventoryWidget(window, conn)
	}
	inv.widget.Show()
}
