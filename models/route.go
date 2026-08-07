package models

import "time"

type Route struct {
	Id           int        `json:"id" db:"id"`
	TradePointId int        `json:"trade_point_id" db:"trade_point_id"`
	TradeAgentId int        `json:"trade_agent_id" db:"trade_agent_id"`
	EkspeditorId int        `json:"ekspeditor_id" db:"ekspeditor_id"`
	DayId        int        `json:"day_id" db:"day_id"`
	Status       bool       `json:"status" db:"status"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	TradeAgent   TradeAgent `json:"trade_agent" gorm:"foreignKey:TradeAgentId;constraint:OnDelete:SET NULL;"`
	Ekspeditor   Ekspeditor `json:"ekspeditor" gorm:"foreignkey:EkspeditorId;constraint:OnDelete:SET NULL;"`
	TradePoint   TradePoint `json:"trade_point" gorm:"foreignKey:TradePointId;constraint:OnDelete:SET NULL;"`
	Day          Day        `json:"day" gorm:"foreignKey:DayId;constraint:OnDelete:SET NULL;"`
}

type CreateRoute struct {
	TradePointId int  `json:"trade_point_id"`
	TradeAgentId int  `json:"trade_agent_id"`
	EkspeditorId int  `json:"ekspeditor_id"`
	DayId        int  `json:"day_id"`
	Status       bool `json:"status"`
}

type ReadRoute struct {
	Id           int    `json:"id" db:"id"`
	Code         string `json:"code" db:"code"`
	TradePointId int    `json:"trade_point_id"`
	TradeAgentId int    `json:"trade_agent_id"`
	EkspeditorId int    `json:"ekspeditor_id"`
	DayId        int    `json:"day_id"`
	CityId       int    `json:"city_id"`
	DistrictId   int    `json:"district_id"`
	Status       bool   `json:"status"`
	City         string `json:"city"`
	District     string `json:"district"`
	TradePoint   string `json:"trade_point"`
	TradeAgent   string `json:"trade_agent"`
	Ekspeditor   string `json:"ekspeditor"`
	Day          string `json:"day"`
	Updatedat    string `json:"updated_at"`
	Location     string `json:"location"`
	Orientir     string `json:"orientir"`
}

type RoutePage struct {
	Items      []ReadRoute `json:"items"`
	PageNumber int         `json:"pageNumber"`
	PageSize   int         `json:"pageSize"`
	PageCount  int         `json:"pageCount"`
}

type PointSelect struct {
	Id           int    `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	City         string `json:"city"`
	Location     string `json:"location"`
	TradeAgentId int    `json:"trade_agent_id"`
	TradeAgent   string `json:"trade_agent"`
}
