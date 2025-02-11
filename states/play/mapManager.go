package play

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

type MapManager struct {
	window  fyne.Window
	conn    *net.Connection
	handler *messages.MessageHandler

	mb *multiBoard
}

func NewMapManager() *MapManager {
	return &MapManager{}
}

func (mm *MapManager) Init(window fyne.Window, conn *net.Connection, handler *messages.MessageHandler) {
	mm.window = window
	mm.conn = conn
	mm.handler = handler

	// Multiboard seutp.
	faceset := data.CurrentFaceSet()
	mm.mb = newMultiBoard(11, 11, 10, faceset.Width, faceset.Height)
	mm.mb.onSizeChanged = func(rows, cols int) {
		mm.conn.Send(&messages.MessageSetup{
			MapSize: struct {
				Use   bool
				Value string
			}{Use: true, Value: fmt.Sprintf("%dx%d", rows, cols)},
		})
	}

	// Map setup handler.
	mm.handler.On(&messages.MessageSetup{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageSetup)
		if msg.MapSize.Use {
			parts := strings.Split(msg.MapSize.Value, "x")
			if len(parts) != 2 {
				fmt.Println("Invalid map size:", msg.MapSize.Value)
				return
			}
			rows, err := strconv.Atoi(parts[0])
			if err != nil {
				fmt.Println("Invalid map size:", msg.MapSize.Value)
				return
			}
			cols, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid map size:", msg.MapSize.Value)
			}
			mm.mb.SetBoardSize(rows+1, cols+1)
		}
	})

	// Map update handlers.

	mm.handler.On(&messages.MessageMap2{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageMap2)

		for _, m := range msg.Coords {
			if len(m.Data) == 0 {
				// TODO ???
				continue
			}
			for _, c := range m.Data {
				switch d := c.(type) {
				case messages.MessageMap2CoordDataDarkness:
					// TODO
				case messages.MessageMap2CoordDataAnim:
					// TODO
				case messages.MessageMap2CoordDataClear:
					mm.mb.SetCells(m.X, m.Y, nil)
				case messages.MessageMap2CoordDataClearLayer:
					mm.mb.SetCell(m.X, m.Y, int(d.Layer), nil)
				case messages.MessageMap2CoordDataImage:
					if d.FaceNum == 0 {
						mm.mb.SetCell(m.X, m.Y, int(d.Layer), nil)
						continue
					}
					faceImage, ok := data.GetFace(int(d.FaceNum))
					if !ok {
						// FIXME: this is dumb.
						pendingImages = append(pendingImages, boardPendingImage{X: m.X, Y: m.Y, Z: int(d.Layer), Num: int16(d.FaceNum)})
						continue
					}
					mm.mb.SetCell(m.X, m.Y, int(d.Layer), &faceImage)
				}
			}
		}

	})
	mm.handler.On(&messages.MessageNewMap{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		mm.mb.Clear()
	})
}
func (mm *MapManager) CanvasObject() fyne.CanvasObject {
	return mm.mb.container
}
