package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/afrozalm/minimess/server/connHandler"
	"github.com/afrozalm/minimess/server/supervisor"
)

func main() {
	port := flag.String("port", "8080", "listening to client on this port")
	flag.Parse()
	log.SetLevel(log.TraceLevel)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	s := supervisor.NewSupervisor()
	go s.Run()
	addr := "127.0.0.1:" + *port
	go run(s, addr)

	<-sigs
	s.Close()
}

func run(s *supervisor.Supervisor, addr string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connHandler.ServeWSConn(s, w, r)
	})
	log.Info("going to run server at ", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe failed with:", err)
	}
}
