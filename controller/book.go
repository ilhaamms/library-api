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

type BookController interface {
	CreateBook(c *gin.Context)
	GetAllBook(c *gin.Context)
	GetBookById(c *gin.Context)
	DeleteBookById(c *gin.Context)
	Update(c *gin.Context)
}

type bookController struct {
	bookService service.BookService
}

func NewBookController(bookService service.BookService) BookController {
	return &bookController{bookService: bookService}
}

func (bc *bookController) CreateBook(c *gin.Context) {

	var book request.CreateBook

	err := c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	_, err = bc.bookService.Save(book)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusCreated, response.WebResponseBook{
		StatusCode: http.StatusCreated,
		Message:    "Berhasil menyimpan data book",
		Data:       book,
	})
}

func (bc *bookController) GetAllBook(c *gin.Context) {
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

	books, totalPages, err := bc.bookService.FindAll(page, limit)
	if err != nil {

		if err.Error() == "data book kosong" {
			c.JSON(http.StatusOK, response.WebResponseBook{
				StatusCode: http.StatusOK,
				Message:    "Data book kosong",
				Data:       books,
			})
			return
		}

		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, response.WebResponseBooks{
		StatusCode: http.StatusOK,
		Message:    "Data buku berhasil diambil",
		Pagination: response.Pagination{
			CurrentPage: page,
			TotalPage:   totalPages,
			Limit:       limit,
		},
		Data: books,
	})
}

func (bc *bookController) GetBookById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	book, err := bc.bookService.FindById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, response.WebResponseBook{
		StatusCode: http.StatusOK,
		Message:    "Data buku berhasil diambil",
		Data:       book,
	})
}

func (bc *bookController) DeleteBookById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	book, err := bc.bookService.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, response.WebResponseBook{
		StatusCode: http.StatusOK,
		Message:    "Data buku berhasil dihapus",
		Data:       book,
	})
}

func (bc *bookController) Update(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	var book request.UpdateBook

	err = c.ShouldBind(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	_, err = bc.bookService.Update(id, book)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, response.WebResponseBook{
		StatusCode: http.StatusOK,
		Message:    "Data buku berhasil diupdate",
		Data:       book,
	})
}
