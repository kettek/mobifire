package layouts

import "fyne.io/fyne/v2"

type Dialog struct {
	window fyne.Window
	Full   bool
}

func NewDialog(window fyne.Window) *Dialog {
	return &Dialog{window: window}
}

func (d *Dialog) MinSize(objects []fyne.CanvasObject) fyne.Size {
	size := d.window.Canvas().Size()

	// Not sure if we have a flag somewhere for landscape vs. portrait, but...
	padding := float32(0)
	if d.Full {
		if size.Width > size.Height {
			padding = size.Height / 6
		} else {
			padding = size.Width / 6
		}
	} else {
		if size.Width > size.Height {
			padding = size.Height / 2
		} else {
			padding = size.Width / 2
		}
	}
	return fyne.NewSize(size.Width-padding, size.Height-padding)
}

func (d *Dialog) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(size)
	}
}
