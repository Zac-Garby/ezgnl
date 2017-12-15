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
	errors             chan error
	closeChan          chan bool

	handlers map[string]MessageHandler
}

// New constructs a new Client, which will connect to the
// given address. But, New doesn't actually connect to the
// server yet -- use .Connect() for that.
func New(addr string) *Client {
	return &Client{
		address:   addr,
		closed:    true,
		incoming:  make(chan *message.Message),
		outgoing:  make(chan *message.Message),
		errors:    make(chan error),
		closeChan: make(chan bool),
		handlers:  make(map[string]MessageHandler),
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

	go c.sendLoop()
	go c.receiveLoop()

	for !c.closed {
		select {
		case err := <-c.errors:
			c.conn.Close()
			return err

		case _ = <-c.closeChan:
			c.closed = true
			break

		case msg := <-c.incoming:
			c.handleMessage(msg)

		case msg := <-c.outgoing:
			message.Send(msg.Type, msg.Data, c.conn)
		}
	}

	return c.conn.Close()
}

// Close closes the client's connection.
func (c *Client) Close() {
	c.closeChan <- true
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

func (c *Client) sendLoop() {
	for {
		out := <-c.outgoing
		message.Send(out.Type, out.Data, c.conn)
	}
}

func (c *Client) receiveLoop() {
	for {
		msg, err := message.Receive(c.conn)
		if err != nil {
			// HANDLE
		}

		c.incoming <- msg
	}
}
