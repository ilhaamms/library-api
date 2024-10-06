package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ilhaamms/library-api/entity/data"
	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"github.com/ilhaamms/library-api/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Save(user request.User) (*response.CreateUser, error)
	CheckUsername(username string) (bool, error)
	Login(user request.User) (bool, *response.ResponseUserLogin, error)
}

type userService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{UserRepository: userRepository}
}

func (s *userService) Save(user request.User) (*response.CreateUser, error) {

	if user.Username == "" || user.Password == "" {
		return nil, errors.New("username dan password wajib diisi")
	}

	isUsername, err := s.CheckUsername(user.Username)
	if err != nil {
		return nil, err
	}

	if isUsername {
		return nil, errors.New("username sudah digunakan oleh user lain")
	}

	if len(user.Password) < 8 {
		return nil, errors.New("harap masukkan password minimal 8 karakter")
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(bcryptPassword)

	err = s.UserRepository.Save(user)
	if err != nil {
		return nil, err
	}

	return &response.CreateUser{
		Username: user.Username,
		Password: user.Password,
	}, nil
}

func (s *userService) CheckUsername(username string) (bool, error) {
	if len(username) < 5 {
		return false, errors.New("username minimal 5 karakter")
	}

	if len(username) > 20 {
		return false, errors.New("username maksimal 20 karakter")
	}

	dataUsername, err := s.UserRepository.CheckUsername(username)
	if err != nil {
		return false, err
	}

	return dataUsername, nil
}

func (s *userService) Login(user request.User) (bool, *response.ResponseUserLogin, error) {

	if user.Username == "" || user.Password == "" {
		return false, nil, errors.New("username dan password wajib diisi")
	}

	dataUser, err := s.UserRepository.GetUserByUsername(user.Username)
	if err != nil {
		return false, nil, errors.New("username atau password salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(user.Password))
	if err != nil {
		return false, nil, errors.New("username atau password salah")
	}

	claims := &data.Claims{
		Username: dataUser.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(data.JwtKey))
	if err != nil {
		return false, nil, err
	}

	return true, &response.ResponseUserLogin{
		Username: dataUser.Username,
		Password: dataUser.Password,
		Token:    tokenString,
	}, nil
}
