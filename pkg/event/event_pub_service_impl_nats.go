package event

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/pkg/logger"
)

type PubServiceNatsImpl struct {
	nc         *nats.Conn
	logService logger.Service
}

func NewPubService(appConfig *config.AppConfig, logService logger.Service) (PubService, error) {
	url := appConfig.NATS_URL
	if url == "" {
		url = nats.DefaultURL
		logService.Debug(fmt.Sprintf("NATS url is not provided, so using default url of %s", url))
	}

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		return nil, fmt.Errorf("nats.NewPublisherService(): %w", err)
	}

	return &PubServiceNatsImpl{nc: nc, logService: logService}, nil
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

func (e *PubServiceNatsImpl) Subscribe(topic string) error {
	_, err := e.nc.Subscribe(topic, func(msg *nats.Msg) {
		log.Printf("Message received: %v", msg)
	})

	if err != nil {
		return fmt.Errorf("event.EventSub.Subscribe(): %w", err)
	}

	e.logService.Debug(fmt.Sprintf("subscribed to event: %s", topic))

	return nil
}
