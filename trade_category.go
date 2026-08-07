package tagma

import "time"

type Trade_category struct {
	Id        int       `json:"id" db:"id"`
	NameRu    string    `json:"name_ru" db:"name_ru" binding:"required"`
	NameTm    string    `json:"name_tm" db:"name_tm" binding:"required"`
	Min_Sum   float32   `json:"min_sum" db:"min_sum" binding:"required"`
	Max_Sum   float32   `json:"max_sum" db:"max_sum" binding:"required"`
	Status    bool      `json:"status" db:"status"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTrade_category struct {
	Name    string  `json:"name"`
	Min_Sum float32 `json:"min_sum"`
	Max_Sum float32 `json:"max_sum"`
	Status  bool    `json:"status"`
}

type TradeCategoryPage struct {
	Items      []ReadTrade_category `json:"items"`
	PageNumber int                  `json:"pageNumber"`
	PageSize   int                  `json:"pageSize"`
	PageCount  int                  `json:"pageCount"`
}

type ReadTrade_category struct {
	Id        int     `json:"id"`
	NameRu    string  `json:"name_ru"`
	NameTm    string  `json:"name_tm"`
	Min_Sum   float32 `json:"min_sum"`
	Max_Sum   float32 `json:"max_sum"`
	Status    bool    `json:"status"`
	UpdatedAt string  `json:"updated_at"`
}
