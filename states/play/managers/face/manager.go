package face

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play/managers"
	"github.com/kettek/termfire/messages"
)

type Manager struct {
	window        fyne.Window
	conn          *net.Connection
	handler       *messages.MessageHandler
	pendingImages []int32
	managers      *managers.Managers
}

func NewManager() *Manager {
	return &Manager{}
}

func (fm *Manager) Init() {
	fm.handler.On(&messages.MessageFace2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageFace2)
		if _, ok := data.GetFace(int(msg.Num)); !ok {
			fm.conn.Send(&messages.MessageAskFace{Face: int32(msg.Num)})
		}
	})

	fm.handler.On(&messages.MessageImage2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageImage2)
		data.AddFaceImage(*msg)
		img, _ := data.GetFace(int(msg.Face))
		for _, manager := range fm.managers.GetFaceReceivers() {
			manager.(managers.FaceReceiver).OnFaceLoaded(int16(msg.Face), img)
		}
	})

	fm.handler.On(&messages.MessageAnim{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageAnim)
		data.AddAnim(*msg)
	})
}

func (fm *Manager) SetConnection(conn *net.Connection) {
	fm.conn = conn
}

func (fm *Manager) SetHandler(handler *messages.MessageHandler) {
	fm.handler = handler
}

func (fm *Manager) SetManagers(managers *managers.Managers) {
	fm.managers = managers
}
