package login

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/data/settings"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/chars"
	"github.com/kettek/rebui"
	"github.com/kettek/rebui/widgets"

	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

// TODO: Move this to some common pkg.
type serverSettings struct {
	Username         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	RememberPassword bool   `json:"rememberPassword,omitempty"`
}

// State provides username + account login management. If successful, sends to chars, otherwise will remain in the login state.
type State struct {
	messages.MessageHandler
	Hostname  string
	conn      *net.Connection
	faces     []messages.MessageFace2
	layout    rebui.Layout
	userNode  *rebui.Node
	passNode  *rebui.Node
	loginNode *rebui.Node
	rulesNode *rebui.Node
	settings  serverSettings
}

// NewState returns a State from the given connection.
func NewState(conn *net.Connection, hostname string) *State {
	return &State{
		conn:     conn,
		Hostname: hostname,
	}
}

// Enter sets up all the necessary logic for logging in.
func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)

	s.settings = settings.GetWithFallback(s.Hostname, serverSettings{})

	layout, err := data.GetLayout("login")
	if err != nil {
		panic(err)
	}
	for _, node := range layout.Nodes {
		s.layout.AddNode(node)
	}

	s.userNode = s.layout.GetByID("username")
	s.userNode.Widget.(*widgets.TextInput).AssignText(s.settings.Username)

	s.passNode = s.layout.GetByID("password")
	s.passNode.Widget.(*widgets.TextInput).AssignText(s.settings.Password)

	// TODO: Create remember me and/or remember password checkboxes

	s.loginNode = s.layout.GetByID("login")

	s.rulesNode = s.layout.GetByID("rules")

	var currentImageSet int
	var imageSets []messages.MessageReplyInfoDataImageInfoSet
	// TODO: Create image set combo box.
	/*imageSetCombo = widget.NewSelect([]string{}, func(_ string) {
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
	})*/

	s.On(&messages.MessageSetup{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		if failure != nil {
			// TODO: Dialog error?
			return
		}
		msg := m.(*messages.MessageSetup)
		if msg.FaceSet.Use {
			currentImageSet = int(msg.FaceSet.Value)
			// TODO: Set combo box index.
		}
	})

	s.On(&messages.MessageAccountLogin{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		if mf != nil {
			fmt.Println("Failed to login: ", mf.Reason)
			// TODO: Show error dialog.
			return
		}
	})

	s.On(&messages.MessageAccountPlayers{}, &messages.MessageAccountLogin{}, func(msg messages.Message, failure *messages.MessageFailure) {
		if failure != nil {
			// TODO: Show error dialog.
			return
		}

		// TODO: Save username/pass if requested, otherwise delete them.

		// Create our image set to use.
		imageSet := imageSets[currentImageSet]
		data.AddFaceSet(imageSet.Index, imageSet.Width, imageSet.Height)
		data.SetCurrentFaceSet(imageSet.Index)

		m := msg.(*messages.MessageAccountPlayers)
		next(chars.NewState(s.conn, m.Characters, s.faces))
	})

	s.On(&messages.MessageReplyInfo{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageReplyInfo)
		switch d := m.Data.(type) {
		case messages.MessageReplyInfoDataImageInfo:
			slices.SortStableFunc(d.Sets, func(a, b messages.MessageReplyInfoDataImageInfoSet) int {
				return a.Index - b.Index
			})
			imageSets = d.Sets
			// Clear image set combo box.
			/*for _, set := range imageSets {
				imageSetCombo.Options = append(imageSetCombo.Options, set.Name)
			}
			imageSetCombo.SetSelected(imageSets[0].Name)*/
		case messages.MessageReplyInfoDataRules:
			s.rulesNode.Widget.(*widgets.Text).AssignText(string(d))
			// Update our rules element with the rules text.
		}
	})
	// Request the server's image info -- this is used for properly setting face images.
	s.conn.Send(&messages.MessageRequestInfo{Data: messages.MessageRequestInfoDataImageInfo{}})
	// Request rules... seems reasonable enough.
	s.conn.Send(&messages.MessageRequestInfo{Data: messages.MessageRequestInfoRules{}})

	s.On(&messages.MessageFace2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m, ok := msg.(*messages.MessageFace2)
		if !ok {
			return
		}
		s.faces = append(s.faces, *m)
	})

	// TODO: Create username, password, and remember me layout form(?)
	// On submit, issue:
	//s.conn.Send(&messages.MessageAccountLogin{Account: usernameEntry.Text, Password: passwordEntry.Text})

	return nil
}

func (s *State) Update() error {
	s.layout.Update()
	return nil
}

func (s *State) Draw(screen *ebiten.Image) {
	s.layout.Draw(screen)
	// TODO
}
