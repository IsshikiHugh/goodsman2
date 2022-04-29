package model

type BorrowGoodsReq struct {
	Eid string `json:"employee_id" binding:"required"`
	Gid string `json:"goods_id" binding:"required"`
	Num int    `json:"goods_num" binding:"required"`
}

type ReturnGoodsReq struct {
	Eid string `json:"employee_id" binding:"required"`
	Gid string `json:"goods_id" binding:"required"`
	Num int    `json:"goods_num" binding:"required"`
	Msg string `json:"msg" binding:"required"`
}
