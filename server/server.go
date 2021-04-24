package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func serveTimestamps(w http.ResponseWriter, r *http.Request) {
	log.Println("request url: ", r.URL, "request method: ", r.Method)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Could not upgrade to ws connection")
	}
	go handleConn(conn)
}

func handleConn(conn *websocket.Conn) {
	ticker := time.NewTicker(5 * time.Second)

	defer conn.Close()
	defer ticker.Stop()

	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		log.Println("pong received")
		return nil
	})

	for {
		select {
		case <-ticker.C:
			conn.WriteMessage(websocket.TextMessage, []byte("I'm alive"))
			log.Println("sent message")
			conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
			w, err := conn.NextWriter(websocket.PingMessage)
			if err != nil {
				log.Println("going to close", err)
				return
			}

			w.Write(nil)
			log.Println("ping sent")

			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("received: %s", message)
		}
	}
}
