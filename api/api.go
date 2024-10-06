package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ilhaamms/library-api/controller"
	"github.com/ilhaamms/library-api/middleware"
)

type API struct {
	authorController controller.AuthorController
	userController   controller.UserController
	bookController   controller.BookController
}

func NewAPI(
	authorController controller.AuthorController,
	userController controller.UserController,
	bookController controller.BookController,
) *API {
	return &API{
		authorController: authorController,
		userController:   userController,
		bookController:   bookController,
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
	r.GET("/authors", middleware.Auth(), a.authorController.GetAllAuthor)
	r.GET("/authors/:id", middleware.Auth(), a.authorController.GetAuthorsById)
	r.DELETE("/authors/:id", middleware.Auth(), a.authorController.DeleteAuthorsById)
	r.PUT("/authors/:id", middleware.Auth(), a.authorController.UpdateAuthorsById)

	r.POST("/books", middleware.Auth(), a.bookController.CreateBook)
	r.GET("/books", middleware.Auth(), a.bookController.GetAllBook)
	r.GET("/books/:id", middleware.Auth(), a.bookController.GetBookById)
	r.DELETE("/books/:id", middleware.Auth(), a.bookController.DeleteBookById)
	r.PUT("/books/:id", middleware.Auth(), a.bookController.Update)

	return r
}

func (a *API) Run() {
	r := a.RegisterRoutes()
	r.Run(":8080")
}
