package play

import (
	"fmt"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/states/play/layouts"
	"github.com/kettek/termfire/messages"
)

// inventory may be a player, a container, or the ground.
type inventory struct {
	name          string
	items         []int32
	sortedItems   []int32
	selectedIndex int
	OnRequest     func(messages.Message)
	list          *widget.List
	info          *widget.RichText
	infoScroll    *container.Scroll
}

func (inv *inventory) showDialog(window fyne.Window) {
	inv.list = widget.NewList(
		func() int {
			return len(inv.sortItems())
		},
		func() fyne.CanvasObject {
			img := canvas.NewImageFromResource(data.GetResource("blank.png"))
			img.FillMode = canvas.ImageFillContain
			img.ScaleMode = canvas.ImageScalePixels
			/*img2 := canvas.NewImageFromResource(resourceBlankPng)
			img2.FillMode = canvas.ImageFillContain
			img2.ScaleMode = canvas.ImageScalePixels*/
			flags := container.NewHBox(widget.NewLabel(""))
			return container.New(&layouts.InventoryEntry{IconSize: data.CurrentFaceSet().Width}, img /*img2,*/, widget.NewLabel(""), flags, widget.NewLabel(""))
			//return container.NewBorder(nil, nil, container.New(&layouts.Height2Width{}, img), otherInfo, widget.NewLabel(""))
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			container := item.(*fyne.Container)
			itemTag := inv.sortItems()[i]
			invItem := GetObject(itemTag)
			if invItem == nil {
				return
			}
			if face, ok := data.GetFace(int(invItem.Face)); ok {
				container.Objects[0].(*canvas.Image).Resource = &face
				container.Objects[0].(*canvas.Image).Refresh()
			}

			/*typeIcon := container.Objects[1].(*canvas.Image)
			if invItem.Type.IsAmmo() {
				typeIcon.Resource = resourceAmmoPng
			} else if invItem.Type.IsRangedWeapon() {
				typeIcon.Resource = resourceRangedPng
			} else if invItem.Type.IsMeleeWeapon() {
				typeIcon.Resource = resourceWeaponPng
			} else if invItem.Type.IsContainer() {
				typeIcon.Resource = resourceContainerPng
			} else if invItem.Type.IsBodyArmor() {
				typeIcon.Resource = resourceBodyarmorPng
			} else if invItem.Type.IsShield() {
				typeIcon.Resource = resourceShieldPng
			} else if invItem.Type.IsCloak() {
				typeIcon.Resource = resourceCloakPng
			} else {
				typeIcon.Resource = resourceBlankPng
			}
			typeIcon.Refresh()*/

			label := container.Objects[1].(*widget.Label)
			label.Importance = widget.MediumImportance
			otherContainer := container.Objects[2].(*fyne.Container)
			otherContainer.RemoveAll()
			if invItem.Flags.Unpaid() {
				img := canvas.NewImageFromResource(data.GetResource("icon_unpaid.png"))
				img.FillMode = canvas.ImageFillOriginal
				img.ScaleMode = canvas.ImageScalePixels
				otherContainer.Objects = append(otherContainer.Objects, img)
				label.Importance = widget.WarningImportance
			}
			if invItem.Flags.Unidentified() {
				img := canvas.NewImageFromResource(data.GetResource("icon_unidentified.png"))
				img.FillMode = canvas.ImageFillOriginal
				img.ScaleMode = canvas.ImageScalePixels
				otherContainer.Objects = append(otherContainer.Objects, img)
				label.Importance = widget.LowImportance
				label.TextStyle.Italic = true
			} else {
				label.TextStyle.Italic = false
			}
			if invItem.Flags.Magic() {
				/*img := canvas.NewImageFromResource(resourceMagicPng)
				img.FillMode = canvas.ImageFillOriginal
				img.ScaleMode = canvas.ImageScalePixels
				otherContainer.Objects = append(otherContainer.Objects, img)*/
				label.Importance = widget.HighImportance
			}
			if invItem.Flags.Damned() {
				img := canvas.NewImageFromResource(data.GetResource("icon_damned.png"))
				img.FillMode = canvas.ImageFillOriginal
				img.ScaleMode = canvas.ImageScalePixels
				otherContainer.Objects = append(otherContainer.Objects, img)
				label.Importance = widget.DangerImportance
			}
			if invItem.Flags.Cursed() {
				img := canvas.NewImageFromResource(data.GetResource("icon_cursed.png"))
				img.FillMode = canvas.ImageFillOriginal
				img.ScaleMode = canvas.ImageScalePixels
				otherContainer.Objects = append(otherContainer.Objects, img)
				label.Importance = widget.DangerImportance
			}
			if invItem.Flags.Blessed() {
				img := canvas.NewImageFromResource(data.GetResource("icon_blessed.png"))
				img.FillMode = canvas.ImageFillOriginal
				img.ScaleMode = canvas.ImageScalePixels
				otherContainer.Objects = append(otherContainer.Objects, img)
				label.Importance = widget.SuccessImportance
			}
			if invItem.Flags.Applied() {
				img := canvas.NewImageFromResource(data.GetResource("icon_applied.png"))
				img.FillMode = canvas.ImageFillOriginal
				img.ScaleMode = canvas.ImageScalePixels
				otherContainer.Objects = append(otherContainer.Objects, img)
				label.TextStyle.Bold = true
			} else {
				label.TextStyle.Bold = false
			}
			if invItem.Flags.Locked() {
				img := canvas.NewImageFromResource(data.GetResource("icon_locked.png"))
				img.FillMode = canvas.ImageFillOriginal
				img.ScaleMode = canvas.ImageScalePixels
				otherContainer.Objects = append(otherContainer.Objects, img)
			}
			if invItem.TotalWeight > 0 {
				kg := float64(invItem.Weight) / 1000
				container.Objects[3].(*widget.Label).SetText(fmt.Sprintf("%.2fkg", kg))
			}

			// SetText after because we adjust styling with the flags checks.
			label.SetText(invItem.GetName())
		},
	)
	inv.list.OnSelected = func(i widget.ListItemID) {
		inv.selectedIndex = i
		tag := inv.getSelectedTag()
		item := GetObject(tag)
		if item != nil {
			inv.request(&messages.MessageExamine{
				Tag: tag,
			})
			inv.info.Segments = data.TextToRichTextSegments(item.examineInfo)
		}
	}
	invToolbar := widget.NewToolbar(
		widget.NewToolbarAction(data.GetResource("icon_apply.png"), func() {
			inv.request(&messages.MessageApply{
				Tag: inv.getSelectedTag(),
			})
		}),
		widget.NewToolbarAction(data.GetResource("icon_marked.png"), func() {
			inv.request(&messages.MessageMark{
				Tag: inv.getSelectedTag(),
			})
		}),
		widget.NewToolbarAction(data.GetResource("icon_locked.png"), func() {
			if item := GetObject(inv.getSelectedTag()); item != nil {
				inv.request(&messages.MessageLock{
					Tag:  inv.getSelectedTag(),
					Lock: !item.Flags.Locked(),
				})
			}
		}),
		widget.NewToolbarAction(data.GetResource("icon_get.png"), func() {
			inv.request(&messages.MessageMove{
				To:   0, // The ground, for now.
				Tag:  inv.getSelectedTag(),
				Nrof: 0, // Perhaps add a drop amount prompt?
			})
		}),
		widget.NewToolbarAction(data.GetResource("icon_get.png"), func() {
			if item := GetObject(inv.getSelectedTag()); item != nil {
				var countEntry *widget.Entry

				countEntry = widget.NewEntry()
				if item.Nrof > 1 {
					countEntry.SetText(fmt.Sprintf("%d", item.Nrof/2))
				} else {
					countEntry.SetText(fmt.Sprintf("%d", item.Nrof))
				}
				dialog.ShowForm("Drop Item", "Drop", "Cancel", []*widget.FormItem{
					widget.NewFormItem("Amount", countEntry),
				}, func(b bool) {
					if !b {
						return
					}
					count, _ := strconv.Atoi(countEntry.Text)
					inv.request(&messages.MessageMove{
						To:   0, // The ground, for now.
						Tag:  inv.getSelectedTag(),
						Nrof: int32(count), // Perhaps add a drop amount prompt?
					})
				}, window)
			}
		}),
	)

	inv.info = widget.NewRichTextWithText("...")
	inv.info.Wrapping = fyne.TextWrapWord
	inv.infoScroll = container.NewVScroll(inv.info)

	invContainer := container.NewBorder(widget.NewLabel(inv.name), invToolbar, nil, nil, container.New(&layouts.Inventory{}, inv.list, inv.infoScroll))

	dialog := layouts.NewDialog(window)
	dialog.Full = true
	popup := widget.NewPopUp(container.New(dialog, invContainer), window.Canvas())
	ps := popup.MinSize()
	ws := window.Canvas().Size()
	x := (ws.Width - ps.Width) / 2
	y := (ws.Height - ps.Height) / 2
	popup.ShowAtPosition(fyne.NewPos(x, y))
}

