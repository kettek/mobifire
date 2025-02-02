package play

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type boardPendingImage struct {
	X   int
	Y   int
	Z   int
	Num uint16
}

type multiBoard struct {
	fyne.Layout
	container *fyne.Container
	boards    []*board
}

func newMultiBoard(w, h, count int, cellWidth int, cellHeight int) *multiBoard {
	b := &multiBoard{}

	b.Layout = layout.NewStackLayout()

	var boardContainers []fyne.CanvasObject
	for i := 0; i < count; i++ {
		board := newBoard(w, h, cellWidth, cellHeight)
		boardContainers = append(boardContainers, board.Container)
		b.boards = append(b.boards, board)

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				board.SetImage(x, y, nil)
			}
		}
	}

	b.container = container.New(layout.NewStackLayout(), boardContainers...)

	return b
}

func (b *multiBoard) SetCell(x, y, z int, img fyne.Resource) {
	b.boards[z].SetImage(x, y, img)
	b.container.Refresh()
}

func (b *multiBoard) SetCells(x, y int, img fyne.Resource) {
	for _, board := range b.boards {
		board.SetImage(x, y, img)
	}
	b.container.Refresh()
}

func (b *multiBoard) ClearBoard(z int) {
	for y := 0; y < b.boards[z].Height; y++ {
		for x := 0; x < b.boards[z].Width; x++ {
			b.SetCell(x, y, z, nil)
		}
	}
	b.container.Refresh()
}

func (b *multiBoard) CalculateCells(size fyne.Size) (int, int) {
	// ... I don't know if size is scaled or not.
	rows := size.Width / float32(b.boards[0].CellWidth)
	cols := size.Height / float32(b.boards[0].CellHeight)
	return int(math.Round(float64(rows))), int(math.Round(float64(cols)))
}

type board struct {
	Container  *fyne.Container
	Tiles      [][]*tile
	Width      int
	Height     int
	CellWidth  int
	CellHeight int
}

func newBoard(w, h, cellWidth, cellHeight int) *board {
	b := &board{
		Width:      w,
		Height:     h,
		CellWidth:  cellWidth,
		CellHeight: cellHeight,
	}

	for i := 0; i < h; i++ {
		row := make([]*tile, w)
		for j := 0; j < w; j++ {
			row[j] = newTile()
		}
		b.Tiles = append(b.Tiles, row)
	}

	b.Container = container.New(b, b.Flatten()...)

	return b
}

func (b *board) SetImage(x, y int, img fyne.Resource) {
	if img == nil {
		b.SetHidden(x, y, true)
		return
	}
	b.Tiles[y][x].SetResource(img)
	// Automatically unhide.
	b.SetHidden(x, y, false)
}

func (b *board) SetHidden(x, y int, hide bool) {
	b.Tiles[y][x].Hidden = hide
}

func (b *board) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(float32(b.CellWidth*b.Width), float32(b.CellHeight*b.Height))
}

func (b *board) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
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

func (b *board) Flatten() []fyne.CanvasObject {
	var objects []fyne.CanvasObject
	for _, row := range b.Tiles {
		for _, tile := range row {
			objects = append(objects, tile)
		}
	}
	return objects
}

type tile struct {
	widget.Icon
}

// NewTile creates a new tile of the given type
func newTile() *tile {
	t := &tile{}
	t.ExtendBaseWidget(t)
	t.Hide()
	return t
}
