package net

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/kettek/termfire/messages"
)

// Connection is a connection to a server.
type Connection struct {
	net.Conn
	OnMessage      func(messages.Message)
	queuedMessages []messages.Message
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
		var length [2]byte
		n, err := c.Read(length[:])
		if err != nil || n != 2 {
			fmt.Println("Error reading message length:", err)
			c.Close()
			return
		}
		buf := make([]byte, (int(length[0])<<8)|int(length[1]))
		n, err = c.Read(buf)
		if err != nil || n != len(buf) {
			fmt.Println("Error reading message:", err)
			c.Close()
			return
		}
		message, err := messages.UnmarshalMessage(buf)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			c.Close()
			return
		}

		fmt.Printf("msg %+v\n", message)
		if c.OnMessage != nil {
			c.OnMessage(message)
		} else {
			c.queuedMessages = append(c.queuedMessages, message)
		}
	}
}

func (c *Connection) SetMessageHandler(handler func(messages.Message)) {
	c.OnMessage = handler
	for _, message := range c.queuedMessages {
		handler(message)
	}
	c.queuedMessages = nil
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
