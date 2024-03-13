package goodsHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hezzl/internal/domain/service/goodsService"
)

type Handler interface {
	Register(gin.RouterGroup)
}

type handler struct {
	log logrus.Entry

	goodsService goodsService.GoodsService
}

func NewHandler(log logrus.Entry, goodsService goodsService.GoodsService) Handler {
	return &handler{log: log, goodsService: goodsService}
}

func (h *handler) Register(group gin.RouterGroup) {
	group.POST("/create", h.createGood)
	group.PATCH("/update", h.updateGood)
	group.DELETE("/remove", h.removeGood)
	group.GET("/list", h.goodsList)
	//group.PATCH("/reprioritize", h.reprioritizeGood)
}
