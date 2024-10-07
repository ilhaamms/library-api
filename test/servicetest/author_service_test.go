package servicetest

import (
	"errors"
	"testing"
	"time"

	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"github.com/ilhaamms/library-api/service"
	"github.com/ilhaamms/library-api/test/repomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthorService_SaveFailedNameEmpty(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	author := request.CreateAuthor{
		Name:      "",
		Birthdate: "2000-01-01",
	}

	authorRepositoryMock.Mock.On("Save", author).Return(nil)

	result, err := authorService.Save(author)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "", author.Name)
	assert.Equal(t, "nama dan tanggal lahir tidak boleh kosong", err.Error())
	assert.Equal(t, "2000-01-01", author.Birthdate)
}

func TestAuthorService_SaveFailedBirthdateEmpty(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	author := request.CreateAuthor{
		Name:      "Ilhaam",
		Birthdate: "",
	}

	authorRepositoryMock.Mock.On("Save", author).Return(nil)

	result, err := authorService.Save(author)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "", author.Birthdate)
	assert.Equal(t, "nama dan tanggal lahir tidak boleh kosong", err.Error())
}

func TestAuthorService_SaveFailedNameLength(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	author := request.CreateAuthor{
		Name:      "Il",
		Birthdate: "2000-01-01",
	}

	authorRepositoryMock.Mock.On("Save", author).Return(nil)

	result, err := authorService.Save(author)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Il", author.Name)
	assert.Equal(t, "nama minimal 3 karakter", err.Error())
	assert.Equal(t, "2000-01-01", author.Birthdate)
}

func TestAuthorService_SaveFailedBirthdateFormat(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	author := request.CreateAuthor{
		Name:      "Ilhaam",
		Birthdate: "2000-01-",
	}

	authorRepositoryMock.Mock.On("Save", author).Return(nil)

	result, err := authorService.Save(author)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Ilhaam", author.Name)
	assert.Equal(t, "format bithdate salah, format harus YYYY-MM-DD atau tanggal, bulan anda tidak valid", err.Error())
	assert.Equal(t, "2000-01-", author.Birthdate)
}

func TestAuthorService_SaveFailedSaveAuthor(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	author := request.CreateAuthor{
		Name:      "Ilhaam",
		Birthdate: "2000-01-01",
	}

	authorRepositoryMock.Mock.On("Save", author).Return(errors.New("gagal menyimpan data author"))

	result, err := authorService.Save(author)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Ilhaam", author.Name)
	assert.Equal(t, "gagal menyimpan data author", err.Error())
	assert.Equal(t, "2000-01-01", author.Birthdate)
}

func TestAuthorService_SaveSuccess(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	author := request.CreateAuthor{
		Name:      "Syifa",
		Birthdate: "2000-06-01",
	}

	authorRepositoryMock.Mock.On("Save", author).Return(nil)

	result, err := authorService.Save(author)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Syifa", author.Name)
	assert.Equal(t, "2000-06-01", author.Birthdate)
}

func TestAuthorService_FindAllDataEmpty(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	authorRepositoryMock.Mock.On("FindAll").Return(nil, errors.New("data author kosong"))

	result, _, err := authorService.FindAll(1, 10)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "data author kosong", err.Error())
}

func TestAuthorService_FindAllSuccess(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	birthDate, _ := time.Parse("2006-01-02", "2000-06-11")
	expectDate, _ := time.Parse("2006-01-02", "2000-06-11")

	authors := []response.Author{
		{
			ID:        1,
			Name:      "Ilhaam Sidiq",
			BirthDate: birthDate,
		},
	}

	authorRepositoryMock.Mock.On("FindAll").Return(authors, nil)

	result, _, err := authorService.FindAll(1, 10)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, authors[0].ID)
	assert.Equal(t, "Ilhaam Sidiq", authors[0].Name)
	assert.Equal(t, expectDate, authors[0].BirthDate)
}

