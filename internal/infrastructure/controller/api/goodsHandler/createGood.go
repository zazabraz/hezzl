package goodsHandler

import (
	"github.com/gin-gonic/gin"
	"hezzl/internal/infrastructure/controller/api"
	"net/http"
	"strconv"
	"time"
)

type createGoodRequest struct {
	Name string `json:"name" binding:"required"`
}

type createGoodResponse struct {
	ID          int       `json:"id" binding:"required"`
	ProjectID   int       `json:"projectId" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Priority    int       `json:"priority" binding:"required"`
	Removed     bool      `json:"removed" binding:"required"`
	CreatedAt   time.Time `json:"createdAt" binding:"required"`
}

func (h *handler) createGood(c *gin.Context) {
	projectIdS, ok := c.GetQuery("projectId")
	if !ok {
		h.log.Errorln(api.MissingRequiredQuery(c, "projectId"))
		return
	}
	projectIdInt, err := strconv.Atoi(projectIdS)
	if err != nil {
		h.log.Errorln(err)
		api.ErrorString(c, http.StatusBadRequest, 0, err.Error(), err)
		return
	}

	body := &createGoodRequest{}
	err = c.ShouldBindQuery(body)
	if err != nil {
		h.log.Errorln(err)
		api.ErrorString(c, http.StatusBadRequest, 0, err.Error(), err)
		return
	}

	good, err := h.goodsService.Create(c, projectIdInt, body.Name)
	if err != nil {
		h.log.Errorln(err)
		api.ErrorString(c, http.StatusBadRequest, 0, err.Error(), err)
		return
	}

	res := &createGoodResponse{
		good.ID,
		good.ProjectID,
		good.Name,
		good.Description,
		good.Priority,
		good.Removed,
		good.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}
