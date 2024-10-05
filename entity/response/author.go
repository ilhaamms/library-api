package response

import "time"

type CreateAuthor struct {
	Name      string    `json:"name"`
	Birthdate time.Time `json:"birth_date" gorm:"column:birth_date"`
}

type Author struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Birthdate time.Time `json:"birth_date" gorm:"column:birth_date"`
}

type UpdateAuthor struct {
	Name      string    `json:"name"`
	Birthdate time.Time `json:"birth_date" gorm:"column:birth_date"`
}

type WebResponseAuthor struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type WebResponseAuthors struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Pagination Pagination  `json:"pagination"`
	Data       interface{} `json:"data"`
}
