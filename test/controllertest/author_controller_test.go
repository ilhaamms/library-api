package controllertest

import (
	"encoding/json"
	"io"
	"log"
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

func TruncateAuthorTable(db *gorm.DB) {
	db.Exec("DELETE FROM author")

	db.Exec("DELETE FROM sqlite_sequence WHERE name='author'")
}

func SetupRouterAuthor() *gin.Engine {

	gin.SetMode(gin.TestMode)

	db, err := config.InitDbSQLite()
	if err != nil {
		panic(err)
	}

	TruncateAuthorTable(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	authorRepo := repository.NewAuthorRepository(db)
	authorService := service.NewAuthorService(authorRepo)
	authorController := controller.NewAuthorController(authorService)

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

	return r
}

func RequestRegisterUser(r *gin.Engine, reqBody string) *httptest.ResponseRecorder {

	TruncateUserTable()

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/register", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	return recorder
}

func RequestLoginUser(r *gin.Engine, reqBody string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/login", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	return recorder
}

func RequestCreateAuthor(r *gin.Engine, reqBody, token string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/authors", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	return recorder
}

func TestCreateAuthorUnauthenticated(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"name": "ilham",
		"birth_date": "2000-06-11"
	}`

	recorder := RequestCreateAuthor(r, reqBody, "")

	response := recorder.Result()

	var responseBody map[string]interface{}
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "invalid token", responseBody["message"])
}

func TestCreateAuthorSuccess(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	err := json.Unmarshal(body, &responseLogin)
	assert.Nil(t, err)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	reqBody = `{
		"name": "ilham",
		"birth_date": "2000-06-11"
	}`

	recorder := RequestCreateAuthor(r, reqBody, token)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	response = recorder.Result()

	bodyCreateAuthors, _ := io.ReadAll(response.Body)

	var responseAuthor map[string]interface{}
	err = json.Unmarshal(bodyCreateAuthors, &responseAuthor)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, "Berhasil menyimpan data author", responseAuthor["message"])

	TruncateUserTable()
}

func TestCreateAuthorFailedNameEmpty(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	recorderCreateAuthor := RequestCreateAuthor(r, `{"name": "", "birth_date": "2000-06-11"}`, token)

	response = recorderCreateAuthor.Result()

	body, _ = io.ReadAll(response.Body)

	var responseAuthor map[string]interface{}
	json.Unmarshal(body, &responseAuthor)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : nama dan tanggal lahir tidak boleh kosong", responseAuthor["error"])
}

func TestCreateAuthorFailedNameLength(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	reqBody = `{
		"name": "il",
		"birth_date": "2000-06-11"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqBody, token)
	assert.Equal(t, http.StatusBadRequest, recorderCreateAuthor.Code)

	response = recorderCreateAuthor.Result()

	body, _ = io.ReadAll(response.Body)

	var responseAuthor map[string]interface{}
	json.Unmarshal(body, &responseAuthor)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : nama minimal 3 karakter", responseAuthor["error"])
}

func TestCreateAuthorFailedBirthdateFormat(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	reqBody = `{
		"name": "ilham",
		"birth_date": "06-12-2000"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqBody, token)
	assert.Equal(t, http.StatusBadRequest, recorderCreateAuthor.Code)

	response = recorderCreateAuthor.Result()

	body, _ = io.ReadAll(response.Body)

	var responseAuthor map[string]interface{}
	json.Unmarshal(body, &responseAuthor)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : format bithdate salah, format harus YYYY-MM-DD atau tanggal, bulan anda tidak valid", responseAuthor["error"])
}

func TestGetAllUnauthenticated(t *testing.T) {
	r := SetupRouterAuthor()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/authors", nil)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	var responseBody map[string]interface{}
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "authorization required", responseBody["message"])
}

func TestGetAllAuthorSuccess(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	err := json.Unmarshal(body, &responseLogin)
	assert.Nil(t, err)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	createAuthor := RequestCreateAuthor(r, `{"name": "ilham", "birth_date": "2000-06-11"}`, token)
	assert.Equal(t, http.StatusCreated, createAuthor.Code)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/authors", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseAuthor map[string]interface{}
	err = json.Unmarshal(body, &responseAuthor)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Berhasil mengambil data list author", responseAuthor["message"])
}

func TestGetAllAuthorFailedPage(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	recorderCreateAuthor := RequestCreateAuthor(r, `{"name": "ilham", "birth_date": "2000-06-11"}`, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	recorderCreateAuthor = RequestCreateAuthor(r, `{"name": "ilham", "birth_date": "2000-06-11"}`, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/authors?page=300", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseAuthor map[string]interface{}
	json.Unmarshal(body, &responseAuthor)

	log.Println("responseAuthor : ", responseAuthor)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : page sudah melebihi total page", responseAuthor["error"])
}

func TestFindByIdUnauthenticated(t *testing.T) {
	r := SetupRouterAuthor()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/authors/1", nil)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	var responseBody map[string]interface{}
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "authorization required", responseBody["message"])
}

func TestFindByIdInvalidId(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/authors/0", nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	log.Println("responseBody : ", responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : id tidak valid", responseBody["error"])
}

func TestFindByIdAuthorNotFound(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/authors/300", nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : record not found", responseBody["error"])
}

func TestDeleteByIdInvalidId(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/authors/0", nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	log.Println("responseBody : ", responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : id tidak valid", responseBody["error"])
}

func TestDeleteByIdUnauthenticated(t *testing.T) {
	r := SetupRouterAuthor()

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/authors/1", nil)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	var responseBody map[string]interface{}
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "authorization required", responseBody["message"])
}

func TestDeleteByIdAuthorNotFound(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/authors/1", nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : gagal menghapus data author, author tidak ditemukan", responseBody["error"])
}

func TestUpdateUnauthenticated(t *testing.T) {
	r := SetupRouterAuthor()

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/authors/1", nil)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	var responseBody map[string]interface{}
	body, _ := io.ReadAll(response.Body)
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "authorization required", responseBody["message"])
}

func TestUpdateByIdInvalidId(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	reqBody = `{
		"name": "ilham",
		"birth_date": "2000-06-11"
	}`

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/authors/0", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	log.Println("responseBody : ", responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : id tidak valid", responseBody["error"])
}

func TestUpdateByIdAuthorNotFound(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	reqBody = `{
		"name": "ilham",
		"birth_date": "2000-06-11"
	}`

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/authors/300", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : gagal mengupdate data author : record not found", responseBody["error"])
}

func TestUpdateByIdNameBirthdateEmpty(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	reqBody = `{
		"name": "",
		"birth_date": ""
	}`

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/authors/1", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : field name dan birthdate tidak boleh kosong", responseBody["error"])
}

func TestUpdateByIdNameLength(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	reqBody = `{
		"name": "il",
		"birth_date": "2000-06-11"
	}`

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/authors/1", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : harap masukan nama minimal 3 karakter", responseBody["error"])
}

func TestUpdateByIdSuccess(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	reqBody = `{
		"name": "ilham",
		"birth_date": "2000-06-11"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqBody, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	response = recorderCreateAuthor.Result()

	bodyCreateAuthors, _ := io.ReadAll(response.Body)

	var responseAuthor map[string]interface{}
	json.Unmarshal(bodyCreateAuthors, &responseAuthor)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, "Berhasil menyimpan data author", responseAuthor["message"])

	reqBody = `{
		"name": "ilham sidiq",
		"birth_date": "2000-06-11"
	}`

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/authors/1", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Berhasil mengupdate data author", responseBody["message"])
}

