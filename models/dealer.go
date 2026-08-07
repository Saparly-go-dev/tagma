package models

import (
	"github.com/Saparly-go-dev/tagma"
	"time"
)

type Dealer struct {
	Id          int        `json:"id" db:"id"`
	NameRu      string     `json:"name_ru" db:"name_ru"`
	NameTm      string     `json:"name_tm" db:"name_tm"`
	AddressRu   string     `json:"address_ru" db:"address_ru"`
	AddressTm   string     `json:"address_tm" db:"address_tm"`
	Company     string     `json:"company" db:"company"`
	PhoneNumber string     `json:"phone_number" db:"phone_number"`
	Email       string     `json:"email" db:"email"`
	CityId      int        `json:"city_id" db:"city_id"`
	Account     float64    `json:"account" db:"account"`
	Debt        float64    `json:"debt" db:"debt"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	Status      bool       `json:"status" db:"status"`
	City        tagma.City `json:"city" gorm:"foreignKey:CityId;constraint:OnDelete:SET NULL;"`
}

type CreateDealer struct {
	NameRu      string `json:"name_ru"`
	NameTm      string `json:"name_tm"`
	AddressRu   string `json:"address_ru"`
	AddressTm   string `json:"address_tm"`
	Company     string `json:"company"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	CityId      int    `json:"city_id"`
	Status      bool   `json:"status"`
}

type ReadDealer struct {
	Id          int     `json:"id"`
	NameRu      string  `json:"name_ru"`
	NameTm      string  `json:"name_tm"`
	AddressRu   string  `json:"address_ru"`
	AddressTm   string  `json:"address_tm"`
	Company     string  `json:"company"`
	PhoneNumber string  `json:"phone_number"`
	Email       string  `json:"email"`
	Account     float64 `json:"account"`
	Debt        float64 `json:"debt"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	CityId      int     `json:"city_id"`
	City        string  `json:"city"`
	Status      bool    `json:"status"`
}

type DealerPage struct {
	Items      []ReadDealer `json:"items"`
	PageNumber int          `json:"pageNumber"`
	PageSize   int          `json:"pageSize"`
	PageCount  int          `json:"pageCount"`
	Total      int          `json:"total"`
}
