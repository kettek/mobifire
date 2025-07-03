package face

import (
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play/managers"
	"github.com/kettek/termfire/messages"
)

// Manager manages adding faces, animations, and requesting face images.
type Manager struct {
	conn     *net.Connection
	handler  *messages.MessageHandler
	managers *managers.Managers
}

// NewManager creates a new face manager.
func NewManager() *Manager {
	return &Manager{}
}

// Init sets up message handlers.
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
			manager.OnFaceLoaded(int16(msg.Face), img)
		}
	})

	fm.handler.On(&messages.MessageAnim{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageAnim)
		data.AddAnim(*msg)
	})
}

// SetConnection sets the connection for the manager.
func (fm *Manager) SetConnection(conn *net.Connection) {
	fm.conn = conn
}

// SetHandler sets the message handler for the manager.
func (fm *Manager) SetHandler(handler *messages.MessageHandler) {
	fm.handler = handler
}

// SetManagers sets the managers for the face manager.
func (fm *Manager) SetManagers(managers *managers.Managers) {
	fm.managers = managers
}
