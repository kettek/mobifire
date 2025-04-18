package action

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/states/play/cfwidgets"
	"github.com/kettek/termfire/messages"
)

// EntryApplyKind represents an action to equip, unequip, or fire an object.
type EntryApplyKind struct {
	ObjectName      string // The name of the object. This is used to do a lookup.
	objectTag       int32  // The ID of the object. This is cached and used when possible. If the ObjectTag does not exist, the Name will be used to do a lookup.
	OnlyIfUnapplied bool
	Fire            bool // Whether to fire the object or not
	exists          bool // This is set to true once the action has been triggered and the object tag is known to exist.
}

// EntrySkillKind represents an action to use or ready a skill.
type EntrySkillKind struct {
	Skill int32
	Name  string
	Ready bool // Whether to ready or use the skill
}

// EntrySpellKind represents an action to cast or invoke a spell.
type EntrySpellKind struct {
	Spell  int32
	Name   string
	Extra  string // Extra string to pass into the spell -- used for create food, etc.
	Ready  bool   // Whether to ready or cast the spell
	exists bool   // This is set to true once the action has been triggered and the spell tag is known to exist.
}

// EntryCommandKind represents an action to execute a command.
type EntryCommandKind struct {
	Command string
	Repeat  int
}

// EntryStepForwardKind represents an action to step the player forward.
type EntryStepForwardKind struct {
}

// EntryFaceKind represents an action to face a direction.
type EntryFaceKind struct {
	Dir string
}

// Entry represents an action entry in the action manager.
type Entry struct {
	Image  fyne.Resource
	widget *cfwidgets.AssignableButton
	Kind   interface{}
	Next   *Entry `json:",omitempty"`
}

// TypeString returns a string representation of the entry type.
func (e Entry) TypeString() string {
	str := ""
	switch k := e.Kind.(type) {
	case EntryApplyKind:
		str = "apply"
		str += " " + k.ObjectName
		if k.OnlyIfUnapplied {
			str += " (if unapplied)"
		}
		if k.Fire {
			str += " (fire)"
		}
	case EntrySkillKind:
		str = "skill"
		str += " " + k.Name
		if k.Ready {
			str += " (ready)"
		}
	case EntrySpellKind:
		str = "spell"
		str += " " + k.Name
		if k.Ready {
			str += " (ready)"
		}
	case EntryCommandKind:
		str = "command"
		str += " " + k.Command
	case EntryStepForwardKind:
		str = "step forward"
	case EntryFaceKind:
		if k.Dir == "" {
			str = "face"
		} else {
			str = "face " + k.Dir
		}
	}
	return str
}

// NewEntryFromString creates a new Entry from a JSON string.
func NewEntryFromString(data string) *Entry {
	var entry Entry
	json.Unmarshal([]byte(data), &entry)
	return &entry
}

type kindWrapper struct {
	Kind  string
	Value json.RawMessage
}

type entryWrapper struct {
	Image string // base64 data of image
	Data  kindWrapper
	Next  json.RawMessage
}

// MarshalJSON marshals the entry to JSON.
func (e *Entry) MarshalJSON() ([]byte, error) {
	var kind string
	var value []byte

	var err error
	switch e.Kind.(type) {
	case EntryApplyKind:
		kind = "apply"
		value, err = json.Marshal(e.Kind.(EntryApplyKind))
	case EntrySkillKind:
		kind = "skill"
		value, err = json.Marshal(e.Kind.(EntrySkillKind))
	case EntrySpellKind:
		kind = "spell"
		value, err = json.Marshal(e.Kind.(EntrySpellKind))
	case EntryCommandKind:
		kind = "command"
		value, err = json.Marshal(e.Kind.(EntryCommandKind))
	case EntryStepForwardKind:
		kind = "step_forward"
		value, err = json.Marshal(e.Kind.(EntryStepForwardKind))
	case EntryFaceKind:
		kind = "face"
		value, err = json.Marshal(e.Kind.(EntryFaceKind))
	}
	if err != nil {
		return nil, err
	}

	if e.Image == nil {
		e.Image = data.GetResource("blank.png")
	}
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(e.Image.Content())))
	base64.StdEncoding.Encode(dst, e.Image.Content())

	wrapper := entryWrapper{
		Image: string(dst),
		Data: kindWrapper{
			Kind:  kind,
			Value: json.RawMessage(value),
		},
	}
	if e.Next != nil {
		next, err := json.Marshal(e.Next)
		if err != nil {
			fmt.Println("error marshalling next entry:", err)
		} else {
			wrapper.Next = next
		}
	}

	return json.Marshal(wrapper)
}

