package models

import (
	"time"

	"github.com/Saparly-go-dev/tagma"
)

type TradePoint struct {
	Id              int                  `json:"id" db:"id" binding:"required"`
	Code            string               `json:"code" db:"code" binding:"required"`
	NameRu          string               `json:"name_ru" db:"name_ru" binding:"required"`
	NameTm          string               `json:"name_tm" db:"name_tm" binding:"required"`
	LocationRu      string               `json:"location_ru" db:"location_ru" binding:"required"`
	LocationTm      string               `json:"location_tm" db:"location_tm" binding:"required"`
	OrientirRu      string               `json:"orientir_ru" db:"orientir_ru" binding:"required"`
	OrientirTm      string               `json:"orientir_tm" db:"orientir_tm" binding:"required"`
	Status          bool                 `json:"status" db:"status" binding:"required"`
	CityId          int                  `json:"city_id" db:"city_id" binding:"required"`
	DistrictId      int                  `json:"district_id" db:"district_id" binding:"required"`
	TradeAgentId    int                  `json:"trade_agent_id" db:"trade_agent_id" binding:"required"`
	TradeChannelId  int                  `json:"trade_channel_id" db:"trade_channel_id" binding:"required"`
	TradeCategoryId int                  `json:"trade_category_id" db:"trade_category_id" binding:"required"`
	DayId           int                  `json:"day_id" db:"day_id" binding:"required"`
	EkspeditorId    int                  `json:"ekspeditor_id" db:"ekspeditor_id" binding:"required"`
	UpdatedAt       time.Time            `json:"updated_at" db:"updated_at"`
	City            tagma.City           `json:"city" gorm:"foreignKey:CityId;constraint:OnDelete:SET NULL;"`
	District        District             `json:"district" gorm:"foreignKey:DistrictId;constraint:OnDelete:SET NULL;"`
	TradeAgent      TradeAgent           `json:"trade_agent" gorm:"foreignKey:TradeAgentId;constraint:OnDelete:SET NULL;"`
	TradeChannel    TradeChannel         `json:"trade_channel" gorm:"foreignKey:TradeChannelId;constraint:OnDelete:SET NULL;"`
	TradeCategory   tagma.Trade_category `json:"trade_category" gorm:"foreignKey:TradeCategoryId;constraint:OnDelete:SET NULL;"`
	Day             Day                  `json:"day" gorm:"foreignKey:DayId;constraint:OnDelete:SET NULL;"`
	Ekspeditor      Ekspeditor           `json:"ekspeditor" gorm:"foreignKey:EkspeditorId;constraint:OnDelete:SET NULL;"`
}

type CreateTradePoint struct {
	Code           string `json:"code"`
	NameRu         string `json:"name_ru"`
	NameTm         string `json:"name_tm"`
	LocationRu     string `json:"location_ru"`
	LocationTm     string `json:"location_tm"`
	OrientirRu     string `json:"orientir_ru"`
	OrientirTm     string `json:"orientir_tm"`
	Status         bool   `json:"status"`
	CityId         int    `json:"city_id"`
	DistrictId     int    `json:"district_id"`
	TradeAgentId   int    `json:"trade_agent_id"`
	TradeChannelId int    `json:"trade_channel_id"`
	DayId          int    `json:"day_id"`
	EkspeditorId   int    `json:"ekspeditor_id"`
}

type ReadTradePoint struct {
	Id              int                      `json:"id"`
	Code            string                   `json:"code"`
	Name            string                   `json:"name"`
	Location        string                   `json:"location"`
	Orientir        string                   `json:"orientir"`
	Status          bool                     `json:"status"`
	CityId          int                      `json:"city_id"`
	DistrictId      int                      `json:"district_id"`
	TradeAgentId    int                      `json:"trade_agent_id"`
	TradeChannelId  int                      `json:"trade_channel_id"`
	UpdatedAt       string                   `json:"updated_at"`
	TradeCategoryId int                      `json:"trade_category_id"`
	DayId           int                      `json:"day_id"`
	EkspeditorId    int                      `json:"ekspeditor_id"`
	FurnitureList   []ReadFurnitureForMobile `json:"furniture_list"`
	TradCategory    string                   `json:"trade_category"`
	City            string                   `json:"city"`
	District        string                   `json:"district"`
	TradeChannel    string                   `json:"trade_channel"`
	Debt            float64                  `json:"debt"`
	TradeAgent      string                   `json:"trade_agent"`
	Day             string                   `json:"day"`
	Ekspeditor      string                   `json:"ekspeditor"`
}

type TradePointMobile struct {
	Id       int     `json:"id"`
	City     string  `json:"city"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Location string  `json:"location"`
	Orientir string  `json:"orientir"`
	Status   bool    `json:"status"`
	Contact  string  `json:"contact"`
	Number   string  `json:"number"`
	Note     string  `json:"note" default:""`
	Debt     float64 `json:"debt" default:"0"`
}

type TradePointPage struct {
	Items         []ReadTradePoint         `json:"items"`
	PageNumber    int                      `json:"pageNumber"`
	PageSize      int                      `json:"pageSize"`
	PageCount     int                      `json:"pageCount"`
	Total         int                      `json:"total"`
	FurnitureList []ReadFurnitureForMobile `json:"furniture_list"`
}

type GeneralTradePoint struct {
	Id            int         `json:"id"`
	Code          string      `json:"code"`
	CityId        int         `json:"city_id"`
	DistrictId    int         `json:"district_id"`
	DayId         int         `json:"day_id"`
	EkspeditorId  int         `json:"ekspeditor_id"`
	City          string      `json:"city"`
	District      string      `json:"district"`
	Location      string      `json:"location"`
	Orientir      string      `json:"orientir"`
	Name          string      `json:"name"`
	TradeChannel  string      `json:"trade_channel"`
	TradeCategory string      `json:"trade_category"`
	TradeAgent    string      `json:"trade_agent"`
	Day           string      `json:"day"`
	Ekspeditor    string      `json:"ekspeditor"`
	Debt          float64     `json:"debt"`
	Cash          float64     `json:"cash"`
	Terminal      float64     `json:"terminal"`
	Transfer      float64     `json:"transfer"`
	Images        []ReadImage `json:"images"`
}

type TradePointDebtResponseForMobile struct {
	TradePointId int     `json:"trade_point_id"`
	Currency     float64 `json:"currency"`
	CreatedAt    string  `json:"created_at"`
	PaymentType  string  `json:"payment_type"`
}
