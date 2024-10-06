package repository

import (
	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"gorm.io/gorm"
)

type BookRepository interface {
	Save(book request.CreateBook) error
	FindBookByIsbn(isbn string) (response.Book, error)
	FindAll() ([]response.Book, error)
	FindById(id int) (response.Book, error)
	Delete(id int) (*response.ResultBook, error)
	Update(id int, book request.UpdateBook) (*response.ResultBook, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Save(book request.CreateBook) error {
	bookData := request.CreateBook{
		Title:    book.Title,
		Isbn:     book.Isbn,
		AuthorId: book.AuthorId,
	}

	err := r.db.Table("book").Create(&bookData).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *bookRepository) FindBookByIsbn(isbn string) (response.Book, error) {
	var book response.Book

	err := r.db.Table("book").Where("isbn = ?", isbn).First(&book).Error
	if err != nil {
		return book, err
	}

	return book, nil
}

func (r *bookRepository) FindAll() ([]response.Book, error) {
	var books []response.Book

	err := r.db.Table("book AS b").
		Select("b.id, b.title, b.isbn, a.id AS author_id, a.name AS author_name, a.birth_date").
		Joins("INNER JOIN author AS a on b.author_id = a.id").
		Find(&books).Error

	if err != nil {
		return nil, err
	}

	return books, nil
}

func (r *bookRepository) FindById(id int) (response.Book, error) {
	var book response.Book

	err := r.db.Table("book AS b").
		Select("b.id, b.title, b.isbn, a.id AS author_id, a.name AS author_name, a.birth_date").
		Joins("INNER JOIN author AS a on b.author_id = a.id").
		Where("b.id = ?", id).
		First(&book).Error

	if err != nil {
		return book, err
	}

	return book, nil
}

func (r *bookRepository) Delete(id int) (*response.ResultBook, error) {

	var book response.Book
	err := r.db.Table("book AS b").
		Select("b.id, b.title, b.isbn, a.id AS author_id, a.name AS author_name, a.birth_date").
		Joins("INNER JOIN author AS a on b.author_id = a.id").
		Where("b.id = ?", id).
		First(&book).Error

	if err != nil {
		return nil, err
	}

	err = r.db.Table("book").Where("id = ?", id).Delete(&response.Book{}).Error
	if err != nil {
		return nil, err
	}

	return &response.ResultBook{
		Id:    book.Id,
		Title: book.Title,
		Isbn:  book.Isbn,
		AuthorBook: response.AuthorBook{
			ID:        book.AuthorId,
			Name:      book.AuthorName,
			BirthDate: book.BirthDate,
		},
	}, nil
}

func (r *bookRepository) Update(id int, book request.UpdateBook) (*response.ResultBook, error) {
	var bookResponse response.Book

	err := r.db.Table("book").Where("id = ?", id).Updates(&book).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Table("book AS b").
		Select("b.id, b.title, b.isbn, a.id AS author_id, a.name AS author_name, a.birth_date").
		Joins("INNER JOIN author AS a on b.author_id = a.id").
		Where("b.id = ?", id).
		First(&bookResponse).Error

	if err != nil {
		return nil, err
	}

	return &response.ResultBook{
		Id:    bookResponse.Id,
		Title: bookResponse.Title,
		Isbn:  bookResponse.Isbn,
		AuthorBook: response.AuthorBook{
			ID:        bookResponse.AuthorId,
			Name:      bookResponse.AuthorName,
			BirthDate: bookResponse.BirthDate,
		},
	}, nil
}
