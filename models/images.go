package models

type Image struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	TradePointId int    `json:"trade_point_id"`
}

type ReadImage struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}
