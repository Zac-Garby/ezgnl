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
	newConns  chan net.Conn
	errors    chan error
	closeChan chan bool

	// connHandler is called for each
	connHandler ConnectionHandler
	handlers    map[string]MessageHandler

	// whether the server is closed or not
	closed bool
}

// New constructs a new server on the given port, but doesn't
// start listening yet.
func New(port string) *Server {
	return &Server{
		Port:        port,
		connections: make(map[UUID]Conn),
		incomings:   make(chan incoming),
		outgoings:   make(chan outgoing),
		newConns:    make(chan net.Conn),
		errors:      make(chan error),
		closeChan:   make(chan bool),
		closed:      true,
	}
}

// Listen starts listening for incoming connections. The network
// parameter is the type of network to use, for example "tcp" or
// "udp".
func (s *Server) Listen(network string) error {
	go s.waitForConnections()

	s.closed = false

	for !s.closed {
		select {
		case err := <-s.errors:
			s.closeAll()
			return err

		case _ = <-s.closeChan:
			break

		case msg := <-s.incomings:
			go s.handleMessage(msg.sender, msg.Message)

		case msg := <-s.outgoings:
			go message.Send(msg.Type, msg.Data, s.connections[msg.receiver])

		case conn := <-s.newConns:
			go s.handleConnection(conn)
		}
	}

	return s.closeAll()
}

// Close stops the server listening and closes the connection.
func (s *Server) Close() {
	s.closeChan <- true
}

// closeAll closes all open connections.
func (s *Server) closeAll() error {
	for _, conn := range s.connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) waitForConnections() {
	ln, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		s.errors <- err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.errors <- err
		}

		s.newConns <- conn
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	var (
		id = uuid.NewV4()

		c = &connection{
			Conn:     conn,
			handlers: make(map[string]MessageHandler),
		}
	)

	s.connHandler(c, id)
	s.connections[id] = c

	go s.awaitMessages(id, c)
}

func (s *Server) awaitMessages(id UUID, conn Conn) {
	for {
		msg, err := message.Receive(conn)

		// An error was most likely caused by the client
		// disconnecting, but either way, end the connection.
		if err != nil {
			if err := s.Disconnect(id); err != nil {
				s.errors <- err
			}

			break
		}

		s.incomings <- incoming{
			Message: msg,
			sender:  id,
		}
	}
}

func (s *Server) handleMessage(id UUID, msg *message.Message) {
	s.connections[id].handle(msg)
}

// Disconnect disconnects a client from the server
func (s *Server) Disconnect(id UUID) error {
	if conn, ok := s.connections[id]; ok {
		delete(s.connections, id)
		return conn.Close()
	}

	return nil
}
