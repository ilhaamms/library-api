package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ilhaamms/library-api/controller"
	"github.com/ilhaamms/library-api/middleware"
)

type API struct {
	authorController controller.AuthorController
	userController   controller.UserController
}

func NewAPI(
	authorController controller.AuthorController,
	userController controller.UserController,
) *API {
	return &API{
		authorController: authorController,
		userController:   userController,
	}
}

func (a *API) RegisterRoutes() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", a.userController.Register)
		auth.POST("/login", a.userController.Login)
	}

	r.POST("/authors", middleware.Auth(), a.authorController.CreateAuthor)
	r.GET("/authors", middleware.Auth(), a.authorController.GetAuthors)
	r.GET("/authors/:id", middleware.Auth(), a.authorController.GetAuthorsById)
	r.DELETE("/authors/:id", middleware.Auth(), a.authorController.DeleteAuthorsById)
	r.PUT("/authors/:id", middleware.Auth(), a.authorController.UpdateAuthorsById)

	return r
}

func (a *API) Run() {
	r := a.RegisterRoutes()
	r.Run(":8080")
}
