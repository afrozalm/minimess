package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/afrozalm/minimess/domain/connHandler"
	"github.com/afrozalm/minimess/domain/server"
)

func main() {
	port := flag.String("port", "8080", "http service address")
	flag.Parse()
	s := server.NewServer()
	addr := "127.0.0.1:" + *port
	run(s, addr)
}

func run(s *server.Server, addr string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connHandler.ServeWSConn(s, w, r)
	})
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe failed with:", err)
	}
}
