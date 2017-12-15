package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Zac-Garby/ezgnl/client"
)

func main() {
	c := client.New("localhost:8080")
	defer c.Close()

	reader := bufio.NewReader(os.Stdin)

	c.Handle("reply", func(msg interface{}) {
		log.Println("got reply:", msg)
	})

	go func() {
		if err := c.Listen("tcp"); err != nil {
			log.Println("listen:", err)
		}
	}()

	for {
		getLine(reader, c)
	}
}

func getLine(r *bufio.Reader, c *client.Client) {
	fmt.Print(">> ")

	line, err := r.ReadString('\n')
	if err != nil {
		log.Println("read:", err)
		return
	}

	c.Send("message", strings.TrimSpace(line))
}
