package metaserver

import (
	"io"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/mobifire/states/play"
	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/messages"
)

var metaservers = []string{
	"http://crossfire.real-time.com/metaserver2/meta_client.php",
	"http://metaserver.eu.cross-fire.org/meta_client.php",
	"http://metaserver.us.cross-fire.org/meta_client.php",
}

type State struct {
	container  *fyne.Container
	serverList *fyne.Container
}

func (s *State) Enter(next func(states.State)) (leave func()) {
	button := widget.NewButton("nextie", func() {
		next(&play.State{})
	})

	s.serverList = container.New(layout.NewVBoxLayout())

	s.container = container.New(layout.NewVBoxLayout(), s.serverList, button)

	// Load servers on load, obv.
	s.refreshMetaservers()

	return nil
}

func (s *State) refreshMetaservers() {
	s.serverList.RemoveAll()
	var serverEntries messages.ServerEntries
	for _, m := range metaservers {
		entries, err := s.requestServers(m)
		if err != nil {
			debug.Debug("Failed to get server list from metaserver: ", err)
			continue
		}
		for _, e := range entries {
			found := false
			for _, server := range serverEntries {
				if server.Hostname == e.Hostname && server.Port == e.Port {
					found = true
					break
				}
			}
			if !found {
				serverEntries = append(serverEntries, e)
			}
		}
	}

	accordion := widget.NewAccordion()
	for _, e := range serverEntries {
		acc := widget.NewAccordionItem(e.Hostname, widget.NewLabel(e.TextComment))
		accordion.Append(acc)
	}
	s.serverList.Add(accordion)
}

func (s *State) Container() *fyne.Container {
	return s.container
}

func (s *State) requestServers(metaserver string) (messages.ServerEntries, error) {
	resp, err := http.Get(metaserver)
	http.DefaultClient.Timeout = 5 * time.Second
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	serverEntries := messages.ServerEntries{}

	err = serverEntries.UnmarshalBinary(body)
	if err != nil {
		return nil, err
	}

	return serverEntries, nil
}
