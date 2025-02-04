package layouts

import (
	"fyne.io/fyne/v2"
)

type Game struct {
	Board    fyne.CanvasObject
	Left     fyne.CanvasObject
	Right    fyne.CanvasObject
	Messages fyne.CanvasObject
}

func (l *Game) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(0, 0)
}

func (l *Game) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	centerSize := l.Board.MinSize()
	//remainingWidth := size.Width - centerSize.Width
	remainingWidth := size.Width / 2
	leftWidth := remainingWidth / 2
	rightWidth := remainingWidth - leftWidth

	l.Left.Resize(fyne.NewSize(leftWidth, size.Height))
	l.Left.Move(fyne.NewPos(0, 0))
	l.Board.Resize(fyne.NewSize(size.Width, size.Height))
	l.Board.Move(fyne.NewPos((size.Width-centerSize.Width)/2, (size.Height-centerSize.Height)/2))
	l.Right.Resize(fyne.NewSize(rightWidth, size.Height))
	l.Right.Move(fyne.NewPos(size.Width-rightWidth, 0))
	l.Messages.Resize(fyne.NewSize(remainingWidth-8, size.Height/6))
	l.Messages.Move(fyne.NewPos((size.Width-remainingWidth)/2+4, size.Height-size.Height/6))
}
