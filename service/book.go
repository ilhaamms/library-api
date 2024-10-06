package service

import (
	"errors"
	"math"

	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"github.com/ilhaamms/library-api/repository"
)

type BookService interface {
	Save(book request.CreateBook) (*response.CreateBook, error)
	FindAll(page, limit int) (*[]response.ResultBook, int, error)
	FindById(id int) (*response.ResultBook, error)
	DeleteById(id int) (*response.ResultBook, error)
	Update(id int, book request.UpdateBook) (*response.ResultBook, error)
}

type bookService struct {
	BookRepository repository.BookRepository
}

func NewBookService(bookRepository repository.BookRepository) BookService {
	return &bookService{BookRepository: bookRepository}
}

func (s *bookService) Save(book request.CreateBook) (*response.CreateBook, error) {

	if book.Title == "" || book.Isbn == "" || book.AuthorId == 0 {
		return nil, errors.New("judul, isbn, dan author_id tidak boleh kosong")
	}

	if len(book.Title) < 3 {
		return nil, errors.New("judul minimal 3 karakter")
	}

	if len(book.Isbn) < 10 {
		return nil, errors.New("isbn minimal 10 karakter")
	}

	if len(book.Isbn) > 13 {
		return nil, errors.New("isbn maksimal 13 karakter")
	}

	if book.AuthorId < 0 {
		return nil, errors.New("author_id tidak boleh negatif")
	}

	_, err := s.BookRepository.FindBookByIsbn(book.Isbn)
	if err == nil {
		return nil, errors.New("isbn sudah digunakan oleh buku lain")
	}

	err = s.BookRepository.Save(book)
	if err != nil {
		return nil, err
	}

	return &response.CreateBook{
		Title: book.Title,
		Isbn:  book.Isbn,
	}, nil

}

func (s *bookService) FindAll(page, limit int) (*[]response.ResultBook, int, error) {
	books, err := s.BookRepository.FindAll()
	if err != nil {
		return nil, 0, errors.New("gagal mengambil data book : " + err.Error())
	}

	var listBook []response.ResultBook
	for _, book := range books {
		dataBook := response.ResultBook{
			Id:    book.Id,
			Title: book.Title,
			Isbn:  book.Isbn,
			AuthorBook: response.AuthorBook{
				ID:        book.AuthorId,
				Name:      book.AuthorName,
				BirthDate: book.BirthDate,
			},
		}

		listBook = append(listBook, dataBook)
	}

	if len(listBook) == 0 {
		return nil, 0, errors.New("data book kosong")
	}

	startIndex := (page - 1) * limit
	endIndex := int(math.Min(float64(startIndex+limit), float64(len(listBook))))
	totalPages := int(math.Ceil(float64(len(listBook)) / float64(limit)))

	if page > totalPages {
		return nil, 0, errors.New("page sudah melebihi total page")
	}

	listBook = listBook[startIndex:endIndex]

	return &listBook, totalPages, nil
}

func (s *bookService) FindById(id int) (*response.ResultBook, error) {

	if id <= 0 {
		return nil, errors.New("id tidak boleh negatif atau 0")
	}

	book, err := s.BookRepository.FindById(id)
	if err != nil {
		return nil, errors.New("book tidak ditemukan")
	}

	dataBook := response.ResultBook{
		Id:    book.Id,
		Title: book.Title,
		Isbn:  book.Isbn,
		AuthorBook: response.AuthorBook{
			ID:        book.AuthorId,
			Name:      book.AuthorName,
			BirthDate: book.BirthDate,
		},
	}

	return &dataBook, nil
}

func (s *bookService) DeleteById(id int) (*response.ResultBook, error) {

	if id <= 0 {
		return nil, errors.New("id tidak boleh negatif atau 0")
	}

	book, err := s.BookRepository.Delete(id)
	if err != nil {
		return nil, errors.New("gagal menghapus data book, book tidak ditemukan")
	}

	return book, nil
}

func (s *bookService) Update(id int, book request.UpdateBook) (*response.ResultBook, error) {

	if id <= 0 {
		return nil, errors.New("id tidak boleh negatif atau 0")
	}

	if book.Title == "" || book.Isbn == "" || book.AuthorId == 0 {
		return nil, errors.New("judul, isbn, dan author_id tidak boleh kosong")
	}

	if len(book.Title) < 3 {
		return nil, errors.New("judul minimal 3 karakter")
	}

	if len(book.Isbn) < 10 {
		return nil, errors.New("isbn minimal 10 karakter")
	}

	if len(book.Isbn) > 13 {
		return nil, errors.New("isbn maksimal 13 karakter")
	}

	if book.AuthorId < 0 {
		return nil, errors.New("author_id tidak boleh negatif")
	}

	_, err := s.BookRepository.FindBookByIsbn(book.Isbn)
	if err == nil {
		return nil, errors.New("isbn sudah digunakan oleh buku lain")
	}

	bookUpdate, err := s.BookRepository.Update(id, book)
	if err != nil {
		return nil, errors.New("gagal mengupdate data book, book tidak ditemukan")
	}

	dataBook := response.ResultBook{
		Id:    bookUpdate.Id,
		Title: bookUpdate.Title,
		Isbn:  bookUpdate.Isbn,
		AuthorBook: response.AuthorBook{
			ID:        bookUpdate.Id,
			Name:      bookUpdate.Name,
			BirthDate: bookUpdate.BirthDate,
		},
	}

	return &dataBook, nil
}
