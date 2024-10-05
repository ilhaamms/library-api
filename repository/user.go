package repository

import (
	"github.com/ilhaamms/library-api/entity/request"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user request.User) error
	CheckUsername(username string) (bool, error)
	GetUserByUsername(username string) (request.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Save(user request.User) error {
	err := r.db.Table("user").Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) CheckUsername(username string) (bool, error) {
	var user request.User
	err := r.db.Table("user").Where("username = ?", username).First(&user).Error
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (r *userRepository) GetUserByUsername(username string) (request.User, error) {
	var user request.User
	err := r.db.Table("user").Where("username = ?", username).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
