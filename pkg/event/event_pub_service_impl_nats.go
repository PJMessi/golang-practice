package event

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type PubServiceNatsImpl struct {
	nc *nats.Conn
}

func NewPubService() (PubService, error) {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		return nil, fmt.Errorf("nats.NewPublisherService(): %w", err)
	}

	return &PubServiceNatsImpl{nc: nc}, nil
}

func (p *PubServiceNatsImpl) Close() error {
	p.nc.Close()
	return nil
}

func (p *PubServiceNatsImpl) Publish(topic string, payload []byte) error {
	err := p.nc.Publish(topic, payload)
	if err != nil {
		return fmt.Errorf("event.PubServiceNatsImpl.Publish(): %w", err)
	}
	return nil
}
