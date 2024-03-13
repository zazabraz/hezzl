package product

import (
	"github.com/sirupsen/logrus"
	"hezzl/_pkg/natsPkg"
	"hezzl/internal/domain/event"
	"hezzl/internal/infrastructure/controller/api"
	"hezzl/internal/infrastructure/controller/monitor"
)

type Controllers struct {
	API     api.API
	Monitor monitor.Monitor
}

func NewControllers(
	logger logrus.Entry,

	services Services,
	APIHost, APIPort string,
	goodsChangeQueue natsPkg.Queue[event.GoodsChange],
) *Controllers {
	restapi := api.New(*logger.WithField("loc", "api"), services.GoodsService, APIHost, APIPort)
	return &Controllers{
		API: restapi,
		Monitor: monitor.NewMonitor(
			*logger.WithField("loc", "monitor"),
			goodsChangeQueue.SubQueue,
			services.GoodsEventService,
		),
	}
}
