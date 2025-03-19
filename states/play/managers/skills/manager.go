package skills

import (
	"fmt"
	"image/color"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play/cfwidgets"
	"github.com/kettek/mobifire/states/play/layouts"
	"github.com/kettek/termfire/messages"
)

// Manager provides storage and handling of player skills.
type Manager struct {
	window           fyne.Window
	conn             *net.Connection
	handler          *messages.MessageHandler
	skills           map[uint16]Skill
	knownSkills      map[uint16]messages.MessageStatSkill
	knownSkillsSlice []uint8
	exp              []uint64
	sortMode         SortMode
	sortAsc          bool
}

// SortMode defines how the skills list should be sorted.
type SortMode int

// Sort modes for skills.
const (
	SortByLevel SortMode = iota
	SortByExp
	SortByName
)

// Skill is a convenience merger around SkillInfo and SkillExtraInfo.
type Skill struct {
	messages.SkillInfo
	messages.SkillExtraInfo
}

// NewManager creates a new skill manager.
func NewManager() *Manager {
	return &Manager{}
}

// SetWindow sets the window for the manager.
func (s *Manager) SetWindow(window fyne.Window) {
	s.window = window
}

// SetConnection sets the connection for the manager.
func (s *Manager) SetConnection(conn *net.Connection) {
	s.conn = conn
}

// SetHandler sets the message handler for the manager.
func (s *Manager) SetHandler(handler *messages.MessageHandler) {
	s.handler = handler
}

// Init sets up message handling for skills as well as sending the initial requests for skills and exp.
func (s *Manager) Init() {
	s.skills = make(map[uint16]Skill)
	s.knownSkills = make(map[uint16]messages.MessageStatSkill)

	s.handler.On(&messages.MessageStats{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageStats)
		for _, stat := range msg.Stats {
			switch stat := stat.(type) {
			case *messages.MessageStatSkill:
				s.knownSkills[uint16(stat.Skill)] = *stat
				s.syncKnownSkills()
			}
		}
	})

	s.handler.On(&messages.MessageReplyInfo{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageReplyInfo)
		switch data := msg.Data.(type) {
		case messages.MessageReplyInfoDataSkillInfo:
			for num, skill := range data.Skills {
				if sk, ok := s.skills[uint16(num)]; ok {
					sk.SkillInfo = skill
					s.skills[uint16(num)] = sk
				} else {
					s.skills[uint16(num)] = Skill{
						SkillInfo: skill,
					}
				}
			}
			s.syncKnownSkills()
		case messages.MessageReplyInfoDataSkillExtra:
			for num, skill := range data.Skills {
				if sk, ok := s.skills[uint16(num)]; ok {
					sk.SkillExtraInfo = skill
					s.skills[uint16(num)] = sk
				} else {
					s.skills[uint16(num)] = Skill{
						SkillExtraInfo: skill,
					}
				}
			}
			s.syncKnownSkills()
		case messages.MessageReplyInfoDataExpTable:
			s.exp = data
		}
	})

	// Request our skills and our extra info.
	s.conn.Send(&messages.MessageRequestInfo{
		Data: messages.MessageRequestInfoSkillInfo(true),
	})
	s.conn.Send(&messages.MessageRequestInfo{
		Data: messages.MessageRequestInfoSkillExtra(1),
	})
	// Also request exp table.
	s.conn.Send(&messages.MessageRequestInfo{
		Data: messages.MessageRequestInfoExpTable{},
	})
}

// Skill returns a skill by its number.
func (m *Manager) Skill(num uint16) Skill {
	return m.skills[num]
}

// ExpToNextLevel returns the amount of exp needed to reach the next level for a given skill.
func (m *Manager) ExpToNextLevel(skill uint16) uint64 {
	if m.knownSkills[skill].Level >= int8(len(m.exp)) {
		return 0
	}
	exp := m.exp[m.knownSkills[skill].Level]
	return exp - uint64(m.knownSkills[skill].Exp)
}

