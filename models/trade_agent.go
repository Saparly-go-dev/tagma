package models

import (
	"time"
)

type TradeAgent struct {
	Id           int            `json:"id" db:"id"`
	Code         string         `json:"code" db:"code" binding:"required"`
	NameRu       string         `json:"name_ru" db:"name_ru" binding:"required"`
	NameTm       string         `json:"name_tm" db:"name_tm" binding:"required"`
	Number       string         `json:"number" db:"number" binding:"required"`
	Status       bool           `json:"status" db:"status" binding:"required"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
	Driver       []Driver       `json:"drivers" gorm:"many2many:agents_drivers;constraint:OnDelete:Set NULL;"`
	Ekspeditor   []Ekspeditor   `json:"ekspeditors" gorm:"many2many:agents_ekspeditors;constraint:OnDelete:Set NULL;"`
	Merchandiser []Merchandiser `json:"merchandisers" gorm:"many2many:agents_merchandisers;constraint:OnDelete:Set NULL;"`
}

type CreateTradeAgent struct {
	Code            string `json:"code"`
	NameRu          string `json:"name_ru" `
	NameTm          string `json:"name_tm"`
	Number          string `json:"number"`
	Status          bool   `json:"status" `
	DriverIds       []int  `json:"driver_ids" `
	EkspeditorIds   []int  `json:"ekspeditor_ids" `
	MerchandiserIds []int  `json:"merchandiser_ids"`
}

type ReadTradeAgent struct {
	Id            int            `json:"id"`
	Code          string         `json:"code"`
	NameRu        string         `json:"name_ru" `
	NameTm        string         `json:"name_tm" `
	Number        string         `json:"number"`
	Status        bool           `json:"status"`
	Drivers       []SelectObject `json:"drivers"`
	Ekspeditors   []SelectObject `json:"ekspeditors"`
	Merchandisers []SelectObject `json:"merchandisers"`
	Debt          float64        `json:"debt"`
	UpdatedAt     string         `json:"updated_at"`
}

type SelectObject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type AgentDebt struct {
	Name     string          `json:"name"`
	Number   string          `json:"number"`
	Sum      float64         `json:"sum"`
	Self     float64         `json:"self"`
	Cash     float64         `json:"cash"`
	Terminal float64         `json:"terminal"`
	Transfer float64         `json:"transfer"`
	List     []AgentDebtItem `json:"list"`
}

type AgentDebtItem struct {
	Name      string  `json:"name"`
	City      string  `json:"city"`
	District  string  `json:"district"`
	Sum       float64 `json:"sum"`
	Paid      float64 `json:"paid"`
	NotPaid   float64 `json:"not_paid"`
	Status    bool    `json:"status"`
	CreatedAt string  `json:"created_at"`
}

type TradeAgentPage struct {
	Items      []ReadTradeAgent `json:"items"`
	PageNumber int              `json:"pageNumber"`
	PageSize   int              `json:"pageSize"`
	PageCount  int              `json:"pageCount"`
	Total      int              `json:"total"`
}

type PointsMerchendising_List struct {
	List      []PointsMerchendising_Item `json:"list"`
	Furniture []FurnitureInformation     `json:"furniture"`
}

type PointsMerchendising_Item struct {
	Id          int                    `json:"id"`
	City        string                 `json:"city"`
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Location    string                 `json:"location"`
	Orientir    string                 `json:"orientir"`
	Furniture   []FurnitureInformation `json:"furniture"`
	Planogramma float64                `json:"planogramma"`
	Ishave      bool                   `json:"is_have"`
}

type FurnitureInformation struct {
	Name                  string  `json:"name"`
	Count                 int     `json:"count"`
	Percentage            float64 `json:"percentage"`
	Item_percentage_value float64 `json:"item_percentage_value"`
}
