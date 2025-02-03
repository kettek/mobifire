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
	dragging  bool
	startPos  fyne.Position
	deltaPos  fyne.Position
	onCommand func(cmd string)
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

func (r *thumbpadWidget) command(cmd string) {
	if r.onCommand != nil {
		r.onCommand(cmd)
	}
}

func (r *thumbpadWidget) Dragged(event *fyne.DragEvent) {
	if !r.dragging {
		r.startPos = event.Position
		r.dragging = true
	} else {
		r.deltaPos = fyne.NewPos(event.Position.X-r.startPos.X, event.Position.Y-r.startPos.Y)
		dirX := 0
		dirY := 0
		if r.deltaPos.X < -20 {
			dirX = -1
		} else if r.deltaPos.X > 20 {
			dirX = 1
		}
		if r.deltaPos.Y < -20 {
			dirY = -1
		} else if r.deltaPos.Y > 20 {
			dirY = 1
		}
		// TODO: commands need to have some sort of throttling, as we probably don't want to spam the server per drag event.
		if dirX == -1 && dirY == -1 {
			r.command("northwest")
		} else if dirX == 0 && dirY == -1 {
			r.command("north")
		} else if dirX == 1 && dirY == -1 {
			r.command("northeast")
		} else if dirX == -1 && dirY == 0 {
			r.command("west")
		} else if dirX == 1 && dirY == 0 {
			r.command("east")
		} else if dirX == -1 && dirY == 1 {
			r.command("southwest")
		} else if dirX == 0 && dirY == 1 {
			r.command("south")
		} else if dirX == 1 && dirY == 1 {
			r.command("southeast")
		}
		//r.Refresh() ???
	}
}

func (r *thumbpadWidget) DragEnd() {
	r.dragging = false
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
