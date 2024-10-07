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
	"github.com/ilhaamms/library-api/repository"
	"github.com/ilhaamms/library-api/service"
	"github.com/stretchr/testify/assert"
)

func SetupRouterUser() *gin.Engine {

	gin.SetMode(gin.TestMode)

	db, err := config.InitDbSQLite()
	if err != nil {
		panic(err)
	}

	TruncateUserTable()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
	}

	return r
}

func TruncateUserTable() {

	db, err := config.InitDbSQLite()
	if err != nil {
		panic(err)
	}

	db.Exec("DELETE FROM user")
}

func TestRegisterUserSuccess(t *testing.T) {
	r := SetupRouterUser()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/register", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, "registrasi user berhasil", responseBody["message"])
	assert.Equal(t, "ilhamm.ms", responseBody["data"].(map[string]interface{})["username"])
}

func TestRegisterUserFailedUsernameEmpty(t *testing.T) {
	r := SetupRouterUser()

	reqBody := `{
		"username": "",
		"password": "ilhamsidiq"
	}`

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/register", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, "error : username dan password wajib diisi", responseBody["error"])
}

func TestRegisterUserFailedPasswordEmpty(t *testing.T) {
	r := SetupRouterUser()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": ""
	}`

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/register", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, "error : username dan password wajib diisi", responseBody["error"])
}

func TestRegisterUserFailedMinUsername(t *testing.T) {
	r := SetupRouterUser()

	reqBody := `{
		"username": "il",
		"password": "ilhamsidiq"
	}`

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/register", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, "error : username minimal 5 karakter", responseBody["error"])
}

func TestRegisterUserFailedMaxUsername(t *testing.T) {
	r := SetupRouterUser()

	reqBody := `{
		"username": "ilham muhammad sidiq sedang bermain bola",
		"password": "ilhamsidiq"
	}`

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/register", strings.NewReader(reqBody))

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, "error : username maksimal 20 karakter", responseBody["error"])
}

func TestRegisterUserFailedMinPassword(t *testing.T) {
	r := SetupRouterUser()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilham"
	}`

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/register", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, "error : harap masukkan password minimal 8 karakter", responseBody["error"])
}

func TestLoginFailed(t *testing.T) {
	r := SetupRouterUser()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/login", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, "error : username atau password salah", responseBody["error"])
}

func TestLoginSuccess(t *testing.T) {
	r := SetupRouterUser()

	reqBody := `{
		"username": "ilhamm.ms",
		"password": "ilhamsidiq"
	}`

	request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/register", strings.NewReader(reqBody))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	request = httptest.NewRequest(http.MethodPost, "http://localhost:8080/auth/login", strings.NewReader(reqBody))

	request.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	r.ServeHTTP(recorder, request)

	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", response.Header.Get("Content-Type"))
	assert.Equal(t, "login berhasil", responseBody["message"])
	assert.Equal(t, "ilhamm.ms", responseBody["data"].(map[string]interface{})["username"])
}
