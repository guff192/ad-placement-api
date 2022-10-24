package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/guff192/ad-placement-api/pkg/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	placements := router.Group("/placements")
	{
		placements.POST("/request", h.getAds)
	}

	return router
}
