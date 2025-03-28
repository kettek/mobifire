package board

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

// Manager manages the game board and handles incoming messages to update the board state.
type Manager struct {
	window  fyne.Window
	conn    *net.Connection
	handler *messages.MessageHandler

	mb *multiBoard

	localTicker bool

	pendingImages []boardPendingImage
}

// NewManager creates a new board manager.
func NewManager() *Manager {
	return &Manager{}
}

// SetConnection sets the connection for the manager.
func (mm *Manager) SetConnection(conn *net.Connection) {
	mm.conn = conn
}

// SetHandler sets the message handler for the manager.
func (mm *Manager) SetHandler(handler *messages.MessageHandler) {
	mm.handler = handler
}

// SetWindow sets the window for the manager.
func (mm *Manager) SetWindow(window fyne.Window) {
	mm.window = window
}

// OnFaceLoaded handles the loading of a face image and updates the board accordingly.
func (mm *Manager) OnFaceLoaded(faceID int16, faceImage *data.FaceImage) {
	for i := len(mm.pendingImages) - 1; i >= 0; i-- {
		if mm.pendingImages[i].Num == faceID {
			mm.mb.SetCell(mm.pendingImages[i].X, mm.pendingImages[i].Y, mm.pendingImages[i].Z, faceImage)
			mm.pendingImages = append(mm.pendingImages[:i], mm.pendingImages[i+1:]...)
		}
	}
}

// PreInit sets up the board and sends a setup message for map size.
func (mm *Manager) PreInit() {
	// Request a board size of the proper dimensions we want.
	w, h := CalculateBoardSize(mm.window.Canvas().Size(), data.CurrentFaceSet().Width, data.CurrentFaceSet().Height)
	mm.conn.Send(&messages.MessageSetup{
		MapSize: struct {
			Use   bool
			Value string
		}{
			Use:   true,
			Value: fmt.Sprintf("%dx%d", w, h),
		},
	})
}

// Init initializes the board manager and sets up message handling for the board.
func (mm *Manager) Init() {
	// Multiboard setup.
	faceset := data.CurrentFaceSet()
	mm.mb = newMultiBoard(11, 11, 10, faceset.Width, faceset.Height, mm.window.Canvas().Scale())
	mm.mb.onSizeChanged = func(rows, cols int) {
		mm.conn.Send(&messages.MessageSetup{
			MapSize: struct {
				Use   bool
				Value string
			}{Use: true, Value: fmt.Sprintf("%dx%d", rows, cols)},
		})
	}

	// Manager setup handler.
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
			mm.mb.SetBoardSize(rows+2, cols+2) // FIXME: CF can send beyond scope of what we can see... I'm not certain how to fix this with how Fyne does widget rendering... Maybe use canvas.Raster for the board...?
		}
	})

	// Manager update handlers.

	mm.handler.On(&messages.MessageMap2{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageMap2)

		for _, m := range msg.Coords {
			if m.Type == messages.MessageMap2CoordTypeScrollInformation {
				mm.mb.Shift(int(m.X), int(m.Y))
			}

			if len(m.Data) == 0 {
				// TODO ???
				continue
			}
			for _, c := range m.Data {
				switch d := c.(type) {
				case messages.MessageMap2CoordDataDarkness:
					mm.mb.SetDarkness(m.X, m.Y, uint8(d.Darkness))
				case messages.MessageMap2CoordDataAnim:
					anim := data.GetAnim(int(d.Anim))
					mm.mb.SetAnim(m.X, m.Y, int(d.Layer), anim, d.Flags, d.Speed)
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
						mm.pendingImages = append(mm.pendingImages, boardPendingImage{X: m.X, Y: m.Y, Z: int(d.Layer), Num: int16(d.FaceNum)})
						continue
					}
					mm.mb.SetCell(m.X, m.Y, int(d.Layer), faceImage)
				}
			}
		}

	})
	mm.handler.On(&messages.MessageNewMap{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		mm.mb.Clear()
	})

	// Manual ticker.
	if mm.localTicker {
		go func() {
			t := time.NewTicker(time.Microsecond * 120000)
			tick := uint32(0)
			for {
				<-t.C
				mm.mb.Tick(tick)
				tick++
			}
		}()
	} else {
		tick := messages.MessageTick(0)
		mm.handler.On(&tick, nil, func(m messages.Message, mf *messages.MessageFailure) {
			mm.mb.Tick(uint32(*(m.(*messages.MessageTick))))
		})
	}
}

// CanvasObject returns the canvas object for the board manager.
func (mm *Manager) CanvasObject() fyne.CanvasObject {
	return mm.mb.container
}
