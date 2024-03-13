package goodsHandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"hezzl/internal/domain/storage"
	"hezzl/internal/infrastructure/controller/api"
	"net/http"
	"strconv"
	"time"
)

type updateGoodRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"omitnil"`
}

type updateGoodResponse struct {
	ID          int       `json:"id" binding:"required"`
	ProjectID   int       `json:"projectId" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Priority    int       `json:"priority" binding:"required"`
	Removed     bool      `json:"removed" binding:"required"`
	CreatedAt   time.Time `json:"createdAt" binding:"required"`
}

func (h *handler) updateGood(c *gin.Context) {
	goodIdS, ok := c.GetQuery("id")
	if !ok {
		h.log.Errorln(api.MissingRequiredQuery(c, "id"))
		return
	}
	goodIdInt, err := strconv.Atoi(goodIdS)
	if err != nil {
		h.log.Errorln(err)
		api.ErrorString(c, http.StatusBadRequest, 0, err.Error(), err)
		return
	}

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

	body := &updateGoodRequest{}
	err = c.ShouldBindQuery(body)
	if err != nil {
		h.log.Errorln(err)
		api.ErrorString(c, http.StatusBadRequest, 0, err.Error(), err)
		return
	}

	if len(body.Name) <= 0 {

	}

	updated, err := h.goodsService.Update(c, goodIdInt, projectIdInt, body.Name, body.Description)
	if err != nil {
		h.log.Errorln(err)
		api.ErrorString(c, http.StatusBadRequest, 0, err.Error(), err)
		return
	}
	if errors.Is(err, storage.PgErrNoEffect) {
		h.log.Errorln(err)
		err := api.ErrorGoodNotFound
		api.ErrorString(c, http.StatusNotFound, 3, err.Error(), err)
		return
	}

	res := &createGoodResponse{
		updated.ID,
		updated.ProjectID,
		updated.Name,
		updated.Description,
		updated.Priority,
		updated.Removed,
		updated.CreatedAt,
	}

	c.JSON(http.StatusOK, res)

}
