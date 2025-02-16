package items

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

type Manager struct {
	window  fyne.Window
	conn    *net.Connection
	handler *messages.MessageHandler

	// inventories []*Inventory
	// items []*Item
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Init(window fyne.Window, conn *net.Connection, handler *messages.MessageHandler) {
	m.window = window
	m.conn = conn
	m.handler = handler

	m.handler.On(&messages.MessageItem2{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
	})
	m.handler.On(&messages.MessageUpdateItem{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
	})
	m.handler.On(&messages.MessageDeleteItem{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
	})
}
