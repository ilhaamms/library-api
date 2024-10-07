package controllertest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ilhaamms/library-api/config"
	"github.com/ilhaamms/library-api/controller"
	"github.com/ilhaamms/library-api/middleware"
	"github.com/ilhaamms/library-api/repository"
	"github.com/ilhaamms/library-api/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TruncateTableBook(db *gorm.DB) {
	db.Exec("DELETE FROM book")

	db.Exec("DELETE FROM sqlite_sequence WHERE name='book'")
}

func SetupRouterBook() *gin.Engine {

	gin.SetMode(gin.TestMode)

	db, err := config.InitDbSQLite()
	if err != nil {
		panic(err)
	}

	TruncateTableBook(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	authorRepo := repository.NewAuthorRepository(db)
	authorService := service.NewAuthorService(authorRepo)
	authorController := controller.NewAuthorController(authorService)

	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookController := controller.NewBookController(bookService)

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
	}

	r.POST("/authors", middleware.Auth(), authorController.CreateAuthor)
	r.GET("/authors", middleware.Auth(), authorController.GetAllAuthor)
	r.GET("/authors/:id", middleware.Auth(), authorController.GetAuthorsById)
	r.DELETE("/authors/:id", middleware.Auth(), authorController.DeleteAuthorsById)
	r.PUT("/authors/:id", middleware.Auth(), authorController.UpdateAuthorsById)

	r.POST("/books", middleware.Auth(), bookController.CreateBook)
	r.GET("/books", middleware.Auth(), bookController.GetAllBook)
	r.GET("/books/:id", middleware.Auth(), bookController.GetBookById)
	r.DELETE("/books/:id", middleware.Auth(), bookController.DeleteBookById)
	r.PUT("/books/:id", middleware.Auth(), bookController.Update)

	return r
}

func RequestCreateBook(r *gin.Engine, reqBody string, token string) *httptest.ResponseRecorder {

	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	return recorder
}

func TestSaveBookUnAuthorized(t *testing.T) {
	r := SetupRouterBook()

	reqBody := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorder := RequestCreateBook(r, reqBody, "")

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "invalid token", responseBody["message"])
}

func TestSaveBookSuccess(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	responseCreateBook := recorderCreateBook.Result()

	bodyCreateBook, _ := io.ReadAll(responseCreateBook.Body)

	var responseBodyCreateBook map[string]interface{}
	json.Unmarshal(bodyCreateBook, &responseBodyCreateBook)

	assert.Equal(t, http.StatusCreated, responseCreateBook.StatusCode)
	assert.Equal(t, "Berhasil menyimpan data book", responseBodyCreateBook["message"])
}

func TestSaveBookFailedTitleEmpty(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`
	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusBadRequest, recorderCreateBook.Code)

	responseCreateBook := recorderCreateBook.Result()

	bodyCreateBook, _ := io.ReadAll(responseCreateBook.Body)

	var responseBodyCreateBook map[string]interface{}
	json.Unmarshal(bodyCreateBook, &responseBodyCreateBook)

	assert.Equal(t, http.StatusBadRequest, responseCreateBook.StatusCode)
	assert.Equal(t, "error : judul, isbn, dan author_id tidak boleh kosong", responseBodyCreateBook["error"])
}

func TestSaveBookFailedIsbnEmpty(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusBadRequest, recorderCreateBook.Code)

	responseCreateBook := recorderCreateBook.Result()

	bodyCreateBook, _ := io.ReadAll(responseCreateBook.Body)

	var responseBodyCreateBook map[string]interface{}
	json.Unmarshal(bodyCreateBook, &responseBodyCreateBook)

	assert.Equal(t, http.StatusBadRequest, responseCreateBook.StatusCode)
	assert.Equal(t, "error : judul, isbn, dan author_id tidak boleh kosong", responseBodyCreateBook["error"])
}

func TestSaveBookFailedAuthorIdEmpty(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 0
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusBadRequest, recorderCreateBook.Code)

	responseCreateBook := recorderCreateBook.Result()

	bodyCreateBook, _ := io.ReadAll(responseCreateBook.Body)

	var responseBodyCreateBook map[string]interface{}
	json.Unmarshal(bodyCreateBook, &responseBodyCreateBook)

	assert.Equal(t, http.StatusBadRequest, responseCreateBook.StatusCode)
	assert.Equal(t, "error : judul, isbn, dan author_id tidak boleh kosong", responseBodyCreateBook["error"])
}

func TestSaveMinTitleBook(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Go",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusBadRequest, recorderCreateBook.Code)

	responseCreateBook := recorderCreateBook.Result()

	bodyCreateBook, _ := io.ReadAll(responseCreateBook.Body)

	var responseBodyCreateBook map[string]interface{}
	json.Unmarshal(bodyCreateBook, &responseBodyCreateBook)

	assert.Equal(t, http.StatusBadRequest, responseCreateBook.StatusCode)
	assert.Equal(t, "error : judul minimal 3 karakter", responseBodyCreateBook["error"])
}

func TestSaveMinIsbnBook(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "123",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusBadRequest, recorderCreateBook.Code)

	responseCreateBook := recorderCreateBook.Result()

	bodyCreateBook, _ := io.ReadAll(responseCreateBook.Body)

	var responseBodyCreateBook map[string]interface{}
	json.Unmarshal(bodyCreateBook, &responseBodyCreateBook)

	assert.Equal(t, http.StatusBadRequest, responseCreateBook.StatusCode)
	assert.Equal(t, "error : isbn minimal 10 karakter", responseBodyCreateBook["error"])
}

func TestSaveMaxIsbnBook(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "12345678901234",
		"author_id": 1
	}`
	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusBadRequest, recorderCreateBook.Code)

	responseCreateBook := recorderCreateBook.Result()

	bodyCreateBook, _ := io.ReadAll(responseCreateBook.Body)

	var responseBodyCreateBook map[string]interface{}
	json.Unmarshal(bodyCreateBook, &responseBodyCreateBook)

	assert.Equal(t, http.StatusBadRequest, responseCreateBook.StatusCode)
	assert.Equal(t, "error : isbn maksimal 13 karakter", responseBodyCreateBook["error"])
}

