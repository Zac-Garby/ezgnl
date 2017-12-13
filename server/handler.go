package server

// A ConnectionHandler handles new connections to the server.
type ConnectionHandler func(conn Conn, id UUID)

// A MessageHandler handles incoming messages to the server.
type MessageHandler func(data interface{})

// Accept sets the server's connection handling function.
func (s *Server) Accept(fn ConnectionHandler) {
	s.connHandler = fn
}
