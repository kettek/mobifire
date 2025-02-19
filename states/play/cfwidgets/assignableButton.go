package cfwidgets

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type AssignableButton struct {
	widget.Button
	longAction func()
}

func NewAssignableButton(icon fyne.Resource, action func(), longAction func()) *AssignableButton {
	btn := &AssignableButton{}
	btn.ExtendBaseWidget(btn)
	btn.SetIcon(icon)
	btn.Button.OnTapped = action
	btn.longAction = longAction

	return btn
}

func (a *AssignableButton) TappedSecondary(e *fyne.PointEvent) {
	fmt.Println("wuh")
	if a.longAction != nil {
		a.longAction()
	}
}
