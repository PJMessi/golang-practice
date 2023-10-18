package event

import "github.com/stretchr/testify/mock"

type PubServiceMock struct {
	mock.Mock
}

func (p *PubServiceMock) Close() error {
	args := p.Called()
	return args.Error(0)
}

func (p *PubServiceMock) Publish(topic string, payload []byte) error {
	args := p.Called(topic, payload)
	return args.Error(0)
}
