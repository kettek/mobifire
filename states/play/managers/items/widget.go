package items

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play/cfwidgets"
	"github.com/kettek/mobifire/states/play/layouts"
	"github.com/kettek/termfire/messages"
)

const (
	actionApply = iota
	actionDrop
	actionDropSome
	actionLock
	actionMark
)

type InventoryWidget struct {
	inv    *Inventory
	window fyne.Window
	conn   *net.Connection

	itemList             *widget.List
	itemListScroll       *container.Scroll
	itemInfo             *widget.RichText
	toolbar              *widget.Toolbar
	toolbarActions       [5]*widget.ToolbarAction
	fullContentContainer *fyne.Container
	minContentContainer  *fyne.Container
	popup                *cfwidgets.PopUp
	minimal              bool

	selectedIndex int
}

func newInventoryWidget(inv *Inventory, window fyne.Window, conn *net.Connection) *InventoryWidget {
	iw := &InventoryWidget{
		inv:    inv,
		window: window,
		conn:   conn,
	}

	// Create our item list.
	iw.itemList = widget.NewList(
		func() int {
			return len(inv.Items)
		},
		func() fyne.CanvasObject {
			return iw.makeEntryTemplate()
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			item := inv.Items[i]
			iw.updateEntry(o.(*fyne.Container), item)
		},
	)
	iw.itemList.OnSelected = func(id widget.ListItemID) {
		iw.selectedIndex = id
		iw.refreshSelected()
	}

	// Create our item info.
	iw.itemInfo = widget.NewRichText()
	iw.itemInfo.Wrapping = fyne.TextWrapWord

	// Create our toolbar.
	iw.toolbarActions[actionApply] = widget.NewToolbarAction(data.GetResource("icon_apply.png"), func() {
		conn.Send(&messages.MessageApply{
			Tag: iw.selectedTag(),
		})
	})
	iw.toolbarActions[actionDrop] = widget.NewToolbarAction(data.GetResource("icon_drop.png"), func() {
		conn.Send(&messages.MessageMove{
			To:   0, // da ground
			Tag:  iw.selectedTag(),
			Nrof: 0, // all
		})
	})
	iw.toolbarActions[actionDropSome] = widget.NewToolbarAction(data.GetResource("icon_dropsome.png"), func() {
		item := iw.inv.Items[iw.selectedIndex]

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
			conn.Send(&messages.MessageMove{
				To:   0, // The ground, for now.
				Tag:  item.Tag,
				Nrof: int32(count), // Perhaps add a drop amount prompt?
			})
		}, window)
	})
	iw.toolbarActions[actionLock] = widget.NewToolbarAction(data.GetResource("icon_locked.png"), func() {
		item := iw.inv.Items[iw.selectedIndex]
		conn.Send(&messages.MessageLock{
			Tag:  iw.selectedTag(),
			Lock: !item.Flags.Locked(),
		})
	})
	iw.toolbarActions[actionMark] = widget.NewToolbarAction(data.GetResource("icon_marked.png"), func() {
		conn.Send(&messages.MessageMark{
			Tag: iw.selectedTag(),
		})
	})
	// ... Why do I have to convert this? Am I dumb?
	var actions []widget.ToolbarItem
	for _, action := range iw.toolbarActions {
		actions = append(actions, action)
	}
	iw.toolbar = widget.NewToolbar(actions...)

	// Kinda messy setup...
	iw.itemListScroll = container.NewVScroll(iw.itemInfo)
	listInfoContainer := container.New(&layouts.Inventory{}, iw.itemList, iw.itemListScroll)

	iw.fullContentContainer = container.NewBorder(widget.NewLabel(inv.Item.Name), iw.toolbar, nil, nil, listInfoContainer)
	iw.minContentContainer = container.NewBorder(nil, nil, nil, nil, iw.itemList)

	return iw
}
func (iw *InventoryWidget) makeEntryTemplate() *fyne.Container {
	img := canvas.NewImageFromResource(data.GetResource("blank.png"))
	img.FillMode = canvas.ImageFillContain
	img.ScaleMode = canvas.ImageScalePixels
	/*img2 := canvas.NewImageFromResource(resourceBlankPng)
	img2.FillMode = canvas.ImageFillContain
	img2.ScaleMode = canvas.ImageScalePixels*/
	flags := container.NewHBox(widget.NewLabel(""))
	return container.New(&layouts.InventoryEntry{IconSize: data.CurrentFaceSet().Width}, img /*img2,*/, widget.NewLabel(""), flags, widget.NewLabel(""))
}

