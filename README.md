# ezgnl

--- Easy Game Networking Library

## An example

**The server**

```go
// A simple echo server
func main() {
	s := server.New("8080")
	defer s.Close()

	s.Accept(func(conn server.Conn, id server.UUID) {
		conn.Handle("message", func(msg interface{}) {
			log.Println("received message:", msg)
			conn.Send("reply", msg)
		})

		conn.Disconnect(func(interface{}) {
			log.Println("a user left")
		})
	})

	log.Println("listening...")

	if err := s.Listen("tcp"); err != nil {
		log.Println("listen:", err)
	}
}

```