func (inv *inventory) getSelectedTag() int32 {
	items := inv.sortItems()
	if inv.selectedIndex < 0 || inv.selectedIndex >= len(items) {
		return -1
	}
	return items[inv.selectedIndex]
}

func (inv *inventory) request(msg messages.Message) {
	if inv.OnRequest != nil {
		inv.OnRequest(msg)
	}
}

func (inv *inventory) RefreshInfo() {
	if inv.infoScroll == nil {
		return
	}
	inv.infoScroll.ScrollToTop()
	if item := GetObject(inv.getSelectedTag()); item != nil {
		inv.info.Segments = data.TextToRichTextSegments(item.examineInfo)
	}
	inv.info.Refresh()
}

func (inv *inventory) RefreshList() {
	if inv.list != nil {
		inv.list.Refresh()
	}
}

func (inv *inventory) RefreshItem(tag int32) {
	items := inv.sortItems()
	for i, item := range items {
		if item == tag {
			if inv.selectedIndex == i {
				// FIXME: When an item is identified by an examine, we end up receiving 2 MessageUpdateItem(s) (the "You discover mystic forces on those items")... maybe we can clear the info if that text is discovered.
				// It's the selected one... let's do an examine request.
				/*inv.request(&messages.MessageExamine{
					Tag: inv.getSelectedTag(),
				})*/
			}
			break
		}
	}
	// Update the list, ofc. TODO: We could check if in view, maybe?
	inv.RefreshList()
}

