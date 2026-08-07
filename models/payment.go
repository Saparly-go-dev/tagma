package models

import "time"

type TradePointForPayment struct {
	Id       int     `json:"id"`
	City     string  `json:"city"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Location string  `json:"location"`
	Orientir string  `json:"orientir"`
	Status   bool    `json:"status"`
	OrderId  int     `json:"order_id"`
	Sum      float64 `json:"sum"`
	Paid     float64 `json:"paid"`
	NotPaid  float64 `json:"not_paid"`
	Contact  string  `josn:"contact"`
	Number   string  `json:"number"`
}

type OrderInformationForPayment struct {
	Id      int     `json:"id"`
	Sum     float64 `json:"sum"`
	Paid    float64 `json:"paid"`
	NotPaid float64 `json:"not_paid"`
}

type Payment struct {
	Id            int         `json:"id" db:"id"`
	OrderId       int         `json:"order_id" db:"order_id"`
	Currency      float64     `json:"currency" db:"currency"`
	PaymentTypeId int         `json:"payment_type_id" db:"payment_type_id"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	TradeAgentId  int         `json:"trade_agent_id" db:"trade_agent_id"`
	Status        bool        `json:"status" db:"status"`
	Count         int         `json:"count" db:"count"`
	IsAgent       bool        `json:"is_agent" db:"is_agent"`
	TradeAgent    TradeAgent  `json:"trade_agent" gorm:"foreignKey:TradeAgentId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Order         Order       `json:"order" gorm:"foreignKey:OrderId;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
	PaymentType   PaymentType `json:"payment_type" gorm:"foreignKey:PaymentTypeId;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
}

type PaymentType struct {
	Id     int    `json:"id" db:"id"`
	NameRu string `json:"name_ru" db:"name_ru"`
	NameTm string `json:"name_tm" db:"name_tm"`
}

type CreatePayment struct {
	OrderId       int     `json:"order_id"`
	Currency      float64 `json:"currency"`
	PaymentTypeId int     `json:"payment_type_id"`
	Is_Agent      bool    `json:"Is_Agent"`
}

type ReadPayment struct {
	Id            int     `json:"id" db:"id"`
	Code          string  `json:"code"`
	OrderId       int     `json:"order_id"`
	Currency      float64 `json:"currency"`
	PaymentTypeId int     `json:"payment_type_id"`
	Status        bool    `json:"status"`
	CreatedAt     string  `json:"created_at"`
	PaymentType   string  `json:"payment_type"`
	TradeAgent    string  `json:"trade_agent"`
	TradePoint    string  `json:"trade_point"`
	IsAgent       bool    `json:"is_agent"`
	Count         int     `json:"count" db:"count"`
}

type ObjectWithIdNameAndCurrency struct {
	OrdersValue       float64 `json:"order_value"`
	AgentCash         float64 `json:"agent_cash"`
	EkspeditorCash    float64 `json:"ekspeditor_cash"`
	SelfDebt          float64 `json:"self_debt"`
	CashSum           float64 `json:"cash_sum"`
	AgentTerminal     float64 `json:"agent_terminal"`
	AgentTransfer     float64 `json:"agent_transfer"`
	EkspeditorLed     float64 `json:"ekspeditor_led"`
	TradePointAllDebt float64 `json:"trade_point_all_debt"`
}
