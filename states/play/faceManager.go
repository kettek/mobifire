package play

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

type FaceManager struct {
	window        fyne.Window
	conn          *net.Connection
	handler       *messages.MessageHandler
	pendingImages []int32
	managers      *Managers
}

func NewFaceManager() *FaceManager {
	return &FaceManager{}
}

func (fm *FaceManager) Init() {
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
			manager.(FaceReceiver).OnFaceLoaded(int16(msg.Face), &img)
		}
	})
}

func (fm *FaceManager) SetConnection(conn *net.Connection) {
	fm.conn = conn
}

func (fm *FaceManager) SetHandler(handler *messages.MessageHandler) {
	fm.handler = handler
}

func (fm *FaceManager) SetManagers(managers *Managers) {
	fm.managers = managers
}