func (iw *InventoryWidget) updateEntry(entry *fyne.Container, item *Item) {
	img := entry.Objects[0].(*canvas.Image)
	label := entry.Objects[1].(*widget.Label)
	flagsContainer := entry.Objects[2].(*fyne.Container)
	weightLabel := entry.Objects[3].(*widget.Label)

	if face, ok := data.GetFace(int(item.Face)); ok {
		img.Resource = &face
		img.Refresh()
	}

	label.Importance = widget.MediumImportance
	flagsContainer.RemoveAll()
	if item.Flags.Unpaid() {
		img := canvas.NewImageFromResource(data.GetResource("icon_unpaid.png"))
		img.FillMode = canvas.ImageFillOriginal
		img.ScaleMode = canvas.ImageScalePixels
		flagsContainer.Objects = append(flagsContainer.Objects, img)
		label.Importance = widget.WarningImportance
	}
	if item.Flags.Unidentified() {
		img := canvas.NewImageFromResource(data.GetResource("icon_unidentified.png"))
		img.FillMode = canvas.ImageFillOriginal
		img.ScaleMode = canvas.ImageScalePixels
		flagsContainer.Objects = append(flagsContainer.Objects, img)
		label.Importance = widget.LowImportance
		label.TextStyle.Italic = true
	} else {
		label.TextStyle.Italic = false
	}
	if item.Flags.Magic() {
		/*img := canvas.NewImageFromResource(resourceMagicPng)
		img.FillMode = canvas.ImageFillOriginal
		img.ScaleMode = canvas.ImageScalePixels
		flagsContainer.Objects = append(flagsContainer.Objects, img)*/
		label.Importance = widget.HighImportance
	}
	if item.Flags.Damned() {
		img := canvas.NewImageFromResource(data.GetResource("icon_damned.png"))
		img.FillMode = canvas.ImageFillOriginal
		img.ScaleMode = canvas.ImageScalePixels
		flagsContainer.Objects = append(flagsContainer.Objects, img)
		label.Importance = widget.DangerImportance
	}
	if item.Flags.Cursed() {
		img := canvas.NewImageFromResource(data.GetResource("icon_cursed.png"))
		img.FillMode = canvas.ImageFillOriginal
		img.ScaleMode = canvas.ImageScalePixels
		flagsContainer.Objects = append(flagsContainer.Objects, img)
		label.Importance = widget.DangerImportance
	}
	if item.Flags.Blessed() {
		img := canvas.NewImageFromResource(data.GetResource("icon_blessed.png"))
		img.FillMode = canvas.ImageFillOriginal
		img.ScaleMode = canvas.ImageScalePixels
		flagsContainer.Objects = append(flagsContainer.Objects, img)
		label.Importance = widget.SuccessImportance
	}
	if item.Flags.Applied() {
		img := canvas.NewImageFromResource(data.GetResource("icon_applied.png"))
		img.FillMode = canvas.ImageFillOriginal
		img.ScaleMode = canvas.ImageScalePixels
		flagsContainer.Objects = append(flagsContainer.Objects, img)
		label.TextStyle.Bold = true
	} else {
		label.TextStyle.Bold = false
	}
	if item.Flags.Locked() {
		img := canvas.NewImageFromResource(data.GetResource("icon_locked.png"))
		img.FillMode = canvas.ImageFillOriginal
		img.ScaleMode = canvas.ImageScalePixels
		flagsContainer.Objects = append(flagsContainer.Objects, img)
	}
	if item.TotalWeight > 0 {
		kg := float64(item.Weight) / 1000
		weightLabel.SetText(fmt.Sprintf("%.2fkg", kg))
	}

	// SetText after because we adjust styling with the flags checks.
	label.SetText(item.GetName())
}

