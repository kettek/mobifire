package cfwidgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// PopUp is a widget that can be shown centered on the screen and dismissed by tapping outside of it.
type PopUp struct {
	*widget.PopUp
	onHide       func()
	overlayShown bool
}

// ShowCentered shows the popup centered on the given canvas.
func (p *PopUp) ShowCentered(canvas fyne.Canvas) {
	ps := p.MinSize()
	ws := canvas.Size()
	x := (ws.Width - ps.Width) / 2
	y := (ws.Height - ps.Height) / 2
	p.ShowAtPosition(fyne.NewPos(x, y))
}

// Show shows the popup.
func (p *PopUp) Show() {
	if !p.overlayShown {
		p.Canvas.Overlays().Add(p)
		p.overlayShown = true
	}
	p.Refresh()
	p.BaseWidget.Show()
}

// Hide hides the popup.
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

// NewPopUp creates a new PopUp widget with the given content and canvas.
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
