package models

type DProduct struct {
	Id          int     `json:"id"`
	ProductId   int     `json:"product_id"`
	DiscountId  int     `json:"discount_id"`
	Count       int     `json:"count"`
	Coefficient int     `json:"coefficient"`
	Type        bool    `json:"type"`
	Product     Product `json:"product" gorm:"foreignkey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreateDProduct struct {
	ProductId   int `json:"product_id"`
	Count       int `json:"count"`
	Coefficient int `json:"coefficient"`
}

type ReadDProduct struct {
	ProductId   int `json:"product_id"`
	Count       int `json:"count"`
	Coefficient int `json:"coefficient"`
}
