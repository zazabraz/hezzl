package natsPkg

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type Callback[MessageType Message] func(MessageType)
type CancelSubscription func()

// Error is a function that encapsulates error delivery
type Error func(error)

type SubQueue[MessageType Message] struct {
	nc                   *nats.Conn
	errTube              Error
	key                  string
	fetchTimeMillisecond int
}

func NewSubQueue[MessageType Message](nc *nats.Conn, errTube Error, key string, fetchTimeMillisecond int) *SubQueue[MessageType] {
	return &SubQueue[MessageType]{nc: nc, errTube: errTube, key: key, fetchTimeMillisecond: fetchTimeMillisecond}
}

func (s *SubQueue[MessageType]) Sub(callback Callback[MessageType]) {
	sub, err := s.nc.SubscribeSync(s.key)
	if err != nil {
		s.errTube(fmt.Errorf("subscribe: %w", err))
		return
	}

	go func() {
		for {
			msgNats, err := sub.NextMsg(time.Duration(s.fetchTimeMillisecond) * time.Millisecond)
			if err != nil {
				s.errTube(fmt.Errorf("fetch message: %w", err))
				return
			} else {
				var mes MessageType
				mesI, err := mes.Unmarshal(msgNats.Data)
				if err != nil {
					s.errTube(fmt.Errorf("scan message's value: %w", err))
					return
				}
				mes = mesI.(MessageType)

				callback(mes)
			}
		}
	}()

	return
}
