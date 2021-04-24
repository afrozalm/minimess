package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "127.0.0.1:8080", "the port at which the server listenes to client requests")
	flag.Parse()
	// fmt.Println("got port as ", *port)
	// fmt.Println("Current time in unix epoch seconds is", time.Now().Unix())

	http.HandleFunc("/ws", serveTimestamps)

	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Println("ListenAndServe failed with:", err)
	}
}
