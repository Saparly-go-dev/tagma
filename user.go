package tagma

import "time"

type User struct {
	Id        int       `json:"-" db:"id"`
	NameTm    string    `json:"name_tm" db:"name_tm" binding:"required"`
	NameRu    string    `json:"name_ru" db:"name_ru" binding:"required"`
	Username  string    `json:"username" db:"username" binding:"required"`
	Password  string    `json:"password" db:"password" binding:"required"`
	Type      string    `json:"type" db:"type"`
	Status    bool      `json:"status" db:"status"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

type UserTradeAgent struct {
	Id           int `json:"id" db:"id"`
	UserId       int `json:"user_id" db:"user_id"`
	TradeAgentId int `json:"trade_agent_id" db:"trade_agent_id"`
}

type CreateUser struct {
	Id              int    `json:"id"`
	Type            string `json:"type"`
	NameRu          string `json:"name_ru"`
	NameTm          string `json:"name_tm"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Username        string `json:"username"`
	EkspeditorId    int    `json:"ekspeditor_id"`
	AgentId         int    `json:"agent_id"`
	MerchandiserId  int    `json:"merchandiser_id"`
}

type Change_Password struct {
	UserId          int    `json:"user_id"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ReadUser struct {
	Id             int    `json:"id"`
	Type           string `json:"type"`
	NameTm         string `json:"name_tm" `
	NameRu         string `json:"name_ru"`
	Username       string `json:"username"`
	CreatedAt      string `json:"created_at"`
	TypeName       string `json:"type_name"`
	Status         bool   `json:"status"`
	AgentId        int    `json:"agent_id"`
	EkspeditorId   int    `json:"ekspeditor_id"`
	MerchandiserId int    `json:"merchandiser_id"`
}

type UserPage struct {
	Items      []ReadUser `json:"items"`
	PageNumber int        `json:"pageNumber"`
	PageSize   int        `json:"pageSize"`
	PageCount  int        `json:"pageCount"`
}