func TestSaveAuthorIdNegativeBook(t *testing.T) {

	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": -1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusBadRequest, recorderCreateBook.Code)

	responseCreateBook := recorderCreateBook.Result()

	bodyCreateBook, _ := io.ReadAll(responseCreateBook.Body)

	var responseBodyCreateBook map[string]interface{}
	json.Unmarshal(bodyCreateBook, &responseBodyCreateBook)

	assert.Equal(t, http.StatusBadRequest, responseCreateBook.StatusCode)
	assert.Equal(t, "error : author_id tidak boleh negatif", responseBodyCreateBook["error"])
}

func TestSaveBookIsbnExist(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	reqCreateBook = `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook = RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusBadRequest, recorderCreateBook.Code)

	responseCreateBook := recorderCreateBook.Result()

	bodyCreateBook, _ := io.ReadAll(responseCreateBook.Body)

	var responseBodyCreateBook map[string]interface{}
	json.Unmarshal(bodyCreateBook, &responseBodyCreateBook)

	assert.Equal(t, http.StatusBadRequest, responseCreateBook.StatusCode)
	assert.Equal(t, "error : isbn sudah digunakan oleh buku lain", responseBodyCreateBook["error"])
}

func TestGetAllBookUnAuthorized(t *testing.T) {
	r := SetupRouterBook()

	req := httptest.NewRequest(http.MethodGet, "/books", nil)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "authorization required", responseBody["message"])
}

func TestSuccessGetAllBook(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Data buku berhasil diambil", responseBodyGetAllBook["message"])
}

func TestGetAllBookEmptyData(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Data book kosong", responseBodyGetAllBook["message"])
}

func TestGetAllBookPageMorethanTotalPagge(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	req := httptest.NewRequest(http.MethodGet, "/books?page=2", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : page sudah melebihi total page", responseBodyGetAllBook["error"])
}

func TestFindByIdBookUnauthorized(t *testing.T) {
	r := SetupRouterBook()

	req := httptest.NewRequest(http.MethodGet, "http://localhost/books/1", nil)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "authorization required", responseBody["message"])
}

func TestFindByBookSuccess(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/books/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Data buku berhasil diambil", responseBodyGetAllBook["message"])
}

func TestFindByIdBookNegative(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	req := httptest.NewRequest(http.MethodGet, "http://localhost/books/-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : id tidak boleh negatif atau 0", responseBodyGetAllBook["error"])
}

func TestFindByIdBookNotFound(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	req := httptest.NewRequest(http.MethodGet, "/books/2", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : book tidak ditemukan", responseBodyGetAllBook["error"])
}

func TestDeleteBookUnAuthorized(t *testing.T) {
	r := SetupRouterBook()

	req := httptest.NewRequest(http.MethodDelete, "http://localhost/books/1", nil)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "authorization required", responseBody["message"])
}

func TestDeleteBookSuccess(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	req := httptest.NewRequest(http.MethodDelete, "http://localhost/books/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Data buku berhasil dihapus", responseBodyGetAllBook["message"])
}

func TestDeleteBookNegativeId(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	req := httptest.NewRequest(http.MethodDelete, "http://localhost/books/-1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : id tidak boleh negatif atau 0", responseBodyGetAllBook["error"])
}

func TestDeleteBookNotfoundId(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	req := httptest.NewRequest(http.MethodDelete, "http://localhost/books/300", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : gagal menghapus data book, book tidak ditemukan", responseBodyGetAllBook["error"])
}

func TestUpdateBookUnauthorized(t *testing.T) {
	r := SetupRouterBook()

	req := httptest.NewRequest(http.MethodPut, "http://localhost/books/1", nil)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "authorization required", responseBody["message"])
}

func TestUpdateBookSuccess(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	reqUpdateBook := `{
		"title": "Belajar Golang",
		"isbn": "12345678910",
		"author_id": 1
	}`

	req := httptest.NewRequest(http.MethodPut, "http://localhost/books/1", strings.NewReader(reqUpdateBook))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Data buku berhasil diupdate", responseBodyGetAllBook["message"])
}

func TestUpdateBookFailedIdNegative(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	reqUpdateBook := `{
		"title": "Belajar Golang",
		"isbn": "12345678910",
		"author_id": 1
	}`

	req := httptest.NewRequest(http.MethodPut, "http://localhost/books/-1", strings.NewReader(reqUpdateBook))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : id tidak boleh negatif atau 0", responseBodyGetAllBook["error"])
}

func TestUpdateBookIdNotFound(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqCreateAuthor := `{
		"name": "Ilham Sidiq",
		"birth_date": "1996-01-01"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqCreateAuthor, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	reqCreateBook := `{
		"title": "Belajar Golang",
		"isbn": "1234567890",
		"author_id": 1
	}`

	recorderCreateBook := RequestCreateBook(r, reqCreateBook, token)
	assert.Equal(t, http.StatusCreated, recorderCreateBook.Code)

	reqUpdateBook := `{
		"title": "Belajar Golang",
		"isbn": "12345678910",
		"author_id": 1
	}`

	req := httptest.NewRequest(http.MethodPut, "http://localhost/books/300", strings.NewReader(reqUpdateBook))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : gagal mengupdate data book, book tidak ditemukan", responseBodyGetAllBook["error"])
}

func TestUpdateBookFailedTitleEmpty(t *testing.T) {
	r := SetupRouterBook()

	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)
	reqUpdateBook := `{
		"title": "",
		"isbn": "12345678910",
		"author_id": 1
	}`

	req := httptest.NewRequest(http.MethodPut, "http://localhost/books/1", strings.NewReader(reqUpdateBook))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : judul, isbn, dan author_id tidak boleh kosong", responseBodyGetAllBook["error"])
}

