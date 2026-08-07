package models

import "time"

type Product struct {
	Id            int         `json:"id" db:"id"`
	Code          string      `json:"code" db:"code" binding:"required"`
	TasteRu       string      `json:"taste_ru" db:"taste_ru" binding:"required"`
	TasteTm       string      `json:"taste_tm" db:"taste_tm" binding:"required"`
	Price         float64     `json:"price" db:"price" binding:"required"`
	DPrice        float64     `json:"d_price" db:"d_price" binding:"required"`
	Status        bool        `json:"status" db:"status"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	BrandId       int         `json:"brand_id" db:"brand_id" binding:"required"`
	Name          Brand       `json:"name" gorm:"foreignKey:BrandId;constraint:OnDelete:SET NULL;"`
	ProductTypeId int         `json:"product_type_id" db:"product_type_id" binding:"required"`
	ProductType   ProductType `json:"product_type" gorm:"foreignKey:ProductTypeId;constraint:OnDelete:SET NULL;"`
}

type ReadProduct struct {
	Id            int     `json:"id" db:"id"`
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" db:"name"`
	TasteRu       string  `json:"taste_ru"  binding:"required"`
	TasteTm       string  `json:"taste_tm" binding:"required"`
	Price         float64 `json:"price"  binding:"required"`
	DPrice        float64 `json:"d_price" binding:"required"`
	Status        bool    `json:"status"`
	BlockPrice    float64 `json:"block_price" `
	ProductType   string  `json:"product_type"`
	ProductTypeId int     `json:"product_type_id"`
	Volume        float64 `json:"volume"`
	Updatedat     string  `json:"updated_at"`
	BrandId       int     `json:"brand_id"`
}

type CreateProduct struct {
	BrandId       int    `json:"brand_id" binding:"required"`
	TasteRu       string `json:"taste_ru"  binding:"required"`
	TasteTm       string `json:"taste_tm" binding:"required"`
	Price1        int    `json:"price1"`
	Price2        int    `json:"price2"`
	DPrice1       int    `json:"d_price1"`
	DPrice2       int    `json:"d_price2"`
	Status        bool   `json:"status" `
	ProductTypeId int    `json:"product_type_id" db:"product_type_id" binding:"required"`
}

type ProductPage struct {
	Items      []ReadProduct `json:"items"`
	PageNumber int           `json:"pageNumber"`
	PageSize   int           `json:"pageSize"`
	PageCount  int           `json:"pageCount"`
	Total      int           `json:"total"`
}

type ReadProductForTradePoint struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Taste         string  `json:"taste"`
	ProductTypeId int     `json:"product_type_id"`
	ProductType   string  `json:"product_type"`
	Volume        float64 `json:"volume"`
	Count         int     `json:"count"`
	Free          int     `json:"free"`
}

type GetGalyndyInArray struct {
	Id            int                        `json:"id"`
	Name          string                     `json:"name"`
	Taste         string                     `json:"taste"`
	ProductTypeId int                        `json:"product_type_id"`
	ProductType   string                     `json:"product_type"`
	Volume        float64                    `json:"volume"`
	List          []EkspeditorsGalyndyObject `json:"list"`
}

type EkspeditorsGalyndyObject struct {
	EkspeditorId int `json:"ekspeditor_id"`
	Free         int `json:"free"`
	Count        int `json:"count"`
}

type SaleHistoryTradePoint struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Taste         string  `json:"taste"`
	ProductTypeId int     `json:"product_type_id"`
	ProductType   string  `json:"product_type"`
	Volume        float64 `json:"volume"`
	First         int     `json:"first"`
	Second        int     `json:"second"`
	Third         int     `json:"third"`
	Fourth        int     `json:"fourth"`
}

type ProductDistinct struct {
	Id              int
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Taste           string  `json:"taste"`
	ProductTypeName string  `json:"product_type_name"`
	Volume          float64 `json:"volume"`
	GroupId         int     `json:"group_id"`
	Block_Count     int     `json:"block_count"`
	Price           float64 `json:"price"`
}
