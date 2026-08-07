package models

import "time"

type EkspData struct {
	Id           int        `json:"id" db:"id"`
	EkspeditorId int        `json:"ekspeditor_id"`
	ProductId    int        `json:"product_id"`
	Count        int        `json:"count"`
	CreatedAt    time.Time  `json:"created_at"`
	Ekspeditor   Ekspeditor `json:"ekspeditor" gorm:"foreignkey:EkspeditorId;constraint:OnDelete:SET NULL;"`
	Product      Product    `json:"product" gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
