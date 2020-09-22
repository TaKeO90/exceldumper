package main

import (
	"log"
	"os"

	"github.com/TaKeO90/exceldumper/server/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	err := os.Setenv("web", "true")
	port := os.Getenv("PORT")
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()

	router.POST("/file/upload", handler.HandleRequest)
	router.OPTIONS("/file/upload", handler.HandleRequest)
	router.POST("/file/download", handler.HandleRequest)
	router.OPTIONS("/file/download", handler.HandleRequest)

	router.Run(":" + port)
}
