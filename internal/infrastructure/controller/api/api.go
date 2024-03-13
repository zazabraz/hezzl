package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hezzl/internal/domain/service/goodsService"
	"hezzl/internal/infrastructure/controller/api/goodsHandler"
	"net/http"
)

type errorResponse struct {
	Code    int    `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
	Details string `json:"details" binding:"required"`
}

var (
	ErrorInvalidBodyFormat = errors.New("invalid body format")
	ErrorGoodNotFound      = errors.New("errors.good.notFound")
)

func ErrorString(c *gin.Context, httpCode int, internalCode int, message string, err error) {
	if err == nil {
		return
	}
	c.AbortWithStatusJSON(
		httpCode,
		errorResponse{
			Code:    internalCode,
			Message: message,
			Details: err.Error(),
		},
	)
	c.Error(err)
}

func MissingRequiredQuery(c *gin.Context, queryName string) error {
	err := fmt.Errorf("error %s query is required", queryName)
	ErrorString(c, http.StatusBadRequest, 1, err.Error(), err)
	return err
}

type API interface {
	Run(ctx context.Context) error
}

type api struct {
	log          logrus.Entry
	goodsHandler goodsHandler.Handler
	host         string
	port         string
}

func New(
	log logrus.Entry,
	goodsService goodsService.GoodsService,
	host string,
	port string,
) API {
	return &api{
		log:          log,
		goodsHandler: goodsHandler.NewHandler(*log.WithField("handler", "goods"), goodsService),
		host:         host,
		port:         port,
	}
}

func (h *api) Run(_ context.Context) error {
	r := gin.Default()

	apiGroup := r.Group("/good")

	h.goodsHandler.Register(*apiGroup)

	err := r.Run(fmt.Sprintf("%s:%s", h.host, h.port))
	if err != nil {
		return fmt.Errorf("api - Run: run router: %w", err)
	}
	return nil
}
