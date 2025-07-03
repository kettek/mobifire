package layouts

import (
	"fyne.io/fyne/v2"
)

type Left struct {
}

func (l *Left) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}

func (l *Left) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) != 3 {
		return
	}
	topSize := objects[0].MinSize()
	botSize := objects[2].MinSize()
	centerHeight := size.Height - topSize.Height - botSize.Height

	objects[0].Resize(fyne.NewSize(size.Width, topSize.Height))
	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Resize(fyne.NewSize(size.Width, centerHeight))
	objects[1].Move(fyne.NewPos(0, topSize.Height))
	objects[2].Resize(fyne.NewSize(size.Width, botSize.Height))
	objects[2].Move(fyne.NewPos(0, topSize.Height+centerHeight))
}