func TestAuthorService_FindAllPageFailed(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	birthDate, _ := time.Parse("2006-01-02", "2000-06-11")

	authors := []response.Author{
		{
			ID:        1,
			Name:      "Ilhaam Sidiq",
			BirthDate: birthDate,
		},
	}

	authorRepositoryMock.Mock.On("FindAll").Return(authors, nil)

	result, _, err := authorService.FindAll(2, 10)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "page sudah melebihi total page", err.Error())
}

func TestAuthorService_FindByIdFailedIdInvalid(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	authorRepositoryMock.Mock.On("FindById", 0).Return(response.Author{}, nil)

	result, err := authorService.FindById(0)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "id tidak valid", err.Error())
}

func TestAuthorService_FindByIdSuccess(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	birthDate, _ := time.Parse("2006-01-02", "2000-06-11")
	expectDate, _ := time.Parse("2006-01-02", "2000-06-11")

	author := response.Author{
		ID:        1,
		Name:      "Ilhaam Sidiq",
		BirthDate: birthDate,
	}

	authorRepositoryMock.Mock.On("FindById", 1).Return(author, nil)

	result, err := authorService.FindById(1)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, author.ID)
	assert.Equal(t, "Ilhaam Sidiq", author.Name)
	assert.Equal(t, expectDate, author.BirthDate)
}

func TestAuthorService_DeleteByIdFailedIdInvalid(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	authorRepositoryMock.Mock.On("DeleteById", 0).Return(response.Author{}, nil)

	result, err := authorService.DeleteById(0)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "id tidak valid", err.Error())
}

func TestAuthorService_DeleteByIdSuccess(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	birthDate, _ := time.Parse("2006-01-02", "2000-06-11")
	expectDate, _ := time.Parse("2006-01-02", "2000-06-11")

	author := response.Author{
		ID:        1,
		Name:      "Ilhaam Sidiq",
		BirthDate: birthDate,
	}

	authorRepositoryMock.Mock.On("DeleteById", 1).Return(&author, nil)

	result, err := authorService.DeleteById(1)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, author.ID)
	assert.Equal(t, "Ilhaam Sidiq", author.Name)
	assert.Equal(t, expectDate, author.BirthDate)
}

func TestAuthorService_UpdateByIdFailedIdInvalid(t *testing.T) {

	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	author := request.UpdateAuthor{
		Name:      "Ilhaam Sidiq",
		Birthdate: "2000-06-11",
	}

	authorRepositoryMock.Mock.On("UpdateById", 0, author).Return(response.Author{}, nil)

	result, err := authorService.UpdateById(0, author)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "id tidak valid", err.Error())
}

func TestAuthorService_UpdateByIdFailedNameLength(t *testing.T) {
	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	author := request.UpdateAuthor{
		Name:      "Il",
		Birthdate: "2000-06-11",
	}

	authorRepositoryMock.Mock.On("UpdateById", 1, author).Return(&response.Author{}, nil)

	result, err := authorService.UpdateById(1, author)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "harap masukan nama minimal 3 karakter", err.Error())
}

func TestAuthorService_UpdateByIdFailedNameBirthdateEmpty(t *testing.T) {
	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	author := request.UpdateAuthor{
		Name:      "",
		Birthdate: "",
	}

	authorRepositoryMock.Mock.On("UpdateById", 1, author).Return(&response.Author{}, nil)

	result, err := authorService.UpdateById(1, author)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "field name dan birthdate tidak boleh kosong", err.Error())
}

func TestAuthorService_UpdateByIdFormat(t *testing.T) {
	var authorRepositoryMock = repomock.AuthorRepositoryMock{Mock: mock.Mock{}}
	var authorService = service.AuthorServices{AuthorRepo: &authorRepositoryMock}

	author := request.UpdateAuthor{
		Name:      "Ilhaam",
		Birthdate: "2000-06-",
	}

	authorRepositoryMock.Mock.On("UpdateById", 1, author).Return(&response.Author{}, nil)

	result, err := authorService.UpdateById(1, author)

	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "format bithdate salah, format harus YYYY-MM-DD atau tanggal, bulan anda tidak valid", err.Error())
}
