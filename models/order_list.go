package models

type OrderList struct {
	Id          int      `json:"id" db:"id"`
	ProductId   int      `json:"product_id" db:"product_id"`
	DiscountId  int      `json:"discount_id" db:"discount_id"`
	OrderId     int      `json:"order_id" db:"order_id"`
	Price       float64  `json:"price" db:"price"`
	SalePrice   float64  `json:"sale_price" db:"sale_price"`
	Coefficient int      `json:"coefficient" db:"coefficient"`
	Count       int      `json:"count" db:"count"`
	Status      bool     `json:"status" db:"status"`
	Product     Product  `json:"product" gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Discount    Discount `json:"discount" gorm:"foreignKey:DiscountId"`
}

type ReadOrderList struct {
	Id            int     `json:"id"`
	ProductId     int     `json:"product_id"`
	Name          string  `json:"name"`
	Taste         string  `json:"taste"`
	Volume        float64 `json:"volume"`
	Count         int     `json:"count"`
	Price         float64 `json:"price"`
	DiscountPrice float64 `json:"discount_price"`
	IsProcent     bool    `json:"is_procent"`
	Start         string  `json:"start"`
	End           string  `json:"end"`
	Status        bool    `json:"status"`
	Coefficient   int     `json:"coefficient"`
}

type ProductOrder struct {
	Id              int     `json:"id"`
	ProductId       int     `json:"product_id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Taste           string  `json:"taste"`
	Volume          float64 `json:"volume"`
	Count           int     `json:"count"`
	Price           float64 `json:"price"`
	SalePrice       float64 `json:"sale_price"`
	Status          bool    `json:"status"`
	ProductTypeName string  `json:"product_type_name"`
}

type GetOrderObject struct {
	Products    []ProductOrder `json:"products"`
	Price       float64        `json:"price"`
	SalePrice   float64        `json:"sale_price"`
	Discount    float64        `json:"discount"`
	Promo       float64        `json:"promo"`
	Is_Credit   bool           `json:"is_credit"`
	PaymentType string         `json:"payment_type"`
	Images      []ReadImage    `json:"images"`
}