// ExpToNextLevelPercentage returns the exp to next level as a percentage for a given skill.
func (m *Manager) ExpToNextLevelPercentage(skill uint16) float64 {
	if m.knownSkills[skill].Level >= int8(len(m.exp)) {
		return 0
	}
	exp := m.exp[m.knownSkills[skill].Level]
	return float64(m.knownSkills[skill].Exp) / float64(exp)
}

// KnownSkillsSlice returns a slice of known skills sorted by the current sort mode.
func (m *Manager) KnownSkillsSlice() []uint8 {
	var skills []uint8
	for _, skill := range m.knownSkills {
		skills = append(skills, skill.Skill)
	}
	// I guess just sort by name.
	/*slices.SortFunc(skills, func(a, b uint8) int {
		return strings.Compare(m.skills[uint16(a)].Name, m.skills[uint16(b)].Name)
	})*/
	// Actually, sort by exp.
	slices.SortFunc(skills, func(a, b uint8) int {
		return int(m.knownSkills[uint16(b)].Exp - m.knownSkills[uint16(a)].Exp)
	})
	return skills
}

func (m *Manager) syncKnownSkills() {
	m.knownSkillsSlice = m.KnownSkillsSlice()
	m.sortKnownSkills()
}

func (m *Manager) sortKnownSkills() {
	switch m.sortMode {
	case SortByExp:
		if m.sortAsc {
			slices.SortFunc(m.knownSkillsSlice, func(a, b uint8) int {
				return int((m.ExpToNextLevelPercentage(uint16(a)) - m.ExpToNextLevelPercentage(uint16(b))) * 100000)
			})
		} else {
			slices.SortFunc(m.knownSkillsSlice, func(a, b uint8) int {
				return int((m.ExpToNextLevelPercentage(uint16(b)) - m.ExpToNextLevelPercentage(uint16(a))) * 100000)
			})
		}
	case SortByLevel:
		if m.sortAsc {
			slices.SortFunc(m.knownSkillsSlice, func(a, b uint8) int {
				return int(m.knownSkills[uint16(a)].Level - m.knownSkills[uint16(b)].Level)
			})
		} else {
			slices.SortFunc(m.knownSkillsSlice, func(a, b uint8) int {
				return int(m.knownSkills[uint16(b)].Level - m.knownSkills[uint16(a)].Level)
			})
		}
	case SortByName:
		if m.sortAsc {
			slices.SortFunc(m.knownSkillsSlice, func(a, b uint8) int {
				return strings.Compare(m.skills[uint16(a)].Name, m.skills[uint16(b)].Name)
			})
		} else {
			slices.SortFunc(m.knownSkillsSlice, func(a, b uint8) int {
				return strings.Compare(m.skills[uint16(b)].Name, m.skills[uint16(a)].Name)
			})
		}
	}
}

