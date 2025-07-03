package chars

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"

	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

// State provides the character selection and creation screen.
type State struct {
	messages.MessageHandler
	conn       *net.Connection
	characters []messages.Character
	faces      []messages.MessageFace2
}

// NewState provides a new State from a connection, Character, and Face messages.
func NewState(conn *net.Connection, characters []messages.Character, faces []messages.MessageFace2) *State {
	return &State{
		conn:       conn,
		characters: characters,
		faces:      faces,
	}
}

// Enter sets up the necessary state.
func (s *State) Enter(next func(states.State)) (leave func()) {
	s.conn.SetMessageHandler(s.OnMessage)

	// Request faces sent during login.
	for _, face := range s.faces {
		s.conn.Send(&messages.MessageAskFace{Face: int32(face.Num)})
	}

	s.On(&messages.MessageImage2{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageImage2)
		data.AddFaceImage(*msg)
	})

	// We also handle/override characters here, as when the player leaves the game, it resends the characters again.
	s.On(&messages.MessageAccountPlayers{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageAccountPlayers)
		s.refreshCharacters(msg.Characters, next)
	})

	// Selection
	// TODO: Make character list element.
	s.refreshCharacters(s.characters, next)

	// Creation
	//creationContainer := s.setupCreation()

	// Tabs
	// TODO: Make tabs to switch between character creation and selection.
	if len(s.characters) > 1 {
		//tabs.SelectIndex(1)
	}

	//s.container = container.New(layout.NewVBoxLayout(), characterList)

	return nil
}

func (s *State) refreshCharacters(characters []messages.Character, next func(states.State)) {
	// TODO: Clear character list.

	for _, character := range characters {
		if character.Name == "" {
			// Skip the weird bogus empty char.
			continue
		}
		// TODO: Add button to:
		//next(play.NewState(s.conn, character.Name))
		// TODO: Show selectable cards with  character.Name, fmt.Sprintf("%d %s %s", character.Level, character.Race, character.Class)
		// TODO: Add the c ard to the character list.
	}

}

func (s *State) setupCreation() {
	// TODO: Setup the UI for character creation.
	// Race
	var races []messages.MessageReplyInfoDataRaceInfo

	// TOOD: This will be a combo that updates a race description from races[index].Description

	// TODO: Same as above for class.
	// Class
	var classes []messages.MessageReplyInfoDataClassInfo

	// TODO: Name and stat selection elements.
	// Name + Stats
	// TODO: Some sort of list of stats... we need to also get this from request_info.

	// TODO: The above we could put in tabs at the toppe.

	// Handle stuff

	s.On(&messages.MessageReplyInfo{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageReplyInfo)
		switch d := msg.Data.(type) {
		case messages.MessageReplyInfoDataRaceList:
			races = nil
			// TODO: Refresh races combo.
			// It feels excessive, but we do want to have the race info, so I guess we just spam for each. (We'll replace the given options with their matching one)
			for _, r := range d {
				s.conn.Send(&messages.MessageRequestInfo{Data: messages.MessageRequestInfoRaceInfo(r)})
				// We also queue up races to be filled here -- might as well also re-use the message structure.
				races = append(races, messages.MessageReplyInfoDataRaceInfo{
					Arch: r,
				})
			}
		case messages.MessageReplyInfoDataRaceInfo:
			for i, r := range races {
				if r.Arch == d.Arch {
					//caser := cases.Title(language.English)
					// TODO: Make the race selection nicer with caser.String(d.Name)
					// Eh... let's capitalize each starting letter in Name.
					races[i] = d // Store the full race as well.
					break
				}
			}
			// TODO: Refresh combo.
		case messages.MessageReplyInfoDataClassList:
			classes = nil
			// TODO: Refresh class combo.
			// Do the same as for races.
			for _, c := range d {
				s.conn.Send(&messages.MessageRequestInfo{Data: messages.MessageRequestInfoClassInfo(c)})
				classes = append(classes, messages.MessageReplyInfoDataClassInfo{
					Arch: c,
				})
			}
		case messages.MessageReplyInfoDataClassInfo:
			for i, c := range classes {
				if c.Arch == d.Arch {
					//caser := cases.Title(language.English)
					// Use caser as well.
					classes[i] = d
					break
				}
			}
			// TODO: Refresh class combo.
		}
	})

	// Send our requesties.
	s.conn.Send(&messages.MessageRequestInfo{Data: messages.MessageRequestInfoRaceList{}})
	s.conn.Send(&messages.MessageRequestInfo{Data: messages.MessageRequestInfoClassList{}})

	// TODO: Return somethin'
	return
}

func (s *State) Update() error {
	return nil
}

func (s *State) Draw(screen *ebiten.Image) {
	// TODO
}
