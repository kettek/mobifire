package play

import (
	"fyne.io/fyne/v2"
	"github.com/kettek/termfire/messages"
)

type command struct {
	Name       string
	OnActivate func()
	OnMessage  func(msg *messages.MessageDrawExtInfo) bool
}

type commandsManager struct {
	commands []command
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
