package foo

import (
	"log"
	"time"

	"fmt"

	"golang.org/x/net/websocket"
)

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func fooSocket(ws *websocket.Conn) {
	done := make(chan struct{})
	fmt.Println("New websocket connection established")

	go func(c *websocket.Conn) {
		for {
			var msg message

			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println(err)
				break
			}

			fmt.Printf("received message %s\n", msg.Data)
		}

		close(done)
	}(ws)

loop:
	for {
		select {
		case <-done:
			fmt.Println("Connection was closed, lets break out of here")
			break loop

		default:
			foos, err := getToptenFoos()

			if err != nil {
				log.Println(err)
				break
			}

			if err := websocket.JSON.Send(ws, foos); err != nil {
				log.Println(err)
				break
			}

			time.Sleep(10 * time.Second)
		}
	}

	fmt.Println("Closing the websocket")
	defer ws.Close()
}
