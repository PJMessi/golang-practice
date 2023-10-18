package event

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/pkg/logger"
)

type PubServiceNatsImpl struct {
	nc         *nats.Conn
	loggerUtil logger.Util
}

func NewPubService(appConfig *config.AppConfig, loggerUtil logger.Util) (PubService, error) {
	url := appConfig.NATS_URL
	if url == "" {
		url = nats.DefaultURL
		loggerUtil.Debug(fmt.Sprintf("NATS url is not provided, so using default url of %s", url))
	}

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		return nil, fmt.Errorf("nats.NewPublisherService(): %w", err)
	}

	return &PubServiceNatsImpl{nc: nc, loggerUtil: loggerUtil}, nil
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
