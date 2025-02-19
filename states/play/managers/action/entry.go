package action

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/termfire/messages"
)

type EntryApplyKind struct {
	ObjectName      string // The name of the object. This is used to do a lookup.
	objectTag       int32  // The ID of the object. This is cached and used when possible. If the ObjectTag does not exist, the Name will be used to do a lookup.
	OnlyIfUnapplied bool
	exists          bool // This is set to true once the action has been triggered and the object tag is known to exist.
}

type EntrySkillKind struct {
	Skill int32
	Name  string
	Ready bool // Whether to ready or use the skill
}

type EntrySpellKind struct {
	Spell  int32
	Name   string
	Ready  bool // Whether to ready or cast the spell
	exists bool // This is set to true once the action has been triggered and the spell tag is known to exist.
}

type EntryCommandKind struct {
	Command string
	Repeat  int
}

type Entry struct {
	Image fyne.Resource
	Kind  interface{}
}

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
}

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
	}
	if err != nil {
		return nil, err
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

	return json.Marshal(wrapper)
}

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
	}

	b64, err := base64.StdEncoding.DecodeString(wrapper.Image)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		e.Image = data.GetResource("blank.png")
		return nil // Don't return error, as we shouldn't explode if something wack happened to the embedded image.
	}
	e.Image = fyne.NewStaticResource(fmt.Sprintf("entry_%d", index), b64)

	return nil
}

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
			m.conn.SendCommand(fmt.Sprintf("cast %d", k.Spell), 1)
		} else {
			m.conn.SendCommand(fmt.Sprintf("invoke %d", k.Spell), 1)
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
	}
}
