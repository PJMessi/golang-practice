package nats

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

type ServiceImpl struct {
	natsCon    *nats.Conn
	logService logger.Service
	jetStream  jetstream.JetStream
	subjects   []string
}

func NewPubService(appConfig *config.AppConfig, logService logger.Service) (Service, error) {
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

	return &ServiceImpl{natsCon: nc, logService: logService, jetStream: js, subjects: []string{
		appConfig.NATS_EVENT_USER_REGISTRATION,
	}}, nil
}

func (s *ServiceImpl) Close() {
	s.natsCon.Close()
	s.logService.Debug("NATS connection closed")
}

func (s *ServiceImpl) Publish(topic string, payload []byte) error {
	err := s.natsCon.Publish(topic, payload)
	if err != nil {
		return fmt.Errorf("nats.ServiceImpl.Publish(): %w", err)
	}
	return nil
}

func (s *ServiceImpl) Subscribe(topic string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := s.createStream(ctx, topic)
	if err != nil {
		return err
	}

	con, err := s.createConsumer(ctx, stream)
	if err != nil {
		return err
	}

	// we don't need to stop the consuming if we run s.Close() func
	_, err = con.Consume(func(msg jetstream.Msg) {
		log.Printf("Topic: %s, Message: %s", msg.Subject(), msg.Data())
		err := msg.Ack()
		if err != nil {
			s.logService.Error(fmt.Sprintf("nats.ServiceImpl.Subscribe(): error fetching messages : %s", err))
		}
	})

	if err != nil {
		return fmt.Errorf("nats.ServiceImpl.Subscribe(): error setting consumer handler: %w", err)
	}

	s.logService.Debug("NATS consumer started")
	return nil
}

func (s *ServiceImpl) createConsumer(ctx context.Context, stream jetstream.Stream) (jetstream.Consumer, error) {
	con, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:      "database_usage_service",
		AckPolicy:    jetstream.AckExplicitPolicy,
		ReplayPolicy: jetstream.ReplayInstantPolicy,
	})

	if err != nil {
		return nil, fmt.Errorf("nats.ServiceImpl.Subscribe(): error creating consumer : %w", err)
	}

	s.logService.Debug("NATS consumer created")

	return con, nil
}

func (s *ServiceImpl) createStream(ctx context.Context, name string) (jetstream.Stream, error) {
	stream, err := s.jetStream.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     name,
		Subjects: s.subjects,
	})

	if err != nil {
		return nil, fmt.Errorf("nats.ServiceImpl.createStream(): %w", err)
	}

	s.logService.Debug(fmt.Sprintf("NATS stream created: '%s'", name))

	return stream, nil
}
