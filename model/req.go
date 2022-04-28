package model

type BorrowGoodsReq struct {
	Eid string `json:"employee_id"`
	Gid string `json:"goods_id"`
	Num int    `json:"goods_num"`
}
