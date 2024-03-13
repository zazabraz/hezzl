package goodsHandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"hezzl/internal/domain/storage"
	"hezzl/internal/infrastructure/controller/api"
	"net/http"
	"strconv"
)

type removeGoodResponse struct {
	ID         int  `json:"id" binding:"required"`
	CampaignID int  `json:"campaignID" binding:"required"`
	Removed    bool `json:"removed" binding:"required"` //always TRUE
}

func (h *handler) removeGood(c *gin.Context) {
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

	deleted, err := h.goodsService.DeleteByIDAndProjectID(c, goodIdInt, projectIdInt)
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
	res := &removeGoodResponse{
		ID:         deleted.ID,
		CampaignID: deleted.ProjectID,
		Removed:    true,
	}

	c.JSON(http.StatusOK, res)
}
