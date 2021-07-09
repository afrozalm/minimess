package interfaces

type Storage interface {
	// a method to add topic-host entry
	AddHostTopic(string, string) error

	// a method to remove topic-host
	RemoveHostTopic(string, string) error

	// a method to retrieve list of hosts subscribed to topic
	GetHostsForTopic(string) []string

	// a method to retrieve list of topics subscribed by host
	GetTopicsForHost(string) []string
}
