package request

type Book struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Isbn     string `json:"isbn"`
	AuthorId int    `json:"author_id"`
}

type CreateBook struct {
	Title    string `json:"title" form:"title"`
	Isbn     string `json:"isbn" form:"isbn"`
	AuthorId int    `json:"author_id" form:"author_id"`
}

type UpdateBook struct {
	Title    string `json:"title" form:"title"`
	Isbn     string `json:"isbn" form:"isbn"`
	AuthorId int    `json:"author_id" form:"author_id"`
}
