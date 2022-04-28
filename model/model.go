package model

type Good struct {
	Id    string
	Name  string
	Lore  string
	Msg   string
	Num   int
	Price float64
	Auth  int
	Image string
}

type Employee struct {
	Id    string
	Name  string
	Auth  int
	Money float64
}

type Record_H struct {
	Id   string
	Eid  string
	Gid  string
	Num  int
	Date string //TODO:
}

type Record_D struct {
	*Record_H
	Origin string
}
