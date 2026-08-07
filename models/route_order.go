package models

type RouteOrder struct {
	Order_Number int        `json:"order_number"`
	TradePointId int        `json:"trade_point_id" db:"trade_point_id"`
	TradeAgentId int        `json:"trade_agent_id" db:"trade_agent_id"`
	DayId        int        `json:"day_id" db:"day_id"`
	TradeAgent   TradeAgent `json:"trade_agent" gorm:"foreignKey:TradeAgentId;constraint:OnDelete:SET NULL;"`
	TradePoint   TradePoint `json:"trade_point" gorm:"foreignKey:TradePointId;constraint:OnDelete:SET NULL;"`
	Day          Day        `json:"day" gorm:"foreignKey:DayId;constraint:OnDelete:SET NULL;"`
}
