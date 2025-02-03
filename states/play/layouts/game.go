package layouts

import (
	"fyne.io/fyne/v2"
)

type Game struct {
}

func (l *Game) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}

func (l *Game) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) != 3 {
		return
	}
	centerSize := objects[0].MinSize()
	//remainingWidth := size.Width - centerSize.Width
	remainingWidth := size.Width / 2
	leftWidth := remainingWidth / 2
	rightWidth := remainingWidth - leftWidth

	objects[1].Resize(fyne.NewSize(leftWidth, size.Height))
	objects[1].Move(fyne.NewPos(0, 0))
	objects[0].Resize(fyne.NewSize(size.Width, size.Height))
	objects[0].Move(fyne.NewPos((size.Width-centerSize.Width)/2, (size.Height-centerSize.Height)/2))
	objects[2].Resize(fyne.NewSize(rightWidth, size.Height))
	objects[2].Move(fyne.NewPos(size.Width-rightWidth, 0))
}
