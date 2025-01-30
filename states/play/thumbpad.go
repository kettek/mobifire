package play

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type thumbpadWidget struct {
	widget.BaseWidget
}

func (r *thumbpadWidget) CreateRenderer() fyne.WidgetRenderer {
	r.ExtendBaseWidget(r)
	return &thumbpadWidgetRenderer{
		rect: &canvas.Rectangle{
			StrokeColor: theme.Color(theme.ColorNameForeground),
			StrokeWidth: 1,
		},
	}
}

func (r *thumbpadWidget) MinSize() fyne.Size {
	r.ExtendBaseWidget(r)
	return r.BaseWidget.MinSize()
}

func (r *thumbpadWidget) Tapped(event *fyne.PointEvent) {
	log.Println("Tapped", event)
}

func (r *thumbpadWidget) TappedSecondary(event *fyne.PointEvent) {
	log.Println("TappedSecondary", event)
}

func (r *thumbpadWidget) Dragged(event *fyne.DragEvent) {
	log.Println("Dragged", event)
}

func (r *thumbpadWidget) DragEnd() {
	log.Println("DragEnd")
}

var _ fyne.WidgetRenderer = (*thumbpadWidgetRenderer)(nil)

type thumbpadWidgetRenderer struct {
	rect *canvas.Rectangle
}

func (r *thumbpadWidgetRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (r *thumbpadWidgetRenderer) Destroy() {
}

func (r *thumbpadWidgetRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
}

func (r *thumbpadWidgetRenderer) MinSize() fyne.Size {
	return r.rect.MinSize()
}

func (r *thumbpadWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rect}
}

func (r *thumbpadWidgetRenderer) Refresh() {
	r.rect.Refresh()
}
