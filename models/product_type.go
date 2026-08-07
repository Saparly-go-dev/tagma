package models

import "time"

type ProductType struct {
	Id        int       `json:"id" db:"id"`
	Code      string    `json:"code" db:"code" binding:"required"`
	NameRu    string    `json:"name_ru" db:"name_ru" binding:"required"`
	NameTm    string    `json:"name_tm" db:"name_tm" binding:"required"`
	Volume    float64   `json:"volume" db:"volume" binding:"required"`
	Count     int       `json:"count" db:"count" binding:"required"`
	Status    bool      `json:"status" db:"status"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (*ProductType) TableName() string {
	return "product_types"
}

type ReadProductType struct {
	Id        int     `json:"id"`
	Code      string  `json:"code"`
	NameRu    string  `json:"name_ru"`
	NameTm    string  `json:"name_tm"`
	Volume    float64 `json:"volume"`
	Count     int     `json:"count"`
	Status    bool    `json:"status"`
	UpdatedAt string  `json:"updated_at"`
}

type CreateProductType struct {
	Code      string    `json:"code" db:"code"`
	NameRu    string    `json:"name_ru" db:"name_ru" `
	NameTm    string    `json:"name_tm" db:"name_tm" `
	Volume    float64   `json:"volume" db:"volume" `
	Count     int       `json:"count" db:"count" `
	Status    bool      `json:"status" db:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductTypePage struct {
	Items      []ReadProductType `json:"items"`
	PageNumber int               `json:"pageNumber"`
	PageSize   int               `json:"pageSize"`
	PageCount  int               `json:"pageCount"`
}

type Volume struct {
	Id     int     `json:"id" db:"id" binding:"required"`
	Volume float64 `json:"volume" db:"volume" binding:"required"`
}
