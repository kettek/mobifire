package play

import "github.com/kettek/termfire/messages"

type inventoryItem struct {
	messages.ItemObject
	container *inventory
}

// inventory may be a player, a container, or the ground.
type inventory struct {
	items []*inventoryItem
}

func (inv *inventory) addNewItem(itemObject messages.ItemObject) *inventoryItem {
	item := &inventoryItem{ItemObject: itemObject, container: inv}
	inv.addItem(item)
	return item
}

func (inv *inventory) addItem(item *inventoryItem) {
	item.container = inv
	inv.items = append(inv.items, item)
}

func (inv *inventory) removeItem(item *inventoryItem) {
	for i, invItem := range inv.items {
		if invItem == item {
			invItem.container = nil
			inv.items = append(inv.items[:i], inv.items[i+1:]...)
			return
		}
	}
}

func (inv *inventory) clear() {
	for _, item := range inv.items {
		item.container = nil
	}
	inv.items = nil
}

func (inv *inventory) hasItem(id int32) bool {
	for _, item := range inv.items {
		if item.Tag == id {
			return true
		}
	}
	return false
}

func (inv *inventory) getItem(id int32) *inventoryItem {
	for _, item := range inv.items {
		if item.Tag == id {
			return item
		}
	}
	return nil
}

var inventories = map[int32]*inventory{}

func acquireInventory(id int32) *inventory {
	if inv, ok := inventories[id]; ok {
		return inv
	}
	return addNewInventory(id)
}

func addNewInventory(id int32) *inventory {
	inv := &inventory{}
	inventories[id] = inv
	return inv
}

func findInventoryItem(id int32) *inventoryItem {
	for _, inv := range inventories {
		for _, item := range inv.items {
			if item.Tag == id {
				return item
			}
		}
	}
	return nil
}
