package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrozalm/minimess/constants"
	"github.com/afrozalm/minimess/server/fanout/broadcaster"
	fanout "github.com/afrozalm/minimess/server/fanout/fanout_proto"
	"github.com/afrozalm/minimess/server/fanout/register"
	"github.com/afrozalm/minimess/server/fanout/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	offset = flag.Int("offset", 0, "port offset for serving registerer")
)

func main() {
	flag.Parse()
	log.SetLevel(log.TraceLevel)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	cassandraStorage, err := storage.NewCassandraStorage(ctx)
	if err != nil {
		log.Fatal("could not create cassandra storage interface")
	}

	bcaster := broadcaster.NewBroadcaster(ctx, cassandraStorage)
	registerer := register.NewRegister(cassandraStorage)
	go bcaster.Run()
	go runGRPCServer(registerer)

	log.Info("waiting for interrupt signal")
	<-sig
	cancel()
	time.Sleep(time.Second * 5)
}

func runGRPCServer(registerer *register.Register) {
	addr := fmt.Sprintf("127.0.0.1:%d", constants.FanoutBasePort+*offset)
	log.Info("starting grpc registerer at ", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("could not start grpc server at address ", addr)
	}
	var ops []grpc.ServerOption
	grpcServer := grpc.NewServer(ops...)
	fanout.RegisterFanoutServer(grpcServer, registerer)
	grpcServer.Serve(lis)
}
