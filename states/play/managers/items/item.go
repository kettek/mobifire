package items

import "github.com/kettek/termfire/messages"

type Item struct {
	messages.ItemObject
	examineInfo string
}

func (item *Item) Update(msg *messages.MessageUpdateItem) {
	for _, mf := range msg.Fields {
		switch f := mf.(type) {
		case messages.MessageUpdateItemFlags:
			item.Flags = messages.ItemFlags(f)
		case messages.MessageUpdateItemWeight:
			item.Weight = int32(f)
		case messages.MessageUpdateItemFace:
			item.Face = int32(f)
		case messages.MessageUpdateItemName:
			item.Name = f.Name
			item.PluralName = f.PluralName
		case messages.MessageUpdateItemAnim:
			item.Anim = int16(f)
		case messages.MessageUpdateItemAnimSpeed:
			item.AnimSpeed = int8(f)
		case messages.MessageUpdateItemNrof:
			item.Nrof = int32(f)
		}
	}
}
