package models

import "time"

type Discount struct {
	Id           int           `json:"id" db:"id"`
	NameRu       string        `json:"name_ru" db:"name_ru"`
	NameTm       string        `json:"name_tm" db:"name_tm"`
	Start        time.Time     `json:"start" db:"start"`
	End          time.Time     `json:"end" db:"end"`
	SortOrder    int           `json:"sort_order"`
	Status       bool          `json:"status" db:"status"`
	IsProcent    bool          `json:"is_procent" db:"is_procent"`
	IsDirect     bool          `json:"is_direct" db:"is_direct"`
	IsIncrease   bool          `json:"is_increase" db:"is_increase"`
	IsGroup      bool          `json:"is_group"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	ChannelTypes []ChannelType `json:"channel_types" gorm:"many2many:discount_channel_types"`
	TradePoints  []TradePoint  `json:"trade_points" gorm:"many2many:discount_trade_points"`
	TradeAgents  []TradeAgent  `json:"trade_agents" gorm:"many2many:discount_trade_agents"`
}

type DiscountChannelType struct {
	DiscountId    int `json:"discount_id" db:"discount_id" gorm:"column:discount_id;"`
	ChannelTypeId int `json:"channel_type_id" db:"channel_type_id" gorm:"column:channel_type_id;"`
}

type DiscountTradePoint struct {
	DiscountId   int `json:"discount_id" db:"discount_id" gorm:"column:discount_id;"`
	TradePointId int `json:"trade_point_id" db:"trade_point_id" gorm:"column:trade_point_id;"`
}

type DiscountTradeAgent struct {
	DiscountId   int `json:"discount_id" db:"discount_id" gorm:"column:discount_id;"`
	TradeAgentId int `json:"trade_agent_id" db:"trade_agent_id" gorm:"trade_agent_id;"`
}

type CreateDiscount struct {
	NameRu         string           `json:"name_ru"`
	NameTm         string           `json:"name_tm"`
	Start          time.Time        `json:"start"`
	End            time.Time        `json:"end"`
	SortOrder      int              `json:"sort_order"`
	IsProcent      bool             `json:"is_procent"`
	IsDirect       bool             `json:"is_direct"`
	IsIncrease     bool             `json:"is_increase"`
	IsGroup        bool             `json:"is_group"`
	Status         bool             `json:"status"`
	First          []CreateDProduct `json:"first"`
	Second         []CreateDProduct `json:"second"`
	ChannelTypeIds []int            `json:"channel_type_ids"`
	TradePointIds  []int            `json:"trade_point_ids"`
	TradeAgents    []int            `json:"trade_agents"`
}

type ReadDiscount struct {
	Id           int            `json:"id"`
	NameTm       string         `json:"name_tm"`
	NameRu       string         `json:"name_ru"`
	Start        string         `json:"start"`
	End          string         `json:"end"`
	SortOrder    int            `json:"sort_order"`
	IsProscent   bool           `json:"is_procent"`
	IsDirect     bool           `json:"is_direct"`
	IsIncrease   bool           `json:"is_increase"`
	IsGroup      bool           `json:"is_group"`
	Status       bool           `json:"status"`
	First        []ReadDProduct `json:"first"`
	Second       []ReadDProduct `json:"second"`
	ChannelTypes []SelectObject `json:"channel_types"`
	TradePoints  []SelectObject `json:"trade_points"`
	TradeAgent   []SelectObject `json:"trade_agents"`
}

type DiscountPage struct {
	Items      []ReadDiscount `json:"items"`
	PageNumber int            `json:"pageNumber"`
	PageSize   int            `json:"pageSize"`
	PageCount  int            `json:"pageCount"`
	Total      int            `json:"total"`
}

type UpdateDiscountInfo struct {
	NameRu         string    `json:"name_ru"`
	NameTm         string    `json:"name_tm"`
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
	SortOrder      int       `json:"sort_order"`
	ChannelTypeIds []int     `json:"channel_type_ids"`
	TradePointIds  []int     `json:"trade_point_ids"`
	TradeAgents    []int     `json:"trade_agents"`
}
