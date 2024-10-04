package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ilhaamms/library-api/config"
)

func main() {
	db, err := config.InitDbSQLite()
	if err != nil {
		log.Fatal("Error connecting to database : ", err)
	}

	log.Println("Connected to database : ", db)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World from gin gaes",
		})
	})

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Error running server : ", err)
	}
}
