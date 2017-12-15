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
	address string
	conn    net.Conn
	id      UUID
	closed  bool

	incoming, outgoing chan *message.Message

	handlers map[string]MessageHandler
}

// New constructs a new Client, which will connect to the
// given address. But, New doesn't actually connect to the
// server yet -- use .Connect() for that.
func New(addr string) *Client {
	return &Client{
		address:  addr,
		closed:   true,
		incoming: make(chan *message.Message),
		outgoing: make(chan *message.Message),
		handlers: make(map[string]MessageHandler),
	}
}

// Listen dials the address specified in New() and attempts
// to connect to the server.
//
// network is the network to connect to, such as "tcp" or "udp".
func (c *Client) Listen(network string) error {
	conn, err := net.Dial(network, c.address)
	if err != nil {
		return err
	}

	c.conn = conn
	c.closed = false

	for !c.closed {
		msg, err := message.Receive(conn)
		if err != nil {
			return err
		}

		c.handleMessage(msg)
	}

	return nil
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

func (c *Client) handleMessage(msg *message.Message) {
	if handler, ok := c.handlers[msg.Type]; ok {
		handler(msg.Data)
	}
}
