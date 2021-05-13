package main

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/afrozalm/minimess/domain/connHandler"
	"github.com/afrozalm/minimess/domain/server"
)

func main() {
	port := flag.String("port", "8080", "http service address")
	flag.Parse()
	log.SetLevel(log.WarnLevel)
	s := server.NewServer()
	addr := "127.0.0.1:" + *port
	run(s, addr)
}

func run(s *server.Server, addr string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connHandler.ServeWSConn(s, w, r)
	})
	log.Info("going to run server at", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe failed with:", err)
	}
}
