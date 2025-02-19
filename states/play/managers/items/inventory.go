package items

import (
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

// Inventory handles item functionality. It receives handler messages directly and self-filters if the messages do not apply.
type Inventory struct {
	Item     Item // the inventory item itself... this is just used for weight and tag (I think weight is just really for player inventory, as that's the only "floating"/non-contained inventory that uses a weight value afaik).
	Items    []*Item
	handlers []*messages.Handler

	pendingExamineTag int32
	widget            *InventoryWidget

	// I really didn't want to have this field, but whatever, it makes nested calls easier.
	conn *net.Connection

	// This is dynamically set when the main state calls show inventory, so as to provide a callback for actions.
	onSelect func(*Item) bool
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

func (inv *Inventory) setup(handler *messages.MessageHandler, conn *net.Connection) {
	inv.conn = conn
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
	inv.handlers = append(inv.handlers, handler.On(&messages.MessageDrawExtInfo{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageDrawExtInfo)
		if inv.pendingExamineTag == 0 {
			return
		}
		if (msg.Type == messages.MessageTypeCommand && msg.Subtype == messages.SubMessageTypeCommandExamine) || (msg.Type == messages.MessageTypeSpell && msg.Subtype == messages.SubMessageTypeSpellInfo) {
			if item := inv.getItemByTag(inv.pendingExamineTag); item != nil {
				// It's a little hacky, but when we encounter "Examine again", we immediately send an examine request again.
				if strings.HasPrefix(msg.Message, "Examine again") {
					conn.Send(&messages.MessageExamine{
						Tag: inv.pendingExamineTag,
					})
					return
				} else if strings.HasPrefix(msg.Message, "You examine the") {
					// Ignore strings that start with "You examine the", as this is probably a second examine request. This _could_ cause issues if the extra examine information actually contains that string, but until that becomes an issue, so it shall remain as it is.
					return
				}
				item.examineInfo += msg.Message + "\n"
				// Update UI
				if inv.widget != nil && inv.widget.selectedTag() == item.Tag {
					inv.widget.SetExamineInfo(item.examineInfo)
				}
			}
		}
	}))
}

func (inv *Inventory) sortItems() {
	slices.SortStableFunc(inv.Items, func(a, b *Item) int {
		return int(a.Type - b.Type)
	})
}

func (inv *Inventory) handleItem2(msg *messages.MessageItem2) {
	for _, o := range msg.Objects {
		inv.addItem(&Item{ItemObject: o})
	}
	inv.sortItems()
	// Update UI
	if inv.widget != nil {
		inv.widget.itemList.Refresh()
	}
}

func (inv *Inventory) handleDeleteInventory(_ *messages.MessageDeleteInventory) {
	inv.clear()
	// Update UI
	if inv.widget != nil {
		inv.widget.itemList.Refresh()
	}
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
		inv.sortItems()
		// Update UI
		if inv.widget != nil {
			inv.widget.itemList.Refresh()
		}
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
		if inv.widget != nil {
			inv.widget.itemList.Refresh()
		}
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
			selectedTag := inv.widget.selectedTag()
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			if inv.widget != nil && selectedTag == tag {
				inv.widget.refreshSelected()
			}
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

func (inv *Inventory) getItemByName(name string) *Item {
	for _, item := range inv.Items {
		if item.Name == name {
			return item
		}
	}
	return nil
}

func (inv *Inventory) showPopup(window fyne.Window, conn *net.Connection, limited bool) {
	if inv.widget == nil {
		inv.widget = newInventoryWidget(inv, window, conn)
	}
	if limited {
		inv.widget.ShowLimited()
	} else {
		inv.widget.Show()
	}
}

func (inv *Inventory) closePopup() {
	inv.widget.Hide()
}
