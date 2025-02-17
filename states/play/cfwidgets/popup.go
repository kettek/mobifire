package cfwidgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type PopUp struct {
	*widget.PopUp
	onHide       func()
	overlayShown bool
}

func (p *PopUp) ShowCentered(canvas fyne.Canvas) {
	ps := p.MinSize()
	ws := canvas.Size()
	x := (ws.Width - ps.Width) / 2
	y := (ws.Height - ps.Height) / 2
	p.ShowAtPosition(fyne.NewPos(x, y))
}

func (p *PopUp) Show() {
	if !p.overlayShown {
		p.Canvas.Overlays().Add(p)
		p.overlayShown = true
	}
	p.Refresh()
	p.BaseWidget.Show()
}

func (p *PopUp) Hide() {
	if p.onHide != nil {
		p.onHide()
	}
	if p.overlayShown {
		p.Canvas.Overlays().Remove(p)
		p.overlayShown = false
	}
	p.PopUp.Hide()
}

func NewPopUp(content fyne.CanvasObject, canvas fyne.Canvas) *PopUp {
	p := &PopUp{
		PopUp: widget.NewPopUp(content, canvas),
	}
	return p
}

// Tapped is called when the user taps the popUp background - if not modal then dismiss this widget
func (p *PopUp) Tapped(_ *fyne.PointEvent) {
	p.Hide()
}

// TappedSecondary is called when the user right/alt taps the background - if not modal then dismiss this widget
func (p *PopUp) TappedSecondary(_ *fyne.PointEvent) {
	p.Hide()
}
