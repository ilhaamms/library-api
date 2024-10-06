package repomock

import (
	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"github.com/stretchr/testify/mock"
)

type BookRepositoryMock struct {
	Mock mock.Mock
}

func (r *BookRepositoryMock) Save(book request.CreateBook) error {
	args := r.Mock.Called(book)
	if args.Get(0) == nil {
		return nil
	}

	dataBook := args.Get(0).(error)

	return dataBook
}

func (r *BookRepositoryMock) FindAll() ([]response.Book, error) {
	args := r.Mock.Called()
	if args.Get(0) == nil {
		return nil, nil
	}

	dataBooks := args.Get(0).([]response.Book)

	return dataBooks, nil
}

func (r *BookRepositoryMock) FindById(id int) (response.Book, error) {
	args := r.Mock.Called(id)
	if args.Get(0) == nil {
		return response.Book{}, nil
	}

	dataBook := args.Get(0).(response.Book)

	return dataBook, nil
}

func (r *BookRepositoryMock) Delete(id int) (*response.ResultBook, error) {
	args := r.Mock.Called(id)
	if args.Get(0) == nil {
		return nil, nil
	}

	dataBook := args.Get(0).(*response.ResultBook)

	return dataBook, nil
}

func (r *BookRepositoryMock) Update(id int, book request.UpdateBook) (*response.ResultBook, error) {
	args := r.Mock.Called(id, book)
	if args.Get(0) == nil {
		return nil, nil
	}

	dataBook := args.Get(0).(*response.ResultBook)

	return dataBook, nil
}

func (r *BookRepositoryMock) FindBookByIsbn(isbn string) (response.Book, error) {
	args := r.Mock.Called(isbn)
	if args.Get(0) == nil {
		return response.Book{}, nil
	}

	dataBook := args.Get(0).(response.Book)

	return dataBook, nil
}
