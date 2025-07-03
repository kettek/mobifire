package layouts

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

type SpellEntry struct {
	IconSize int
	Rect     *canvas.Rectangle // background for skill.
}

func (e *SpellEntry) MinSize(objects []fyne.CanvasObject) fyne.Size {
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

func (e *SpellEntry) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	width := float32(0.0)
	padding := fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameInnerPadding) / 2

	// 0: icon
	scale := math.Max(1, math.Floor(float64(e.IconSize)/float64(size.Height)))

	objects[1].Resize(fyne.NewSize(float32(e.IconSize)*float32(scale), float32(e.IconSize)*float32(scale)))
	objects[1].Move(fyne.NewPos(float32(math.Round(float64(size.Height/2-objects[1].Size().Width/2))), float32(math.Round(float64(size.Height/2-objects[1].Size().Height/2)))))
	width += size.Height + padding

	fw := size.Width - width // skip the icon
	e.Rect.Resize(fyne.NewSize(fw, size.Height))
	e.Rect.Move(fyne.NewPos(width, 0))

	// 1. Set name pos
	objects[2].Move(fyne.NewPos(width, 0))

	// 3: level
	objects[3].Resize(fyne.NewSize(size.Height, size.Height))
	width += size.Height

	// 4: mana/grace
	objects[4].Resize(fyne.NewSize(size.Height, size.Height))
	width += size.Height

	// 5: casting time
	objects[5].Resize(fyne.NewSize(size.Height, size.Height))
	width += size.Height

	// 2: name (remaining width)
	objects[2].Resize(fyne.NewSize(size.Width-width, size.Height))

	objects[3].Move(fyne.NewPos(objects[2].Position().X+objects[2].Size().Width, 0))
	objects[4].Move(fyne.NewPos(objects[3].Position().X+objects[3].Size().Width, 0))
	objects[5].Move(fyne.NewPos(objects[4].Position().X+objects[4].Size().Width, 0))
}
