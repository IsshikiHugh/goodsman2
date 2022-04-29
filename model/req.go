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

type AddNewGoodsReqGoods struct {
	Name  string  `json:"name" binding:"required"`
	Lore  string  `json:"lore" binding:"required"`
	Num   int     `json:"num" binding:"required"`
	Price float64 `json:"price" binding:"required"`
	Auth  int     `json:"auth" binding:"required"`
	Image string  `json:"image" binding:"required"`
}

type AddNewGoodsReq struct {
	Eid   string              `json:"employee_id" binding:"required"`
	Goods AddNewGoodsReqGoods `json:"goods" binding:"required"`
}

type ChangeGoodsNumberReq struct {
	Eid    string `json:"employee_id" binding:"required"`
	Gid    string `json:"goods_id" binding:"required"`
	DelNum int    `json:"del_num" binding:"required"`
}

type ChangeGoodsPriceReq struct {
	Eid      string  `json:"employee_id" binding:"required"`
	Gid      string  `json:"goods_id" binding:"required"`
	NewPrice float64 `json:"new_price" binding:"required"`
}
