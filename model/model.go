package model

type Goods struct {
	Id    string  `json:"id" bson:"_id"`
	Name  string  `json:"name" bson:"name"`
	Lore  string  `json:"lore" bson:"lore"`
	Msg   string  `json:"msg" bson:"msg"`
	Num   int     `json:"num" bson:"num"`
	Price float64 `json:"price" bson:"price"`
	Auth  int     `json:"auth" bson:"auth"`
	Image string  `json:"image" bson:"string"`
}

type Employee struct {
	Id    string  `json:"id" bson:"_id"`
	Name  string  `json:"name" bson:"name"`
	Auth  int     `json:"auth" bson:"auth"`
	Money float64 `json:"money" bson:"money"`
}

type Record_H struct {
	Id   string `json:"id" bson:"_id"`
	Eid  string `json:"eid" bson:"eid"`
	Gid  string `json:"gid" bson:"gid"`
	Num  int    `json:"num" bson:"num"`
	Date string `json:"date" bson:"date"`
}

type Record_D struct {
	Id   string `json:"id" bson:"_id"`
	Eid  string `json:"eid" bson:"eid"`
	Gid  string `json:"gid" bson:"gid"`
	Num  int    `json:"num" bson:"num"`
	Date string `json:"date" bson:"date"`
}
