package natsPkg

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
)

type PubQueue[MessageType Message] struct {
	nc  *nats.Conn
	key string
}

func NewPubQueue[MessageType Message](
	nc *nats.Conn,
	key string,
) *PubQueue[MessageType] {
	return &PubQueue[MessageType]{nc: nc, key: key}
}

func (p *PubQueue[MessageType]) Pub(ctx context.Context, message MessageType) error {
	rawMes, err := message.Bytes()
	if err != nil {
		return fmt.Errorf("get bytes from message: %w", err)
	}

	err = p.nc.Publish(p.key, rawMes)
	if err != nil {
		return fmt.Errorf("write message to partition: %w", err)
	}
	return nil
}
