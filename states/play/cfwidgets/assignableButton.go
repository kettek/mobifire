package cfwidgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// AssignableButton is a widget that has long-tap/rmb click support.
type AssignableButton struct {
	widget.Button
	longAction func()
}

// NewAssignableButton creates a new button with the given icon and action funcs.
func NewAssignableButton(icon fyne.Resource, action func(), longAction func()) *AssignableButton {
	btn := &AssignableButton{}
	btn.ExtendBaseWidget(btn)
	btn.SetIcon(icon)
	btn.Button.OnTapped = action
	btn.longAction = longAction

	return btn
}

// TappedSecondary is called when the button is long-pressed or RMB clicked.
func (a *AssignableButton) TappedSecondary(e *fyne.PointEvent) {
	if a.longAction != nil {
		a.longAction()
	}
}

// TriggerSecondary is called when the button is long-pressed or RMB clicked.
func (a *AssignableButton) TriggerSecondary() {
	if a.longAction != nil {
		a.longAction()
	}
}
