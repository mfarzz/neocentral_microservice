package broker

import (
	"log"

	"github.com/nats-io/nats.go"
)

type NATSBroker struct {
	nc *nats.Conn
}

func NewNATSBroker(url string) (*NATSBroker, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to NATS at %s", url)
	return &NATSBroker{nc: nc}, nil
}

func (b *NATSBroker) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return b.nc.Subscribe(subject, handler)
}

func (b *NATSBroker) Close() {
	b.nc.Close()
}
