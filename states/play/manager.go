package play

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

type Manager interface {
	Init(window fyne.Window, conn *net.Connection, manager *messages.MessageHandler)
	OnFaceLoaded(faceID int16, faceImage *data.FaceImage)
}
