package model

type Goods struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Lore  string  `json:"lore"`
	Msg   string  `json:"msg"`
	Num   int     `json:"num"`
	Price float64 `json:"price"`
	Auth  int     `json:"auth"`
	Image string  `json:"image"`
}

type Employee struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Auth  int     `json:"auth"`
	Money float64 `json:"money"`
}

type Record_H struct {
	Id   string `json:"id"`
	Eid  string `json:"eid"`
	Gid  string `json:"gid"`
	Num  int    `json:"num"`
	Date string `json:"date"`
}

type Record_D struct {
	Id   string `json:"id"`
	Eid  string `json:"eid"`
	Gid  string `json:"gid"`
	Num  int    `json:"num"`
	Date string `json:"date"`
}
