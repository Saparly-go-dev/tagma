package tagma

import "time"

type City struct {
	Id        int       `json:"id" db:"id"`
	NameRu    string    `json:"name_ru" db:"name_ru"`
	NameTm    string    `json:"name_tm" db:"name_tm"`
	Code      string    `json:"code" db:"code"`
	Status    bool      `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateCity struct {
	NameRu string `json:"name_ru"`
	NameTm string `json:"name_tm"`
	Code   string `json:"code"`
	Status bool   `json:"status"`
}

type CityRead struct {
	Id        int    `json:"id"`
	NameRu    string `json:"name_ru" `
	NameTm    string `json:"name_tm" `
	Code      string `json:"code" `
	Status    bool   `json:"status" `
	UpdatedAt string `json:"updated_at"`
}

type CityPage struct {
	Items      []CityRead `json:"items"`
	PageNumber int        `json:"pageNumber"`
	PageSize   int        `json:"pageSize"`
	PageCount  int        `json:"pageCount"`
}
