package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"github.com/ilhaamms/library-api/service"
)

type AuthorController interface {
	CreateAuthor(c *gin.Context)
	GetAuthors(c *gin.Context)
	GetAuthorsById(c *gin.Context)
	DeleteAuthorsById(c *gin.Context)
	UpdateAuthorsById(c *gin.Context)
}

type authorController struct {
	AuthorService service.AuthorService
}

func NewAuthorController(authorService service.AuthorService) AuthorController {
	return &authorController{AuthorService: authorService}
}

func (ac *authorController) CreateAuthor(c *gin.Context) {

	var author request.CreateAuthor

	err := c.ShouldBind(&author)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	_, err = ac.AuthorService.Save(author)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusCreated, response.WebResponseAuthor{
		StatusCode: http.StatusCreated,
		Message:    "Berhasil menyimpan data author",
		Data:       author,
	})
}

func (ac *authorController) GetAuthors(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	authors, totalPage, err := ac.AuthorService.FindAll(page, limit)
	if err != nil {
		if err.Error() == "data author kosong" {
			c.JSON(http.StatusOK, response.WebResponseAuthor{
				StatusCode: http.StatusOK,
				Message:    "Data author kosong",
				Data:       authors,
			})
			return
		}

		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, response.WebResponseAuthors{
		StatusCode: http.StatusOK,
		Message:    "Berhasil mengambil data list author",
		Pagination: response.Pagination{
			CurrentPage: page,
			TotalPage:   totalPage,
			Limit:       limit,
		},
		Data: authors,
	})
}

func (ac *authorController) GetAuthorsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	author, err := ac.AuthorService.FindById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, response.WebResponseAuthor{
		StatusCode: http.StatusOK,
		Message:    "Berhasil mengambil data author",
		Data:       author,
	})
}

func (ac *authorController) DeleteAuthorsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	author, err := ac.AuthorService.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, response.WebResponseAuthor{
		StatusCode: http.StatusOK,
		Message:    "Berhasil menghapus data author",
		Data:       author,
	})
}

func (ac *authorController) UpdateAuthorsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	var author request.UpdateAuthor

	err = c.ShouldBind(&author)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	authorResponse, err := ac.AuthorService.UpdateById(id, author)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, response.WebResponseAuthor{
		StatusCode: http.StatusOK,
		Message:    "Berhasil mengupdate data author",
		Data:       authorResponse,
	})
}
