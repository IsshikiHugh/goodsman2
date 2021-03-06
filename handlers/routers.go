package handlers

import (
	"goodsman2/utils/feishu"

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
	//NULL              string = ""
	//INVALID_PARAMS    string = "INVALID_PARAMS"
	DB_ERROR          string = "DB_ERROR"
	CONDITION_NOT_MET string = "CONDITION_NOT_MET"
	//FEISHU_ERROR      string = "FEISHU_ERROR"
)

func InitRouter() *gin.Engine {
	//feishu event router
	fr := feishu.NewEventGroup()
	fr.Register(feishu.ReplyEvent, Receive_msg)

	r := gin.Default()
	r.GET("/ping", Ping)
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/user/records/hang", GetRecordsHangOfCertainEmployee)
		apiGroup.POST("/user/update/money", ChangeEmployeeMoney)
		apiGroup.POST("/user/update/auth", ChangeEmployeeAuth)

		apiGroup.POST("/user/login", EmployeeLogin)

		apiGroup.GET("/users/info", GetEmployeeInfo)
		apiGroup.GET("/users/search", GetCertainEmployeeList)
		apiGroup.POST("/users/return_goods", ReturnGoods)

		apiGroup.GET("/goods/msg", GetGoodsInfo)
		apiGroup.GET("/goods/all", GetGoodsBriefInfoList)
		apiGroup.GET("/goods/search", GetCertainGoodsBriefInfoList)
		apiGroup.POST("/goods/borrow", BorrowGoods)
		apiGroup.POST("/goods/new", AddNewGoods)
		apiGroup.POST("/goods/update/num", ChangeGoodsNumber)
		apiGroup.POST("/goods/update/price", ChangeGoodsPrice)

		apiGroup.POST("/records/close", CloseCertainRecordsH)
	}
	eventGroup := r.Group("/event")
	{
		eventGroup.POST("/received/msg", fr.EventListener)
	}
	return r
}
