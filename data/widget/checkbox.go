package widget

import (
	"github.com/kettek/rebui"
	"github.com/kettek/rebui/widgets"
)

type Checkbox struct {
	widgets.Label
	checked  bool
	OnChange func(bool)
}

func (c *Checkbox) HandleGenerate() {
	c.refresh()
}

func (c *Checkbox) HandlePointerPressed(evt rebui.EventPointerPressed) {
	c.Set(!c.checked)
	c.refresh()
}

func (c *Checkbox) Set(v bool) {
	if v == c.checked {
		return
	}
	c.checked = v
	c.refresh()
	if c.OnChange != nil {
		c.OnChange(c.checked)
	}
}

func (c *Checkbox) refresh() {
	if c.checked {
		c.AssignText("☑")
	} else {
		c.AssignText("☐")
	}
}

func (c *Checkbox) Checked() bool {
	return c.checked
}

func init() {
	rebui.RegisterWidget("Checkbox", &Checkbox{})
}
