package nats

import "github.com/stretchr/testify/mock"

type PubServiceMock struct {
	mock.Mock
}

func (p *PubServiceMock) Close() {
	p.Called()
}

func (p *PubServiceMock) Publish(topic string, payload []byte) error {
	args := p.Called(topic, payload)
	return args.Error(0)
}

func (p *PubServiceMock) Subscribe(topic string) error {
	args := p.Called(topic)
	return args.Error(0)
}
