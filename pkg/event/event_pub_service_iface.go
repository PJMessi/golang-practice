package event

type PubService interface {
	Close() error
	Publish(topic string, payload []byte) error
}
