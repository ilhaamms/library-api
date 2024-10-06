package response

type CreateBook struct {
	Title    string `json:"title"`
	Isbn     string `json:"isbn"`
	AuthorId int    `json:"author_id"`
}

type Book struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Isbn       string `json:"isbn"`
	AuthorId   int    `json:"author_id"`
	AuthorName string `json:"author_name"`
	BirthDate  string `json:"birth_date"`
}

type ResultBook struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Isbn       string `json:"isbn"`
	AuthorBook `json:"author"`
}

type WebResponseBook struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type WebResponseBooks struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Pagination Pagination  `json:"pagination"`
	Data       interface{} `json:"data"`
}
