package login

import (
	"errors"
	"fmt"
	"slices"

	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/chars"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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

	var currentImageSet int
	var imageSets []messages.MessageReplyInfoDataImageInfoSet
	var imageSetCombo *widget.Select
	imageSetCombo = widget.NewSelect([]string{}, func(_ string) {
		index := imageSetCombo.SelectedIndex()
		if index == currentImageSet {
			return
		}
		if index < 0 || index >= len(imageSets) {
			return
		}
		s.conn.Send(&messages.MessageSetup{
			FaceSet: struct {
				Use   bool
				Value uint8
			}{
				Use:   true,
				Value: uint8(imageSets[index].Index),
			},
		})
	})

	s.On(&messages.MessageSetup{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		if failure != nil {
			dialog.ShowError(errors.New(failure.Reason), s.window)
			return
		}
		msg := m.(*messages.MessageSetup)
		if msg.FaceSet.Use {
			currentImageSet = int(msg.FaceSet.Value)
			imageSetCombo.SetSelectedIndex(currentImageSet)
		}
	})

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

		// Create our image set to use.
		imageSet := imageSets[currentImageSet]
		data.AddFaceSet(imageSet.Index, imageSet.Width, imageSet.Height)
		data.SetCurrentFaceSet(imageSet.Index)

		m := msg.(*messages.MessageAccountPlayers)
		next(chars.NewState(s.conn, m.Characters, s.faces))
	})

	s.On(&messages.MessageReplyInfo{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageReplyInfo)
		switch data := m.Data.(type) {
		case messages.MessageReplyInfoDataImageInfo:
			slices.SortStableFunc(data.Sets, func(a, b messages.MessageReplyInfoDataImageInfoSet) int {
				return a.Index - b.Index
			})
			imageSets = data.Sets
			imageSetCombo.Options = []string{}
			for _, set := range imageSets {
				imageSetCombo.Options = append(imageSetCombo.Options, set.Name)
			}
			imageSetCombo.SetSelected(imageSets[0].Name)
		}
	})
	// Request the server's image info -- this is used for properly setting face images.
	s.conn.Send(&messages.MessageRequestInfo{Data: messages.MessageRequestInfoDataImageInfo{}})

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
			{Text: "Image Set", Widget: imageSetCombo},
		},
		OnSubmit: func() {
			s.conn.Send(&messages.MessageAccountLogin{Account: usernameEntry.Text, Password: passwordEntry.Text})
		},
	}

	s.container = container.NewBorder(nil, nil, nil, nil, form)

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
