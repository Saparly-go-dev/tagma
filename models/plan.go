package models

import "time"

type Plan struct {
	Id           int        `json:"id" db:"id"`
	ProductId    int        `json:"product_id" db:"product_id"`
	TradeAgentId int        `json:"trade_agent_id" db:"trade_agent_id"`
	Count        int        `json:"count" db:"count"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	Status       bool       `json:"status" db:"status"`
	Product      Product    `json:"product" gorm:"foreignKey:ProductId;constraint:OnDelete:SET NULL;"`
	TradeAgent   TradeAgent `json:"trade_agent" gorm:"foreignKey:TradeAgentId;constraint:OnDelete:SET NULL;"`
}

type CreatePlan struct {
	ProductId int         `json:"product_id"`
	List      map[int]int `json:"list"`
}

type UpdatedPlan struct {
	ProductId int         `json:"product_id"`
	List      map[int]int `json:"list"`
}

//type ReadPlan struct {
//	Id           int    `json:"id"`
//	ProductId    int    `json:"product_id"`
//	TradeAgentId int    `json:"trade_agent_id"`
//	Count        int    `json:"count"`
//	Status       bool   `json:"status"`
//	Product      string `json:"product"`
//	TradeAgent   string `json:"trade_agent"`
//	CreatedAt    string `json:"created_at"`
//}

type ReadPlan struct {
	ProductId int         `json:"product_id"`
	List      map[int]int `json:"list"`
}
type PlanStatisticsMain struct {
	ProductId      int        `json:"product_id"`
	Product        string     `json:"product"`
	Plan           int        `json:"plan"`
	Forecast       int        `json:"forecast"`
	ForecastVsPlan float64    `json:"forecast_vs_plan"`
	Volume         float64    `json:"volume"`
	Name           string     `json:"name"`
	ProductType    string     `json:"product_type"`
	GroupId        int        `json:"group_id"`
	Count          int        `json:"count"`
	Data           []PlanData `json:"data"`
}

type PlanStatisticsForAgent struct {
	ProductId int    `json:"product_id"`
	Name      string `json:"name"`
	Taste     int    `json:"taste"`
	Volume    int    `json:"volume"`
	Count     int    `json:"count"`
	Selled    int    `json:"selled"`
}

//type PlanStatisticsForAgents struct {
//	ProductId    int        `json:"product_id"`
//	TradeAgentId int        `json:"trade_agent_id"`
//	Product      string     `json:"product"`
//	TradeAgent   string     `json:"trade_agent"`
//	Forecast     float64    `json:"forecast"`
//	Data         []PlanData `json:"data"`
//}

type PlanData struct {
	Day   int `json:"day"`
	Count int `json:"count"`
}
