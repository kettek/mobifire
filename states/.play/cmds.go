package play

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

type command struct {
	Name       string
	OnActivate func()
	OnSubmit   func(cmd string)
}

type commandsManager struct {
	conn              *net.Connection
	commands          []command
	pendingQueries    []*queryCommand
	OnCommandComplete func(*queryCommand)
}

func (cm *commandsManager) toMenuItems() []*fyne.MenuItem {
	var items []*fyne.MenuItem
	for i, c := range cm.commands {
		items = append(items, fyne.NewMenuItem(c.Name, func() {
			cm.commands[i].OnActivate()
		}))
	}
	return items
}

func (cm *commandsManager) trigger(name string, args ...string) {
	for _, c := range cm.commands {
		if c.Name == name {
			// TODO: Send command
			return
		}
	}
}

func (cm *commandsManager) checkDrawExtInfo(msg *messages.MessageDrawExtInfo) bool {
	for i, q := range cm.pendingQueries {
		if msg.Type == q.MT && msg.Subtype == q.ST {
			q.Text += msg.Message + "\n"
			cm.pendingQueries[i] = q
			return true
		}
	}
	return false
}

func (cm *commandsManager) checkCommandCompleted(msg *messages.MessageCommandCompleted) bool {
	for i, q := range cm.pendingQueries {
		if msg.Packet == q.PacketID {
			if q.Text != "" {
				// Trim trailing newline.
				q.Text = q.Text[:len(q.Text)-1]
			}
			if cm.OnCommandComplete != nil {
				cm.OnCommandComplete(q)
			}
			cm.pendingQueries = append(cm.pendingQueries[:i], cm.pendingQueries[i+1:]...)
			return true
		}
	}
	return false
}

func (cm *commandsManager) QuerySimpleCommand(cmd string, mt messages.MessageType, st messages.SubMessageType) {
	id, _ := cm.conn.SendCommand(cmd, 0)
	cm.pendingQueries = append(cm.pendingQueries, &queryCommand{
		PacketID:        id,
		Command:         cmd,
		OriginalCommand: cmd,
		MT:              mt,
		ST:              st,
	})
}

func (cm *commandsManager) QuerySimpleCommandWithInput(cmd string, mt messages.MessageType, st messages.SubMessageType) *queryCommand {
	id, _ := cm.conn.SendCommand(cmd, 0)
	cm.pendingQueries = append(cm.pendingQueries, &queryCommand{
		PacketID:        id,
		Command:         cmd,
		OriginalCommand: cmd,
		HasInput:        true,
		MT:              mt,
		ST:              st,
	})
	return cm.pendingQueries[len(cm.pendingQueries)-1]
}

func (cm *commandsManager) QueryComplexCommand(cmd, origCmd string, mt messages.MessageType, st messages.SubMessageType) *queryCommand {
	id, _ := cm.conn.SendCommand(cmd, 0)
	cm.pendingQueries = append(cm.pendingQueries, &queryCommand{
		PacketID:        id,
		Command:         cmd,
		OriginalCommand: origCmd,
		Text:            "",
		MT:              mt,
		ST:              st,
	})
	return cm.pendingQueries[len(cm.pendingQueries)-1]
}

func (cm *commandsManager) QueryComplexCommandWithInput(cmd, origCmd string, mt messages.MessageType, st messages.SubMessageType) *queryCommand {
	id, _ := cm.conn.SendCommand(cmd, 0)
	cm.pendingQueries = append(cm.pendingQueries, &queryCommand{
		PacketID:        id,
		Command:         cmd,
		OriginalCommand: origCmd,
		HasInput:        true,
		Text:            "",
		MT:              mt,
		ST:              st,
	})
	return cm.pendingQueries[len(cm.pendingQueries)-1]
}

func (cm *commandsManager) QueryCommand(cmd queryCommand) *queryCommand {
	id, _ := cm.conn.SendCommand(cmd.Command, 0)
	cmd.PacketID = id
	cmd.Text = ""
	cm.pendingQueries = append(cm.pendingQueries, &cmd)
	return cm.pendingQueries[len(cm.pendingQueries)-1]
}

type queryCommand struct {
	PacketID        uint16
	Command         string
	OriginalCommand string
	Text            string
	SubmitText      string
	Repeat          bool // Can keep requesting input.
	HasInput        bool
	MT              messages.MessageType
	ST              messages.SubMessageType
}
