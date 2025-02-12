package layouts

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

type SkillEntry struct {
	IconSize int
}

func (e *SkillEntry) MinSize(objects []fyne.CanvasObject) fyne.Size {
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

func (e *SkillEntry) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	width := float32(0.0)
	// 0: icon
	scale := math.Max(1, math.Floor(float64(e.IconSize)/float64(size.Height)))

	padding := fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameInnerPadding)

	objects[0].Resize(fyne.NewSize(float32(e.IconSize)*float32(scale), float32(e.IconSize)*float32(scale)))
	objects[0].Move(fyne.NewPos(0, float32(math.Round(float64(size.Height/2-objects[0].Size().Height/2)))))
	width += size.Height + padding
	// 1: name (remaining width)
	objects[1].Resize(fyne.NewSize(size.Width-width, size.Height))
	objects[1].Move(fyne.NewPos(width, 0))
}

type FullSkillEntry struct {
	IconSize int
	Perc     float32
	Rect     *canvas.Rectangle
}

func (e *FullSkillEntry) MinSize(objects []fyne.CanvasObject) fyne.Size {
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

func (e *FullSkillEntry) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	width := float32(0.0)
	padding := fyne.CurrentApp().Settings().Theme().Size(theme.SizeNameInnerPadding)

	// 0: icon
	scale := math.Max(1, math.Floor(float64(e.IconSize)/float64(size.Height)))

	objects[0].Resize(fyne.NewSize(float32(e.IconSize)*float32(scale), float32(e.IconSize)*float32(scale)))
	objects[0].Move(fyne.NewPos(float32(math.Round(float64(size.Height/2-objects[0].Size().Width/2))), float32(math.Round(float64(size.Height/2-objects[0].Size().Height/2)))))
	width += size.Height + padding

	fw := size.Width*e.Perc - width // skip the icon
	e.Rect.Resize(fyne.NewSize(fw, size.Height))
	e.Rect.Move(fyne.NewPos(width, 0))

	// Set name pos
	objects[1].Move(fyne.NewPos(width, 0))

	// 3: level
	objects[2].Resize(fyne.NewSize(size.Height, size.Height))
	width += size.Height

	// 4: exp
	objects[3].Resize(fyne.NewSize(size.Height*2, size.Height))
	width += size.Height * 3

	// 2: name (remaining width)
	objects[1].Resize(fyne.NewSize(size.Width-width, size.Height))

	// Set flags and weight pos
	objects[2].Move(fyne.NewPos(objects[1].Position().X+objects[1].Size().Width, 0))
	objects[3].Move(fyne.NewPos(objects[2].Position().X+objects[2].Size().Width, 0))
}
