package models

type UserMerchandiser struct {
	Id             int `json:"id" db:"id"`
	UserId         int `json:"user_id" db:"user_id"`
	MerchandiserId int `json:"merchandiser_id" db:"merchandiser_id"`
}
