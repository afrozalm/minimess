package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/server/frontend/connHandler"
	frontend "github.com/afrozalm/minimess/server/frontend/frontend_proto"
	"github.com/afrozalm/minimess/server/frontend/supervisor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	offset     = flag.Int("offset", 0, "port offset for serving client and grpc")
	clientAddr string
	grpcAddr   string
)

func main() {
	// client server port :
	// grpc server port :
	flag.Parse()
	log.SetLevel(log.TraceLevel)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	clientAddr = fmt.Sprintf("127.0.0.1:%d", constants.FrontendClientBasePort+*offset)
	grpcAddr = fmt.Sprintf("127.0.0.1:%d", constants.FrontendGRPCBasePort+*offset)
	s := supervisor.NewSupervisor(ctx, grpcAddr)
	go s.Run()
	go runClientServer(s)
	go runGRPCServer(s)

	<-sig
	cancel()
	time.Sleep(time.Second * 5)
}

func runClientServer(s *supervisor.Supervisor) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connHandler.ServeWSConn(s, w, r)
	})

	log.Info("starting client server at ", clientAddr)
	if err := http.ListenAndServe(clientAddr, nil); err != nil {
		log.Fatal("ListenAndServe failed with:", err)
	}
}

func runGRPCServer(s *supervisor.Supervisor) {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal("could not start grpc server at address ", grpcAddr)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	frontend.RegisterFrontEndServer(grpcServer, s)
	log.Info("starting grpc server at ", grpcAddr)
	grpcServer.Serve(lis)
}
