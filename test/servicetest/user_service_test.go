package servicetest

import (
	"testing"

	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/service"
	"github.com/ilhaamms/library-api/test/repomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_SaveFailedNameAndPwEmpty(t *testing.T) {

	var userRepositoryMock = repomock.UserRepositoryMock{Mock: mock.Mock{}}
	var userService = service.UserServices{UserRepository: &userRepositoryMock}

	user := request.User{
		Username: "",
		Password: "",
	}

	_, err := userService.Save(user)

	assert.NotNil(t, err)
	assert.Equal(t, "username dan password wajib diisi", err.Error())
}

func TestUserService_SaveFailedMinUsername(t *testing.T) {

	var userRepositoryMock = repomock.UserRepositoryMock{Mock: mock.Mock{}}
	var userService = service.UserServices{UserRepository: &userRepositoryMock}

	user := request.User{
		Username: "il",
		Password: "12345678",
	}

	_, err := userService.Save(user)

	assert.NotNil(t, err)
	assert.Equal(t, "username minimal 5 karakter", err.Error())
}

func TestUserService_SaveFailedMaxUsername(t *testing.T) {

	var userRepositoryMock = repomock.UserRepositoryMock{Mock: mock.Mock{}}
	var userService = service.UserServices{UserRepository: &userRepositoryMock}

	user := request.User{
		Username: "ilham muhammad sidiq sedang bermain bola",
		Password: "12345678",
	}

	_, err := userService.Save(user)

	assert.NotNil(t, err)
	assert.Equal(t, "username maksimal 20 karakter", err.Error())
}

func TestUserService_SaveFailedUsernameExist(t *testing.T) {

	var userRepositoryMock = repomock.UserRepositoryMock{Mock: mock.Mock{}}
	var userService = service.UserServices{UserRepository: &userRepositoryMock}

	user := request.User{
		Username: "ilham",
		Password: "12345678",
	}

	userRepositoryMock.Mock.On("CheckUsername", user.Username).Return(true, nil)

	_, err := userService.Save(user)

	assert.NotNil(t, err)
	assert.Equal(t, "username sudah digunakan oleh user lain", err.Error())
}

func TestUserService_SaveFailedMinPassword(t *testing.T) {

	var userRepositoryMock = repomock.UserRepositoryMock{Mock: mock.Mock{}}
	var userService = service.UserServices{UserRepository: &userRepositoryMock}

	user := request.User{
		Username: "ilham",
		Password: "123",
	}

	userRepositoryMock.Mock.On("CheckUsername", user.Username).Return(false, nil)

	_, err := userService.Save(user)

	assert.NotNil(t, err)
	assert.Equal(t, "harap masukkan password minimal 8 karakter", err.Error())
}

func TestUserService_SaveUserSuccessRegister(t *testing.T) {

	var userRepositoryMock = repomock.UserRepositoryMock{Mock: mock.Mock{}}
	var userService = service.UserServices{UserRepository: &userRepositoryMock}

	user := request.User{
		Username: "ilham",
		Password: "12345678",
	}

	userRepositoryMock.Mock.On("CheckUsername", user.Username).Return(false, nil)

	userRepositoryMock.Mock.On("Save", mock.MatchedBy(func(u request.User) bool {
		return u.Username == user.Username && bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)) == nil
	})).Return(nil)

	_, err := userService.Save(user)

	assert.Nil(t, err)
}
