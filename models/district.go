package models

import (
	"github.com/Saparly-go-dev/tagma"
	"time"
)

type District struct {
	Id        int        `json:"id" db:"id"`
	NameRu    string     `json:"name_ru" db:"name_ru"`
	NameTm    string     `json:"name_tm" db:"name_tm"`
	CityId    int        `json:"city_id" db:"city_id"`
	Status    bool       `json:"status" db:"status"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	City      tagma.City `json:"city" gorm:"foreignKey:CityId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreateDistrict struct {
	NameRu string `json:"name_ru"`
	NameTm string `json:"name_tm"`
	CityId int    `json:"city_id"`
}

type ReadDistrict struct {
	Id        int    `json:"id"`
	NameRu    string `json:"name_ru"`
	NameTm    string `json:"name_tm"`
	CityId    int    `json:"city_id"`
	Status    bool   `json:"status"`
	UpdatedAt string `json:"updated_at"`
	City      string `json:"city"`
}

type DistrictPage struct {
	Items      []ReadDistrict `json:"items"`
	PageNumber int            `json:"pageNumber"`
	PageSize   int            `json:"pageSize"`
	PageCount  int            `json:"pageCount"`
}
