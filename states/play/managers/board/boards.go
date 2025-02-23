package board

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/data"
)

type boardPendingImage struct {
	X   int
	Y   int
	Z   int
	Num int16
}

type fyneLayout = fyne.Layout

type multiBoard struct {
	container             *fyne.Container
	boards                []*board
	darkness              [][]uint8
	darknessOverlay       *canvas.Raster
	lastWidth, lastHeight float32
	realWidth, realHeight float32
	cellWidth, cellHeight int
	lastRows, lastCols    int
	onSizeChanged         func(rows, cols int)
}

func newMultiBoard(w, h, count int, cellWidth int, cellHeight int) *multiBoard {
	b := &multiBoard{
		cellWidth:  cellWidth,
		cellHeight: cellHeight,
	}

	var boardContainers []fyne.CanvasObject
	for i := 0; i < count; i++ {
		board := newBoard(w, h, cellWidth, cellHeight)
		boardContainers = append(boardContainers, board.Container)
		b.boards = append(b.boards, board)

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				board.SetFace(x, y, nil)
			}
		}
	}
	for y := 0; y < h; y++ {
		b.darkness = append(b.darkness, make([]uint8, w))
	}

	// darkness overlay
	b.darknessOverlay = canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		cellX := x / cellWidth
		cellY := y / cellHeight

		if cellX < 0 || cellX >= b.boards[0].Width || cellY < 0 || cellY >= b.boards[0].Height {
			return color.Black
		}

		clr := color.NRGBA{0, 0, 0, 0}

		darkness := b.darkness[cellY][cellX]

		if darkness != 0 {
			clr.A = 255 - darkness
		}

		return clr
	})

	b.container = container.New(b, boardContainers...)

	return b
}

func (b *multiBoard) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if b.container.Size().Width != b.lastWidth || b.container.Size().Height != b.lastHeight {
		b.lastWidth = b.container.Size().Width
		b.lastHeight = b.container.Size().Height
		rows, cols := CalculateBoardSize(b.container.Size(), b.cellWidth, b.cellHeight)
		if rows != b.lastRows || cols != b.lastCols {
			b.lastRows = rows
			b.lastCols = cols
			if b.onSizeChanged != nil {
				b.onSizeChanged(rows, cols)
			}
		}
	}
	return fyne.NewSize(b.realWidth, b.realHeight)
}

func (b *multiBoard) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Move(fyne.NewPos(0, 0))
		//o.Move(fyne.NewPos((size.Width-b.realWidth)/2, (size.Height-b.realHeight)/2))
		o.Resize(fyne.NewSize(b.realWidth, b.realHeight))
	}
}

func (b *multiBoard) SetCell(x, y, z int, face *data.FaceImage) {
	b.boards[z].SetFace(x, y, face)
	b.container.Refresh()
}

func (b *multiBoard) SetCells(x, y int, face *data.FaceImage) {
	for _, board := range b.boards {
		board.SetFace(x, y, face)
	}
	b.container.Refresh()
}

func (b *multiBoard) SetDarkness(x, y int, darkness uint8) {
	b.darkness[y][x] = darkness
	b.container.Refresh()
}

func (b *multiBoard) Clear() {
	for _, board := range b.boards {
		board.Clear()
	}
	b.container.Refresh()
}

func (b *multiBoard) ClearBoard(z int) {
	b.boards[z].Clear()
	b.container.Refresh()
}

func CalculateBoardSize(size fyne.Size, cellWidth, cellHeight int) (int, int) {
	rows := size.Width / float32(cellWidth)
	cols := size.Height / float32(cellHeight)
	return int(math.Ceil(float64(rows)) + 1), int(math.Ceil(float64(cols)) + 1)
}

func (b *multiBoard) SetBoardSize(rows, cols int) {
	// We can just fully re-create our boards since a new map is sent when map size changes.
	b.container.RemoveAll()
	for i := 0; i < len(b.boards); i++ {
		b.boards[i] = newBoard(rows, cols, b.cellWidth, b.cellHeight)
		b.container.Add(b.boards[i].Container)
	}
	b.darkness = nil
	for y := 0; y < cols; y++ {
		b.darkness = append(b.darkness, make([]uint8, rows))
	}

	b.realWidth = float32(rows * b.cellWidth)
	b.realHeight = float32(cols * b.cellHeight)

	b.container.Add(b.darknessOverlay)

	b.container.Refresh()
}

func (b *multiBoard) Shift(dx, dy int) {
	for _, board := range b.boards {
		board.Shift(dx, dy)
	}

	if dx == 0 && dy == 0 {
		return
	}
	var updates []darknessUpdate

	for y := 0; y < b.boards[0].Height; y++ {
		for x := 0; x < b.boards[0].Width; x++ {
			if x+dx < 0 || x+dx >= b.boards[0].Width || y+dy < 0 || y+dy >= b.boards[0].Height {
				updates = append(updates, darknessUpdate{x, y, 0})
			} else {
				updates = append(updates, darknessUpdate{x, y, b.darkness[y+dy][x+dx]})
			}
		}
	}

	for _, update := range updates {
		b.darkness[update.y][update.x] = update.darkness
	}
	b.container.Refresh()
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

type darknessUpdate struct {
	x, y     int
	darkness uint8
}

type cellUpdate struct {
	x, y int
	Face *data.FaceImage
}

func (b *board) Shift(dx, dy int) {
	if dx == 0 && dy == 0 {
		return
	}
	var updates []cellUpdate

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			if x+dx < 0 || x+dx >= b.Width || y+dy < 0 || y+dy >= b.Height {
				updates = append(updates, cellUpdate{x, y, nil})
			} else {
				updates = append(updates, cellUpdate{x, y, b.Tiles[y+dy][x+dx].Face})
			}
		}
	}

	for _, update := range updates {
		b.SetFace(update.x, update.y, update.Face)
	}
}

func (b *board) SetFace(x, y int, img *data.FaceImage) {
	b.Tiles[y][x].Face = img
	if img != nil {
		b.Tiles[y][x].SetResource(img)
		b.Tiles[y][x].Show()
	} else {
		b.Tiles[y][x].Hide()
	}
}

func (b *board) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(float32(b.CellWidth*b.Width), float32(b.CellHeight*b.Height))
}

func (b *board) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	for y := b.Height - 1; y >= 0; y-- {
		for x := b.Width - 1; x >= 0; x-- {
			o := b.Tiles[y][x]
			px := float32(x * b.CellWidth)
			py := float32(y * b.CellHeight)
			if o.Face != nil {
				o.Resize(fyne.NewSize(float32(o.Face.Width), float32(o.Face.Height)))
				if o.Face.Width > b.CellWidth {
					px -= float32(o.Face.Width) - float32(b.CellWidth)
				}
				if o.Face.Height > b.CellHeight {
					py -= float32(o.Face.Height) - float32(b.CellHeight)
				}
			} else {
				o.Resize(fyne.NewSize(float32(b.CellWidth), float32(b.CellHeight)))
			}
			o.Move(fyne.NewPos(px, py))
		}
	}
}

func (b *board) Flatten() []fyne.CanvasObject {
	var objects []fyne.CanvasObject
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			objects = append(objects, b.Tiles[y][x])
		}
	}
	return objects
}

func (b *board) Clear() {
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			b.SetFace(x, y, nil)
		}
	}
}

type tile struct {
	widget.Icon
	Face *data.FaceImage
}

// NewTile creates a new tile of the given type
func newTile() *tile {
	t := &tile{}
	t.ExtendBaseWidget(t)
	t.Hide()
	return t
}
