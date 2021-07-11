package storage

import (
	"context"

	"github.com/afrozalm/minimess/constants"
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

type CassandraStorage struct {
	session *gocql.Session
	ctx     context.Context
}

func NewCassandraStorage(ctx context.Context) (*CassandraStorage, error) {
	cluster := gocql.NewCluster(constants.CassandraHosts...)
	cluster.Keyspace = "subscriptions"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &CassandraStorage{
		session: session,
		ctx:     ctx,
	}, nil
}

func (s *CassandraStorage) Close() {
	s.session.Close()
}

func (s *CassandraStorage) AddHostTopic(host, topic string) error {
	log.Tracef("adding hosttopic '%s' '%s'", host, topic)
	if err := s.session.Query(`INSERT INTO topic_host (topic, host) VALUES (?, ?) USING TTL 600`,
		topic, host).WithContext(s.ctx).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *CassandraStorage) RemoveHostTopic(host, topic string) error {
	log.Tracef("removing hosttopic '%s' '%s'", host, topic)
	if err := s.session.Query(`DELETE FROM topic_host WHERE topic=? and host=?`,
		topic, host).WithContext(s.ctx).Exec(); err != nil {
		return err
	}
	return nil
}

func (s *CassandraStorage) GetHostsForTopic(topic string) []string {
	log.Tracef("get all hosts for topic '%s'", topic)
	scanner := s.session.Query(`SELECT host FROM topic_host WHERE topic=?`, topic).
		WithContext(s.ctx).
		Iter().
		Scanner()

	hosts := make([]string, 0, 5)
	for scanner.Next() {
		var host string
		err := scanner.Scan(&host)
		if err == nil {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

func (s *CassandraStorage) GetTopicsForHost(host string) []string {
	log.Tracef("get all topics for host '%s'", host)
	scanner := s.session.Query(`SELECT topic FROM topic_host WHERE host=?`, host).
		WithContext(s.ctx).
		Iter().
		Scanner()

	topics := make([]string, 0, 5)
	for scanner.Next() {
		var topic string
		err := scanner.Scan(&topic)
		if err == nil {
			topics = append(topics, topic)
		}
	}
	return topics
}
