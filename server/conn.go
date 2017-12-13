package server

import (
	"net"

	"github.com/Zac-Garby/ezgnl/message"
)

// A Conn is a higher-level wrapper around the net.Conn type.
type Conn interface {
	net.Conn

	// Handle handles a message of type t.
	Handle(t string, fn MessageHandler)

	// Disconnect handles a message of type "disconnect".
	Disconnect(fn MessageHandler)

	// Send sends a message of the given type to the client.
	Send(t string, msg interface{})

	// handle actually calls the handler function.
	handle(msg *message.Message)
}

// connection is the internal Conn implementor.
type connection struct {
	net.Conn

	handlers map[string]MessageHandler
}

// Handle handles a message of type t.
func (c *connection) Handle(t string, fn MessageHandler) {
	c.handlers[t] = fn
}

// Disconnect handles a message of type "disconnect".
func (c *connection) Disconnect(fn MessageHandler) {
	c.handlers["disconnect"] = fn
}

// Send sends a message of the given type to the client.
func (c *connection) Send(t string, msg interface{}) {
	message.Send(t, msg, c.Conn)
}

// handle actually calls the handler function.
func (c *connection) handle(msg *message.Message) {
	if h, ok := c.handlers[msg.Type]; ok {
		h(msg.Data)
	}
}
