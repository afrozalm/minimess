package storage_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/afrozalm/minimess/server/fanout/storage"
)

var (
	store *storage.CassandraStorage
)

func TestAddHostTopicNoThrow(t *testing.T) {
	store.AddHostTopic("localhost:9090", "general")
}

func TestRemoveHostTopicNoThrow(t *testing.T) {
	store.RemoveHostTopic("localhost:9090", "general")
}

func TestGetHostsForTopicNoThrow(t *testing.T) {
	hosts := store.GetHostsForTopic("general")
	fmt.Println("found following hosts for topic 'general'")
	for _, host := range hosts {
		fmt.Println(host)
	}
}

func TestGetTopicsForHostNoThrow(t *testing.T) {
	topics := store.GetTopicsForHost("localhost:9090")
	fmt.Println("found the following topics for host 'localhost:9090'")
	for _, topic := range topics {
		fmt.Println(topic)
	}
}

func TestMain(m *testing.M) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	store, err = storage.NewCassandraStorage(ctx)
	if err != nil {
		log.Fatalf("could not setup storage due to %v", err)
	}
	store.AddHostTopic("localhost:9090", "random")
	store.AddHostTopic("localhost:8080", "general")
	store.AddHostTopic("localhost:9090", "all")
	m.Run()
	store.Close()
	cancel()
}
