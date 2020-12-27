package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8889/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
	for i := 0; i < 10; i++ {
		err = c.WriteMessage(websocket.TextMessage, []byte(`/testservice1/hello/1.1
		{"name":"jack"}`))
		if err != nil {
			log.Println("write:", err)
			return
		}
		time.Sleep(time.Second)
	}
}
