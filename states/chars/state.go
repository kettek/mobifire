package chars

import (
	"fmt"

	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/termfire/messages"
)

// State provides the character selection and creation screen.
type State struct {
	window fyne.Window
	messages.MessageHandler
	container     *fyne.Container
	conn          *net.Connection
	characterList *fyne.Container
	characters    []messages.Character
	faces         []messages.MessageFace2
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
		// It's kind of dumb, but reset our view to 11,11
		s.conn.Send(&messages.MessageSetup{
			MapSize: struct {
				Use   bool
				Value string
			}{
				Use:   true,
				Value: "11x11",
			},
		})
	})

	// Selection
	s.characterList = container.New(layout.NewVBoxLayout())
	s.refreshCharacters(s.characters, next)

	// Creation
	creationContainer := s.setupCreation()

	// Tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Create", creationContainer),
		container.NewTabItem("Select", container.NewVScroll(s.characterList)),
	)
	if len(s.characters) > 1 {
		tabs.SelectIndex(1)
	}

	s.container = container.NewBorder(nil, nil, nil, nil, tabs)

	//s.container = container.New(layout.NewVBoxLayout(), characterList)

	return nil
}

func (s *State) refreshCharacters(characters []messages.Character, next func(states.State)) {
	s.characterList.RemoveAll()

	for _, character := range characters {
		if character.Name == "" {
			// Skip the weird bogus empty char.
			continue
		}
		content := container.New(layout.NewHBoxLayout(), widget.NewLabel(character.Map), widget.NewButton("Play", func() {
			next(play.NewState(s.conn, character.Name))
		}))
		card := widget.NewCard(character.Name, fmt.Sprintf("%d %s %s", character.Level, character.Race, character.Class), content)
		s.characterList.Add(card)
	}

}

func (s *State) setupCreation() fyne.CanvasObject {
	// Race
	var races []messages.MessageReplyInfoDataRaceInfo

	var racesCombo *widget.Select
	var raceDescription *widget.Label

	racesCombo = widget.NewSelect([]string{}, func(r string) {
		// Update the race info!
		raceDescription.SetText(races[racesCombo.SelectedIndex()].Description)
	})

	raceDescription = widget.NewLabel("")
	raceDescription.Wrapping = fyne.TextWrapWord

	raceContainer := container.NewBorder(racesCombo, nil, nil, nil, container.NewVScroll(raceDescription))

	// Class
	var classes []messages.MessageReplyInfoDataClassInfo

	var classCombo *widget.Select
	var classDescription *widget.Label

	classCombo = widget.NewSelect([]string{}, func(r string) {
		classDescription.SetText(classes[classCombo.SelectedIndex()].Description)
	})

	classDescription = widget.NewLabel("")
	classDescription.Wrapping = fyne.TextWrapWord

	classContainer := container.NewBorder(classCombo, nil, nil, nil, container.NewVScroll(classDescription))

	// Name + Stats
	var nameEntry *widget.Entry
	var statsLabel *widget.Label
	// TODO: Some sort of list of stats... we need to also get this from request_info.

	nameEntry = widget.NewEntry()
	nameEntry.PlaceHolder = "Name"

	statsLabel = widget.NewLabel("Stats go here")

	statsContainer := container.NewVScroll(container.New(layout.NewVBoxLayout(), nameEntry, statsLabel))

	// Tabs

	creationTabs := container.NewAppTabs()
	creationTabs.Append(container.NewTabItem("Race", raceContainer))
	creationTabs.Append(container.NewTabItem("Class", classContainer))
	creationTabs.Append(container.NewTabItem("Stats & Name", statsContainer))

	// Handle stuff

	s.On(&messages.MessageReplyInfo{}, nil, func(m messages.Message, failure *messages.MessageFailure) {
		msg := m.(*messages.MessageReplyInfo)
		switch d := msg.Data.(type) {
		case messages.MessageReplyInfoDataRaceList:
			races = nil
			racesCombo.Options = []string(d)
			racesCombo.Refresh()
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
					caser := cases.Title(language.English)
					// Eh... let's capitalize each starting letter in Name.
					racesCombo.Options[i] = caser.String(d.Name)
					races[i] = d // Store the full race as well.
					break
				}
			}
			racesCombo.Refresh()
		case messages.MessageReplyInfoDataClassList:
			classes = nil
			classCombo.Options = []string(d)
			classCombo.Refresh()
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
					caser := cases.Title(language.English)
					classCombo.Options[i] = caser.String(d.Name)
					classes[i] = d
					break
				}
			}
			classCombo.Refresh()
		}
	})

	// Send our requesties.
	s.conn.Send(&messages.MessageRequestInfo{Data: messages.MessageRequestInfoRaceList{}})
	s.conn.Send(&messages.MessageRequestInfo{Data: messages.MessageRequestInfoClassList{}})

	return creationTabs
	//return container.NewVScroll(container.New(layout.NewVBoxLayout(), racesCombo, raceDescription))
}

// SetWindow sets the window for dialog functions.
func (s *State) SetWindow(window fyne.Window) {
	s.window = window
}

// Container returns the container.
func (s *State) Container() *fyne.Container {
	return s.container
}
