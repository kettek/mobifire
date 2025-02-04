package play

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/mobifire/net"
	"github.com/kettek/termfire/messages"
)

type command struct {
	Name       string
	OnActivate func()
}

type commandsManager struct {
	conn              *net.Connection
	commands          []command
	pendingQueries    []queryCommand
	OnCommandComplete func(command string, text string)
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
			if q.Callback != nil {
				q.Callback(q)
			} else if cm.OnCommandComplete != nil {
				cm.OnCommandComplete(q.Command, q.Text)
			}
			cm.pendingQueries = append(cm.pendingQueries[:i], cm.pendingQueries[i+1:]...)
			return true
		}
	}
	return false
}

func (cm *commandsManager) QuerySimpleCommand(cmd string, mt messages.MessageType, st messages.SubMessageType) {
	id, _ := cm.conn.SendCommand(cmd, 0)
	cm.pendingQueries = append(cm.pendingQueries, queryCommand{
		PacketID: id,
		Command:  cmd,
		Text:     "",
		MT:       mt,
		ST:       st,
	})
}

type queryCommand struct {
	PacketID uint16
	Command  string
	Text     string
	MT       messages.MessageType
	ST       messages.SubMessageType
	Callback func(cmd queryCommand)
}
