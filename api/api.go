package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ilhaamms/library-api/controller"
)

type API struct {
	authorController controller.AuthorController
}

func NewAPI(authorController controller.AuthorController) *API {
	return &API{authorController: authorController}
}

func (a *API) RegisterRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/authors", a.authorController.CreateAuthor)
	r.GET("/authors", a.authorController.GetAuthors)
	r.GET("/authors/:id", a.authorController.GetAuthorsById)
	r.DELETE("/authors/:id", a.authorController.DeleteAuthorsById)
	r.PUT("/authors/:id", a.authorController.UpdateAuthorsById)

	return r
}

func (a *API) Run() {
	r := a.RegisterRoutes()
	r.Run(":8080")
}
