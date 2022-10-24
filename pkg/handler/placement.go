package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	placement "github.com/guff192/ad-placement-api"
	"github.com/sirupsen/logrus"
)

type placementRequest struct {
	Id      string            `json:"id"`
	Tiles   []placement.Tile  `json:"tiles"`
	Context placement.Context `json:"context"`
}

type adsResponse struct {
	Id  string                  `json:"id"`
	Imp []placement.ImpResponse `json:"imp"`
}

func (h Handler) getAds(c *gin.Context) {
	var input placementRequest

	logrus.Info("Got http request")

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

	ads, err := h.service.GetAdsForPlacements(input.Id, input.Tiles, input.Context)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "error getting ads: "+err.Error())
		return
	}

	resp := &adsResponse{
		Id:  input.Id,
		Imp: ads,
	}

	c.JSON(http.StatusOK, *resp)
}
