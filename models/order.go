package models

import (
	"time"
)

type CreateOrder struct {
	Id    int `json:"id"`
	Count int `json:"count"`
}

type Order struct {
	Id            int         `json:"id" db:"id"`
	TradePointId  int         `json:"trade_point_id" db:"trade_point_id"`
	TradeAgentId  int         `json:"trade_agent_id" db:"trade_agent_id"`
	Sum           float64     `json:"sum" db:"sum"`
	SaleSum       float64     `json:"sale_sum" db:"sale_sum"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	List          []OrderList `json:"list" gorm:"OrderId"`
	Status        bool        `json:"status" db:"status"`
	IsClosed      bool        `json:"is_closed" db:"is_closed"`
	IsCredit      bool        `json:"is_credit" db:"is_credit"`
	PaymentTypeId int         `json:"payment_type_id" db:"payment_type_id"`
	TradePoint    TradePoint  `json:"trade_point" gorm:"foreignKey:TradePointId;constraint:OnDelete:SET NULL;"`
	PaymentType   PaymentType `json:"payment_type" gorm:"foreignKey:PaymentTypeId"`
}

type ReadOrder struct {
	Id            int     `json:"id"`
	TradePointId  int     `json:"trade_point_id"`
	TradeAgentId  int     `json:"trade_agent_id"`
	TradePoint    string  `json:"trade_point"`
	Location      string  `json:"location"`
	Orientir      string  `json:"orientir"`
	TradeCategory string  `json:"trade_category"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	CreatedAt     string  `json:"created_at"`
	TradeAgent    string  `json:"trade_agent"`
	Status        bool    `json:"status"`
	IsClosed      int     `json:"is_closed"`
	Sum           float64 `json:"sum"`
	Paid          float64 `json:"paid"`
	NotPaid       float64 `json:"not_paid"`
}

type OrderPage struct {
	Items      []ReadOrder `json:"items"`
	PageNumber int         `json:"pageNumber"`
	PageSize   int         `json:"pageSize"`
	PageCount  int         `json:"pageCount"`
	Total      int         `json:"total"`
	Sum        float64     `json:"sum"`
	Paid       float64     `json:"paid"`
	NotPaid    float64     `json:"not_paid"`
}
