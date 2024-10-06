package repomock

import (
	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (r *UserRepositoryMock) Save(user request.User) error {
	args := r.Mock.Called(user)
	if args.Get(0) == nil {
		return nil
	}

	dataUser := args.Get(0).(error)

	return dataUser
}

func (r *UserRepositoryMock) Login(user request.User) (response.User, error) {
	args := r.Mock.Called(user)
	if args.Get(0) == nil {
		return response.User{}, nil
	}

	dataUser := args.Get(0).(response.User)

	return dataUser, nil
}

func (r *UserRepositoryMock) CheckUsername(username string) (bool, error) {
	args := r.Mock.Called(username)
	if args.Get(0) == nil {
		return false, nil
	}

	dataUser := args.Get(0).(bool)

	return dataUser, nil
}

func (r *UserRepositoryMock) GetUserByUsername(username string) (request.User, error) {
	args := r.Mock.Called(username)
	if args.Get(0) == nil {
		return request.User{}, nil
	}

	dataUser := args.Get(0).(request.User)

	return dataUser, nil
}
