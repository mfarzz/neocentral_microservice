package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// NATSClient wraps a NATS connection with JetStream support.
type NATSClient struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

// NewNATSClient connects to NATS and initialises a JetStream context.
func NewNATSClient(url string) (*NATSClient, error) {
	nc, err := nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("⚠️  NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Println("🔄 NATS reconnected")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	log.Println("✅ Connected to NATS JetStream")
	return &NATSClient{conn: nc, js: js}, nil
}

// EnsureStream creates a JetStream stream if it does not already exist.
func (c *NATSClient) EnsureStream(name string, subjects []string) error {
	_, err := c.js.StreamInfo(name)
	if err != nil {
		_, err = c.js.AddStream(&nats.StreamConfig{
			Name:     name,
			Subjects: subjects,
			Storage:  nats.FileStorage,
		})
		if err != nil {
			return fmt.Errorf("failed to create stream %s: %w", name, err)
		}
		log.Printf("📡 Created NATS stream: %s (subjects: %v)", name, subjects)
	}
	return nil
}

// Publish sends a JSON-encoded message to a NATS subject.
func (c *NATSClient) Publish(subject string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	_, err = c.js.Publish(subject, data)
	return err
}

// Subscribe creates a durable push-based consumer.
func (c *NATSClient) Subscribe(subject, durableName string, handler func(msg *nats.Msg)) error {
	_, err := c.js.Subscribe(subject, handler,
		nats.Durable(durableName),
		nats.DeliverAll(),
		nats.AckExplicit(),
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}
	log.Printf("📩 Subscribed to %s (durable: %s)", subject, durableName)
	return nil
}

// Close closes the underlying NATS connection.
func (c *NATSClient) Close() {
	if c.conn != nil {
		c.conn.Drain()
	}
}
