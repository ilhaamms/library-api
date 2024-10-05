package request

import "time"

type Author struct {
	ID        int       `json:"id" form:"id"`
	Name      string    `json:"name" form:"name"`
	Birthdate time.Time `json:"birth_date" form:"birth_date" gorm:"column:birth_date"`
}

type CreateAuthor struct {
	Name      string `json:"name" form:"name"`
	Birthdate string `json:"birth_date" form:"birth_date" gorm:"column:birth_date"`
}

type UpdateAuthor struct {
	Name      string `json:"name" form:"name"`
	Birthdate string `json:"birth_date" form:"birth_date" gorm:"column:birth_date"`
}
