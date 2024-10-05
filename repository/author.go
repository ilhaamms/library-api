package repository

import (
	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"gorm.io/gorm"
)

type AuthorRepository interface {
	Save(author request.CreateAuthor) error
	FindAll() ([]response.Author, error)
	FindById(id int) (response.Author, error)
	DeleteById(id int) (*response.Author, error)
	UpdateById(id int, author request.UpdateAuthor) (*response.Author, error)
}

type authorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) AuthorRepository {
	return &authorRepository{db}
}

func (r *authorRepository) Save(author request.CreateAuthor) error {
	err := r.db.Table("author").Create(&author).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *authorRepository) FindAll() ([]response.Author, error) {
	var authors []response.Author
	err := r.db.Table("author").Find(&authors).Error
	if err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *authorRepository) FindById(id int) (response.Author, error) {
	var author response.Author
	err := r.db.Table("author").Where("id = ?", id).First(&author).Error
	if err != nil {
		return author, err
	}

	return author, nil
}

func (r *authorRepository) DeleteById(id int) (*response.Author, error) {

	var author response.Author

	err := r.db.Table("author").Where("id = ?", id).First(&author).Error
	if err != nil {
		return &author, err
	}

	err = r.db.Table("author").Where("id = ?", id).Delete(&response.Author{}).Error
	if err != nil {
		return &author, err
	}

	return &author, nil
}

func (r *authorRepository) UpdateById(id int, author request.UpdateAuthor) (*response.Author, error) {
	var authorResponse response.Author

	err := r.db.Table("author").Where("id = ?", id).Updates(&author).Error
	if err != nil {
		return &authorResponse, err
	}

	err = r.db.Table("author").Where("id = ?", id).First(&authorResponse).Error
	if err != nil {
		return &authorResponse, err
	}

	return &authorResponse, nil
}
