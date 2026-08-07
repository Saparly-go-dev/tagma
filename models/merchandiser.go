package models

import "time"

type Merchandiser struct {
	Id         int          `json:"id" db:"id"`
	NameRu     string       `json:"name_ru" db:"name_ru" binding:"required"`
	NameTm     string       `json:"name_tm" db:"name_tm" binding:"required"`
	Number     string       `json:"number" db:"number" binding:"required"`
	Status     bool         `json:"status" db:"status" binding:"required"`
	UpdatedAt  time.Time    `json:"updated_at" db:"updated_at"`
	TradeAgent []TradeAgent `json:"trade_agent" gorm:"many2many:agents_merchandisers;constraint:OnDelete:Set NULL;"`
}

type AgentsMerchandisers struct {
	TradeAgentId   int `gorm:"column:trade_agent_id;foreignkey:trade_agent_id"`
	MerchandiserId int `gorm:"column:merchandiser_id;foreignkey:merchandiser_id"`
}

func (*Merchandiser) TableName() string {
	return "merchandisers"
}

type CreateMerchandiser struct {
	NameRu string `json:"name_ru"  binding:"required"`
	NameTm string `json:"name_tm"  binding:"required"`
	Number string `json:"number" binding:"required"`
	Status bool   `json:"status"`
}

type ReadMerchandiser struct {
	Id         int    `json:"id"`
	Code       string `json:"code"`
	TradeAgent string `json:"trade_agent"`
	NameRu     string `json:"name_ru"`
	NameTm     string `json:"name_tm"`
	Number     string `json:"number"`
	Status     bool   `json:"status"`
	UpdatedAt  string `json:"updated_at"`
}

type MerchandiserPage struct {
	Items      []ReadMerchandiser `json:"items"`
	PageNumber int                `json:"pageNumber"`
	PageSize   int                `json:"pageSize"`
	PageCount  int                `json:"pageCount"`
}
