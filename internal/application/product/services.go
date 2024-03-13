package product

import (
	"github.com/sirupsen/logrus"
	"hezzl/_pkg/natsPkg"
	"hezzl/internal/domain/event"
	"hezzl/internal/domain/service/goodsCacheService"
	"hezzl/internal/domain/service/goodsEventService"
	"hezzl/internal/domain/service/goodsService"
)

type Services struct {
	GoodsService      goodsService.GoodsService
	GoodsCacheService goodsCacheService.GoodsCacheService

	GoodsEventService goodsEventService.GoodsEventService //works with clickhouse
}

func NewServices(
	log logrus.Entry,
	storage *Storages,
	goodsChangeQueue natsPkg.Queue[event.GoodsChange],
) (*Services, error) {
	goodsCacheServiceServ := goodsCacheService.NewGoodsCacheService(
		*log.WithField("loc", "GoodsCacheService"),
		storage.GoodsCache,
	)
	return &Services{
		GoodsService: goodsService.NewGoodsService(
			*log.WithField("loc", "GoodsService"),
			storage.Goods,
			goodsCacheServiceServ,
			goodsChangeQueue.PubQueue,
		),
		GoodsEventService: goodsEventService.NewGoodsCacheService(
			*log.WithField("loc", "GoodsEventService"),
			storage.GoodsEvent,
		),
	}, nil
}
