package models

import "time"

type Promo struct {
	Id           int           `json:"id" db:"id"`
	NameRu       string        `json:"name_ru" db:"name_ru"`
	NameTm       string        `json:"name_tm" db:"name_tm"`
	Start        time.Time     `json:"start" db:"start"`
	End          time.Time     `json:"end" db:"end"`
	NeedCount    int           `json:"need_count" db:"need_count"`
	FreeCount    int           `json:"free_count" db:"free_count"`
	Status       bool          `json:"status" db:"status"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
	ChannelTypes []ChannelType `json:"channel_types" gorm:"many2many:promo_channel_types"`
}

type PromoChannelType struct {
	PromoId       int `json:"promo_id" db:"discount_id" gorm:"column:promo_id"`
	ChannelTypeId int `json:"channel_type_id" db:"channel_type_id" gorm:"column:channel_type_id;"`
}

type UpdatePromo struct {
	NameRu         string    `json:"name_ru"`
	NameTm         string    `json:"name_tm"`
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
	NeedCount      int       `json:"need_count"`
	FreeCount      int       `json:"free_count"`
	NeedIds        []int     `json:"need_ids"`
	FreeIds        []int     `json:"free_ids"`
	ChannelTypeIds []int     `json:"channel_type_ids"`
}

type ReadPromo struct {
	Id           int            `json:"id"`
	NameRu       string         `json:"name_ru"`
	NameTm       string         `json:"name_tm"`
	Start        time.Time      `json:"start"`
	End          time.Time      `json:"end"`
	NeedCount    int            `json:"need_count"`
	FreeCount    int            `json:"free_count"`
	Status       bool           `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	ChannelTypes []SelectObject `json:"channel_types"`
	NeedIds      []int          `json:"need_ids"`
	FreeIds      []int          `json:"free_ids"`
}

type PromoProduct struct {
	Id        int     `json:"id"`
	PromoId   int     `json:"promo_id"`
	ProductId int     `json:"product_id"`
	IsNeeded  bool    `json:"is_needed"`
	Promo     Promo   `json:"promo" gorm:"foreignKey:PromoId;constraint:OnDelete:SET NULL;"`
	Product   Product `json:"product" gorm:"foreignKey:ProductId;constraint:OnDelete:SET NULL"`
}

type PromoPage struct {
	Items      []ReadPromo `json:"items"`
	PageNumber int         `json:"pageNumber"`
	PageSize   int         `json:"pageSize"`
	PageCount  int         `json:"pageCount"`
	Total      int         `json:"total"`
}
