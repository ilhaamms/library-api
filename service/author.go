package service

import (
	"errors"
	"math"
	"time"

	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"github.com/ilhaamms/library-api/repository"
)

type AuthorService interface {
	Save(author request.CreateAuthor) (*response.CreateAuthor, error)
	FindAll(page, limit int) (*[]response.Author, int, error)
	FindById(id int) (*response.Author, error)
	DeleteById(id int) (*response.Author, error)
	UpdateById(id int, author request.UpdateAuthor) (*response.UpdateAuthor, error)
}

type authorService struct {
	AuthorRepo repository.AuthorRepository
}

func NewAuthorService(authorRepo repository.AuthorRepository) AuthorService {
	return &authorService{AuthorRepo: authorRepo}
}

func (s *authorService) Save(author request.CreateAuthor) (*response.CreateAuthor, error) {

	if author.Name == "" || author.Birthdate == "" {
		return nil, errors.New("nama dan tanggal lahir tidak boleh kosong")
	}

	if len(author.Name) < 3 {
		return nil, errors.New("nama minimal 3 karakter")
	}

	birthdate, err := time.Parse("2006-01-02", author.Birthdate)
	if err != nil {
		return nil, errors.New("format bithdate salah, format harus YYYY-MM-DD atau tanggal, bulan anda tidak valid")
	}

	err = s.AuthorRepo.Save(author)
	if err != nil {
		return nil, errors.New("gagal menyimpan data author : " + err.Error())
	}

	return &response.CreateAuthor{
		Name:      author.Name,
		Birthdate: birthdate,
	}, nil
}

func (s *authorService) FindAll(page, limit int) (*[]response.Author, int, error) {
	authors, err := s.AuthorRepo.FindAll()
	if err != nil {
		return nil, 0, errors.New("gagal mengambil data author : " + err.Error())
	}

	if len(authors) == 0 {
		return nil, 0, errors.New("data author kosong")
	}

	startIndex := (page - 1) * limit
	endIndex := int(math.Min(float64(startIndex+limit), float64(len(authors))))
	totalPages := int(math.Ceil(float64(len(authors)) / float64(limit)))

	if page > totalPages {
		return nil, 0, errors.New("page sudah melebihi total page")
	}

	authors = authors[startIndex:endIndex]

	return &authors, totalPages, nil
}

func (s *authorService) FindById(id int) (*response.Author, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	author, err := s.AuthorRepo.FindById(id)
	if err != nil {
		return nil, errors.New("author tidak ditemukan")
	}

	return &author, nil
}

func (s *authorService) DeleteById(id int) (*response.Author, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	author, err := s.AuthorRepo.DeleteById(id)
	if err != nil {
		return nil, errors.New("gagal menghapus data author, author tidak ditemukan")
	}

	return author, nil
}

func (s *authorService) UpdateById(id int, author request.UpdateAuthor) (*response.UpdateAuthor, error) {

	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}

	if author.Name == "" && author.Birthdate == "" {
		return nil, errors.New("field name dan birthdate tidak boleh kosong")
	}

	if author.Name != "" && len(author.Name) < 3 {
		return nil, errors.New("harap masukan nama minimal 3 karakter")
	}

	if author.Birthdate != "" {
		birthdate, err := time.Parse("2006-01-02", author.Birthdate)
		if err != nil {
			return nil, errors.New("format bithdate salah, format harus YYYY-MM-DD atau tanggal, bulan anda tidak valid")
		}

		author.Birthdate = birthdate.Format("2006-01-02")
	}

	authorResponse, err := s.AuthorRepo.UpdateById(id, author)
	if err != nil {
		return nil, errors.New("gagal mengupdate data author : " + err.Error())
	}

	return &response.UpdateAuthor{
		Name:      authorResponse.Name,
		Birthdate: authorResponse.Birthdate,
	}, nil
}
