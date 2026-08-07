package models

import "time"

type DealerProduct struct {
	Id            int         `json:"id" db:"id"`
	Code          string      `json:"code" db:"code" binding:"required"`
	TasteRu       string      `json:"taste_ru" db:"taste_ru" binding:"required"`
	TasteTm       string      `json:"taste_tm" db:"taste_tm" binding:"required"`
	Price         float64     `json:"price" db:"price" binding:"required"`
	Status        bool        `json:"status" db:"status"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	BrandId       int         `json:"brand_id" db:"brand_id" binding:"required"`
	Name          Brand       `json:"name" gorm:"foreignKey:BrandId;constraint:OnDelete:SET NULL;"`
	ProductTypeId int         `json:"product_type_id" db:"product_type_id" binding:"required"`
	ProductType   ProductType `json:"product_type" gorm:"foreignKey:ProductTypeId;constraint:OnDelete:SET NULL;"`
}

type ReadDealerProduct struct {
	Id            int     `json:"id" db:"id"`
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" db:"name"`
	TasteRu       string  `json:"taste_ru"  binding:"required"`
	TasteTm       string  `json:"taste_tm" binding:"required"`
	Price         float64 `json:"price"  binding:"required"`
	Status        bool    `json:"status"`
	BlockPrice    float64 `json:"block_price" `
	ProductType   string  `json:"product_type"`
	ProductTypeId int     `json:"product_type_id" `
	Volume        float64 `json:"volume"`
	Updatedat     string  `json:"updated_at"`
	BrandId       int     `json:"brand_id"`
}

type CreateDealerProduct struct {
	BrandId       int    `json:"brand_id" binding:"required"`
	TasteRu       string `json:"taste_ru"  binding:"required"`
	TasteTm       string `json:"taste_tm" binding:"required"`
	Price1        int    `json:"price1"`
	Price2        int    `json:"price2"`
	Status        bool   `json:"status" `
	ProductTypeId int    `json:"product_type_id"`
}

type DealerProductPage struct {
	Items      []ReadDealerProduct `json:"items"`
	PageNumber int                 `json:"pageNumber"`
	PageSize   int                 `json:"pageSize"`
	PageCount  int                 `json:"pageCount"`
	Total      int                 `json:"total"`
}
