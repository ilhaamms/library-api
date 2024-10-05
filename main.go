package main

import (
	"log"

	"github.com/ilhaamms/library-api/api"
	"github.com/ilhaamms/library-api/config"
	"github.com/ilhaamms/library-api/controller"
	"github.com/ilhaamms/library-api/repository"
	"github.com/ilhaamms/library-api/service"
)

func main() {
	db, err := config.InitDbSQLite()
	if err != nil {
		log.Fatal("Error connecting to database : ", err)
	}

	authorRepo := repository.NewAuthorRepository(db)

	authorService := service.NewAuthorService(authorRepo)

	authorController := controller.NewAuthorController(authorService)

	api := api.NewAPI(authorController)
	api.Run()
}