func TestUpdateBookFailedIsbnEmpty(t *testing.T) {
	r := SetupRouterBook()
	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)

	reqUpdateBook := `{
		"title": "Belajar Golang",
		"isbn": "",
		"author_id": 1
	}`

	req := httptest.NewRequest(http.MethodPut, "http://localhost/books/1", strings.NewReader(reqUpdateBook))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, req)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : judul, isbn, dan author_id tidak boleh kosong", responseBodyGetAllBook["error"])
}

func TestUpdateBookFailedAuthorIdEmpty(t *testing.T) {
	r := SetupRouterBook()
	reqBodyRegister := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBodyRegister)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	reqBodyLogin := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderLogin := RequestLoginUser(r, reqBodyLogin)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()
	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	token := responseBody["data"].(map[string]interface{})["token"].(string)
	reqUpdateBook := `{
		"title": "Belajar Golang",
		"isbn": "12345678910",
		"author_id": 0
	}`

	req := httptest.NewRequest(http.MethodPut, "http://localhost/books/1", strings.NewReader(reqUpdateBook))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)
	var responseBodyGetAllBook map[string]interface{}
	json.Unmarshal(body, &responseBodyGetAllBook)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : judul, isbn, dan author_id tidak boleh kosong", responseBodyGetAllBook["error"])
}
