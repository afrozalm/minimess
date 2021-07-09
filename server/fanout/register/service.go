package register

import (
	"context"

	fanout "github.com/afrozalm/minimess/server/fanout/fanout_proto"
	"github.com/afrozalm/minimess/server/fanout/interfaces"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Register struct {
	fanout.UnimplementedFanoutServer
	storage interfaces.Storage
}

func NewRegister(storage interfaces.Storage) *Register {
	return &Register{storage: storage}
}

func (r *Register) Subscribe(ctx context.Context, ht *fanout.HostTopic) (*emptypb.Empty, error) {
	host := ht.GetHost()
	topic := ht.GetTopic()
	if err := r.storage.AddHostTopic(host, topic); err != nil {
		return new(emptypb.Empty), err
	}
	return new(emptypb.Empty), nil
}

func (r *Register) Unsubscribe(ctx context.Context, ht *fanout.HostTopic) (*emptypb.Empty, error) {
	host := ht.GetHost()
	topic := ht.GetTopic()
	if err := r.storage.RemoveHostTopic(host, topic); err != nil {
		return new(emptypb.Empty), err
	}
	return new(emptypb.Empty), nil
}

func (r *Register) UnsubscribeAll(ctx context.Context, host *fanout.Host) (*emptypb.Empty, error) {
	hostname := host.GetHost()
	topics := r.storage.GetTopicsForHost(hostname)
	for _, topic := range topics {
		if err := r.storage.RemoveHostTopic(hostname, topic); err != nil {
			return new(emptypb.Empty), err
		}
	}
	return new(emptypb.Empty), nil
}
