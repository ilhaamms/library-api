package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilhaamms/library-api/entity/request"
	"github.com/ilhaamms/library-api/entity/response"
	"github.com/ilhaamms/library-api/service"
)

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{userService: userService}
}

func (uc *userController) Register(ctx *gin.Context) {

	var user request.User

	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	dataUser, err := uc.userService.Save(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response.WebResponseUser{
		StatusCode: http.StatusCreated,
		Message:    "registrasi user berhasil",
		Data:       dataUser,
	})

}

func (uc *userController) Login(ctx *gin.Context) {

	var user request.User

	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	isLogin, dataUser, err := uc.userService.Login(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Sprintf("error : %v", err.Error()),
		})
		return
	}

	if !isLogin {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      "username atau password salah",
		})
		return
	}

	cookie := &http.Cookie{
		Name:  "username",
		Value: dataUser.Username,
	}

	http.SetCookie(ctx.Writer, cookie)

	ctx.JSON(http.StatusOK, response.WebResponseUser{
		StatusCode: http.StatusOK,
		Message:    "login berhasil",
		Data:       dataUser,
	})
}
