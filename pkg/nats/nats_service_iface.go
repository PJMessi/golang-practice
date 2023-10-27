package nats

type Service interface {
	Close()
	Publish(topic string, payload []byte) error
	Subscribe(topic string) error
}
