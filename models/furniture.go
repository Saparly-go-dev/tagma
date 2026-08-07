package models

import "time"

type Furniture struct {
	Id        int       `json:"id" db:"id"`
	NameRu    string    `json:"name_ru" db:"name_ru"`
	NameTm    string    `json:"name_tm" db:"name_tm"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type TradePointsFurniture struct {
	Id           int        `json:"id" db:"id"`
	TradePointId int        `json:"trade_point_id" db:"trade_point_id"`
	FurnitureId  int        `json:"furniture_id" db:"furniture_id"`
	Count        int        `json:"count" db:"count"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	TradePoint   TradePoint `json:"trade_point" gorm:"foregnKey:TradePointId;constraint:OnDelete:CASCADE;"`
	Furniture    Furniture  `json:"furniture" gorm:"foreignKey:FurnitureId;constraint:OnDelete:CASCADE;"`
}

type TradePointFurnitureRead struct {
	Id           int    `json:"id"`
	TradePointId int    `json:"trade_point_id"`
	FurnitureId  int    `json:"furniture_id"`
	Name         string `json:"name"`
	Count        int    `json:"count"`
}

type CreateFurniture struct {
	NameRu string `json:"name_ru" db:"name_ru"`
	NameTm string `json:"name_tm" db:"name_tm"`
}

type ReadFurniture struct {
	Id        int    `json:"id" db:"id"`
	NameRu    string `json:"name_ru" db:"name_ru"`
	NameTm    string `json:"name_tm" db:"name_tm"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type ReadFurnitureForMobile struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type FurniturePage struct {
	Items      []ReadFurniture `json:"items"`
	PageNumber int             `json:"pageNumber"`
	PageSize   int             `json:"pageSize"`
	PageCount  int             `json:"pageCount"`
	Total      int             `json:"total"`
}
