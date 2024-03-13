package goodsHandler

import (
	"github.com/gin-gonic/gin"
	"hezzl/internal/infrastructure/controller/api"
	"net/http"
	"strconv"
	"time"
)

type goodInList struct {
	ID          int       `json:"id" binding:"required"`
	ProjectID   int       `json:"projectId" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Priority    int       `json:"priority" binding:"required"`
	Removed     bool      `json:"removed" binding:"required"`
	CreatedAt   time.Time `json:"createdAt" binding:"required"`
}

type goodsListMeta struct {
	Total   int `json:"total" binding:"required"`
	Removed int `json:"removed" binding:"required"`
	Limit   int `json:"limit" binding:"required"`
	Offset  int `json:"offset" binding:"required"`
}

type goodsListResponse struct {
	Meta      goodsListMeta `json:"meta" binding:"required"`
	GoodsList []goodInList  `json:"goods" binding:"required"`
}

func (h *handler) goodsList(c *gin.Context) {
	var limitInt int
	limitS, limitOK := c.GetQuery("limit")
	if !limitOK {
		limitInt = 10
	} else {
		limitConv, err := strconv.Atoi(limitS)
		if err != nil {
			h.log.Errorln(err)
			api.ErrorString(c, http.StatusBadRequest, 0, err.Error(), err)
			return
		}
		limitInt = limitConv
	}

	var offsetInt int
	offsetS, offsetOK := c.GetQuery("offset")
	if !offsetOK {
		offsetInt = 1
	} else {
		offsetConv, err := strconv.Atoi(offsetS)
		if err != nil {
			h.log.Errorln(err)
			api.ErrorString(c, http.StatusBadRequest, 0, err.Error(), err)
			return
		}
		offsetInt = offsetConv
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

	list, err := h.goodsService.GetList(c, limitInt, offsetInt, projectIdInt)
	if err != nil {
		h.log.Errorln(err)
		api.ErrorString(c, http.StatusBadRequest, 0, err.Error(), err)
		return
	}

	res := &goodsListResponse{
		GoodsList: make([]goodInList, 0),
		Meta: goodsListMeta{
			Offset: offsetInt,
			Limit:  limitInt,
		},
	}
	for _, good := range list {
		res.Meta.Total++
		if good.Removed {
			res.Meta.Removed++
		}
		res.GoodsList = append(res.GoodsList, goodInList{
			good.ID,
			good.ProjectID,
			good.Name,
			good.Description,
			good.Priority,
			good.Removed,
			good.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}
