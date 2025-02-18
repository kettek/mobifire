package items

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play/cfwidgets"
	"github.com/kettek/mobifire/states/play/layouts"
	"github.com/kettek/termfire/messages"
)

type InventoryWidget struct {
	inv    *Inventory
	window fyne.Window
	conn   *net.Connection

	itemList             *widget.List
	itemInfo             *widget.RichText
	toolbar              *widget.Toolbar
	toolbarActions       []*widget.ToolbarAction
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
	// TODO: Setup UI.
	iw.itemList = widget.NewList(
		func() int {
			return len(inv.Items)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Item")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			label.SetText(inv.Items[i].GetName())
		},
	)
	iw.itemList.OnSelected = func(id widget.ListItemID) {
		item := inv.Items[id]
		iw.selectedIndex = id
		if inv.onSelect != nil && inv.onSelect(item) {
			// Bail if select is set and it says to stop processing.
			return
		}
		if iw.minimal {
			// Bail early if minimal, as we don't need no requests.
			return
		}
		iw.inv.pendingExamineTag = item.Tag

		// Set to our existing exmaine
		iw.itemInfo.Segments = data.TextToRichTextSegments(item.examineInfo)

		// Clear out old. Maybe we should keep this and only reset (and send request) if invalidated by an update or a failed examine (if we can even get that info -- maybe check for a fail string, if one exists)
		item.examineInfo = ""

		// Send request for item.
		conn.Send(&messages.MessageExamine{
			Tag: item.Tag,
		})
	}

	iw.itemInfo = widget.NewRichText()
	iw.itemInfo.Wrapping = fyne.TextWrapWord

	// Kinda messy setup...
	listInfoScroll := container.NewVScroll(iw.itemInfo)
	listInfoContainer := container.New(&layouts.Inventory{}, iw.itemList, listInfoScroll)

	iw.fullContentContainer = container.NewBorder(widget.NewLabel(inv.Item.Name), nil, nil, nil, listInfoContainer)
	iw.minContentContainer = container.NewBorder(nil, nil, nil, nil, iw.itemList)

	return iw
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
