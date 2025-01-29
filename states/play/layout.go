package play

import (
	"fyne.io/fyne/v2"
)

type gameLayout struct {
	window fyne.Window
}

func (l *gameLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}

func (l *gameLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) != 3 {
		return
	}
	centerSize := objects[1].MinSize()
	remainingWidth := size.Width - centerSize.Width
	leftWidth := remainingWidth / 2
	rightWidth := remainingWidth - leftWidth

	objects[0].Resize(fyne.NewSize(leftWidth, size.Height))
	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Resize(centerSize)
	objects[1].Move(fyne.NewPos(leftWidth, (size.Height-centerSize.Height)/2))
	objects[2].Resize(fyne.NewSize(rightWidth, size.Height))
	objects[2].Move(fyne.NewPos(leftWidth+centerSize.Width, 0))
}