func TestUpdateByIdBirthdateFormatFailed(t *testing.T) {
	r := SetupRouterAuthor()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	recorderRegister := RequestRegisterUser(r, reqBody)
	assert.Equal(t, http.StatusCreated, recorderRegister.Code)

	recorderLogin := RequestLoginUser(r, reqBody)
	assert.Equal(t, http.StatusOK, recorderLogin.Code)

	response := recorderLogin.Result()

	body, _ := io.ReadAll(response.Body)

	var responseLogin map[string]interface{}
	json.Unmarshal(body, &responseLogin)

	token := responseLogin["data"].(map[string]interface{})["token"].(string)

	reqBody = `{
		"name": "ilham",
		"birth_date": "2000-06-11"
	}`

	recorderCreateAuthor := RequestCreateAuthor(r, reqBody, token)
	assert.Equal(t, http.StatusCreated, recorderCreateAuthor.Code)

	response = recorderCreateAuthor.Result()

	bodyCreateAuthors, _ := io.ReadAll(response.Body)

	var responseAuthor map[string]interface{}
	json.Unmarshal(bodyCreateAuthors, &responseAuthor)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, "Berhasil menyimpan data author", responseAuthor["message"])

	reqBody = `{
		"name": "ilham sidiq",
		"birth_date": "11-06-2000"
	}`

	request := httptest.NewRequest(http.MethodPut, "http://localhost:8080/authors/1", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	r.ServeHTTP(recorder, request)

	response = recorder.Result()

	body, _ = io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "error : format bithdate salah, format harus YYYY-MM-DD atau tanggal, bulan anda tidak valid", responseBody["error"])
}