// Unmarshal unmarshals the entry from JSON.
func (e *Entry) Unmarshal(index int, b []byte) error {
	var wrapper entryWrapper
	if err := json.Unmarshal(b, &wrapper); err != nil {
		// TODO: Adjust entry to show it is erroneous.
		return nil
	}

	if wrapper.Data.Kind == "apply" {
		var kind EntryApplyKind
		if err := json.Unmarshal(wrapper.Data.Value, &kind); err != nil {
			return err
		}
		e.Kind = kind
	} else if wrapper.Data.Kind == "skill" {
		var kind EntrySkillKind
		if err := json.Unmarshal(wrapper.Data.Value, &kind); err != nil {
			return err
		}
		e.Kind = kind
	} else if wrapper.Data.Kind == "spell" {
		var kind EntrySpellKind
		if err := json.Unmarshal(wrapper.Data.Value, &kind); err != nil {
			return err
		}
		e.Kind = kind
	} else if wrapper.Data.Kind == "command" {
		var kind EntryCommandKind
		if err := json.Unmarshal(wrapper.Data.Value, &kind); err != nil {
			return err
		}
		e.Kind = kind
	} else if wrapper.Data.Kind == "step_forward" {
		var kind EntryStepForwardKind
		if err := json.Unmarshal(wrapper.Data.Value, &kind); err != nil {
			return err
		}
		e.Kind = kind
	} else if wrapper.Data.Kind == "face" {
		var kind EntryFaceKind
		if err := json.Unmarshal(wrapper.Data.Value, &kind); err != nil {
			return err
		}
		e.Kind = kind
	}

	b64, err := base64.StdEncoding.DecodeString(wrapper.Image)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		e.Image = data.GetResource("blank.png")
		return nil // Don't return error, as we shouldn't explode if something wack happened to the embedded image.
	}
	e.Image = fyne.NewStaticResource(fmt.Sprintf("entry_%d", index), b64)

	if len(wrapper.Next) > 0 && string(wrapper.Next) != "null" {
		e.Next = &Entry{}
		if err := e.Next.Unmarshal(index+100 /* fixme: this is dumb */, wrapper.Next); err != nil {
			return err
		}
	}

	return nil
}

// Trigger triggers the entry action.
func (e Entry) Trigger(m *Manager) {
	switch k := e.Kind.(type) {
	case EntryApplyKind:
		if !k.exists {
			// Lookup the given object by name... I guess search _all_ inventories for now.
			if item := m.itemsManager.GetItemByName(k.ObjectName); item != nil {
				k.objectTag = item.Tag
				k.exists = true
			}
		}
		if item := m.itemsManager.GetItemByTag(k.objectTag); item != nil {
			if k.Fire {
				if !item.Flags.Applied() {
					m.conn.Send(&messages.MessageApply{
						Tag: k.objectTag,
					})
				}
				m.conn.SendCommand(fmt.Sprintf("fire %d", m.lastDir), 1)
				m.conn.SendCommand("fire_stop", 1)
				return
			}
			if k.OnlyIfUnapplied && item.Flags.Applied() {
				return
			}
			m.conn.Send(&messages.MessageApply{
				Tag: k.objectTag,
			})
		}
	case EntrySpellKind:
		if !k.exists {
			// Lookup the given spell by name... I guess search _all_ spells for now.
			if spell := m.spellsManager.GetSpellByName(k.Name); spell != nil {
				k.Spell = spell.Tag
				k.exists = true
			}
		}
		if k.Ready {
			if k.Extra != "" {
				m.conn.SendCommand(fmt.Sprintf("cast %d %s", k.Spell, k.Extra), 1)
			} else {
				m.conn.SendCommand(fmt.Sprintf("cast %d", k.Spell), 1)
			}
		} else {
			if k.Extra != "" {
				m.conn.SendCommand(fmt.Sprintf("invoke %d %s", k.Spell, k.Extra), 1)
			} else {
				m.conn.SendCommand(fmt.Sprintf("invoke %d", k.Spell), 1)
			}
		}
	case EntrySkillKind:
		// TODO: Maybe add ready and use skill option? This would ensure that a talisman or holy symbol gets equipped before using the skill.
		if k.Ready {
			m.conn.SendCommand("ready_skill "+k.Name, 1)
		} else {
			m.conn.SendCommand("use_skill "+k.Name, 1)
		}
	case EntryCommandKind:
		m.conn.SendCommand(k.Command, uint32(k.Repeat))
	case EntryStepForwardKind:
		m.conn.SendCommand(m.GetStringFromDirection(), 1)
	case EntryFaceKind:
		if k.Dir == "" {
			// If blank, use last known direction used by the player.
			m.conn.SendCommand("face "+m.GetStringFromDirection(), 1)
		} else {
			// Otherwise whatever was set.
			m.conn.SendCommand("face "+k.Dir, 1)
		}
	}
	if e.Next != nil {
		e.Next.Trigger(m)
	}
}
