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
	packetId       uint16
	OnLoss         func(error)
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
	c.packetId = 1 // Skip 0 for default value sanity

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
			c.Close()
			err = errors.Join(err, errors.New("failed to read message length"))
			fmt.Println(err)
			if c.OnLoss != nil {
				c.OnLoss(err)
			}
			return
		}
		size := (int(length[0]) << 8) | int(length[1])
		buf := make([]byte, size)
		err = c.ReadBytes(buf, size)
		if err != nil {
			c.Close()
			err = errors.Join(err, errors.New("failed to read message"))
			fmt.Println(err)
			if c.OnLoss != nil {
				c.OnLoss(err)
			}
			return
		}
		message, err := messages.UnmarshalMessage(buf)
		if err != nil {
			c.Close()
			err = errors.Join(err, errors.New("failed to unmarshal message"))
			fmt.Println(err)
			if c.OnLoss != nil {
				c.OnLoss(err)
			}
			return
		}

		//fmt.Printf("msg %+v\n", message)
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

func (c *Connection) ReadBytes(buf []byte, size int) error {
	pos := 0
	for {
		n, err := c.Read(buf[pos:])
		if err != nil {
			return err
		}
		pos += n
		if pos == size {
			return nil
		}
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

func (c *Connection) SendCommand(command string, repeat uint32) (uint16, error) {
	msg := messages.MessageCommand{Command: command, Repeat: repeat, Packet: c.packetId}
	c.packetId++
	return c.packetId - 1, c.Send(&msg)
}
