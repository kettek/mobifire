package layouts

import "fyne.io/fyne/v2"

type Inventory struct {
}

func (l *Inventory) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}

func (l *Inventory) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	// 3:2 ratio
	infoWidth := size.Width / 2
	listWidth := size.Width - infoWidth
	objects[0].Resize(fyne.NewSize(listWidth, size.Height))
	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Resize(fyne.NewSize(infoWidth, size.Height))
	objects[1].Move(fyne.NewPos(listWidth, 0))
}
