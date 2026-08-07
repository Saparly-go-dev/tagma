package models

import "time"

type Manual struct {
	Id        int       `json:"id" db:"id"`
	TitleRu   string    `json:"title_ru" db:"title_ru"`
	ContentRu string    `json:"content_ru" db:"content_ru"`
	TitleTm   string    `json:"title_tm" db:"title_tm"`
	ContentTm string    `json:"content_tm" db:"content_tm"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ManualImage struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ManualId int    `json:"manual_id"`
	Manual   Manual `json:"manual" gorm:"foreignkey:ManualId;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
}

type CreateManual struct {
	TitleRu   string `json:"title_ru"`
	ContentRu string `json:"content_ru"`
	TitleTm   string `json:"title_tm"`
	ContentTm string `json:"content_tm"`
}

type ReadManual struct {
	Id        int           `json:"id"`
	TitleRu   string        `json:"title_ru"`
	ContentRu string        `json:"content_ru"`
	TitleTm   string        `json:"title_tm"`
	ContentTm string        `json:"content_tm"`
	CreatedAt string        `json:"created_at"`
	Images    []ManualImage `json:"images"`
}

type ReadManualForMobile struct {
	Id      int      `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}
