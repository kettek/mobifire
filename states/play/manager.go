package play

import (
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
		manager.Init(window, conn, handler)
	}
}

func (m *Managers) GetByType(manager Manager) Manager {
	for _, v := range *m {
		if reflect.TypeOf(v) == reflect.TypeOf(manager) {
			return v
		}
	}
	return nil
}

func (m *Managers) GetFaceLoadedManagers() []ManagerWithFaceLoaded {
	var managers []ManagerWithFaceLoaded
	for _, v := range *m {
		if manager, ok := v.(ManagerWithFaceLoaded); ok {
			managers = append(managers, manager)
		}
	}
	return managers
}

type Manager interface {
	Init(window fyne.Window, conn *net.Connection, manager *messages.MessageHandler)
}

type ManagerWithFaceLoaded interface {
	Manager
	OnFaceLoaded(faceID int16, faceImage *data.FaceImage)
}
