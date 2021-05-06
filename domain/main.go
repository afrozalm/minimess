package main

import (
	"flag"

	"github.com/afrozalm/minimess/domain/server"
)

func main() {
	port := flag.String("port", "8080", "http service address")
	flag.Parse()
	s := server.Create()
	s.Run("127.0.0.1", *port)
}
