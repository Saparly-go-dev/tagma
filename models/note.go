package models

import "time"

type Note struct {
	Id          int       `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateNote struct {
	Description string `json:"description"`
}

type ReadNote struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	CreateAt    string `json:"create_at"`
}

type NotePage struct {
	Items      []ReadNote `json:"items"`
	PageNumber int        `json:"pageNumber"`
	PageSize   int        `json:"pageSize"`
	PageCount  int        `json:"pageCount"`
	Total      int        `json:"total"`
}
