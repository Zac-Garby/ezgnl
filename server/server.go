package server

import (
	"net"

	"github.com/Zac-Garby/ezgnl/message"
	"github.com/satori/go.uuid"
)

// A UUID is a unique ID given to each connection.
type UUID = uuid.UUID

type incoming struct {
	*message.Message
	sender UUID
}

type outgoing struct {
	*message.Message
	receiver UUID
}

// A Server is a server which can receive/send messages
// from/to clients.
type Server struct {
	// Port is the port to listen for connections on.
	Port string

	// connections maps the clients' UUIDs to their connections.
	connections map[uuid.UUID]Conn

	// incoming and outgoing are the message queues for
	// received/to-be-sent messages.
	incomings chan incoming
	outgoings chan outgoing

	// connHandler is called for each
	connHandler ConnectionHandler
	handlers    map[string]MessageHandler
}

// New constructs a new server on the given port, but doesn't
// start listening yet.
func New(port string) *Server {
	return &Server{
		Port:        port,
		connections: make(map[UUID]Conn),
		incomings:   make(chan incoming),
		outgoings:   make(chan outgoing),
	}
}

// Listen starts listening for incoming connections. The network
// parameter is the type of network to use, for example "tcp" or
// "udp".
func (s *Server) Listen(network string) error {
	ln, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		return err
	}

	go s.processIncoming()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err // TODO: It might be worth closing all connections here
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	id := uuid.NewV4()

	s.connections[id] = &connection{
		Conn:     conn,
		handlers: make(map[string]MessageHandler),
	}

	go s.awaitMessages(id, s.connections[id])
}

func (s *Server) awaitMessages(id UUID, conn Conn) {
	for {
		msg, err := message.Receive(conn)

		// An error was most likely caused by the client
		// disconnecting, but either way, end the connection.
		if err != nil {
			s.Disconnect(id)
			break
		}

		s.incomings <- incoming{
			Message: msg,
			sender:  id,
		}
	}
}

func (s *Server) processIncoming() {
	for {
		msg := <-s.incomings
		go s.handleMessage(msg.sender, msg.Message)
	}
}

func (s *Server) handleMessage(id UUID, msg *message.Message) {
	s.connections[id].handle(msg)
}

// Disconnect disconnects a client from the server
func (s *Server) Disconnect(id UUID) {

}
