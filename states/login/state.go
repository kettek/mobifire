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
	Username         string `yaml:"username,omitempty"`
	Password         string `yaml:"password,omitempty"`
	RememberPassword bool   `yaml:"rememberPassword,omitempty"`
	ImageSet         string `yaml:"imageset,omitempty"`
}

// State provides username + account login management. If successful, sends to chars, otherwise will remain in the login state.
type State struct {
	messages.MessageHandler
	Hostname     string
	conn         *net.Connection
	faces        []messages.MessageFace2
	layout       rebui.Layout
	userNode     *rebui.Node
	passNode     *rebui.Node
	loginNode    *rebui.Node
	rulesNode    *rebui.Node
	imagesetNode *rebui.Node
	settings     serverSettings
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

	carouselLeft := s.layout.GetByID("carousel_left").Widget.(*widgets.Button)
	carouselRight := s.layout.GetByID("carousel_right").Widget.(*widgets.Button)
	carouselContent := s.layout.GetByID("carousel_content").Widget.(*widgets.Text)
	carouselContent.AssignText(s.settings.ImageSet)

	// TODO: Create remember me and/or remember password checkboxes

	s.loginNode = s.layout.GetByID("login")

	s.rulesNode = s.layout.GetByID("rules")

	var currentImageSet int
	var imageSets []messages.MessageReplyInfoDataImageInfoSet

	carouselItems := []string{}
	refreshCarousel := func(items []string) {
		carouselItems = items
		if currentImageSet >= len(items) {
			currentImageSet = len(items) - 1
		}
		if currentImageSet < 0 {
			currentImageSet = 0
		}
		// Refresh arrows
		if currentImageSet > 0 {
			carouselLeft.AssignDisabled(false)
			carouselLeft.AssignText("<")
		} else {
			carouselLeft.AssignDisabled(true)
			carouselLeft.AssignText("")
		}
		if currentImageSet == len(items)-1 {
			carouselRight.AssignDisabled(true)
			carouselRight.AssignText("")
		} else {
			carouselRight.AssignDisabled(false)
			carouselRight.AssignText(">")
		}
		if len(items) == 0 {
			carouselContent.AssignText("")
			return
		}
		carouselContent.AssignText(items[currentImageSet])
	}
	s.layout.GetByID("carousel_left").OnPointerPressed = func(epp rebui.EventPointerPressed) {
		s.conn.Send(&messages.MessageSetup{
			FaceSet: struct {
				Use   bool
				Value uint8
			}{
				Use:   true,
				Value: uint8(currentImageSet - 1),
			},
		})
	}
	s.layout.GetByID("carousel_right").OnPointerPressed = func(epp rebui.EventPointerPressed) {
		s.conn.Send(&messages.MessageSetup{
			FaceSet: struct {
				Use   bool
				Value uint8
			}{
				Use:   true,
				Value: uint8(currentImageSet + 1),
			},
		})
	}

	s.On(&messages.MessageSetup{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		if failure != nil {
			// TODO: Dialog error?
			return
		}
		msg := m.(*messages.MessageSetup)
		if msg.FaceSet.Use {
			currentImageSet = int(msg.FaceSet.Value)
			refreshCarousel(carouselItems)
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
		// Might as well remember our image set.
		s.settings.ImageSet = imageSets[currentImageSet].Name

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
			// Populate our carousel with the image set names.
			{
				items := []string{}
				for _, s := range imageSets {
					items = append(items, s.Name)
				}
				refreshCarousel(items)
			}
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
	// s.conn.Send(&messages.MessageAccountLogin{Account: usernameEntry.Text, Password: passwordEntry.Text})

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
