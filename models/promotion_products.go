package models

type PromotionProducts struct {
	Id          int     `json:"id"`
	ProductId   int     `json:"product_id"`
	PromotionId int     `json:"promotion_id"`
	Product     Product `json:"product" gorm:"foreignkey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
