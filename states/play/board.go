package play

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type MultiBoard struct {
	fyne.Layout
	container *fyne.Container
	boards    []*Board
}

func NewMultiBoard(w, h, count int) *MultiBoard {
	b := &MultiBoard{}

	b.Layout = layout.NewStackLayout()

	var boardContainers []fyne.CanvasObject
	for i := 0; i < count; i++ {
		board := NewBoard(w, h)
		boardContainers = append(boardContainers, board.Container)
		b.boards = append(b.boards, board)

		if i == 0 {
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					board.SetImage(x, y, resourceBlankPng)
					board.SetHidden(x, y, false)
				}
			}
		} else if i == 1 {
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					if x == 0 || y == 0 || x == w-1 || y == h-1 {
						board.SetImage(x, y, resourceMarkPng)
						board.SetHidden(x, y, false)
					}
				}
			}
		}
	}

	b.container = container.New(layout.NewStackLayout(), boardContainers...)

	return b
}

type Board struct {
	Container  *fyne.Container
	Tiles      [][]*Tile
	Width      int
	Height     int
	CellWidth  int
	CellHeight int
}

func NewBoard(w, h int) *Board {
	b := &Board{
		Width:      w,
		Height:     h,
		CellWidth:  32,
		CellHeight: 32,
	}

	for i := 0; i < h; i++ {
		row := make([]*Tile, w)
		for j := 0; j < w; j++ {
			row[j] = NewTile()
		}
		b.Tiles = append(b.Tiles, row)
	}

	b.Container = container.New(b, b.Flatten()...)

	return b
}

func (b *Board) SetImage(x, y int, img *fyne.StaticResource) {
	b.Tiles[y][x].SetResource(img)
	// Automatically unhide.
	b.SetHidden(x, y, false)
}

func (b *Board) SetHidden(x, y int, hide bool) {
	b.Tiles[y][x].Hidden = hide
}

func (b *Board) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(float32(b.CellWidth*b.Width), float32(b.CellHeight*b.Height))
}

func (b *Board) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for i, o := range objects {
		o.Resize(fyne.NewSize(float32(b.CellWidth), float32(b.CellHeight)))
		o.Move(pos)

		if i%b.Width == b.Width-1 {
			_, h := pos.Components()
			pos = fyne.NewPos(0, h).Add(fyne.NewPos(0, float32(b.CellHeight)))
		} else {
			pos = pos.Add(fyne.NewPos(float32(b.CellWidth), 0))
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
}

// NewTile creates a new tile of the given type
func NewTile() *Tile {
	t := &Tile{}
	t.ExtendBaseWidget(t)
	t.Hide()
	return t
}
