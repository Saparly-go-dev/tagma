package models

import "time"

type PointNote struct {
	Id           int        `json:"id" db:"id"`
	Description  string     `json:"description" db:"description"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	TradePointId int        `json:"trade_point_id" db:"trade_point_id"`
	TradePoint   TradePoint `json:"trade_point" gorm:"foreignKey:TradePointId;constraint:OnDelete:SET NULL;"`
}

type CreatePointNote struct {
	Description  string `json:"description"`
	TradePointId int    `json:"trade_point_id"`
}

type ReadPointNote struct {
	Id             int    `json:"id"`
	CityId         int    `json:"city_id"`
	DisctrictId    int    `json:"disctrict_id"`
	TradePointCode string `json:"trade_point_code"`
	Description    string `json:"description"`
	TradePointId   int    `json:"trade_point_id"`
	City           string `json:"city"`
	Disctrict      string `json:"disctrict"`
	Location       string `json:"location"`
	CreatedAt      string `json:"created_at"`
	TradePoint     string `json:"trade_point"`
}

type PointNotePage struct {
	Items      []ReadPointNote `json:"items"`
	PageNumber int             `json:"pageNumber"`
	PageSize   int             `json:"pageSize"`
	PageCount  int             `json:"pageCount"`
	Total      int             `json:"total"`
}
