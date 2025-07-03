package board

import (
	"image/color"
	"math"
	"math/rand"

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

type multiBoard struct {
	container             *fyne.Container
	boards                []*board
	darkness              [][]uint8
	darknessOverlay       *canvas.Raster
	scale                 float32
	lastWidth, lastHeight float32
	realWidth, realHeight float32
	cellWidth, cellHeight int
	lastRows, lastCols    int
	onSizeChanged         func(rows, cols int)
}

func newMultiBoard(w, h, count int, cellWidth int, cellHeight int, scale float32) *multiBoard {
	b := &multiBoard{
		cellWidth:  cellWidth,
		cellHeight: cellHeight,
		scale:      scale,
	}

	var boardContainers []fyne.CanvasObject
	for range count {
		board := newBoard(w, h, cellWidth, cellHeight)
		boardContainers = append(boardContainers, board.Container)
		b.boards = append(b.boards, board)

		for y := range h {
			for x := range w {
				board.SetFace(x, y, nil)
			}
		}
	}
	for range h {
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

	// A lil rasterizer test...
	/*raster := canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		cellX := int(float32(x) / (float32(cellWidth) * scale))
		cellY := int(float32(y) / (float32(cellHeight) * scale))

		if cellX < 0 || cellX >= b.boards[0].Width || cellY < 0 || cellY >= b.boards[0].Height {
			return color.Black
		}

		clr := color.NRGBA{0, 0, 0, 0}

		//darkness := b.darkness[cellY][cellX]

		// Draw from top-down so we can skip checks if pixel alpha is found.
		for i := len(b.boards) - 1; i >= 0; i-- {
			board := b.boards[i]
			if cellX >= 0 && cellX < board.Width && cellY >= 0 && cellY < board.Height {
				tile := board.Tiles[cellY][cellX]
				if tile.Face != nil {
					faceX := int(float32(x)/scale) - cellX*cellWidth
					faceY := int(float32(y)/scale) - cellY*cellHeight
					if faceX >= 0 && faceX < tile.Face.Width && faceY >= 0 && faceY < tile.Face.Height {
						bclr := tile.Face.Image.At(faceX, faceY)
						r, g, b, a := bclr.RGBA()
						if a > 0 {
							clr.R = uint8(r >> 8)
							clr.G = uint8(g >> 8)
							clr.B = uint8(b >> 8)
							clr.A = uint8(a >> 8)
							break
						}
					}
				} else {
					clr.A = 0
				}
			}
		}

			clr.A = darkness

		return clr
	})*/

	//b.container = container.New(b, raster)
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

func (b *multiBoard) Tick(tick uint32) {
	for _, board := range b.boards {
		board.Tick(tick)
	}
}

func (b *multiBoard) SetAnim(x, y, z int, anim *data.Anim, flags int8, speed int8) {
	b.boards[z].SetAnim(x, y, anim, flags, speed)
	b.container.Refresh()
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
	for i := range len(b.boards) {
		b.boards[i] = newBoard(rows, cols, b.cellWidth, b.cellHeight)
		b.container.Add(b.boards[i].Container)
	}
	b.darkness = nil
	for range cols {
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

	for y := range b.boards[0].Height {
		for x := range b.boards[0].Width {
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
	lastTick   uint64
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

	for range h {
		row := make([]*tile, w)
		for j := range w {
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
	// FIXME: Redo all this.
	Anim    *data.Anim
	Frame   int
	Speed   int8
	Counter int
	Flags   int8
}

func (b *board) Shift(dx, dy int) {
	if dx == 0 && dy == 0 {
		return
	}
	var updates []cellUpdate

	for y := range b.Height {
		for x := range b.Width {
			if x+dx < 0 || x+dx >= b.Width || y+dy < 0 || y+dy >= b.Height {
				updates = append(updates, cellUpdate{x, y, nil, nil, 0, 0, 0, 0})
			} else {
				tile := b.Tiles[y+dy][x+dx]
				updates = append(updates, cellUpdate{x, y, tile.Face, tile.Anim, tile.Frame, tile.Speed, tile.Counter, tile.Flags})
			}
		}
	}

	for _, update := range updates {
		t := b.Tiles[update.y][update.x]
		t.Anim = update.Anim
		t.Frame = update.Frame
		t.Speed = update.Speed
		t.Counter = update.Counter
		t.Flags = update.Flags
		b.SetFace(update.x, update.y, update.Face)
	}
}

func (b *board) Tick(t uint32) {
	delta := uint64(t) - b.lastTick
	b.lastTick = uint64(t)
	// FIXME: This is unnecessarily heavy, use some sort of cached animated coords.
	for y, row := range b.Tiles {
		for x, t := range row {
			if t.Anim == nil {
				continue
			}
			t.Counter += int(delta)
			if t.Counter >= int(t.Speed) {
				if t.Flags == 1 { // Randomize
					t.Counter = 0
					t.Frame = rand.Intn(len(t.Anim.Faces))
				} else {
					for t.Counter >= int(t.Speed) {
						t.Counter -= int(t.Speed)
						t.Frame++
						if t.Frame >= len(t.Anim.Faces) {
							t.Frame = 0
						}
					}
				}
				if face, ok := data.GetFace(t.Anim.Faces[t.Frame]); ok {
					b.SetFace(x, y, face)
				}
			}
		}
	}
}

func (b *board) SetFace(x, y int, img *data.FaceImage) {
	if len(b.Tiles) <= y || len(b.Tiles[y]) <= x {
		return
	}
	b.Tiles[y][x].Face = img
	if img != nil {
		b.Tiles[y][x].SetResource(img)
		b.Tiles[y][x].Show()
	} else {
		b.Tiles[y][x].Anim = nil // Clear anim if face is nil, as this _should_ signify a clear.
		b.Tiles[y][x].Hide()
	}
}

func (b *board) SetAnim(x, y int, anim *data.Anim, flags int8, speed int8) {
	b.Tiles[y][x].Anim = anim
	b.Tiles[y][x].Speed = speed
	b.Tiles[y][x].Flags = flags
	// I guess set the face if we can.
	if anim != nil {
		if face, ok := data.GetFace(anim.Faces[0]); ok {
			b.SetFace(x, y, face)
		}
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
	for y := range b.Height {
		for x := range b.Width {
			objects = append(objects, b.Tiles[y][x])
		}
	}
	return objects
}

func (b *board) Clear() {
	for y := range b.Height {
		for x := range b.Width {
			b.SetFace(x, y, nil)
		}
	}
}

type tile struct {
	widget.Icon
	Face    *data.FaceImage
	Anim    *data.Anim
	Frame   int
	Speed   int8
	Counter int
	Flags   int8
}

// NewTile creates a new tile of the given type
func newTile() *tile {
	t := &tile{}
	t.ExtendBaseWidget(t)
	t.Hide()
	return t
}
