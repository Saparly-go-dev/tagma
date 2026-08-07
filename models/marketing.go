package models

import "time"

type Marketing struct {
	Id           int        `json:"id" db:"id"`
	TradePointId int        `json:"trade_point_id" db:"trade_point_id"`
	ProductId    int        `json:"product_id" db:"product_id"`
	Have         bool       `json:"have" db:"have"`
	Exposure     bool       `json:"exposure" db:"exposure"`
	Saled        bool       `json:"saled" db:"saled"`
	CreatedAt    time.Time  `json:"created_at"`
	TradePoint   TradePoint `json:"trade_point" gorm:"foreignKey:TradePointId;constraint:OnDelete:SET NULL;"`
}

type MarketingPost struct {
	ProductId int  `json:"product_id"`
	Have      bool `json:"have"`
	Exposure  bool `json:"exposure"`
	Saled     bool `json:"saled"`
}

type MarketingMainPage struct {
	ProductId         int     `json:"product_id"`
	Name              string  `json:"name"`
	Taste             string  `json:"taste"`
	ProductTypeId     int     `json:"product_type_id"`
	ProductType       string  `json:"product_type"`
	Volume            float64 `json:"volume"`
	Distribution      int     `json:"distribution"`
	Have              int     `json:"have"`
	Wystawlaymost     int     `json:"wystawlaymost"`
	Rasprodano        int     `json:"rasprodano"`
	Trade_point_Count int     `json:"trade_point_count"`
}
type ReadMarketingById struct {
	ProductId     int     `json:"product_id"`
	ProductTypeId int     `json:"product_type_id"`
	Name          string  `json:"name"`
	Taste         string  `json:"taste"`
	ProductType   string  `json:"product_type"`
	Volume        float64 `json:"volume"`
	Status        bool    `json:"status"`
}
