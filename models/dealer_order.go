package models

import "time"

type DealerOrder struct {
	Id             int               `json:"id" db:"id"`
	DealerId       int               `json:"dealer_id" db:"dealer_id"`
	Sum            float64           `json:"sum" db:"sum"`
	SaleSum        float64           `json:"sale_sum" db:"sale_sum"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	List           []DealerOrderList `json:"list" gorm:"DealerOrderId"`
	Status         bool              `json:"status" db:"status"`
	IsClosed       bool              `json:"is_closed" db:"is_closed"`
	Completed_date time.Time         `json:"completed_date" db:"completed_date"`
	Paid_date      time.Time         `json:"paid_date" db:"paid_date"`
	Dealer         Dealer            `json:"dealer" gorm:"foreignKey:DealerId;constraint:OnDelete:CASCADE;"`
}

type ReadDealerOrder struct {
	Id            int     `json:"id"`
	CreatedAt     string  `json:"created_at"`
	CompletedDate string  `json:"completed_date"`
	PaidDate      string  `json:"paid_date"`
	DealerId      int     `json:"dealer_id"`
	Dealer        string  `json:"dealer"`
	Status        bool    `json:"status"`
	IsClosed      bool    `json:"is_closed"`
	Sum           float64 `json:"sum"`
}

type DealerOrderPage struct {
	Items      []ReadDealerOrder `json:"items"`
	PageNumber int               `json:"pageNumber"`
	PageSize   int               `json:"pageSize"`
	PageCount  int               `json:"pageCount"`
	Total      int               `json:"total"`
}

type DealerOrderList struct {
	Id            int      `json:"id" db:"id"`
	ProductId     int      `json:"product_id" db:"product_id"`
	DiscountId    int      `json:"discount_id" db:"discount_id"`
	DealerOrderId int      `json:"dealer_order_id" db:"dealer_order_id"`
	Price         float64  `json:"price" db:"price"`
	SalePrice     float64  `json:"sale_price" db:"sale_price"`
	Coefficient   int      `json:"coefficient" db:"coefficient"`
	Count         int      `json:"count" db:"count"`
	Status        bool     `json:"status" db:"status"`
	Product       Product  `json:"product" gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Discount      Discount `json:"discount" gorm:"foreignKey:DiscountId"`
}
