package models

type Debt struct {
	Id           int        `json:"id" db:"id"`
	TradeAgentId int        `json:"trade_agent_id" db:"trade_agent_id"`
	Total        float64    `json:"total" db:"total"`
	TradeAgent   TradeAgent `json:"trade_agent" gorm:"foreignKey:TradeAgentId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
