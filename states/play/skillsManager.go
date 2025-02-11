package play

import (
	"fmt"
	"image/color"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/mobifire/states/play/layouts"
	"github.com/kettek/termfire/messages"
)

// SkillsManager provides storage and handling of player skills.
type SkillsManager struct {
	window      fyne.Window
	conn        *net.Connection
	handler     *messages.MessageHandler
	skills      map[uint16]Skill
	knownSkills map[uint16]messages.MessageStatSkill
	exp         []uint64
}

// Skill is a convenience merger around SkillInfo and SkillExtraInfo.
type Skill struct {
	messages.SkillInfo
	messages.SkillExtraInfo
}

func NewSkillsManager() *SkillsManager {
	return &SkillsManager{}
}

func (s *SkillsManager) Init(window fyne.Window, conn *net.Connection, handler *messages.MessageHandler) {
	s.window = window
	s.conn = conn
	s.handler = handler

	s.skills = make(map[uint16]Skill)
	s.knownSkills = make(map[uint16]messages.MessageStatSkill)

	s.handler.On(&messages.MessageStats{}, nil, func(m messages.Message, mf *messages.MessageFailure) {
		msg := m.(*messages.MessageStats)
		for _, stat := range msg.Stats {
			switch stat := stat.(type) {
			case *messages.MessageStatSkill:
				s.knownSkills[uint16(stat.Skill)] = *stat
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

func (m *SkillsManager) Skill(num uint16) Skill {
	return m.skills[num]
}

func (m *SkillsManager) ExpToNextLevel(skill uint16) uint64 {
	if m.knownSkills[skill].Level >= int8(len(m.exp)) {
		return 0
	}
	exp := m.exp[m.knownSkills[skill].Level]
	return exp - uint64(m.knownSkills[skill].Exp)
}

func (m *SkillsManager) ExpToNextLevelPercentage(skill uint16) float64 {
	if m.knownSkills[skill].Level >= int8(len(m.exp)) {
		return 0
	}
	exp := m.exp[m.knownSkills[skill].Level]
	return float64(m.knownSkills[skill].Exp) / float64(exp)
}

func (m *SkillsManager) KnownSkillsSlice() []uint8 {
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

func (m *SkillsManager) ShowSimpleSkillsList(cb func(id int)) {
	var popup *widget.PopUp
	skillIDs := m.KnownSkillsSlice()
	list := widget.NewList(
		func() int {
			return len(skillIDs)
		},
		func() fyne.CanvasObject {
			return container.New(&layouts.SkillEntry{IconSize: data.CurrentFaceSet().Width}, &canvas.Image{}, widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			skill := m.skills[uint16(skillIDs[i])]
			if face, ok := data.GetFace(int(skill.Face)); ok {
				o.(*fyne.Container).Objects[0].(*canvas.Image).Resource = &face
			} else {
				o.(*fyne.Container).Objects[0].(*canvas.Image).Resource = data.GetResource("blank.png")
			}
			o.(*fyne.Container).Objects[0].(*canvas.Image).Refresh()
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(skill.Name)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		cb(int(skillIDs[id]))
		popup.Hide()
	}

	dialog := layouts.NewDialog(m.window)
	dialog.Full = true

	popup = widget.NewPopUp(container.New(dialog, list), m.window.Canvas())

	ps := popup.MinSize()
	ws := m.window.Canvas().Size()
	x := (ws.Width - ps.Width) / 2
	y := (ws.Height - ps.Height) / 2
	popup.ShowAtPosition(fyne.NewPos(x, y))
}

func (m *SkillsManager) ShowSkillsList() {
	var popup *widget.PopUp
	skillIDs := m.KnownSkillsSlice()

	info := widget.NewRichTextWithText("...")
	info.Wrapping = fyne.TextWrapWord

	list := widget.NewList(
		func() int {
			return len(skillIDs)
		},
		func() fyne.CanvasObject {
			rect := canvas.NewRectangle(color.NRGBA{255, 255, 255, 100})
			return container.New(&layouts.FullSkillEntry{IconSize: data.CurrentFaceSet().Width, Rect: rect}, &canvas.Image{}, widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""), rect)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// Set the percentage until next level (used for BG)
			o.(*fyne.Container).Layout.(*layouts.FullSkillEntry).Perc = float32(m.ExpToNextLevelPercentage(uint16(skillIDs[i])))
			skill := m.Skill(uint16(skillIDs[i]))
			if face, ok := data.GetFace(int(skill.Face)); ok {
				o.(*fyne.Container).Objects[0].(*canvas.Image).Resource = &face
			} else {
				o.(*fyne.Container).Objects[0].(*canvas.Image).Resource = data.GetResource("blank.png")
			}
			o.(*fyne.Container).Objects[0].(*canvas.Image).Refresh()
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(skill.Name)
			o.(*fyne.Container).Objects[2].(*widget.Label).SetText(fmt.Sprintf("%d", m.knownSkills[uint16(skillIDs[i])].Level))
			v, p := humanize.ComputeSI(float64(m.knownSkills[uint16(skillIDs[i])].Exp))
			f := humanize.SIWithDigits(v, 4, p)
			o.(*fyne.Container).Objects[3].(*widget.Label).SetText(f)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		skill := m.Skill(uint16(skillIDs[id]))
		next := m.ExpToNextLevel(uint16(skillIDs[id]))
		v, p := humanize.ComputeSI(float64(next))
		f := humanize.SIWithDigits(v, 4, p)
		info.Segments = data.TextToRichTextSegments(f + " until level " + fmt.Sprintf("%d", m.knownSkills[uint16(skillIDs[id])].Level+1) + "\n" + skill.Description)
		info.Refresh()
	}

	infoScroll := container.NewVScroll(info)

	cnt := container.New(&layouts.Inventory{}, list, infoScroll)

	dialog := layouts.NewDialog(m.window)
	dialog.Full = true

	popup = widget.NewPopUp(container.New(dialog, cnt), m.window.Canvas())

	ps := popup.MinSize()
	ws := m.window.Canvas().Size()
	x := (ws.Width - ps.Width) / 2
	y := (ws.Height - ps.Height) / 2
	popup.ShowAtPosition(fyne.NewPos(x, y))

}
