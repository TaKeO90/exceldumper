package main

import (
	"log"
	"os"

	"github.com/TaKeO90/exceldumper/server/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	err := os.Setenv("web", "true")
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()

	router.POST("/file/upload", handler.HandleRequest)
	router.POST("/file/download", handler.HandleRequest)
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"origin"},
	}))

	router.Run(":3000")
}
