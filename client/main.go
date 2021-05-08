package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/afrozalm/minimess/client/cli"
	"github.com/afrozalm/minimess/client/connHandler"
)

func main() {
	userID := flag.String("uid", "", "give me your user id")
	flag.Parse()
	if len(*userID) == 0 {
		log.Fatal("please provide a non-empty username")
	}
	fmt.Println("user id given by the user is ", *userID)

	f, err := os.OpenFile(fmt.Sprintf("%s.mess", *userID), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	h := connHandler.NewHandler(*userID)
	h.Run()
	cli.Run(h)
}
