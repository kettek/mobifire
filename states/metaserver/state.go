package metaserver

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/mobifire/data"
	"github.com/kettek/mobifire/data/settings"
	"github.com/kettek/mobifire/states"
	"github.com/kettek/mobifire/states/join"
	"github.com/kettek/rebui"
	_ "github.com/kettek/rebui/defaults/font"
	"github.com/kettek/rebui/widgets"
	_ "github.com/kettek/rebui/widgets"
	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/messages"
)

var metaservers = []string{
	"https://crossfire.real-time.com/metaserver2/meta_client.php",
	"http://metaserver.eu.cross-fire.org/meta_client.php",
	"https://metaserver.us.cross-fire.org/meta_client.php",
}

// State provides a list of servers the user can join.
type State struct {
	next        func(states.State)
	layout      rebui.Layout
	itemNodes   []*rebui.Node
	addressNode *rebui.Node
	joinNode    *rebui.Node
	entryNode   rebui.Node // This is just a template.
}

func (s *State) Update() error {
	s.layout.Update()
	return nil
}

func (s *State) Draw(screen *ebiten.Image) {
	s.layout.Draw(screen)
}

// Enter sets up the base UI containers and loads the server list.
func (s *State) Enter(next func(states.State)) (leave func()) {
	s.next = next

	address := settings.GetWithFallback[string]("address", "localhost:13327")

	layout, err := data.GetLayout("metaserver")
	if err != nil {
		panic(err)
	}
	for _, node := range layout.Nodes {
		s.layout.AddNode(node)
	}

	s.addressNode = s.layout.GetByID("address")
	s.addressNode.Widget.(*widgets.TextInput).AssignText(address)
	s.addressNode.Widget.(*widgets.TextInput).OnChange = func(text string) {
		address = text
	}

	s.joinNode = s.layout.GetByID("connect")
	s.joinNode.OnPointerPressed = func(rebui.EventPointerPressed) {
		var hostname string
		var port int64

		if address == "" {
			fmt.Println("Please enter a valid address.")
			return
		}

		parts := strings.Split(address, ":")
		if len(parts) != 2 {
			fmt.Println("Invalid address format. Please use 'hostname:port'.")
		}
		hostname = parts[0]
		port, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			fmt.Println("Invalid port number. Please enter a valid port.", err)
			return
		}

		// Save it!
		settings.Set("address", address)

		s.next(&join.State{
			Hostname: hostname,
			Port:     int(port),
		})
	}

	// Get our template entry node for reference and then remove it from our layout.
	s.entryNode = *s.layout.GetByID("template__entry")
	s.layout.RemoveNode(s.layout.GetByID("template__entry"))

	s.refreshMetaservers()

	return func() {
		s.layout = rebui.Layout{}
	}
}

// refreshMetaservers iterates thru metaservers and generates non-duplicate servers.
func (s *State) refreshMetaservers() {

	// Generate server entries from the metaservers.
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

	// Create the container with the server list.
	for i, entry := range serverEntries {
		id := fmt.Sprintf("server-%d", i)
		y := "2%"
		if i > 0 {
			y = fmt.Sprintf("after %s", fmt.Sprintf("server-%d", i-1))
		}

		s.entryNode.ID = id
		s.entryNode.Text = fmt.Sprintf("%s:%d", entry.Hostname, entry.Port)
		s.entryNode.Y = y

		node := s.layout.AddNode(s.entryNode)
		node.OnPointerPressed = func(rebui.EventPointerPressed) {
			s.addressNode.Widget.(*widgets.TextInput).AssignText(fmt.Sprintf("%s:%d", entry.Hostname, entry.Port))
		}
		s.itemNodes = append(s.itemNodes, node)
	}
}

// requestServers requests the servers from the given metaserver with a 5 second timeout.
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
