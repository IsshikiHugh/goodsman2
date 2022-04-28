package handlers

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", Ping)
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/users/info", GetEmployeeInfo)
		apiGroup.GET("/goods/msg", GetGoodsInfo)
		apiGroup.POST("/goods/borrow", BorrowGoods)

	}
	// eventGroup := r.Group("/event")
	// {
	// 	eventGroup.POST("/received/msg", handler.ReplyCheck)
	// }
	return r
}
