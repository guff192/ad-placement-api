package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Info("Aborted http request with code " + strconv.Itoa(statusCode))
	logrus.Errorf("%s", message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
