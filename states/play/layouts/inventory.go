package layouts

import (
	"math"

	"fyne.io/fyne/v2"
)

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

type InventoryEntry struct {
	IconSize int
}

func (e *InventoryEntry) MinSize(objects []fyne.CanvasObject) fyne.Size {
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

func (e *InventoryEntry) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	width := float32(0.0)
	// 0: icon
	scale := math.Max(1, math.Floor(float64(e.IconSize)/float64(size.Height)))

	objects[0].Resize(fyne.NewSize(float32(e.IconSize)*float32(scale), size.Height))
	objects[0].Move(fyne.NewPos(0, 0))
	width += objects[0].Size().Width

	// 1: type icon
	/*objects[1].Resize(fyne.NewSize(size.Height/2, size.Height))
	objects[1].Move(fyne.NewPos(width, 0))
	width += size.Height / 2*/

	// Set name pos
	objects[1].Move(fyne.NewPos(width, 0))

	// 3: flags
	objects[2].Resize(fyne.NewSize(size.Height, size.Height))
	width += size.Height

	// 4: weight
	objects[3].Resize(fyne.NewSize(size.Height*2, size.Height))
	width += size.Height * 2

	// 2: name (remaining width)
	objects[1].Resize(fyne.NewSize(size.Width-width, size.Height))

	// Set flags and weight pos
	objects[2].Move(fyne.NewPos(objects[1].Position().X+objects[1].Size().Width, 0))
	objects[3].Move(fyne.NewPos(objects[2].Position().X+objects[2].Size().Width, 0))
}
