package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	userID := flag.String("uid", "", "give me your user id")
	flag.Parse()
	fmt.Println("user id given by the user is ", *userID)
	messagecount := 5
	endpoint := "127.0.0.1:8080"
	u := url.URL{Scheme: "ws", Host: endpoint, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dail:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("got: %s\n", message)
		}
	}()

	err = c.WriteMessage(websocket.TextMessage,
		[]byte("going to send"+string(messagecount)+"messages\n"))
	if err != nil {
		c.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, ""))
	}
	for i := 0; i < messagecount; i++ {
		err = c.WriteMessage(websocket.TextMessage,
			[]byte(fmt.Sprintf("hey there server, I am %s", *userID)))
		if err != nil {
			c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, ""))
		}
		log.Println("talked to the server")
		time.Sleep(time.Second * 5)
	}
	c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	log.Println("exiting")
}
