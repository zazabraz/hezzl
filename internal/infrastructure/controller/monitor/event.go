package monitor

import (
	"context"
	"github.com/sirupsen/logrus"
	"hezzl/_pkg/natsPkg"
	"hezzl/internal/domain/event"
	"hezzl/internal/domain/service/goodsEventService"
)

type Monitor interface {
	Run(ctx context.Context) error
}

type monitor struct {
	log               logrus.Entry
	goodsChangeSub    natsPkg.SubQueue[event.GoodsChange]
	goodsEventService goodsEventService.GoodsEventService
}

func NewMonitor(
	log logrus.Entry,
	goodsChangeSub natsPkg.SubQueue[event.GoodsChange],
	goodsEventService goodsEventService.GoodsEventService,
) Monitor {
	return &monitor{
		log:               log,
		goodsChangeSub:    goodsChangeSub,
		goodsEventService: goodsEventService,
	}
}

func (e monitor) Run(ctx context.Context) error {
	e.goodsChangeSub.Sub(e.handleEvent(ctx))

	<-ctx.Done()
	return ctx.Err()
}

func (e monitor) handleEvent(ctx context.Context) func(goodsChange event.GoodsChange) {
	return func(goodsChange event.GoodsChange) {
		select {
		case <-ctx.Done():
			return
		default:
		}
		err := e.goodsEventService.SaveEvent(ctx, goodsChange)
		if err != nil {
			e.log.Errorf("eventChecker - handleEvent - SaveEvent: %s", err)
			return
		}
	}
}
