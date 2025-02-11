package play

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

type DataManager struct {
	window        fyne.Window
	conn          *net.Connection
	handler       *messages.MessageHandler
	pendingImages []int32
	managers      *[]Manager
}

func NewDataManager(managers *[]Manager) *DataManager {
	return &DataManager{
		managers: managers,
	}
}

func (dm *DataManager) Init(window fyne.Window, conn *net.Connection, handler *messages.MessageHandler) {
	dm.window = window
	dm.conn = conn
	dm.handler = handler

	dm.handler.On(&messages.MessageFace2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageFace2)
		if _, ok := data.GetFace(int(msg.Num)); !ok {
			dm.conn.Send(&messages.MessageAskFace{Face: int32(msg.Num)})
		}
	})

	dm.handler.On(&messages.MessageImage2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageImage2)
		data.AddFaceImage(*msg)
		img, _ := data.GetFace(int(msg.Face))
		for _, manager := range *dm.managers {
			manager.OnFaceLoaded(int16(msg.Face), &img)
		}
	})

}

func (dm *DataManager) OnFaceLoaded(faceID int16, faceImage *data.FaceImage) {
	// Not used, as we cause this.
}
