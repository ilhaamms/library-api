package repomock

import (
	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"github.com/stretchr/testify/mock"
)

type AuthorRepositoryMock struct {
	Mock mock.Mock
}

func (r *AuthorRepositoryMock) Save(author request.CreateAuthor) error {
	args := r.Mock.Called(author)
	if args.Get(0) == nil {
		return nil
	}

	dataAuthor := args.Get(0).(error)

	return dataAuthor
}

func (r *AuthorRepositoryMock) FindAll() ([]response.Author, error) {
	args := r.Mock.Called()
	if args.Get(0) == nil {
		return nil, nil
	}

	dataAuthors := args.Get(0).([]response.Author)

	return dataAuthors, nil
}

func (r *AuthorRepositoryMock) FindById(id int) (response.Author, error) {
	args := r.Mock.Called(id)
	if args.Get(0) == nil {
		return response.Author{}, nil
	}

	dataAuthor := args.Get(0).(response.Author)

	return dataAuthor, nil
}

func (r *AuthorRepositoryMock) DeleteById(id int) (*response.Author, error) {
	args := r.Mock.Called(id)
	if args.Get(0) == nil {
		return nil, nil
	}

	dataAuthor := args.Get(0).(*response.Author)

	return dataAuthor, nil
}

func (r *AuthorRepositoryMock) UpdateById(id int, author request.UpdateAuthor) (*response.Author, error) {
	args := r.Mock.Called(id, author)
	if args.Get(0) == nil {
		return nil, nil
	}

	dataAuthor := args.Get(0).(*response.Author)

	return dataAuthor, nil
}
