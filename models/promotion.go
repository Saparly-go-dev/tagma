package models

import "time"

type Promotion struct {
	Id        int       `json:"id" db:"id"`
	NameRu    string    `json:"name_ru" db:"name_ru"`
	NameTm    string    `json:"name_tm" db:"name_tm"`
	Start     time.Time `json:"start" db:"start"`
	End       time.Time `json:"end" db:"end"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Count     int       `json:"count" db:"count"`
	Status    bool      `json:"status" db:"status"`
}

type CreatePromotion struct {
	NameRu string    `json:"name_ru"`
	NameTm string    `json:"name_tm"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Count  int       `json:"count"`
	List   []int     `json:"list"`
}

type ReadPromotion struct {
	Id        int    `json:"id"`
	NameTm    string `json:"name_tm"`
	NameRu    string `json:"name_ru"`
	Start     string `json:"start"`
	End       string `json:"end"`
	List      []int  `json:"list"`
	CreatedAt string `json:"created_at"`
	Count     int    `json:"count"`
	Status    bool   `json:"status"`
}

type PromotionPage struct {
	Items      []ReadPromotion `json:"items"`
	PageNumber int             `json:"pageNumber"`
	PageSize   int             `json:"pageSize"`
	PageCount  int             `json:"pageCount"`
	Total      int             `json:"total"`
}
