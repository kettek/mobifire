package login

import (
	"fmt"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/data/settings"
	cwidget "github.com/kettek/mobifire/data/widget"
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
	RememberUsername bool   `yaml:"rememberUsername,omitempty"`
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

	statusWidget := s.layout.GetByID("status").Widget.(*widgets.Label)
	passWidget := s.layout.GetByID("password").Widget.(*widgets.TextInput)
	userWidget := s.layout.GetByID("username").Widget.(*widgets.TextInput)
	carouselLeft := s.layout.GetByID("carousel__left").Widget.(*widgets.Button)
	carouselRight := s.layout.GetByID("carousel__right").Widget.(*widgets.Button)
	carouselContent := s.layout.GetByID("carousel__content").Widget.(*widgets.Text)
	rememberUsername := s.layout.GetByID("remember_username__checkbox").Widget.(*cwidget.Checkbox)
	rememberPassword := s.layout.GetByID("remember_password__checkbox").Widget.(*cwidget.Checkbox)
	loginNode := s.layout.GetByID("login")
	rulesNode := s.layout.GetByID("rules")

	userWidget.AssignText(s.settings.Username)
	userWidget.OnChange = func(str string) {
		s.settings.Username = str
	}

	passWidget.AssignText(s.settings.Password)
	passWidget.OnChange = func(str string) {
		s.settings.Password = str
	}

	carouselContent.AssignText(s.settings.ImageSet)

	rememberUsername.Set(s.settings.RememberUsername)
	rememberUsername.OnChange = func(b bool) {
		s.settings.RememberUsername = b
	}
	rememberPassword.Set(s.settings.RememberPassword)
	rememberPassword.OnChange = func(b bool) {
		s.settings.RememberPassword = b
	}

	loginNode.OnPointerPressed = func(epp rebui.EventPointerPressed) {
		s.conn.Send(&messages.MessageAccountLogin{Account: s.settings.Username, Password: s.settings.Password})
	}

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
	s.layout.GetByID("carousel__left").OnPointerPressed = func(epp rebui.EventPointerPressed) {
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
	s.layout.GetByID("carousel__right").OnPointerPressed = func(epp rebui.EventPointerPressed) {
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
			statusWidget.AssignText(mf.Reason)
			return
		}
	})

	s.On(&messages.MessageAccountPlayers{}, &messages.MessageAccountLogin{}, func(msg messages.Message, failure *messages.MessageFailure) {
		if failure != nil {
			// TODO: Show error dialog.
			return
		}

		// TODO: Save username/pass if requested, otherwise delete them.
		if !s.settings.RememberUsername {
			s.settings.Username = ""
		}
		if !s.settings.RememberPassword {
			s.settings.Password = ""
		}
		// Might as well remember our image set.
		s.settings.ImageSet = imageSets[currentImageSet].Name

		// Save settings.
		settings.Set(s.Hostname, s.settings)
		settings.Save()

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
				// See if we can restore our previous image set if we have one stored.
				if s.settings.ImageSet != "" {
					for i := 0; i < len(items); i++ {
						if items[i] == s.settings.ImageSet {
							currentImageSet = i
							break
						}
					}
				}
				refreshCarousel(items)
			}
		case messages.MessageReplyInfoDataRules:
			rulesNode.Widget.(*widgets.Text).AssignText(string(d))
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
