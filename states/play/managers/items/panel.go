package items

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play/layouts"
)

type InventoryPanel struct {
	inv    *Inventory
	window fyne.Window
	conn   *net.Connection

	itemList       *widget.List
	itemListScroll *container.Scroll

	container *fyne.Container
}

func newInventoryPanel(inv *Inventory, window fyne.Window, conn *net.Connection) *InventoryPanel {
	panel := &InventoryPanel{
		inv:    inv,
		window: window,
		conn:   conn,
	}

	panel.itemList = widget.NewList(
		func() int { return len(inv.Items) },
		func() fyne.CanvasObject {
			img := canvas.NewImageFromResource(data.GetResource("blank.png"))
			img.FillMode = canvas.ImageFillContain
			//img.ScaleMode = canvas.ImageScalePixels
			icon := container.New(&layouts.Icon{IconSize: data.CurrentFaceSet().Width}, img)
			text := widget.NewLabel("")
			return container.NewHBox(icon, text)
		},
		func(i int, o fyne.CanvasObject) {
			item := inv.Items[i]

			box := o.(*fyne.Container)
			icon := box.Objects[0].(*fyne.Container)
			img := icon.Objects[0].(*canvas.Image)
			if face, ok := data.GetFace(int(item.Face)); ok {
				img.Resource = face
				img.Refresh()
			}
			icon.Refresh()
			text := box.Objects[1].(*widget.Label)
			text.SetText(item.Name)
			text.Refresh()
		},
	)

	panel.itemListScroll = container.NewScroll(panel.itemList)
	panel.container = container.NewBorder(nil, nil, nil, nil, panel.itemListScroll)

	return panel
}

func (panel *InventoryPanel) Container() *fyne.Container {
	return panel.container
}
