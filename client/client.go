package client

import (
	"net"

	"github.com/Zac-Garby/ezgnl/message"
	"github.com/satori/go.uuid"
)

// A UUID is assigned uniquely to each client.
type UUID = uuid.UUID

// A MessageHandler is a function which is called for
// each received message.
type MessageHandler func(data interface{})

// A Client is used to send messages to the server and
// to receieve and react to messages.
type Client struct {
	conn net.Conn
	id   UUID

	incoming, outgoing chan *message.Message

	handlers map[string]MessageHandler
}

// Handle sets the message handler of a certain message type.
func (c *Client) Handle(t string, handler MessageHandler) {
	c.handlers[t] = handler
}

// Send sends a message to the server.
func (c *Client) Send(t string, data interface{}) {
	c.outgoing <- &message.Message{
		Type: t,
		Data: data,
	}
}
