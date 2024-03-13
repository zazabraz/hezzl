package goodsEventService

import (
	"context"
	"github.com/sirupsen/logrus"
	"hezzl/internal/domain/event"
	"hezzl/internal/domain/storage"
)

type GoodsEventService interface {
	SaveEvent(ctx context.Context, change event.GoodsChange) error
}

type goodsEventService struct {
	log               logrus.Entry
	goodsEventStorage storage.GoodsEventStorage
}

func NewGoodsCacheService(
	log logrus.Entry,
	goodsEventStorage storage.GoodsEventStorage,
) GoodsEventService {
	return &goodsEventService{
		log:               log,
		goodsEventStorage: goodsEventStorage,
	}
}

func (g goodsEventService) SaveEvent(ctx context.Context, change event.GoodsChange) error {
	return g.goodsEventStorage.SaveEvent(ctx, change)
}
