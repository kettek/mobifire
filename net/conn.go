package net

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/kettek/termfire/messages"
)

// Connection is a connection to a server.
type Connection struct {
	net.Conn
}

// Join attempts to join the given server.
func (c *Connection) Join(server string) error {
	if !strings.Contains(server, ":") {
		server += ":13327"
	}

	conn, err := net.DialTimeout("tcp", server, time.Duration(5)*time.Second)
	if err != nil {
		return err
	}
	c.Conn = conn

	go c.readLoop()

	return nil
}

// Close closes the connection to the server.
func (c *Connection) Close() {
	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
	}
}

func (c *Connection) readLoop() {
	for {
		println("TODO: readLoop")
		return
	}
}

// Send send a message.
func (c *Connection) Send(msg messages.Message) error {
	bytes := msg.Bytes()
	if len(bytes) > 0 {
		c.Write([]byte{byte(len(bytes) >> 8), byte(len(bytes))})
		c.Write(bytes)
		return nil
	}
	return errors.New("empty message")
}
