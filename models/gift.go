package models

import "time"

type Gift struct {
	Id         int       `json:"id" db:"id"`
	NameRu     string    `json:"name_ru" db:"name_ru"`
	NameTm     string    `json:"name_tm" db:"name_tm"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	GiftTypeId int       `json:"gift_type_id" db:"gift_type_id"`
	Status     bool      `json:"status" db:"status"`
	GiftType   GiftType  `json:"gift_type" gorm:"foreignKey:GiftTypeId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreateGift struct {
	NameRu     string `json:"name_ru"`
	NameTm     string `json:"name_tm"`
	Status     bool   `json:"status"`
	GiftTypeId int    `json:"gift_type_id"`
}

type ReadGift struct {
	Id         int    `json:"id"`
	NameRu     string `json:"name_ru"`
	NameTm     string `json:"name_tm"`
	GiftTypeId int    `json:"gift_type_id"`
	GiftType   string `json:"gift_type"`
	Status     bool   `json:"status"`
	UpdatedAt  string `json:"updated_at"`
}

type GiftPage struct {
	Items      []ReadGift `json:"items"`
	PageNumber int        `json:"pageNumber"`
	PageSize   int        `json:"pageSize"`
	PageCount  int        `json:"pageCount"`
}