func (inv *inventory) sortItems() []int32 {
	inv.sortedItems = make([]int32, len(inv.items))
	copy(inv.sortedItems, inv.items)
	slices.SortStableFunc(inv.sortedItems, func(a, b int32) int {
		ai := GetObject(a)
		if ai == nil {
			return 1
		}
		bi := GetObject(b)
		if bi == nil {
			return -1
		}
		return int(ai.Type - bi.Type)
	})
	return inv.sortedItems
}

func (inv *inventory) addItem(tag int32) {
	inv.items = append(inv.items, tag)
}

func (inv *inventory) removeItem(tag int32) {
	lastTag := inv.getSelectedTag()
	for i, item := range inv.items {
		if item == tag {
			inv.items = append(inv.items[:i], inv.items[i+1:]...)
			break
		}
	}

	if inv.list != nil {
		if inv.selectedIndex >= len(inv.items) {
			inv.list.Select(len(inv.items) - 1)
			return // Just return, as selecting will cause an examine.
		}
	}

	if lastTag == tag {
		// Request info for new item.
		inv.request(&messages.MessageExamine{
			Tag: inv.getSelectedTag(),
		})
	}
}

func (inv *inventory) clear() {
	inv.items = nil
}

func (inv *inventory) hasItem(tag int32) bool {
	for _, item := range inv.items {
		if item == tag {
			return true
		}
	}
	return false
}

var inventories = map[int32]*inventory{}

func getInventory(id int32) *inventory {
	if inv, ok := inventories[id]; ok {
		return inv
	}
	return nil
}

func acquireInventory(id int32) (*inventory, bool) {
	if inv, ok := inventories[id]; ok {
		return inv, true
	}
	return addNewInventory(id), false
}

func addNewInventory(id int32) *inventory {
	inv := &inventory{}
	inventories[id] = inv
	return inv
}
