package handler

import "github.com/gin-gonic/gin"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	placements := router.Group("/placemnets")
	{
		placements.POST("request/", h.placementRequest)
	}

	return router
}
