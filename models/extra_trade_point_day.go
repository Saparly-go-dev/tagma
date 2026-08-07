package models

type ExtraTradePointDay struct {
	Id           int        `json:"id" db:"id"`
	TradePointId int        `json:"trade_point_id" db:"trade_point_id"`
	DayId        int        `json:"day_id" db:"day_id"`
	Day          Day        `json:"day" gorm:"foreignKey:DayId;constraint:OnDelete:CASCADE;"`
	TradePoint   TradePoint `json:"trade_point" gorm:"foreignKey:TradePointId;constraint:OnDelete:CASCADE;"`
}
