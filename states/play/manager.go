package play

import (
	"fmt"
	"reflect"

	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

type Managers []Manager

func (m *Managers) Add(manager Manager) {
	*m = append(*m, manager)
}

func (m *Managers) Remove(manager Manager) {
	for i, v := range *m {
		if v == manager {
			*m = append((*m)[:i], (*m)[i+1:]...)
			return
		}
	}
}

func (m *Managers) Init(window fyne.Window, conn *net.Connection, handler *messages.MessageHandler) {
	for _, manager := range *m {
		if manager, ok := manager.(ManagersAccessor); ok {
			manager.SetManagers(m)
		}
		if manager, ok := manager.(WindowAccessor); ok {
			manager.SetWindow(window)
		}
		if manager, ok := manager.(ConnectionAccessor); ok {
			manager.SetConnection(conn)
		}
		if manager, ok := manager.(HandlerAccessor); ok {
			manager.SetHandler(handler)
		}
		manager.Init()
	}
}

func (m *Managers) GetByType(manager Manager) Manager {
	for _, v := range *m {
		fmt.Printf("%T %T %+v %+v\n", v, manager, v, manager)
		if reflect.TypeOf(v) == reflect.TypeOf(manager) {
			return v
		}
	}
	return nil
}

func (m *Managers) GetFaceReceivers() []FaceReceiver {
	var managers []FaceReceiver
	for _, v := range *m {
		if v, ok := v.(FaceReceiver); ok {
			managers = append(managers, v)
		}
	}
	return managers
}

type Manager interface {
	Init()
}

type FaceReceiver interface {
	Manager
	OnFaceLoaded(faceID int16, faceImage *data.FaceImage)
}

type ManagersAccessor interface {
	SetManagers(managers *Managers)
}

type WindowAccessor interface {
	SetWindow(window fyne.Window)
}

type ConnectionAccessor interface {
	SetConnection(conn *net.Connection)
}

type HandlerAccessor interface {
	SetHandler(handler *messages.MessageHandler)
}
