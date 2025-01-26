package main

import (
	_ "embed"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

//go:embed blank.base.111.png
var sourceBlankPng []byte
var resourceBlankPng = &fyne.StaticResource{
	StaticName:    "blank",
	StaticContent: sourceBlankPng,
}

func main() {
	a := app.New()
	w := a.NewWindow(("Crossfire Mobile"))

	w.SetContent(widget.NewLabel("Crossfire Mobile"))
	w.Resize(fyne.NewSize(360, 800))

	board := &Board{
		Width:      11,
		Height:     11,
		CellWidth:  32,
		CellHeight: 32,
	}
	for i := 0; i < 11; i++ {
		row := make([]*Tile, 11)
		for j := 0; j < 11; j++ {
			row[j] = NewTile()
		}
		board.Tiles = append(board.Tiles, row)
	}

	grid := container.New(board, board.Flatten()...)

	center := container.New(layout.NewCenterLayout(), grid)

	w.SetContent(center)

	w.ShowAndRun()
}

type Board struct {
	Tiles      [][]*Tile
	Width      int
	Height     int
	CellWidth  int
	CellHeight int
}

func (b *Board) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(float32(b.CellWidth*b.Width), float32(b.CellHeight*b.Height))
}

func (b *Board) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	//pos := fyne.NewPos(0, containerSize.Height-b.MinSize(objects).Height)
	pos := fyne.NewPos(0, 0)
	for i, o := range objects {
		o.Resize(fyne.NewSize(float32(b.Width), float32(b.Height)))
		o.Move(pos)

		if i%b.Width == b.Width-1 {
			_, h := pos.Components()
			pos = fyne.NewPos(0, h).Add(fyne.NewPos(0, float32(b.Height)))
		} else {
			pos = pos.Add(fyne.NewPos(float32(b.Width*2), 0))
		}
	}
}

func (b *Board) Flatten() []fyne.CanvasObject {
	var objects []fyne.CanvasObject
	for _, row := range b.Tiles {
		for _, tile := range row {
			objects = append(objects, tile)
		}
	}
	return objects
}

type Tile struct {
	widget.Icon
	canvas.Image
}

// NewTile creates a new tile of the given type
func NewTile() *Tile {
	t := &Tile{}
	t.ExtendBaseWidget(t)
	t.SetResource(resourceBlankPng)
	return t
}
