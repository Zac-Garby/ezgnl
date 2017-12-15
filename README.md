# ezgnl

--- Easy Game Networking Library

ezgnl is a networking library in Go suitable for making game servers/clients - or,
at least, it will be. It works as a generic networking library (a bit like socket.io),
but I've yet to add features useful for games. Things I'm planning on are:

 - Lobbies
 - Improved speed
 - Listing local servers
 - Rejection of connects from the server
 - Kicking clients
 - **Make disconnect handlers work**

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

**The client**

```go
// Sends a message to the echo server
func main() {
	c := client.New("localhost:8080")
	defer c.Close()

	reader := bufio.NewReader(os.Stdin)

	c.Handle("reply", func(msg interface{}) {
		log.Println("got reply:", msg)
	})

	// Listen concurrently but still report errors
	go func() {
		if err := c.Listen("tcp"); err != nil {
			log.Println("listen:", err)
		}
	}()

	c.Send("message", "hello, world!")
}
```

As you can see, much easier than using the `net` package directly.
