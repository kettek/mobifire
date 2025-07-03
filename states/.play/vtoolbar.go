package play

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Toolbar widget creates a vertical list of tool buttons
type Toolbar struct {
	widget.BaseWidget
	Items []widget.ToolbarItem
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (t *Toolbar) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	r := &toolbarRenderer{toolbar: t, layout: layout.NewVBoxLayout()}
	r.resetObjects()
	return r
}

// Append a new widget.ToolbarItem to the end of this Toolbar
func (t *Toolbar) Append(item widget.ToolbarItem) {
	t.Items = append(t.Items, item)
	t.Refresh()
}

// Prepend a new widget.ToolbarItem to the start of this Toolbar
func (t *Toolbar) Prepend(item widget.ToolbarItem) {
	t.Items = append([]widget.ToolbarItem{item}, t.Items...)
	t.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (t *Toolbar) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

// NewToolbar creates a new toolbar widget.
func NewToolbar(items ...widget.ToolbarItem) *Toolbar {
	t := &Toolbar{Items: items}
	t.ExtendBaseWidget(t)

	t.Refresh()
	return t
}

type toolbarRenderer struct {
	layout  fyne.Layout
	items   []fyne.CanvasObject
	toolbar *Toolbar
}

func (r *toolbarRenderer) MinSize() fyne.Size {
	return r.layout.MinSize(r.items)
}

func (r *toolbarRenderer) Layout(size fyne.Size) {
	r.layout.Layout(r.items, size)
}

func (r *toolbarRenderer) Refresh() {
	r.resetObjects()
	canvas.Refresh(r.toolbar)
}

func (r *toolbarRenderer) resetObjects() {
	r.items = make([]fyne.CanvasObject, 0, len(r.toolbar.Items))
	for _, item := range r.toolbar.Items {
		r.items = append(r.items, item.ToolbarObject())
	}
}

func (r *toolbarRenderer) Destroy() {
}

func (r *toolbarRenderer) Objects() []fyne.CanvasObject {
	return r.items
}
