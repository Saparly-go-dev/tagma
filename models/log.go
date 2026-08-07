package models

import "time"

type Log struct {
	Id       int       `json:"id" db:"id"`
	InfoTm   string    `json:"info_tm" db:"info_tm"`
	InfoRu   string    `json:"info_ru" db:"info_ru"`
	LogType  string    `json:"log_type" db:"log_type"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
	User     string    `json:"user" db:"user"`
}

type ReadLog struct {
	Id       int    `json:"id"`
	Info     string `json:"info"`
	LogType  string `json:"log_type"`
	CreateAt string `json:"create_at"`
	User     string `json:"user"`
}

type LogPage struct {
	Items      []ReadLog `json:"items"`
	PageNumber int       `json:"pageNumber"`
	PageSize   int       `json:"pageSize"`
	PageCount  int       `json:"pageCount"`
}
