package models

type Upr struct {
	Id       int       `json:"id" db:"id"`
	NameRu   string    `json:"name_ru" db:"name_ru"`
	NameTm   string    `json:"name_tm" db:"name_tm"`
	Contacts []Contact `gorm:"many2many:contact_uprs;constraint:OnUpdate:SET NULL;"`
}