// ShowSimpleSkillsList shows a simple list of skills with their icons and names.
func (m *Manager) ShowSimpleSkillsList(cb func(id int)) {
	var popup *cfwidgets.PopUp
	list := widget.NewList(
		func() int {
			return len(m.knownSkillsSlice)
		},
		func() fyne.CanvasObject {
			return container.New(&layouts.SkillEntry{IconSize: data.CurrentFaceSet().Width}, &canvas.Image{}, widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			skill := m.skills[uint16(m.knownSkillsSlice[i])]
			if face, ok := data.GetFace(int(skill.Face)); ok {
				o.(*fyne.Container).Objects[0].(*canvas.Image).Resource = face
			} else {
				o.(*fyne.Container).Objects[0].(*canvas.Image).Resource = data.GetResource("blank.png")
			}
			o.(*fyne.Container).Objects[0].(*canvas.Image).Refresh()
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(skill.Name)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		cb(int(m.knownSkillsSlice[id]))
		popup.Hide()
	}

	dialog := layouts.NewDialog(m.window)
	dialog.Full = true

	popup = cfwidgets.NewPopUp(container.New(dialog, list), m.window.Canvas())

	popup.ShowCentered(m.window.Canvas())
}

// ShowSkillsList shows a detailed list of skills.
func (m *Manager) ShowSkillsList() {
	var popup *cfwidgets.PopUp

	info := widget.NewRichTextWithText("...")
	info.Wrapping = fyne.TextWrapWord

	list := widget.NewList(
		func() int {
			return len(m.knownSkillsSlice)
		},
		func() fyne.CanvasObject {
			rect := canvas.NewRectangle(color.NRGBA{255, 255, 255, 100})
			return container.New(&layouts.FullSkillEntry{IconSize: data.CurrentFaceSet().Width, Rect: rect}, &canvas.Image{}, widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), rect)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// Set the percentage until next level (used for BG)
			o.(*fyne.Container).Layout.(*layouts.FullSkillEntry).Perc = float32(m.ExpToNextLevelPercentage(uint16(m.knownSkillsSlice[i])))
			skill := m.Skill(uint16(m.knownSkillsSlice[i]))
			if face, ok := data.GetFace(int(skill.Face)); ok {
				o.(*fyne.Container).Objects[0].(*canvas.Image).Resource = face
			} else {
				o.(*fyne.Container).Objects[0].(*canvas.Image).Resource = data.GetResource("blank.png")
			}
			o.(*fyne.Container).Objects[0].(*canvas.Image).Refresh()
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(skill.Name)
			o.(*fyne.Container).Objects[2].(*widget.Label).SetText(fmt.Sprintf("%d", m.knownSkills[uint16(m.knownSkillsSlice[i])].Level))
			v, p := humanize.ComputeSI(float64(m.knownSkills[uint16(m.knownSkillsSlice[i])].Exp))
			f := humanize.SIWithDigits(v, 4, p)
			o.(*fyne.Container).Objects[3].(*widget.Label).SetText(f)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		skill := m.Skill(uint16(m.knownSkillsSlice[id]))
		next := m.ExpToNextLevel(uint16(m.knownSkillsSlice[id]))
		v, p := humanize.ComputeSI(float64(next))
		f := humanize.SIWithDigits(v, 4, p)
		info.Segments = data.TextToRichTextSegments(f + " until level " + fmt.Sprintf("%d", m.knownSkills[uint16(m.knownSkillsSlice[id])].Level+1) + "\n" + skill.Description)
		info.Refresh()
	}

	infoScroll := container.NewVScroll(info)

	cnt := container.New(&layouts.Inventory{}, list, infoScroll)

	var actionSortDir, actionSortExp, actionSortLevel, actionSortName *widget.ToolbarAction
	actionSortDir = widget.NewToolbarAction(data.GetResource("icon_descending.png"), func() {
		m.sortAsc = !m.sortAsc
		if m.sortAsc {
			actionSortDir.SetIcon(data.GetResource("icon_ascending.png"))
		} else {
			actionSortDir.SetIcon(data.GetResource("icon_descending.png"))
		}
		m.sortKnownSkills()
		list.Refresh()
	})
	actionSortExp = widget.NewToolbarAction(data.GetResource("icon_exp.png"), func() {
		m.sortMode = SortByExp
		m.sortKnownSkills()
		list.Refresh()
	})
	actionSortLevel = widget.NewToolbarAction(data.GetResource("icon_level.png"), func() {
		m.sortMode = SortByLevel
		m.sortKnownSkills()
		list.Refresh()
	})
	actionSortName = widget.NewToolbarAction(data.GetResource("icon_name.png"), func() {
		m.sortMode = SortByName
		m.sortKnownSkills()
		list.Refresh()
	})
	topControls := widget.NewToolbar(
		actionSortLevel,
		actionSortExp,
		actionSortName,
		widget.NewToolbarSeparator(),
		actionSortDir,
	)
	topBar := container.NewHBox(widget.NewLabel("Skills"), topControls)

	blep := container.NewBorder(topBar, nil, nil, nil, cnt)

	dialog := layouts.NewDialog(m.window)
	dialog.Full = true

	popup = cfwidgets.NewPopUp(container.New(dialog, blep), m.window.Canvas())

	popup.ShowCentered(m.window.Canvas())
}
