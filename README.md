# ezgnl

--- Easy Game Networking Library

## An example

**The server**

```go
package main

import (
    "log"

    "github.com/Zac-Garby/ezgnl/server"
)

type MoveResult struct {
    Collided bool
}

func main() {
    s, _ := server.New(":8080")
    defer s.Close()

    s.Accept(func(conn server.Conn, id server.UUID) {
        conn.Handle("move", func() {
            dir := conn.Data("direction")
            log.Println("user", id, "moved in direction:", dir)

            conn.Send("result", MoveResult{ Collided: true })
        })

        conn.Disconnect(func() {
            log.Println("user", id, "has left the server")
        })
    })

    _ = s.Listen()
}
```

**The client**

```go
package main

import (
    "log"

    "github.com/Zac-Garby/ezgnl/client"
)

type MoveMessage struct {
    Direction string
}

func main() {
    c, _ := client.New("localhost:8080")
    defer c.Close()

    time.Sleep(time.Second * 2)
    c.Send("move", MoveMessage{ Direction: "right" })

    time.Sleep(time.Second)
    c.Send("move", MoveMessage{ Direction: "up" })
}
```
