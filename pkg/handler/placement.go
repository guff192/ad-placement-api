package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	placement "github.com/guff192/ad-placement-api"
)

type Context struct {
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

type PlacementRequest struct {
	Id      string `json:"id"`
	Tiles   []Tile `json:"tiles"`
	Context `json:"context"`
}

func (h Handler) getAds(c *gin.Context) {
	var input placement.PlacementRequest
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "WRONG_SCHEMA")
		return
	}
	if len(input.Tiles) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "EMPTY_TILES")
		return
	}
	if input.Ip == "" || input.UserAgent == "" || input.Id == "" {
		newErrorResponse(c, http.StatusBadRequest, "EMPTY_FIELD")
		return
	}

	h.service.GetImps()
}
