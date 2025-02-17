package items

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play/cfwidgets"
)

type InventoryWidget struct {
	window fyne.Window
	conn   *net.Connection

	itemList       *widget.List
	itemInfo       *widget.RichText
	toolbar        *widget.Toolbar
	toolbarActions []*widget.ToolbarAction
	popup          *cfwidgets.PopUp

	selectedTag int32
}

func newInventoryWidget(window fyne.Window, conn *net.Connection) *InventoryWidget {
	iw := &InventoryWidget{
		window: window,
		conn:   conn,
	}
	// TODO: Setup UI.
	return iw
}

func (iw *InventoryWidget) Show() {
	iw.popup.ShowCentered(iw.window.Canvas())
}
