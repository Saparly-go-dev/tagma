package tagma

type Test struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name_ru" db:"name_ru"`
}
