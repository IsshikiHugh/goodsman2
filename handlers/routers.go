package handlers

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", Ping)
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/user/records/hang", GetRecordsHangOfCertainEmployee)

		apiGroup.GET("/users/info", GetEmployeeInfo)
		apiGroup.GET("/users/search", GetCertainEmployeeList)
		apiGroup.POST("/users/return_goods", ReturnGoods)

		apiGroup.GET("/goods/msg", GetGoodsInfo)
		apiGroup.GET("/goods/search", GetCertainGoodsBriefInfoList)
		apiGroup.POST("/goods/borrow", BorrowGoods)

	}
	// eventGroup := r.Group("/event")
	// {
	// 	eventGroup.POST("/received/msg", handlers.ReplyCheck)
	// }
	return r
}
