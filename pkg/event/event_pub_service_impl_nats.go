package event

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/pkg/logger"
)

type PubServiceNatsImpl struct {
	nc         *nats.Conn
	logService logger.Service
	js         jetstream.JetStream
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

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("nats.NewPublisherService(): %w", err)
	}

	return &PubServiceNatsImpl{nc: nc, logService: logService, js: js}, nil
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := e.createStream(ctx, topic)
	if err != nil {
		return err
	}

	con, err := e.createConsumer(ctx, stream)
	if err != nil {
		return err
	}

	_, err = con.Consume(func(msg jetstream.Msg) {
		log.Printf("Topic: %s, Message: %s", msg.Subject(), msg.Data())
		err := msg.Ack()
		if err != nil {
			e.logService.Error(fmt.Sprintf("event.PubServiceNatsImpl.Subscribe(): error fetching messages : %s", err))
		}
	})

	if err != nil {
		return fmt.Errorf("event.PubServiceNatsImpl.Subscribe(): error setting consumer handler: %w", err)
	}

	e.logService.Debug("consumer handler set")
	return nil
}

func (e *PubServiceNatsImpl) createConsumer(ctx context.Context, stream jetstream.Stream) (jetstream.Consumer, error) {
	con, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:      "database_usage_service",
		AckPolicy:    jetstream.AckExplicitPolicy,
		ReplayPolicy: jetstream.ReplayInstantPolicy,
	})

	if err != nil {
		return nil, fmt.Errorf("event.PubServiceNatsImpl.Subscribe(): error creating consumer : %w", err)
	}

	e.logService.Debug("consumer created")

	return con, nil
}

func (e *PubServiceNatsImpl) createStream(ctx context.Context, name string) (jetstream.Stream, error) {
	stream, err := e.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     name,
		Subjects: []string{fmt.Sprintf("%s.*", name)},
	})

	if err != nil {
		return nil, fmt.Errorf("event.PubServiceNatsImpl.createStream(): %w", err)
	}

	e.logService.Debug(fmt.Sprintf("stream created with name: %s", name))

	return stream, nil
}
