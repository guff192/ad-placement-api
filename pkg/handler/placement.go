package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	placement "github.com/guff192/ad-placement-api"
)

type placementRequest struct {
	Id      string            `json:"id"`
	Tiles   []placement.Tile  `json:"tiles"`
	Context placement.Context `json:"context"`
}

func (h Handler) getAds(c *gin.Context) {
	var input placementRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "WRONG_SCHEMA")
		return
	}
	if len(input.Tiles) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "EMPTY_TILES")
		return
	}
	if input.Id == "" || input.Context.UserAgent == "" || input.Id == "" {
		newErrorResponse(c, http.StatusBadRequest, "EMPTY_FIELD")
		return
	}

	h.service.GetAllImps(input.Id, input.Tiles, input.Context)
}
