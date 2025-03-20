package layouts

import (
	"math"

	"fyne.io/fyne/v2"
)

type Icon struct {
	IconSize int
}

func (e *Icon) MinSize(objects []fyne.CanvasObject) fyne.Size {
	size := fyne.NewSize(0, float32(e.IconSize))
	for _, object := range objects {
		min := object.MinSize()
		if min.Height > size.Height {
			size.Height = min.Height
		}
		size.Width += min.Height
	}
	return size
}

func (e *Icon) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	scale := math.Max(1, math.Floor(float64(e.IconSize)/float64(size.Height)))
	objects[0].Resize(fyne.NewSize(float32(e.IconSize)*float32(scale), size.Height))
	objects[0].Move(fyne.NewPos(0, 0))
}
