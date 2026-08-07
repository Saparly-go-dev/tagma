package models

type UserEkspeditor struct {
	Id           int `json:"id" db:"id"`
	UserId       int `json:"userId" db:"user_id"`
	EkspeditorId int `json:"ekspeditorId" db:"ekspeditor_id"`
}
