package login

import (
	"errors"
	"fmt"

	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/chars"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

// State provides username + account login management. If successful, sends to chars, otherwise will remain in the login state.
type State struct {
	messages.MessageHandler
	app       fyne.App
	window    fyne.Window
	container *fyne.Container
	conn      *net.Connection
	faces     []messages.MessageFace2
}

// NewState returns a State from the given connection.
func NewState(conn *net.Connection) *State {
	return &State{
		conn: conn,
	}
}

// Enter sets up all the necessary logic for logging in.
func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)

	// Variables used for storing username and password.
	host := s.app.Preferences().String("lastServer")
	port := s.app.Preferences().Int("lastPort")
	key := fmt.Sprintf("%s-%d", host, port)

	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(s.app.Preferences().StringWithFallback(key+"-account", ""))
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetText(s.app.Preferences().StringWithFallback(key+"-password", ""))
	rememberCheck := widget.NewCheck("", func(remember bool) {
		s.app.Preferences().SetBool(key+"-remember", remember)
	})
	rememberCheck.SetChecked(s.app.Preferences().Bool(key + "-remember"))

	s.On(&messages.MessageAccountLogin{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		if mf != nil {
			fmt.Println("Failed to login: ", mf.Reason)
			dialog.ShowError(errors.New(mf.Reason), s.window)
			return
		}
	})

	s.On(&messages.MessageAccountPlayers{}, &messages.MessageAccountLogin{}, func(msg messages.Message, failure *messages.MessageFailure) {
		if failure != nil {
			dialog.ShowError(errors.New(failure.Reason), s.window)
			return
		}

		if s.app.Preferences().Bool(key + "-remember") {
			s.app.Preferences().SetString(key+"-account", usernameEntry.Text)
			s.app.Preferences().SetString(key+"-password", passwordEntry.Text)
		} else {
			// Clear it out.
			s.app.Preferences().SetString(key+"-account", "")
			s.app.Preferences().SetString(key+"-password", "")
		}

		m := msg.(*messages.MessageAccountPlayers)
		next(chars.NewState(s.conn, m.Characters, s.faces))
	})

	s.On(&messages.MessageFace2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m, ok := msg.(*messages.MessageFace2)
		if !ok {
			return
		}
		s.faces = append(s.faces, *m)
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
			{Text: "Remember", Widget: rememberCheck},
		},
		OnSubmit: func() {
			s.conn.Send(&messages.MessageAccountLogin{Account: usernameEntry.Text, Password: passwordEntry.Text})
		},
	}

	s.container = container.New(layout.NewVBoxLayout(), form)

	return nil
}

// SetWindow sets the window -- used for showing errors.
func (s *State) SetWindow(window fyne.Window) {
	s.window = window
}

// SetApp sets the app.
func (s *State) SetApp(app fyne.App) {
	s.app = app
}

// Container returns the container.
func (s *State) Container() *fyne.Container {
	return s.container
}
