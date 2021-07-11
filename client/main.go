package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	h := connHandler.NewHandler(*userID)
	go h.Run(ctx)
	go cli.Run(h, ctx)

	<-sig
	cancel()
	time.Sleep(5 * time.Second)
}
