package models

type Contact struct {
	Id           int        `json:"id" db:"id"`
	NameRu       string     `json:"name_ru" db:"name_ru" binding:"required"`
	NameTm       string     `json:"name_tm" db:"name_tm" binding:"required"`
	Number       string     `json:"number" db:"number" binding:"required"`
	KindId       int        `json:"kind_id" db:"kind_id"`
	PostId       int        `json:"post_id" db:"post_id"`
	TradePointId int        `json:"trade_point_id" db:"trade_point_id"`
	TradePoint   TradePoint `json:"trade_point" gorm:"foreignkey:TradePointId;constraint:OnDelete:SET NULL;"`
	Kind         Kind       `json:"kind" gorm:"foreignKey:KindId;constraint:OnDelete:SET NULL;"`
	Post         Post       `json:"post" gorm:"foreignKey:PostId;constraint:OnDelete:SET NULL;"`
	Uprs         []Upr      `json:"uprs" gorm:"many2many:contact_uprs;constraint:OnUpdate:SET NULL;"`
}

type ContactUprs struct {
	ContactId int `gorm:"column:contact_id;foreignkey:contact_id"`
	UprId     int `gorm:"column:upr_id;foreignkey:upr_id"`
}

func (*ContactUprs) TableName() string {
	return "contact_uprs"
}

type CreateContact struct {
	NameRu       string `json:"name_ru"`
	NameTm       string `json:"name_tm"`
	Number       string `json:"number"`
	TradePointId int    `json:"trade_point_id"`
	KindId       int    `json:"kind_id"`
	PostId       int    `json:"post_id"`
	UprsIds      []int  `json:"uprs_ids"`
}

type ReadContact struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Number string   `json:"number"`
	KindId int      `json:"kind_id"`
	PostId int      `json:"post_id"`
	Kind   string   `json:"kind"`
	Post   string   `json:"post"`
	Uprs   []string `json:"uprs"`
}

type ContactPage struct {
	Items      []ReadContact `json:"items"`
	PageNumber int           `json:"pageNumber"`
	PageSize   int           `json:"pageSize"`
	PageCount  int           `json:"pageCount"`
}