func (iw *InventoryWidget) Show() {
	iw.minimal = false
	// I wish we could re-use popup... but for now, we'll just recreate it.
	dialog := layouts.NewDialog(iw.window)
	dialog.Full = true
	if iw.popup != nil {
		iw.popup.Hide()
		iw.popup = nil
	}
	iw.popup = cfwidgets.NewPopUp(container.New(dialog, iw.fullContentContainer), iw.window.Canvas())
	iw.popup.ShowCentered(iw.window.Canvas())
}

func (iw *InventoryWidget) ShowLimited() {
	iw.minimal = true
	// I wish we could re-use popup... but for now, we'll just recreate it.
	dialog := layouts.NewDialog(iw.window)
	dialog.Full = true
	if iw.popup != nil {
		iw.popup.Hide()
		iw.popup = nil
	}
	iw.popup = cfwidgets.NewPopUp(container.New(dialog, iw.minContentContainer), iw.window.Canvas())
	iw.popup.ShowCentered(iw.window.Canvas())
}

func (iw *InventoryWidget) selectedTag() int32 {
	if iw.selectedIndex < 0 || iw.selectedIndex >= len(iw.inv.Items) {
		return -1
	}
	return iw.inv.Items[iw.selectedIndex].Tag
}

func (iw *InventoryWidget) SetExamineInfo(info string) {
	iw.itemInfo.Segments = data.TextToRichTextSegments(info)
	iw.itemInfo.Refresh()
}

func (iw *InventoryWidget) refreshSelected() {
	if iw.selectedIndex < 0 || iw.selectedIndex >= len(iw.inv.Items) {
		return
	}
	item := iw.inv.Items[iw.selectedIndex]
	if iw.inv.onSelect != nil && iw.inv.onSelect(item) {
		// Bail if select is set and it says to stop processing.
		return
	}
	if iw.minimal {
		// Bail early if minimal, as we don't need no requests.
		return
	}

	// Set icons to be specific to the given item type.
	icon := iw.toolbarActions[actionApply].Icon
	if item.Flags.Applied() {
		if item.Flags.Open() {
			// Show "close container" icon
		} else {
			if item.Type.IsContainer() {
				// Show "open container" icon
			} else {
				// Show "unapply" icon
			}
		}
	} else {
		if item.Type.IsDrink() || item.Type.IsPotion() {
			icon = data.GetResource("icon_drink.png")
		} else if item.Type.IsFood() || item.Type.IsFlesh() {
			icon = data.GetResource("icon_eat.png")
		} else if item.Type.IsSpellCastingConsumable() {
			icon = data.GetResource("icon_cast.png")
		} else if item.Type.IsReadable() {
			icon = data.GetResource("icon_read.png")
		} else {
			icon = data.GetResource("icon_apply.png")
		}
	}
	if icon != iw.toolbarActions[actionApply].Icon {
		iw.toolbarActions[actionApply].SetIcon(icon)
	}

	// Set our pending tag so we can re-acquire examineInfo
	iw.inv.pendingExamineTag = item.Tag

	// Set to our existing examine
	iw.itemInfo.Segments = data.TextToRichTextSegments(item.examineInfo)
	iw.itemInfo.Refresh()
	iw.itemListScroll.ScrollToTop() // And scroll back up

	// Clear out old. Maybe we should keep this and only reset (and send request) if invalidated by an update or a failed examine (if we can even get that info -- maybe check for a fail string, if one exists)
	item.examineInfo = ""

	// Send request for item.
	iw.inv.conn.Send(&messages.MessageExamine{
		Tag: item.Tag,
	})

}
