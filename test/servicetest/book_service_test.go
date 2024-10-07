package servicetest

import (
	"errors"
	"testing"

	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"github.com/ilhaamms/library-api/service"
	"github.com/ilhaamms/library-api/test/repomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookService_SaveFailedAllDataRequestEmpty(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book := request.CreateBook{
		Title:    "",
		Isbn:     "",
		AuthorId: 0,
	}

	_, err := bookService.Save(book)

	assert.NotNil(t, err)
	assert.Equal(t, "judul, isbn, dan author_id tidak boleh kosong", err.Error())

}

func TestBookService_SaveFailedMinTitle(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book := request.CreateBook{
		Title:    "il",
		Isbn:     "1234567890",
		AuthorId: 1,
	}

	_, err := bookService.Save(book)

	assert.NotNil(t, err)
	assert.Equal(t, "judul minimal 3 karakter", err.Error())

}

func TestBookService_SaveFailedMinIsbn(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book := request.CreateBook{
		Title:    "ilham",
		Isbn:     "123456789",
		AuthorId: 1,
	}

	_, err := bookService.Save(book)

	assert.NotNil(t, err)
	assert.Equal(t, "isbn minimal 10 karakter", err.Error())

}

func TestBookService_SaveFailedMaxIsbn(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book := request.CreateBook{
		Title:    "ilham",
		Isbn:     "12345678901234",
		AuthorId: 1,
	}

	_, err := bookService.Save(book)

	assert.NotNil(t, err)
	assert.Equal(t, "isbn maksimal 13 karakter", err.Error())

}

func TestBookService_SaveFailedAuthorIdNegative(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book := request.CreateBook{
		Title:    "ilham",
		Isbn:     "1234567890",
		AuthorId: -1,
	}

	_, err := bookService.Save(book)

	assert.NotNil(t, err)
	assert.Equal(t, "author_id tidak boleh negatif", err.Error())

}

func TestBookService_SaveFailedIsbnExist(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book := request.CreateBook{
		Title:    "ilham",
		Isbn:     "1234567890",
		AuthorId: 1,
	}

	bookRepositoryMock.Mock.On("FindBookByIsbn", book.Isbn).Return(nil, errors.New("isbn sudah digunakan oleh buku lain"))

	_, err := bookService.Save(book)

	assert.NotNil(t, err)
	assert.Equal(t, "isbn sudah digunakan oleh buku lain", err.Error())

}

func TestBookService_FindAllFailed(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	bookRepositoryMock.Mock.On("FindAll").Return(nil, errors.New("gagal mengambil data book"))

	book, _, err := bookService.FindAll(1, 10)

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_PageMoreThanTotalPages(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	bookRepositoryMock.Mock.On("FindAll").Return([]response.Book{}, nil)

	book, _, err := bookService.FindAll(3, 10)

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_InvalidId(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book, err := bookService.FindById(0)

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_DeleteFailedIdNegative(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book, err := bookService.DeleteById(-1)

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_FailedDeleteBookSuccess(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	bookRepositoryMock.Mock.On("Delete", 1).Return(nil, nil)

	book, _ := bookService.DeleteById(1)

	assert.Nil(t, book)
}

func TestBookService_FailedUpdateBookIdNegative(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book, err := bookService.Update(-1, request.UpdateBook{})

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_FailedUpdateBookEmptyDataRequest(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book, err := bookService.Update(1, request.UpdateBook{})

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_FailedUpdateBookMinTitle(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book, err := bookService.Update(1, request.UpdateBook{Title: "il"})

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_FailedUpdateBookMinIsbn(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book, err := bookService.Update(1, request.UpdateBook{Isbn: "123456789"})

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_FailedUpdateBookMaxIsbn(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book, err := bookService.Update(1, request.UpdateBook{Isbn: "12345678901234"})

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_FailedUpdateBookAuthorIdNegative(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	book, err := bookService.Update(1, request.UpdateBook{AuthorId: -1})

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_FailedUpdateBookIsbnExist(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	bookRepositoryMock.Mock.On("FindBookByIsbn", "1234567890").Return(response.Book{}, errors.New("isbn sudah digunakan oleh buku lain"))

	book, err := bookService.Update(1, request.UpdateBook{Isbn: "1234567890"})

	assert.Nil(t, book)
	assert.NotNil(t, err)
}

func TestBookService_FailedUpdateBookIdBookBotFound(t *testing.T) {

	var bookRepositoryMock = repomock.BookRepositoryMock{Mock: mock.Mock{}}
	var bookService = service.BookServices{BookRepository: &bookRepositoryMock}

	bookRepositoryMock.Mock.On("Update", 1, request.UpdateBook{}).Return(nil, errors.New("gagal mengupdate data book, book tidak ditemukan"))

	book, err := bookService.Update(1, request.UpdateBook{})

	assert.Nil(t, book)
	assert.NotNil(t, err)
}
