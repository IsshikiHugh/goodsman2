package handlers

import (
	"github.com/gin-gonic/gin"
)

// auth level
var (
	AuthEmplo int = 1
	AuthAdmin int = 2
	AuthSuper int = 3
)

// err type
var (
	NULL              string = "NULL"
	INVALID_PARAMS    string = "INVALID_PARAMS"
	DB_ERROR          string = "DB_ERROR"
	CONDITION_NOT_MET string = "CONDITION_NOT_MET"
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
		apiGroup.POST("/goods/new", AddNewGoods)
		apiGroup.POST("/goods/update/num", ChangeGoodsNumber)
		apiGroup.POST("/goods/update/price", ChangeGoodsPrice)

	}
	// eventGroup := r.Group("/event")
	// {
	// 	eventGroup.POST("/received/msg", handlers.ReplyCheck)
	// }
	return r
}
