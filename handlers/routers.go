package handlers

import (
	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", Ping)
	apiGroup := r.Group("/api")
	{

	}
	eventGroup := r.Group("/event")
	{
		eventGroup.POST("/received/msg", handler.ReplyCheck)
	}
	return r
}
