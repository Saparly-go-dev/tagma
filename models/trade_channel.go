package models

import "time"

type ChannelType struct {
	Id     int    `json:"id" db:"id"`
	NameRu string `json:"name_ru" db:"name_ru"`
	NameTm string `json:"name_tm" db:"name_tm"`
}

type ChannelStructure struct {
	Id     int    `json:"id" db:"id"`
	NameRu string `json:"name_ru" db:"name_ru"`
	NameTm string `json:"name_tm" db:"name_tm"`
}

type ChannelSize struct {
	Id     int    `json:"id" db:"id"`
	NameRu string `json:"name_ru" db:"name_ru"`
	NameTm string `json:"name_tm" db:"name_tm"`
}

type ChannelManagement struct {
	Id     int    `json:"id" db:"id"`
	NameRu string `json:"name_ru" db:"name_ru"`
	NameTm string `json:"name_tm" db:"name_tm"`
}

type ChannelObject struct {
	NameRu string `json:"name_ru"`
	NameTm string `json:"name_tm"`
}

type TradeChannel struct {
	Id           int               `json:"id" db:"id"`
	TypeId       int               `json:"type_id" db:"type_id"`
	StructureId  int               `json:"structure_id" db:"structure_id"`
	SizeId       int               `json:"size_id" db:"size_id"`
	ManagementId int               `json:"management_id" db:"management_id"`
	Status       bool              `json:"status" db:"status"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
	Type         ChannelType       `json:"type" gorm:"foreignKey:TypeId;constraint:OnDelete:SET NULL;"`
	Structure    ChannelStructure  `json:"structure" gorm:"foreignKey:StructureId;constraint:OnDelete:SET NULL;"`
	Size         ChannelSize       `json:"size" gorm:"foreignKey:SizeId;constraint:OnDelete:SET NULL;"`
	Management   ChannelManagement `json:"management" gorm:"foreignKey:ManagementId;constraint:OnDelete:SET NULL;"`
}

type ReadChannel struct {
	Id           int    `json:"id"`
	TypeId       int    `json:"type_id"`
	StructureId  int    `json:"structure_id"`
	SizeId       int    `json:"size_id" `
	ManagementId int    `json:"management_id"`
	Type         string `json:"type"`
	Structure    string `json:"structure"`
	Size         string `json:"size"`
	Management   string `json:"management"`
	Status       bool   `json:"status"`
	UpdatedAt    string `json:"updated_at"`
}

type CreateChannel struct {
	TypeId       int  `json:"type_id"`
	StructureId  int  `json:"structure_id"`
	SizeId       int  `json:"size_id" `
	ManagementId int  `json:"management_id"`
	Status       bool `json:"status"`
}

type TradeChannelPage struct {
	Items      []ReadChannel `json:"items"`
	PageNumber int           `json:"pageNumber"`
	PageSize   int           `json:"pageSize"`
	PageCount  int           `json:"pageCount"`
}
