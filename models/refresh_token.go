package models

type RefreshToken struct {
	Id           int    `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpireDate   string `json:"expire_date"`
	UserId       int    `json:"user_id"`
}
